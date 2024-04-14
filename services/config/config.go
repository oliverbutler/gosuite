package config

import (
	"context"
	"os"
	"path/filepath"
)

func getConfigPath() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(xdgConfigHome, "gosuite", "config.pkl")
}

var emptyPkl = `
databases = new Listing { }
`

func CreateConfigIfMissing(path string) error {
	_, err :=
		os.Stat(path)

	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0700)
		file, err := os.Create(path)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = file.WriteString(emptyPkl)

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func GetConfig() (*AppConfig, error) {

	path := getConfigPath()

	err := CreateConfigIfMissing(path)

	if err != nil {
		return nil, err
	}

	context := context.Background()

	config, err := LoadFromPath(context, path)

	if err != nil {
		return nil, err
	}

	return config, nil
}
