package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/sch8ill/mystprom/api/cryptocompare"
	"github.com/sch8ill/mystprom/api/mystnodes"
	"github.com/sch8ill/mystprom/api/mystnodes/node"
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

func (m *Monitor) monitorNodes() error {
	nodes, err := m.mystApi.Nodes()
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	ids := listNodeIDs(nodes)

	sessions, err := m.getSessions(ids)
	if err != nil {
		return err
	}

	lifetimeEarnings, err := m.getLifetimeEarnings(ids)
	if err != nil {
		return err
	}

	t, err := m.getTotals(ids)
	if err != nil {
		return err
	}

	submitMetrics(nodes, sessions, lifetimeEarnings, t)
	return nil
}

func (m *Monitor) getLifetimeEarnings(ids []string) (map[string]node.LifetimeEarnings, error) {
	earningsMap := make(map[string]node.LifetimeEarnings)

	for _, id := range ids {
		n, err := m.mystApi.Node(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get lifetime earnings for %s: %w", id, err)
		}
		earningsMap[id] = *n.LifetimeEarnings
	}

	return earningsMap, nil
}

func (m *Monitor) getSessions(ids []string) (map[string][]node.Session, error) {
	sessionMap := make(map[string][]node.Session)

	for _, id := range ids {
		sessions, err := m.mystApi.Sessions(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get sessions for %s: %w", id, err)
		}
		sessionMap[id] = sessions
	}

	return sessionMap, nil
}

func (m *Monitor) getTotals(ids []string) (map[string]*totals.Totals, error) {
	totalsMap := make(map[string]*totals.Totals)

	for _, id := range ids {
		t, err := m.mystApi.Totals([]string{id})
		if err != nil {
			return nil, fmt.Errorf("failed to totals for %s: %w", id, err)
		}
		totalsMap[id] = t
	}

	return totalsMap, nil
}

func (m *Monitor) updateMystPrices() error {
	prices, err := m.cryptoCompare.MystPrices()
	if err != nil {
		return err
	}
	metrics.MystPrices(prices)
	return nil
}

func submitMetrics(nodes *node.Nodes, sessions map[string][]node.Session, lifetimeEarnings map[string]node.LifetimeEarnings, t map[string]*totals.Totals) {
	metrics.NodeCount(nodes.Total)

	for _, node := range nodes.Nodes {
		metrics.NodeMetrics(node)
	}

	names := nodeNames(nodes.Nodes)
	for id, s := range sessions {
		metrics.NodeSessions(id, names[id], s)
	}

	for id, earnings := range lifetimeEarnings {
		metrics.NodeLifetimeEarnings(id, names[id], earnings)
	}

	for id, nodeTotals := range t {
		metrics.NodeTotals(id, names[id], nodeTotals)
	}
}

func nodeNames(n []node.Node) map[string]string {
	names := make(map[string]string)
	for _, node := range n {
		names[node.Identity] = node.Name
	}
	return names
}

func listNodeIDs(n *node.Nodes) (ids []string) {
	for _, node := range n.Nodes {
		ids = append(ids, node.Identity)
	}
	return
}
