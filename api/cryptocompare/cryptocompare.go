// Package cryptocompare provides a wrapper for the cryptocompare api.
package cryptocompare

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sch8ill/mystprom/api/client"
)

const (
	BaseURL   = "https://min-api.cryptocompare.com"
	PricePath = "/data/price"

	MystSymbol = "MYST"
)

var Currencies = []string{"EUR", "USD"}

type CryptoCompare struct {
	client *client.HttpClient
}

func New() (*CryptoCompare, error) {
	c := client.New(BaseURL)
	c.SetHeader("Content-Type", "application/json")
	c.SetHeader("Accept", "application/json")

	return &CryptoCompare{client: c}, nil
}

func (c *CryptoCompare) Prices(symbol string, currencies []string) (map[string]float64, error) {
	path := fmt.Sprintf("%s?fsym=%s&tsyms=%s", PricePath, symbol, strings.Join(currencies, ","))
	res, err := c.client.Get(path)
	if err != nil {
		return nil, err
	}

	prices := make(map[string]float64)
	if err := json.NewDecoder(res.Body).Decode(&prices); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	prices[MystSymbol] = 1

	return prices, nil
}

func (c *CryptoCompare) MystPrices() (map[string]float64, error) {
	prices, err := c.Prices(MystSymbol, Currencies)
	if err != nil {
		return nil, err
	}
	return prices, nil
}
