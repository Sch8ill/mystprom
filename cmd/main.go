package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/sch8ill/mystprom/api/cryptocompare"
	"github.com/sch8ill/mystprom/api/mystnodes"
	"github.com/sch8ill/mystprom/config"
	"github.com/sch8ill/mystprom/metrics"
	"github.com/sch8ill/mystprom/monitor"
)

func main() {
	createLogger()

	app := createApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("error encountered")
	}
}

func run(ctx *cli.Context) error {
	config.SetConfig(ctx)
	log.Info().Str("email", config.MystAPIEmail).Bool("password", config.MystAPIPassword != "").Msg("Credentials")
	log.Info().Str("interval", config.ScrapeInterval.String()).Str("metrics_address", config.MetricsAddress).Str("refresh_file", config.RefreshFile).Msg("Config")

	credentials := mystnodes.Credentials{
		Email:    config.MystAPIEmail,
		Password: config.MystAPIPassword,
	}

	refreshToken, err := mystnodes.NewTokenFromFile(config.RefreshFile)
	if err != nil {
		log.Debug().Err(err).Msg("refresh token cache miss (non-critical)")
	}

	mystApi := mystnodes.NewWithRefreshToken(credentials, refreshToken)

	cryptoCompare, err := cryptocompare.New()
	if err != nil {
		return fmt.Errorf("failed to create CryptoCompare api client: %w", err)
	}

	m := monitor.New(mystApi, cryptoCompare, config.ScrapeInterval)
	m.Start()
	defer m.Stop()

	if err := metrics.Listen(); err != nil {
		return fmt.Errorf("failed to start prometheus exporter: %w", err)
	}

	return nil
}

func createApp() *cli.App {
	return &cli.App{
		Name:      "mystprom",
		Usage:     "Monitor your Mysterium Network nodes using prometheus.",
		Copyright: "Copyright (c) 2024 Sch8ill",
		Action:    run,
		Flags:     config.DeclareFlags(),
	}
}

func createLogger() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}
	log.Logger = log.Output(consoleWriter).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}
