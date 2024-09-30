package node

import (
	"encoding/json"
	"time"
)

type Node struct {
	ID               string        `json:"id"`
	UserID           string        `json:"userId"`
	TermsVersion     string        `json:"termsVersion,omitempty"`
	TermsAcceptedAt  time.Time     `json:"termsAcceptedAt,omitempty"`
	Whitelist        string        `json:"whitelist,omitempty"`
	LocalIP          string        `json:"localIp"`
	ExternalIP       string        `json:"externalIp"`
	ISP              string        `json:"isp,omitempty"`
	OS               string        `json:"os"`
	Arch             string        `json:"arch"`
	Version          string        `json:"version"`
	Vendor           string        `json:"vendor"`
	Identity         string        `json:"identity"`
	Malicious        bool          `json:"malicious"`
	Name             string        `json:"name"`
	AvailableAt      time.Time     `json:"availableAt,omitempty"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	Deleted          bool          `json:"deleted"`
	LauncherVersion  string        `json:"launcherVersion,omitempty"`
	IPTagged         bool          `json:"ipTagged"`
	NodeStatus       Status        `json:"nodeStatus"`
	MonitoringStatus string        `json:"monitoringStatus"`
	Earnings         []Earnings    `json:"earnings"`
	UptimeLast24h    time.Duration `json:"uptimeMinLast24H"`
}

func (n *Node) UnmarshalJSON(data []byte) error {
	type Alias Node
	aux := &struct {
		UptimeLast24h int `json:"uptimeMinLast24H"`
		*Alias
	}{
		Alias: (*Alias)(n),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	n.UptimeLast24h = time.Minute * time.Duration(aux.UptimeLast24h)

	return nil
}

type Status struct {
	ID                     string    `json:"id"`
	NodeID                 string    `json:"nodeId"`
	MonitoringFailed       bool      `json:"monitoringFailed"`
	MonitoringFailedLastAt time.Time `json:"monitoringFailedLastAt"`
	Online                 bool      `json:"online"`
	OnlineLastAt           time.Time `json:"onlineLastAt"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	IPCategory             string    `json:"ipCategory"`
	Location               string    `json:"location"`
	Quality                float64   `json:"quality"`
	ServiceTypes           []string  `json:"serviceTypes"`
}

type Earnings struct {
	Service     string  `json:"service"`
	EtherAmount float64 `json:"etherAmount,string"`
}
