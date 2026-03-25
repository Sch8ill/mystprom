package rewards

type Stats struct {
	Data   []float64 `json:"data"`
	Myst   []float64 `json:"myst"`
	Uptime []float64 `json:"uptime"`
	Nodes  []int     `json:"nodes"`
}
