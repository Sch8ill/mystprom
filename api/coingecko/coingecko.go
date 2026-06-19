// Package coingecko provides a wrapper for the coingecko api.
package coingecko

import (
	"encoding/json"
	"fmt"
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

	prices := make(map[string]map[string]float64)
	if err := json.NewDecoder(res.Body).Decode(&prices); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if _, ok := prices[MystID]; !ok {
		return nil, fmt.Errorf("response does not contain: %s", MystID)
	}

	// key to upper
	for key, val := range prices[MystID] {
		prices[MystID][strings.ToUpper(key)] = val
		delete(prices[MystID], key)
	}

	prices[MystID][MystSymbol] = 1

	return prices[MystID], nil
}
