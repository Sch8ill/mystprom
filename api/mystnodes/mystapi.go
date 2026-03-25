// Package mystnodes provides a wrapper for the my.mystnodes.com api.
package mystnodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/sch8ill/mystprom/api/client"
	"github.com/sch8ill/mystprom/api/mystnodes/auth"
	"github.com/sch8ill/mystprom/api/mystnodes/me"
	"github.com/sch8ill/mystprom/api/mystnodes/node"
	"github.com/sch8ill/mystprom/api/mystnodes/notifications"
	"github.com/sch8ill/mystprom/api/mystnodes/rewards"
	"github.com/sch8ill/mystprom/api/mystnodes/totals"
)

const (
	BaseURL           = "https://my.mystnodes.com"
	LoginPath         = "/api/v2/auth/login"
	RefreshPath       = "/api/v2/auth/refresh"
	NodePath          = "/api/v2/node"
	TotalsPath        = "/api/v1/metrics/node-totals"
	NotificationsPath = "/api/v1/me/notifications"
	RewardClaimPath   = "/api/v2/reward-program/reward/claim"
	RewardRanksPath   = "/api/v2/reward-program/ranks"
	RewardPointsPath  = "/api/v2/reward-program/points"
	RewardStatsPath   = "/api/v2/reward-program/stats"
	AccountInfoPath   = "/api/v2/me"
)

type Credentials struct {
	Email    string
	Password string
}

type MystAPI struct {
	client       *client.HttpClient
	credentials  Credentials
	token        *Token
	refreshToken *Token
}

func New(credentials Credentials) *MystAPI {
	return NewWithRefreshToken(credentials, nil)
}

func NewWithRefreshToken(credentials Credentials, refreshToken *Token) *MystAPI {
	c := client.New(BaseURL)
	// required by the api to parse post requests correctly
	c.SetHeader("Content-Type", "application/json")
	c.SetHeader("Accept", "application/json")

	return &MystAPI{
		client:       c,
		credentials:  credentials,
		refreshToken: refreshToken,
	}
}

func (m *MystAPI) Login() error {
	loginPost := auth.LoginRequest{
		Email:      m.credentials.Email,
		Password:   m.credentials.Password,
		RememberMe: true,
	}

	res, err := m.client.PostJSON(LoginPath, loginPost)
	if err != nil {
		return fmt.Errorf("failed to post login credentials: %w", err)
	}

	loginRes := new(auth.LoginResponse)
	if err := m.parseResponse(res, loginRes); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	m.setToken(NewToken(loginRes.AccessToken, loginRes.AccessTokenTTL))
	m.refreshToken = NewToken(loginRes.RefreshToken, loginRes.RefreshTokenTTL)

	return nil
}

func (m *MystAPI) Refresh() error {
	res, err := m.client.PostJSON(RefreshPath, auth.RefreshRequest{RefreshToken: m.refreshToken.Value()})
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	refreshRes := new(auth.RefreshResponse)
	if err := m.parseResponse(res, refreshRes); err != nil {
		return fmt.Errorf("failed to parse token refresh response: %w", err)
	}
	m.setToken(NewToken(refreshRes.AccessToken, refreshRes.AccessTokenTTL))

	return nil
}

func (m *MystAPI) Nodes() (*node.Nodes, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	totalNodes := &node.Nodes{Total: 1}
	nodeCount := 0
	page := 1
	for nodeCount < totalNodes.Total {
		path := fmt.Sprintf("%s?page=%d&itemsPerPage=100", NodePath, page)
		res, err := m.client.Get(path)
		if err != nil {
			return nil, err
		}

		nodeList := new(node.Nodes)
		if err := m.parseResponse(res, nodeList); err != nil {
			return nil, err
		}

		totalNodes.Nodes = append(totalNodes.Nodes, nodeList.Nodes...)
		totalNodes.Total = nodeList.Total
		nodeCount += len(nodeList.Nodes)
		page++
	}

	return totalNodes, nil
}

func (m *MystAPI) Node(identity string) (*node.Node, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", NodePath, identity)
	res, err := m.client.Get(path)
	if err != nil {
		return nil, err
	}

	n := new(node.Node)
	if err := m.parseResponse(res, n); err != nil {
		return nil, err
	}

	return n, nil
}

func (m *MystAPI) Sessions(identity string) ([]node.Session, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s/sessions", NodePath, identity)
	res, err := m.client.Get(path)
	if err != nil {
		return nil, err
	}

	var sessions []node.Session
	if err := m.parseResponse(res, &sessions); err != nil {
		if errors.Is(err, io.EOF) {
			return []node.Session{}, nil
		}
		return nil, err
	}

	return sessions, nil
}

