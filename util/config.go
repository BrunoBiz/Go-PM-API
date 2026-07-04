package util

import "github.com/spf13/viper"

// Stores all configs - Reads with VIPER
type Config struct {
	PVEUrl           string `mapstructure:"PVE_URL"`            // Proxmox API URL
	PVEUser          string `mapstructure:"PVE_USER"`           // Proxmox API User
	PVERealm         string `mapstructure:"PVE_REALM"`          // Proxmox API Realm
	PVEUserRealm     string `mapstructure:"PVE_USER_REALM"`     // PVEUser + "@" + PVERealm -> Ease of use | E.g. go-pm-api@pve or root@pam
	PVETokenID       string `mapstructure:"PVE_TOKEN_ID"`       // "Name" of the token
	PVEToken         string `mapstructure:"PVE_TOKEN"`          // Token itself -> Secret
	PVENodeName      string `mapstructure:"PVE_NODE_NAME"`      // Main node - This API loads only one node and it's containers
	SSHKeyFile       string `mapstructure:"SSH_KEY_FILE"`       // SSH private key file path
	SSHKeyPassphrase string `mapstructure:"SSH_KEY_PASSPHRASE"` // SSH key passphrase
	SSHPveIP         string `mapstructure:"SSH_PVE_IP"`         // PVE IP for SSH Connection (Same IP address as PVE_URL)
	SSHPvePort       string `mapstructure:"SSH_PVE_PORT"`       // PVE PORT for SSH Connection (Default 22)
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
