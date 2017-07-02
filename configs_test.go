package gfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)

func TestConfig(t *testing.T) {
	testPath := "./test-config.json"

	t.Run("Defaults if file doesn't exist", func(t *testing.T) {
		a := assert.New(t)

		configs, err := GetConfigs(testPath)
		if a.NoError(err) {
			a.NotNil(configs)

			t.Run("Default values", func(t *testing.T) {
				a := assert.New(t)

				a.NotEqual("password", configs.Password)
				a.Equal("username", configs.Username)
			})

			t.Run("File should have been created", func(t *testing.T) {
				a := assert.New(t)
				_, err := os.Stat(testPath)
				a.NoError(err)
			})
		}
	})

	// Clean up
	os.Remove(testPath)

	t.Run("Save and read config file", func(t *testing.T) {
		a := assert.New(t)
		testConfig := &Config {
			Password: "superPassword!.!",
			Username:"test",
			Serve:"./somewhere",
		}

		err := SaveConfigs(testPath, testConfig)
		if a.NoError(err) {

			t.Run("Should load file", func(t *testing.T) {
				a := assert.New(t)

				loadedConfig, err := GetConfigs(testPath)
				if a.NoError(err) {
					a.Equal(testConfig.Username, loadedConfig.Username)
					a.Equal(testConfig.Password, loadedConfig.Password)
					a.Equal(testConfig.Serve, loadedConfig.Serve)
				}
			})
		}
	})

	os.Remove(testPath)

}
