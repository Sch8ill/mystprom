package me

import "time"

type AccountInfo struct {
	User      User      `json:"user"`
	UserHash  string    `json:"userHash"`
	GTMEvents []string  `json:"gtmEvents"`
	Locals    Locals    `json:"locals"`
	NodesInfo NodesInfo `json:"nodesInfo"`
}

type User struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	EmailVerifiedAt       time.Time `json:"emailVerifiedAt"`
	ApiKey                string    `json:"apiKey"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	Admin                 bool      `json:"admin"`
	EmailMarketingConsent bool      `json:"emailMarketingConsent"`
	Archived              bool      `json:"archived"`
	WalletAddress         string    `json:"walletAddress"`
}

type Locals struct {
	UserId string `json:"userId"`
}

type NodesInfo struct {
	TotalCount  int `json:"totalCount"`
	OnlineCount int `json:"onlineCount"`
}
