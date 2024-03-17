package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func revertOldViperVal(t *testing.T, key string) {
	t.Helper()
	old := viper.Get(key)
	t.Cleanup(func() {
		viper.Set(key, old)
	})
}

func TestShouldUseColor(t *testing.T) {
	revertOldViperVal(t, ViperDisableColorKey)
	viper.Set(ViperDisableColorKey, false)
	assert.Truef(t, ShouldUseColor(), "should return true when disable color is set to false")

	viper.Set(ViperDisableColorKey, true)
	assert.Falsef(t, ShouldUseColor(), "should return false when disable color is set to true")
}

func TestShouldUseSpinner(t *testing.T) {
	revertOldViperVal(t, ViperLogLevelKey)
	viper.Set(ViperLogLevelKey, nil)
	assert.Truef(t, ShouldUseSpinner(), "should return true when log level is not set")

	viper.Set(ViperLogLevelKey, "off")
	assert.Truef(t, ShouldUseSpinner(), "should return true when log level is off")

	viper.Set(ViperLogLevelKey, "debug")
	assert.Falsef(t, ShouldUseSpinner(), "should return false when log level is debug")
}

func TestInitLogger(t *testing.T) {
	assert.NoError(t, initLogger("debug"))
	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel())

	assert.NoError(t, initLogger("info"))
	assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())

	assert.NoError(t, initLogger("warn"))
	assert.Equal(t, logrus.WarnLevel, logrus.GetLevel())

	assert.NoError(t, initLogger("error"))
	assert.Equal(t, logrus.ErrorLevel, logrus.GetLevel())

	assert.Error(t, initLogger("invalid"))
}
