package cfg_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/crabi/cfg"
	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	tests := map[string]struct {
		configPath     string
		configName     string
		expectedConfig map[string]interface{}
		expectedError  error
	}{
		"TestOverrideHappyPath": {
			configPath: "./testConf/conf/",
			configName: "test",
			expectedConfig: map[string]interface{}{
				"api": map[string]interface{}{
					"host": "127.0.0.1",
					"port": 8080,
				},
				"foo_service": map[string]interface{}{
					"foo_string": "testfoo",
					"foo_int":    42,
					"foo_float":  9.81,
					"foo_bool":   true,
				},
				"bar_service": map[string]interface{}{
					"bar_string": "bar",
					"bar_int":    84,
					"bar_float":  3.141592,
					"bar_bool":   false,
				},
			},
			expectedError: nil,
		},
		"TestNoOverrideHappyPath": {
			configPath: "./testConf/conf/",
			configName: "",
			expectedConfig: map[string]interface{}{
				"api": map[string]interface{}{
					"host": "127.0.0.1",
					"port": 8080,
				},
				"foo_service": map[string]interface{}{
					"foo_string": "foo",
					"foo_int":    42,
					"foo_float":  9.81,
					"foo_bool":   false,
				},
				"bar_service": map[string]interface{}{
					"bar_string": "bar",
					"bar_int":    84,
					"bar_float":  3.141592,
					"bar_bool":   false,
				},
			},
			expectedError: nil,
		},
		"TestDefaultConfigurationNonExistent": {
			configPath:     "./testConf/confMissingDefault/",
			configName:     "",
			expectedConfig: nil,
			expectedError:  errors.New("error loading \"default\" configuration"),
		},
		"TestOverrideNonExistentName": {
			configPath:     "./testConf/conf/",
			configName:     "staging",
			expectedConfig: nil,
			expectedError:  errors.New("error overriding \"staging\" configuration"),
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			viper.Reset()
			config, err := cfg.LoadConfig(testCase.configPath, testCase.configName)
			if !reflect.DeepEqual(config, testCase.expectedConfig) {
				t.Errorf("Expected config %+v, got %+v", testCase.expectedConfig, config)
			}
			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("Expected error %+v, got %+v", testCase.expectedError, err)
			}
		})
	}
}
