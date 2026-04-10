package stats

type Stat struct {
	Name  string  `json:"name"`
	Value *Global `json:"value"`
}

type Global struct {
	TotalNodes     int `json:"totalNodes,string"`
	TotalTraffic   int `json:"totalTraffic"`
	TotalCountries int `json:"totalCountries,string"`
}
