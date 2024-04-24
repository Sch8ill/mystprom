package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slices"

	"github.com/sch8ill/mystprom/api/mystnodes/nodes"
	"github.com/sch8ill/mystprom/api/mystnodes/totals"
)

var nodeCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "myst_node_count",
	Help: "Total number of nodes",
})

var nodeBandwidth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_bandwidth",
	Help: "Internet bandwidth of the node",
}, []string{"id", "name"})

var nodeTraffic = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_traffic",
	Help: "Traffic transferred by the node over the last 30 days",
}, []string{"id", "name"})

var nodeUserID = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_user_id",
	Help: "User ID of user of the node",
}, []string{"id", "name", "user_id"})

var nodeTermsVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_terms_version",
	Help: "Terms version of the node",
}, []string{"id", "name", "version"})

var nodeTermsAcceptedAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_terms_accepted_at",
	Help: "Last time terms were accepted by node",
}, []string{"id", "name"})

var nodeLocalIP = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_local_ip",
	Help: "Local ip address of the node",
}, []string{"id", "name", "ip"})

var nodeExternalIP = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_external_ip",
	Help: "External ip address of the node",
}, []string{"id", "name", "ip"})

var nodeISP = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_isp",
	Help: "Internet Service Provider of the node",
}, []string{"id", "name", "isp"})

var nodeOS = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_os",
	Help: "Operating system the node is running on",
}, []string{"id", "name", "os"})

var nodeArch = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_arch",
	Help: "System architecture of the node",
}, []string{"id", "name", "arch"})

var nodeVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_version",
	Help: "Myst version the node is running on",
}, []string{"id", "name", "version"})

var nodeVendor = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_vendor",
	Help: "Vendor of the node",
}, []string{"id", "name", "vendor"})

var nodeMalicious = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_malicious",
	Help: "whether the node is tagged a malicious",
}, []string{"id", "name"})

var nodeAvailableAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_available_at",
	Help: "Last time the node was available",
}, []string{"id", "name"})

var nodeCreatedAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_created_at",
	Help: "Time the node was created",
}, []string{"id", "name"})

var nodeUpdatedAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_updated_at",
	Help: "Last time the node was updated",
}, []string{"id", "name"})

var nodeDeleted = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_deleted",
	Help: "whether the node is deleted",
}, []string{"id", "name"})

var nodeLauncherVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_launcher_version",
	Help: "Launcher version the node is running on",
}, []string{"id", "name", "version"})

var nodeIPTagged = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_ip_tagged",
	Help: "whether the node is ip tagged",
}, []string{"id", "name"})

var nodeMonitoringFailed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_monitoring_failed",
	Help: "whether monitoring on the node failed",
}, []string{"id", "name"})

var nodeMonitoringFailedLastAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_monitoring_failed_last_at",
	Help: "Last time monitoring failed on node",
}, []string{"id", "name"})

var nodeOnline = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_online",
	Help: "whether the node is online",
}, []string{"id", "name"})

var nodeOnlineLastAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_online_last_at",
	Help: "Last time the node was online",
}, []string{"id", "name"})

var nodeStatusCreatedAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_status_created_at",
	Help: "Time the node monitoring record was created",
}, []string{"id", "name"})

var nodeStatusUpdatedAt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_status_updated_at",
	Help: "Last time the node status was updated",
}, []string{"id", "name"})

var nodeIPCategory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_ip_category",
	Help: "IP category of the node",
}, []string{"id", "name", "category"})

var nodeLocation = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_location",
	Help: "Location of the node",
}, []string{"id", "name", "location"})

var nodeQuality = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_quality",
	Help: "Quality score assigned to the node",
}, []string{"id", "name"})

var nodeService = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_service",
	Help: "whether a service on the node is running",
}, []string{"id", "name", "service"})

var nodeMonitoringStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_monitoring_status",
	Help: "Monitoring status of the node",
}, []string{"id", "name", "status"})

var nodeServiceEarnings = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_node_service_earnings",
	Help: "Earnings by service of node over the last 30 days",
}, []string{"id", "name", "service"})

var mystPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "myst_token_price",
	Help: "Current price of the MYST token",
}, []string{"currency"})

