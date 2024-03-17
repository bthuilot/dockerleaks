package cmd

import (
	"github.com/bthuilot/dockerleaks/cmd/analyze"
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/bthuilot/dockerleaks/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "dockerleaks",
	Short: "Dockerleaks scans docker images for secrets",
	Long: `Scan docker images for secrets leaked via environment variables,
build arguments, or present on the filesystem.

Copyright (C) 2023 Bryce Thuilot.
This program comes with ABSOLUTELY NO WARRANTY; This is free software,
and you are welcome to redistribute it under certain conditions.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cobra.OnInitialize(func() {
		bindFlags()
		if err := config.Init(); err != nil {
			logging.Fatal("error initializing config: %s", err)
		}
	})

	rootCmd.PersistentFlags().StringP("log-level", "l", "off", "log level (off, debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().Bool("disable-color", false, "disable color use")
	rootCmd.PersistentFlags().StringP("config", "c", "./", "path to config file")
	rootCmd.PersistentFlags().BoolP("unmask", "u", false, "secret values should be unmasked")

	rootCmd.AddCommand(analyze.Command)
	if err := rootCmd.MarkPersistentFlagFilename("config", "yaml", "yml"); err != nil {
		log.Fatalf("err marking config as filename %s", err)
	}
}

func bindFlags() {
	_ = viper.BindPFlag(config.ViperLogLevelKey, rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag(config.ViperDisableColorKey, rootCmd.PersistentFlags().Lookup("disable-color"))
	if err := viper.BindPFlag(config.ViperUnmaskKey, rootCmd.PersistentFlags().Lookup("unmask")); err != nil {
		log.Fatalf("err binding reveal %s", err)
	}

	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		log.Fatalf("err binding config %s", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.Fatal(err.Error())
	}
}
