package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	GithubToken string
	GitlabToken string
}

func (c *Config) Set(key, value string) error {
	field := GetField(Config{}, key)
	if field.Tag.Get("kvp") == "" {
		return fmt.Errorf("unknown key: %s", key)
	}

	reflect.ValueOf(c).Elem().FieldByName(field.Name).SetString(value)

	return nil
}

var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}

	config, err := LoadConfig()
	if err != nil {
		if err == ErrConfigNotFound {
			config, err = CreateConfig()
			if err != nil {
				return nil
			}
		} else {
			return nil
		}
	}

	return config
}

func CreateConfig() (*Config, error) {
	configDirectory, err := GetConfigDirectory()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(configDirectory, "fpm", ".fpmrc")

	if _, err := os.Stat(configFile); err == nil {
		return nil, ErrConfigExists
	}

	if _, err := os.Stat(filepath.Join(configDirectory, "fpm")); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(configDirectory, "fpm"), 0755)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Create(configFile)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	config := &Config{}

	for _, field := range GetFields(Config{}) {
		if field.Tag.Get("kvp") == "" {
			continue
		}

		_, err := file.WriteString(field.Name + "=" + field.Tag.Get("default") + "\n")
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func LoadConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	configDirectory, err := GetConfigDirectory()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(configDirectory, "fpm", ".fpmrc")

	if !Exists(configFile) {
		return nil, ErrConfigNotFound
	}

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var line string
	config := &Config{}

	for {
		_, err := fmt.Fscanln(file, &line)
		if err != nil {
			break
		}

		// Ignore comments
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Ignore empty lines
		if line == "" {
			continue
		}

		// Ignore lines that don't have a key=value format
		if !strings.Contains(line, "=") {
			continue
		}

		// Split line into key and value
		split := strings.Split(line, "=")
		key := split[0]
		value := split[1]

		// Find the field that matches the key
		field := GetField(Config{}, key)
		if field.Tag.Get("kvp") == "" {
			continue
		}

		switch field.Tag.Get("type") {
		case "string":
			SetField(reflect.ValueOf(config).Elem().FieldByName(field.Name), value)
			break
		case "int":
			parsedInt, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				parsedInt, err = strconv.ParseInt(field.Tag.Get("default"), 10, 32)
				if err != nil {
					panic(err)
				}
			}
			SetField(reflect.ValueOf(config).Elem().FieldByName(field.Name), int(parsedInt))
			break
		case "int64":
			parsedInt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				parsedInt, err = strconv.ParseInt(field.Tag.Get("default"), 10, 64)
				if err != nil {
					panic(err)
				}
			}

			SetField(reflect.ValueOf(config).Elem().FieldByName(field.Name), parsedInt)
			break
		case "bool":
			parseBool, err := strconv.ParseBool(value)
			if err != nil {
				parseBool, err = strconv.ParseBool(field.Tag.Get("default"))
				if err != nil {
					return nil, err
				}
			}
			SetField(reflect.ValueOf(config).Elem().FieldByName(field.Name), parseBool)
			break
		default:
			panic("Unknown type: " + field.Tag.Get("type"))
		}
	}

	token, err := getGithubToken()
	if err != nil && err != ErrSecretNotFound {
		return nil, err
	}
	config.GithubToken = token

	return config, nil
}
