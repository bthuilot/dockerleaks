package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"math/rand"
	"os"
	"path"
	"testing"
)

func createTestConfigFile(t *testing.T, cfg File) {
	t.Helper()

	fPath := path.Join(t.TempDir(), fmt.Sprintf("dockerleaks-%s.yml", t.Name()))
	f, err := os.Create(fPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	cfgB, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.Write(cfgB)
	if err != nil {
		t.Fatal(err)
	}

	viper.Set(ViperConfigFileKey, fPath)
}

func TestInitViper(t *testing.T) {
	staticRuleName := "testing123"
	staticRulePattern := "(?i)this(is)?a\\s+test"
	staticRuleEntropy := rand.Float64()
	createTestConfigFile(t, File{
		StaticRules: []UserStaticRule{
			{
				Name:       staticRuleName,
				MinEntropy: staticRuleEntropy,
				Pattern:    staticRulePattern,
			},
		},
		IgnoreInvalidRules:         true,
		ExcludeDefaultDynamicRules: false,
	})

	if err := initViper(); err != nil {
		t.Fatalf("Init should not have thrown an error: %s", err)
	}
	var cfg File
	assert.NoError(t, viper.Unmarshal(&cfg))
	// match static rules
	assert.Len(t, cfg.StaticRules, 1)
	assert.Equal(t, staticRuleName, cfg.StaticRules[0].Name)
	assert.Equal(t, staticRulePattern, cfg.StaticRules[0].Pattern)
	assert.Equal(t, staticRuleEntropy, cfg.StaticRules[0].MinEntropy)
	// no dynamic rules
	assert.Len(t, cfg.DynamicRules, 0)
	// set to true
	assert.True(t, cfg.IgnoreInvalidRules)
	// set to false
	assert.False(t, cfg.ExcludeDefaultDynamicRules)
	// default false
	assert.False(t, cfg.ExcludeDefaultStaticRules)
}

func TestInitViperNoFile(t *testing.T) {
	viper.Set(ViperConfigFileKey, "/does/not/exist")
	assert.Error(t, initViper())
}
