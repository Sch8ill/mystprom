# mystprom

[![Release](https://img.shields.io/github/release/sch8ill/mystprom.svg?style=flat-square)](https://github.com/sch8ill/mystprom/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/sch8ill/mystprom)](https://goreportcard.com/report/github.com/sch8ill/mystprom)
![MIT license](https://img.shields.io/badge/license-MIT-green)

---

`mytsprom` is a Prometheus exporter for monitoring Mysterium Network nodes using the my.mystnodes.com api.  
`mystprom` offers all metrics found on my.mystnodes.com and more.  
The monitored metrics include a wide range of vital statistics, including:

- `Node Online Status`: Provides real-time information on whether a node is actively online.
- `Earning Statistics`: Offers insights into earnings generated by each node and service.
- `Traffic Statistics`: Tracks data transfer activities.
- `Bandwidth Statistics`: Delivers detailed metrics on available internet bandwidth.
- `Version Monitoring`: Keeps track of software versioning for nodes.

---

## Installation

### Docker

```bash
docker run -p 9300:9300 -e MYSTPROM_EMAIL="" -e MYSTPROM_PASSWORD="" sch8ill/mystprom:latest
```

### Build

Requires:

```
go >= 1.22
make
```

Build command:

```bash
make build
```

---

## Usage

### Grafana

An [example Grafana panel](https://github.com/sch8ill/mystprom/blob/master/grafana/dashboard.json) can be found in
the [grafana directory](https://github.com/sch8ill/mystprom/blob/master/grafana).

### Prometheus config

Example `prometheus.yml` scrape config:

```yaml
scrape_configs:
  - job_name: mystprom
    scrape_interval: 5m
    static_configs:
      - targets: [ "localhost:9300" ]
```

### Metrics

| name                                | description                                           | labels             | type         |
|-------------------------------------|-------------------------------------------------------|--------------------|--------------|
| myst_node_bandwidth                 | Internet bandwidth of the node                        | id, name           | mbit/s       |
| myst_node_traffic                   | Traffic transferred by the node over the last 30 days | id, name, service, country | gigabytes    |
| myst_node_quality                   | Quality score assigned to the node                    | id, name           | float        |
| myst_node_service                   | whether a service on the node is running              | id, name, service  | boolean      |
| myst_node_earnings                  | Earnings by service of node over the last 30 days     | id, name, service  | MYST         |
| myst_node_earnings_lifetime         | Total lifetime earnings by node                       | id, name           | MYST         |
| myst_node_earnings_settled          | Total settled earnings by node                        | id, name           | MYST         |
| myst_node_earnings_unsettled        | Unsettled earnings by node                            | id, name           | MYST         |
| myst_node_sessions                  | Number of sessions of the node over the last 30 days  | id, name, service, country | int  |
| myst_node_session_earings           | Earnings by node, generated from session log          | id, name, service, country | MYST |
| myst_node_session_durations         | Total duration of sessions over the last 30 days      | id, name, service, country | seconds |
| myst_token_price                    | Current price of the MYST token                       | currency           | EUR/USD      |
| myst_node_location                  | Location of the node                                  | id, name, location | country code |
| myst_node_external_ip               | External ip address of the node                       | id, name, ip       | ip           |
| myst_node_local_ip                  | Local ip address of the node                          | id, name, ip       | ip           |
| myst_node_isp                       | Internet Service Provider of the node                 | id, name, isp      |              |
| myst_node_os                        | Operating system the node is running on               | id, name, os       | os           |
| myst_node_arch                      | System architecture of the node                       | id, name, arch     | architecture |
| myst_node_version                   | Myst version the node is running on                   | id, name, version  | version      |
| myst_node_launcher_version          | Launcher version the node is running on               | id, name, version  | version      |
| myst_node_vendor                    | Vendor of the node                                    | id, name, vendor   |              |
| myst_node_updated_at                | Last time the node was updated                        | id, name           | unix time    |
| myst_node_ip_category               | IP category of the node                               | id, name, category | ip category  |
| myst_node_malicious                 | whether the node is tagged a malicious                | id, name           | boolean      |
| myst_node_ip_tagged                 | whether the node is ip tagged                         | id, name           | boolean      |
| myst_node_online                    | whether the node is online                            | id, name           | boolean      |
| myst_node_online_last_at            | Last time the node was online                         | id, name           | unix time    |
| myst_node_monitoring_status         | Monitoring status of the node                         | id, name, status   |              |
| myst_node_monitoring_failed         | whether monitoring on the node failed                 | id, name           | boolean      |
| myst_node_monitoring_failed_last_at | Last time monitoring failed on node                   | id, name           | unix time    |
| myst_node_available_at              | Last time the node was available                      | id, name           | unix time    |
| myst_node_status_created_at         | Time the node monitoring record was created           | id, name           | unix time    |
| myst_node_status_updated_at         | Last time the node status was updated                 | id, name           | unix time    |
| myst_node_created_at                | Time the node was created                             | id, name           | unix time    |
| myst_node_terms_version             | Terms version of the node                             | id, name, version  | gauge        |
| myst_node_terms_accepted_at         | Last time terms were accepted by node                 | id, name           | unix time    |
| myst_node_count                     | Total number of nodes                                 |                    | gauge        |
| myst_node_user_id                   | User ID of user of the node                           | id, name, user_id  |              |
| myst_node_deleted                   | whether the node is deleted                           | id, name           | boolean      |

### CLI flags

```
   --email value, -m value     email address of the my.mystnodes.com account [$MYSTPROM_EMAIL]
   --password value, -p value  password of the my.mystnodes.com account [$MYSTPROM_PASSWORD]
   --interval value, -i value  interval the Mysterium Network api should be scraped in (default: 10m0s) [$MYSTPROM_INTERVAL]
   --metrics-address value     address the Prometheus metrics exporter listens on (default: ":9300") [$MYSTPROM_METRICS_ADDRESS]
   --refresh-file value        name of the file the refresh token is stored in (default: ".refresh_token.json") [$MYSTPROM_REFRESH_FILE]
   --help, -h                  show help
```

## License

This package is licensed under the [MIT License](LICENSE).
