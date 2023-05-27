package cmd

import (
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
	Run: runDetect,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init()
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("image", "i", "", "the name of the image")
	_ = rootCmd.MarkFlagRequired("image")
	rootCmd.PersistentFlags().StringP("log-level", "l", "off", "log level (off, debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().Bool("disable-color", false, "disable color use")
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config file")
	rootCmd.PersistentFlags().BoolP("pull", "p", false, "image should be pulled from remote")
	// TODO(refine this)
	cobra.OnInitialize(initLog)
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		log.Fatalf("err binding config %s", err)
	}
}

func initLog() {
	_ = viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("pull_image", rootCmd.PersistentFlags().Lookup("pull"))
	_ = viper.BindPFlag("disable_color", rootCmd.PersistentFlags().Lookup("disable-color"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.Fatal(err)
	}
}
