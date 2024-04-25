package config

import (
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

const (
	DefaultScrapeInterval = time.Minute * 10
	DefaultMetricsAddress = ":9300"
	DefaultRefreshFile    = ".refresh_token.json"

	MystAPIEmailFlag    = "email"
	MystAPIPasswordFlag = "password"
	ScrapeIntervalFlag  = "interval"
	MetricsAddressFlag  = "metrics-address"
	RefreshFileFlag     = "refresh-file"
)

var (
	MystAPIEmail    string
	MystAPIPassword string
	ScrapeInterval  time.Duration
	MetricsAddress  string
	RefreshFile     string
)

func DeclareFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     MystAPIEmailFlag,
			Usage:    "email address of the my.mystnodes.com account",
			Required: true,
			Aliases:  []string{"m"},
			EnvVars:  []string{"MYSTPROM_EMAIL"},
		},
		&cli.StringFlag{
			Name:     MystAPIPasswordFlag,
			Usage:    "password of the my.mystnodes.com account",
			Required: true,
			Aliases:  []string{"p"},
			EnvVars:  []string{"MYSTPROM_PASSWORD"},
		},
		&cli.DurationFlag{
			Name:    ScrapeIntervalFlag,
			Usage:   "interval the api should be scraped in",
			Value:   DefaultScrapeInterval,
			Aliases: []string{"i"},
			EnvVars: []string{"MYSTPROM_INTERVAL"},
		},
		&cli.StringFlag{
			Name:    MetricsAddressFlag,
			Usage:   "address the Prometheus metrics exporter listens on",
			Value:   DefaultMetricsAddress,
			EnvVars: []string{"MYSTPROM_METRICS_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    RefreshFileFlag,
			Usage:   "name of the file the refresh token is stored in",
			Value:   DefaultRefreshFile,
			EnvVars: []string{"MYSTPROM_REFRESH_FILE"},
		},
	}
}

func SetConfig(ctx *cli.Context) {
	MystAPIEmail = ctx.String(MystAPIEmailFlag)
	MystAPIPassword = ctx.String(MystAPIPasswordFlag)
	ScrapeInterval = ctx.Duration(ScrapeIntervalFlag)
	MetricsAddress = ctx.String(MetricsAddressFlag)
	RefreshFile = ctx.String(RefreshFileFlag)
}
