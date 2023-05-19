package config

import "github.com/bthuilot/dockerleaks/pkg/detections"

type Config struct {
	Regexp RegexpConfig
}

type RegexpConfig struct {
	Patterns []detections.Pattern
}

type EntropyConfig struct {
}

func Init() error {
	if err := initViper(); err != nil {
		return err
	}
	initLogger()
	return nil
}
