// Package coingecko provides a wrapper for the coingecko api.
package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sch8ill/mystprom/api/client"
)

const (
	BaseURL = "https://api.coingecko.com"
	Path    = "/api/v3/simple/price?ids=mysterium&vs_currencies=EUR,USD"

	MystSymbol = "MYST"
	MystID     = "mysterium"
)

type Coingecko struct {
	client *client.HttpClient
}

func New() (*Coingecko, error) {
	c, err := client.New(BaseURL)
	if err != nil {
		return nil, err
	}

	c.SetHeader("Content-Type", "application/json")
	c.SetHeader("Accept", "application/json")

	return &Coingecko{client: c}, nil
}

func (c *Coingecko) MystPrices() (map[string]float64, error) {
	res, err := c.client.Get(Path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data := make(map[string]map[string]float64)
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if _, ok := data[MystID]; !ok {
		return nil, fmt.Errorf("response does not contain: %s", MystID)
	}

	// key to upper
	prices := make(map[string]float64)
	for key, val := range data[MystID] {
		prices[strings.ToUpper(key)] = val
	}
	prices[MystSymbol] = 1

	return prices, nil
}
