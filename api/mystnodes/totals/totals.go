package totals

type Response struct {
	Totals Totals `json:"totals"`
}

type Totals struct {
	BandwidthTotal float64 `json:"bandwidthTotal"`
	TrafficTotal   float64 `json:"trafficTotal"`
}
