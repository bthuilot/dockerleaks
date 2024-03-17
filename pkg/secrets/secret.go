package secrets

import (
	"github.com/bthuilot/dockerleaks/internal/config"
	"github.com/spf13/viper"
)

type Secret struct {
	// Value is the actual value of the secret
	Value string `json:"value"`
	// Entropy is the shannon entropy of the secret
	Entropy float64 `json:"entropy"`
}

func (s Secret) String() string {
	if s.Value == "" {
		return ""
	}
	if viper.GetBool(config.ViperUnmaskKey) {
		return s.Value
	}
	if len(s.Value) < 6 {
		return "*****"
	}
	return s.Value[:3] + "*****" + s.Value[len(s.Value)-3:]
}
