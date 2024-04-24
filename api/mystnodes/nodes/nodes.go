package nodes

import "time"

type Nodes struct {
	Nodes []Node `json:"nodes"`
	Total int    `json:"total"`
}

type Node struct {
	ID               string     `json:"id"`
	UserID           string     `json:"userId"`
	TermsVersion     string     `json:"termsVersion,omitempty"`
	TermsAcceptedAt  time.Time  `json:"termsAcceptedAt,omitempty"`
	Whitelist        string     `json:"whitelist,omitempty"`
	LocalIP          string     `json:"localIp"`
	ExternalIP       string     `json:"externalIp"`
	ISP              string     `json:"isp,omitempty"`
	OS               string     `json:"os"`
	Arch             string     `json:"arch"`
	Version          string     `json:"version"`
	Vendor           string     `json:"vendor"`
	Identity         string     `json:"identity"`
	Malicious        bool       `json:"malicious"`
	Name             string     `json:"name"`
	AvailableAt      time.Time  `json:"availableAt,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	Deleted          bool       `json:"deleted"`
	LauncherVersion  string     `json:"launcherVersion,omitempty"`
	IPTagged         bool       `json:"ipTagged"`
	NodeStatus       NodeStatus `json:"nodeStatus"`
	MonitoringStatus string     `json:"monitoringStatus"`
	Earnings         []Earnings `json:"earnings"`
}

type NodeStatus struct {
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
