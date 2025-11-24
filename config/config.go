package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDir  = ".uppies"
	ConfigFile = "config.yaml"
)


type Profile struct {
	Name  string `yaml:"name"`
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

type configData struct {
	ActiveProfile string    `yaml:"active_profile"`
	Profiles      []Profile  `yaml:"profiles"`
}

var Token string
var Host string

var defaultConfig = configData {
	ActiveProfile: "default",
	Profiles: []Profile {
		{ Name: "default", Host: "https://api.uppiesplz.com:3000", Token: "" },
		{ Name: "local", Host: "http://localhost:3000", Token: "" },
	},
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ConfigDir, ConfigFile)
}

func LoadConfig() {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)

	var data configData
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		data = defaultConfig
		// Ensure directory exists
		os.MkdirAll(configDir, 0755)
		// Write to file
		file, err := os.Create(configPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		yaml.NewEncoder(file).Encode(data)
	} else {
		// Read from file
		file, err := os.Open(configPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		yaml.NewDecoder(file).Decode(&data)
		// Merge with defaults for missing values
		if data.ActiveProfile == "" {
			data.ActiveProfile = defaultConfig.ActiveProfile
		}
		if len(data.Profiles) == 0 {
			data.Profiles = defaultConfig.Profiles
		}
	}

	// Set active profile values
	for _, profile := range data.Profiles {
		if profile.Name == data.ActiveProfile {
			Token = profile.Token
			Host = profile.Host
			break
		}
	}
}

func GetConfig() (configData, error) {
	configPath := getConfigPath()

	var data configData
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		data = defaultConfig
	} else {
		file, err := os.Open(configPath)
		if err != nil {
			return data, err
		}
		defer file.Close()
		if err := yaml.NewDecoder(file).Decode(&data); err != nil {
			return data, err
		}
		// Merge with defaults for missing values
		if data.ActiveProfile == "" {
			data.ActiveProfile = defaultConfig.ActiveProfile
		}
		if len(data.Profiles) == 0 {
			data.Profiles = defaultConfig.Profiles
		}
	}
	return data, nil
}

func SaveConfig() {
	configPath := getConfigPath()

	// Load current config
	var data configData
	file, err := os.Open(configPath)
	if err == nil {
		yaml.NewDecoder(file).Decode(&data)
		file.Close()
	} else {
		data = defaultConfig
	}

	// Update active profile
	for i, profile := range data.Profiles {
		if profile.Name == data.ActiveProfile {
			data.Profiles[i].Token = Token
			data.Profiles[i].Host = Host
			break
		}
	}

	// Save
	file, err = os.Create(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	yaml.NewEncoder(file).Encode(data)
}

func SaveConfigData(data configData) error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Save
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(data)
}

func SwitchProfile(profileName string) error {
	configPath := getConfigPath()

	// Load current config
	var data configData
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	yaml.NewDecoder(file).Decode(&data)

	// Check if profile exists
	found := false
	for _, profile := range data.Profiles {
		if profile.Name == profileName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("profile %s not found", profileName)
	}

	// Update active profile
	data.ActiveProfile = profileName

	// Save
	file, err = os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	yaml.NewEncoder(file).Encode(data)

	// Reload to set globals
	LoadConfig()
	return nil
}