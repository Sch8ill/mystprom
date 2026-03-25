package rewards

type Ranks struct {
	Items []User
	Limit int
	Total int
}

type User struct {
	Address      string  `json:"address"`
	ActiveNodes  int     `json:"activeNodes,string"`
	PointsData   float64 `json:"pointsData,string"`
	PointsMyst   float64 `json:"pointsMyst,string"`
	PointsTotal  float64 `json:"pointsTotal,string"`
	PointsUptime float64 `json:"pointsUptime,string"`
	Current      bool    `json:"current"`
}
