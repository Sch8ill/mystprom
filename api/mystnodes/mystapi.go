// Package mystnodes provides a wrapper for the my.mystnodes.com api.
package mystnodes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sch8ill/mystprom/api/client"
	"github.com/sch8ill/mystprom/api/mystnodes/auth"
	"github.com/sch8ill/mystprom/api/mystnodes/me"
	"github.com/sch8ill/mystprom/api/mystnodes/nodes"
	"github.com/sch8ill/mystprom/api/mystnodes/notifications"
	"github.com/sch8ill/mystprom/api/mystnodes/totals"
)

const (
	BaseURL           = "https://my.mystnodes.com"
	LoginPath         = "/api/v2/auth/login"
	RefreshPath       = "/api/v2/auth/refresh"
	NodesPath         = "/api/v1/nodes"
	TotalsPath        = "/api/v1/metrics/node-totals"
	NotificationsPath = "/api/v1/me/notifications"
	AccountInfoPath   = "/api/v1/me"
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

func (m *MystAPI) Nodes() (*nodes.Nodes, error) {
	if err := m.authenticate(); err != nil {
		return nil, err
	}

	totalNodes := &nodes.Nodes{Total: 1}
	nodeCount := 0
	page := 1
	for nodeCount < totalNodes.Total {
		path := fmt.Sprintf("%s?page=%d&itemsPerPage=100", NodesPath, page)
		res, err := m.client.Get(path)
		if err != nil {
			return nil, err
		}

		nodeList := new(nodes.Nodes)
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
