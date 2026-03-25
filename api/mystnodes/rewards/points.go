package rewards

type Points struct {
	Items []string `json:"items"`
	Total float64 `json:"total,string"`
}