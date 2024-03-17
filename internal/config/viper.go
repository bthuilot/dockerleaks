package config

import (
	"github.com/spf13/viper"
)

// initViper will initialize the viper configuration
// and read in the config file
func initViper() error {
	viper.SetConfigName("dockerleaks")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dockerleaks/")
	viper.AddConfigPath("$HOME/.dockerleaks")
	viper.AddConfigPath(".")

	if viper.IsSet(ViperConfigFileKey) {
		viper.SetConfigFile(viper.GetString(ViperConfigFileKey))
	}
	//if err := viper.ReadInConfig(); err != nil {
	//var configFileNotFoundError viper.ConfigFileNotFoundError
	//if errors.As(err, &configFileNotFoundError) {
	//	_, _ = fmt.Fprintf(os.Stderr, "unable to read config file, exiting\n")
	//	os.Exit(1)
	//}
	//}

	return viper.ReadInConfig()
}
