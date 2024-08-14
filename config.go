package appctx

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func (app *AppCtx[_, _]) loadConfig(configName string) error {
	if app.noConfig {
		return nil
	}

	f, err := os.Open(configName)
	if err != nil {
		return fmt.Errorf("opening config file \"%v\": %w", configName, err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.DecodeContext(app, &app.cfg)
	if err != nil {
		return fmt.Errorf("decoding YAML: %w", err)
	}

	return nil
}
