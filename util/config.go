package util

import "github.com/spf13/viper"

// Stores all configs - Reads with VIPER
type Config struct {
	PVEUrl       string `mapstructure:"PVE_URL"`
	PVEUser      string `mapstructure:"PVE_USER"`
	PVERealm     string `mapstructure:"PVE_REALM"`
	PVEUserRealm string `mapstructure:"PVE_USER_REALM"`
	PVETokenID   string `mapstructure:"PVE_TOKEN_ID"`
	PVEToken     string `mapstructure:"PVE_TOKEN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("pm")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