func (m *MystAPI) Totals(identities []string) (*totals.Totals, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	// the 'days' parameter is disregarded by the API, it consistently returns metrics
	// for the last 30 days.
	path := fmt.Sprintf("%s?days=30&identities=%s", TotalsPath, strings.Join(identities, "%2C"))
	res, err := m.client.Get(path)
	if err != nil {
		return nil, err
	}

	nodeTotals := new(totals.Response)
	if err := m.parseResponse(res, nodeTotals); err != nil {
		return nil, err
	}

	return &nodeTotals.Totals, nil
}

func (m *MystAPI) AccountInfo() (*me.AccountInfo, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	res, err := m.client.Get(AccountInfoPath)
	if err != nil {
		return nil, err
	}

	accountInfo := new(me.AccountInfo)
	if err := m.parseResponse(res, accountInfo); err != nil {
		return nil, err
	}

	return accountInfo, nil
}

func (m *MystAPI) Notifications() ([]notifications.Notification, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	res, err := m.client.Get(NotificationsPath)
	if err != nil {
		return nil, err
	}

	n := new(notifications.Response)
	if err := m.parseResponse(res, n); err != nil {
		return nil, err
	}

	return n.Notifications, nil
}

func (m *MystAPI) RewardPoints() (*rewards.Points, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	res, err := m.client.Get(RewardPointsPath)
	if err != nil {
		return nil, err
	}

	p := new(rewards.Points)
	if err := m.parseResponse(res, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (m *MystAPI) RewardStats() (*rewards.Stats, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	type stats struct {
		Data   []string `json:"data"`
		Myst   []string `json:"myst"`
		Uptime []string `json:"uptime"`
		Nodes  []int    `json:"nodes"`
	}

	res, err := m.client.Get(RewardStatsPath)
	if err != nil {
		return nil, err
	}

	rawStats := new(stats)
	if err := m.parseResponse(res, rawStats); err != nil {
		return nil, err
	}

	s := new(rewards.Stats)
	s.Data, err = parseFloatList(rawStats.Data)
	if err != nil {
		return nil, err
	}

	s.Myst, err = parseFloatList(rawStats.Myst)
	if err != nil {
		return nil, err
	}

	s.Uptime, err = parseFloatList(rawStats.Uptime)
	if err != nil {
		return nil, err
	}

	s.Nodes = rawStats.Nodes

	return s, nil
}

func (m *MystAPI) RewardRanks() ([]rewards.User, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	const limit int = 100
	var ranks []rewards.User
	for i := 1; ; i++ {
		res, err := m.client.Get(RewardRanksPath + fmt.Sprintf("?page=%d&limit=%d", i, limit))
		if err != nil {
			return nil, err
		}

		r := new(rewards.Ranks)
		if err := m.parseResponse(res, r); err != nil {
			return nil, err
		}

		ranks = append(ranks, r.Items...)

		// total reached
		if r.Total < i*limit {
			break
		}
	}

	return ranks, nil
}

func (m *MystAPI) RefreshToken() *Token {
	return m.refreshToken
}

func (m *MystAPI) authenticate() error {
	if m.refreshToken == nil {
		return m.Login()
	} else if m.refreshToken.Expired() {
		return m.Login()
	}

	if m.token == nil {
		return m.Refresh()
	} else if m.token.Expired() {
		return m.Refresh()
	}

	return nil
}

func (m *MystAPI) setToken(token *Token) {
	m.token = token
	m.client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token.Value()))
}

func (m *MystAPI) parseResponse(res *http.Response, target any) error {
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 300 {
		if err := m.handleApiError(res); err != nil {
			return err
		}
	}

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (m *MystAPI) handleApiError(res *http.Response) error {
	errRes := new(Error)
	if err := json.NewDecoder(res.Body).Decode(errRes); err != nil {
		return fmt.Errorf("failed to parse error message for status code: %d:%w", res.StatusCode, err)
	}

	if errRes.ErrorCode == "expiredAccessToken" {
		if err := m.Refresh(); err != nil {
			return fmt.Errorf("access token expired: %w", err)
		}
	}

	if errRes.ErrorCode == "expiredRefreshToken" {
		if err := m.Login(); err != nil {
			return fmt.Errorf("refresh token expired: %w", err)
		}
	}

	return fmt.Errorf("api error response: %d: %+v", res.StatusCode, errRes)
}

func parseFloatList(list []string) ([]float64, error) {
	parsed := []float64{}

	for _, r := range list {
		v, err := strconv.ParseFloat(r, 64)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, v)
	}

	return parsed, nil
}
