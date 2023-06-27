package cmd

import (
	"fmt"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/config"
	testSuite "github.com/LampardNguyen234/astra-integration-test/test-suite"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

const (
	flagConfig    = "config"
	EnvConfigFile = "CONFIG_FILE_LOCATION"
)

func Run(args []string) error {
	cliApp := &cli.App{
		Name:                 filepath.Base(args[0]),
		Usage:                "Reward Notification Service",
		Version:              "v0.0.1",
		Copyright:            "(c) 2023 stellalab.com",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flagConfig,
				Value: "./default-config.json",
				Usage: "The config file to load from",
			},
		},
		Action: func(ctx *cli.Context) error {
			if args := ctx.Args(); args.Len() > 0 {
				return fmt.Errorf("unexpected arguments: %q", args.Get(0))
			}

			// Prepare FileConfig
			configPath := ctx.String(flagConfig)
			if os.Getenv(EnvConfigFile) != "" {
				configPath = os.Getenv(EnvConfigFile)
			}
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("error load config: %v", err)
			}
			// some configs are missing, load them from the cli.
			if _, err = cfg.IsValid(); err != nil {
				return fmt.Errorf("invalid config: %v", err)
			}

			log := logger.NewZeroLogger(cfg.Logger.LogPath)
			log.SetLogLevel(logger.LogLevel(cfg.Logger.Level))

			mainApp, err := testSuite.NewTestSuite(cfg.TestSuite, log)
			if err != nil {
				return err
			}

			mainApp.RunTest()

			return nil
		},
	}

	err := cliApp.Run(args)
	if err != nil {
		return err
	}

	return nil
}
