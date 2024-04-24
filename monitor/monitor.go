package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/sch8ill/mystprom/api/cryptocompare"
	"github.com/sch8ill/mystprom/api/mystnodes"
	"github.com/sch8ill/mystprom/api/mystnodes/nodes"
	"github.com/sch8ill/mystprom/api/mystnodes/totals"
	"github.com/sch8ill/mystprom/config"
	"github.com/sch8ill/mystprom/metrics"
)

type Monitor struct {
	mystApi       *mystnodes.MystAPI
	cryptoCompare *cryptocompare.CryptoCompare
	interval      time.Duration
	stop          chan struct{}
	wg            sync.WaitGroup
}

func New(mystApi *mystnodes.MystAPI, cryptoCompare *cryptocompare.CryptoCompare, interval time.Duration) *Monitor {
	return &Monitor{
		mystApi:       mystApi,
		cryptoCompare: cryptoCompare,
		interval:      interval,
		stop:          make(chan struct{}),
	}
}

func (m *Monitor) Start() {
	log.Info().Msg("Starting monitor...")
	m.wg.Add(1)
	go m.run()
}

func (m *Monitor) Stop() {
	close(m.stop)
	m.wg.Wait()
}

func (m *Monitor) run() {
	defer m.wg.Done()

	for {
		select {
		case <-m.stop:
			return

		default:
			if err := m.monitorNodes(); err != nil {
				log.Warn().Err(err).Msg("failed to monitor")
			}
			if err := m.mystApi.RefreshToken().Save(config.RefreshFile); err != nil {
				log.Warn().Err(err).Msg("failed to save refresh token")
			}
			if err := m.updateMystPrices(); err != nil {
				log.Warn().Err(err).Msg("")
			}
			time.Sleep(m.interval)
		}
	}
}

func (m *Monitor) updateMystPrices() error {
	prices, err := m.cryptoCompare.MystPrices()
	if err != nil {
		return err
	}
	metrics.MystPrices(prices)
	return nil
}

func (m *Monitor) monitorNodes() error {
	n, err := m.mystApi.Nodes()
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	t, err := m.getTotals(listNodeIDs(n))
	if err != nil {
		return err
	}

	submitMetrics(n, t)
	return nil
}

func (m *Monitor) getTotals(ids []string) (map[string]*totals.Totals, error) {
	tMap := make(map[string]*totals.Totals)

	for _, id := range ids {
		t, err := m.mystApi.Totals([]string{id})
		if err != nil {
			return nil, fmt.Errorf("failed to get node totals for %s: %w", id, err)
		}
		tMap[id] = t
	}
	return tMap, nil
}

func submitMetrics(n *nodes.Nodes, t map[string]*totals.Totals) {
	metrics.NodeCount(n.Total)
	names := nodeNames(n.Nodes)
	for id, total := range t {
		metrics.NodeTotals(id, names[id], total)
	}

	for _, node := range n.Nodes {
		metrics.NodeMetrics(node)
	}
}

func nodeNames(n []nodes.Node) map[string]string {
	names := make(map[string]string)
	for _, node := range n {
		names[node.Identity] = node.Name
	}
	return names
}

func listNodeIDs(n *nodes.Nodes) (ids []string) {
	for _, node := range n.Nodes {
		ids = append(ids, node.Identity)
	}
	return
}
