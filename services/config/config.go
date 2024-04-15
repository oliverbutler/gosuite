package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func getConfigPath() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(xdgConfigHome, "gosuite", "config.yml")
}

const emptyConfig = `
databases:
  - name: "default"
    user: "root"
    password: "password"
    # password_cmd: "op get item database/mysql"
    host: "localhost"
    port: 3306
    database: "my-db"
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

		_, err = file.WriteString(emptyConfig)

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

type DatabaseConfig struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type AppConfig struct {
	Databases []DatabaseConfig `yaml:"databases"`
}

func LoadConfig(path string, config *AppConfig) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(config)

	if err != nil {
		return err
	}

	return nil
}

func GetConfig() (*AppConfig, error) {
	path := getConfigPath()

	err := CreateConfigIfMissing(path)

	if err != nil {
		return nil, err
	}

	config := &AppConfig{}

	err = LoadConfig(path, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