func init() {
	registry.MustRegister(nodeCount, nodeBandwidth, nodeTraffic, nodeUserID, nodeTermsVersion, nodeTermsAcceptedAt,
		nodeLocalIP, nodeExternalIP, nodeISP, nodeOS, nodeArch, nodeVersion, nodeVendor, nodeMalicious,
		nodeAvailableAt, nodeCreatedAt, nodeUpdatedAt, nodeDeleted, nodeLauncherVersion, nodeIPTagged,
		nodeMonitoringFailed, nodeMonitoringFailedLastAt, nodeOnline, nodeOnlineLastAt, nodeStatusCreatedAt,
		nodeStatusUpdatedAt, nodeIPCategory, nodeLocation, nodeQuality, nodeService, nodeMonitoringStatus,
		nodeServiceEarnings, mystPrice)
}

func NodeCount(n int) {
	nodeCount.Set(float64(n))
}

func NodeTotals(id string, name string, t *totals.Totals) {
	nodeBandwidth.WithLabelValues(id, name).Set(t.BandwidthTotal)
	nodeTraffic.WithLabelValues(id, name).Set(t.TrafficTotal)
}

func NodeMetrics(node nodes.Node) {
	nodeUserID.WithLabelValues(node.Identity, node.Name, node.UserID).Set(1)
	nodeTermsVersion.WithLabelValues(node.Identity, node.Name, node.TermsVersion).Set(1)
	nodeTermsAcceptedAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.TermsAcceptedAt.Unix()))

	nodeLocalIP.WithLabelValues(node.Identity, node.Name, node.LocalIP).Set(1)
	nodeExternalIP.WithLabelValues(node.Identity, node.Name, node.ExternalIP).Set(1)
	nodeIPCategory.WithLabelValues(node.Identity, node.Name, node.NodeStatus.IPCategory).Set(1)
	nodeIPTagged.WithLabelValues(node.Identity, node.Name).Set(boolToFloat(node.IPTagged))
	nodeISP.WithLabelValues(node.Identity, node.Name, node.ISP).Set(1)
	nodeLocation.WithLabelValues(node.Identity, node.Name, node.NodeStatus.Location).Set(1)

	nodeOS.WithLabelValues(node.Identity, node.Name, node.OS).Set(1)
	nodeArch.WithLabelValues(node.Identity, node.Name, node.Arch).Set(1)
	nodeVersion.WithLabelValues(node.Identity, node.Name, node.Version).Set(1)
	nodeLauncherVersion.WithLabelValues(node.Identity, node.Name, node.LauncherVersion).Set(1)
	nodeVendor.WithLabelValues(node.Identity, node.Name, node.Vendor).Set(1)
	nodeMalicious.WithLabelValues(node.Identity, node.Name).Set(boolToFloat(node.Malicious))

	nodeAvailableAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.AvailableAt.Unix()))
	nodeCreatedAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.CreatedAt.Unix()))
	nodeUpdatedAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.UpdatedAt.Unix()))
	nodeDeleted.WithLabelValues(node.Identity, node.Name).Set(boolToFloat(node.Deleted))

	nodeMonitoringStatus.WithLabelValues(node.Identity, node.Name, node.MonitoringStatus).Set(1)
	nodeMonitoringFailed.WithLabelValues(node.Identity, node.Name).Set(boolToFloat(node.NodeStatus.MonitoringFailed))
	nodeMonitoringFailedLastAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.NodeStatus.MonitoringFailedLastAt.Unix()))
	nodeOnline.WithLabelValues(node.Identity, node.Name).Set(boolToFloat(node.NodeStatus.Online))
	nodeOnlineLastAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.NodeStatus.OnlineLastAt.Unix()))
	nodeStatusCreatedAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.NodeStatus.CreatedAt.Unix()))
	nodeStatusUpdatedAt.WithLabelValues(node.Identity, node.Name).Set(float64(node.NodeStatus.UpdatedAt.Unix()))
	nodeQuality.WithLabelValues(node.Identity, node.Name).Set(node.NodeStatus.Quality)

	for _, earnings := range node.Earnings {
		nodeServiceEarnings.WithLabelValues(node.Identity, node.Name, earnings.Service).Set(earnings.EtherAmount)
		nodeService.WithLabelValues(node.Identity, node.Name, earnings.Service).Set(
			boolToFloat(slices.Contains(node.NodeStatus.ServiceTypes, earnings.Service)))
	}
}

func MystPrices(prices map[string]float64) {
	for currency, price := range prices {
		mystPrice.WithLabelValues(currency).Set(price)
	}
}

// https://github.com/golang/go/issues/64825
func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
