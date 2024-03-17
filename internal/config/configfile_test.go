package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"testing"
)

func TestInit(t *testing.T) {
	n1, fP1, p1, e1 := "test", ".*[159]?", "testing", rand.Float64()
	n2, fP2, p2, e2 := "another", "\\d+\\w?", "123456", rand.Float64()
	createTestConfigFile(t, File{
		DynamicRules: []UserDynamicRule{
			{
				Name:        n1,
				FilePattern: fP1,
				Pattern:     p1,
				MinEntropy:  e1,
			},
			{
				Name:        n2,
				FilePattern: fP2,
				Pattern:     p2,
				MinEntropy:  e2,
			},
		},
		ExcludeDefaultStaticRules: true,
	})

	if err := Init(); err != nil {
		t.Fatalf("Init should not have thrown an error: %s", err)
	}

	var cfg File
	assert.NoError(t, viper.Unmarshal(&cfg))
	assert.Equal(t, logrus.GetLevel(), logrus.InfoLevel)
	// match dynamic rules
	assert.Len(t, cfg.DynamicRules, 2)
	assert.Equal(t, n1, cfg.DynamicRules[0].Name)
	assert.Equal(t, fP1, cfg.DynamicRules[0].FilePattern)
	assert.Equal(t, p1, cfg.DynamicRules[0].Pattern)
	assert.Equal(t, e1, cfg.DynamicRules[0].MinEntropy)
	assert.Equal(t, n2, cfg.DynamicRules[1].Name)
	assert.Equal(t, fP2, cfg.DynamicRules[1].FilePattern)
	assert.Equal(t, p2, cfg.DynamicRules[1].Pattern)
	assert.Equal(t, e2, cfg.DynamicRules[1].MinEntropy)
	// match exclude default static rules
	assert.True(t, cfg.ExcludeDefaultStaticRules)
	// match defaults
	assert.False(t, cfg.ExcludeDefaultDynamicRules)
	assert.False(t, cfg.IgnoreInvalidRules)
}

func TestInitViperLogLevel(t *testing.T) {
	createTestConfigFile(t, File{})
	viper.Set(ViperLogLevelKey, "debug")
	assert.NoError(t, Init())

	assert.Equal(t, logrus.GetLevel(), logrus.DebugLevel)

	viper.Set(ViperLogLevelKey, "off")
	assert.NoError(t, Init())

	assert.Equal(t, logrus.StandardLogger().Out, io.Discard)
}

func TestInitViperError(t *testing.T) {
	viper.SetConfigFile("nonexistent")
	assert.Error(t, Init())
}
