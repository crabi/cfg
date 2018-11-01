package cfg_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/crabi/cfg"
)

func TestParamsSetDefault(t *testing.T) {
	tests := map[string]struct {
		params         *cfg.Params
		expectedParams *cfg.Params
	}{
		"TestHappyPath": {
			params: &cfg.Params{
				Path:     "path",
				FileType: "type",
			},
			expectedParams: &cfg.Params{
				Path:     "path",
				FileType: "type",
			},
		},
		"TestEmptyPath": {
			params: &cfg.Params{
				FileType: "type",
			},
			expectedParams: &cfg.Params{
				Path:     "./config",
				FileType: "type",
			},
		},
		"TestEmptyFileType": {
			params: &cfg.Params{
				Path: "path",
			},
			expectedParams: &cfg.Params{
				Path:     "path",
				FileType: "yaml",
			},
		},
		"TestEmptyFields": {
			params: &cfg.Params{},
			expectedParams: &cfg.Params{
				Path:     "./config",
				FileType: "yaml",
			},
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			newParams := &(*testCase.params)
			newParams.SetDefaults()
			if !reflect.DeepEqual(newParams, testCase.expectedParams) {
				t.Errorf("Expected %+v, got %+v", testCase.expectedParams, newParams)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := map[string]struct {
		params        *cfg.Params
		envNameValue  string
		expectedError error
	}{
		"TestHappyPathDefaultEnv": {
			params: &cfg.Params{
				Path: "./confs/conf",
			},
			expectedError: nil,
		},
		"TestHappyPathOverrideEnv": {
			params: &cfg.Params{
				Path: "./confs/conf",
			},
			envNameValue:  "production",
			expectedError: nil,
		},
		"TestNilParams": {
			params:        nil,
			expectedError: errors.New("nil params"),
		},
		"TestBadFileType": {
			params: &cfg.Params{
				Path:     "./confs/conf",
				FileType: "xml",
			},
			expectedError: errors.New("error reading \"default\" configuration file"),
		},
		"TestNotFoundDefaultEnv": {
			params: &cfg.Params{
				Path: "./confs/confMissingDefault",
			},
			envNameValue:  "staging",
			expectedError: errors.New("error merging \"staging\" configuration file"),
		},
		"TestNotFoundEnv": {
			params: &cfg.Params{
				Path: "./confs/conf",
			},
			envNameValue:  "staging",
			expectedError: errors.New("error merging \"staging\" configuration file"),
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			os.Setenv("GOENV", testCase.envNameValue)

			_, err := cfg.Load(testCase.params)
			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("Expected error %+v, got %+v", testCase.expectedError, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	os.Setenv("GOENV", "local")
	config, err := cfg.Load(&cfg.Params{
		Path: "./confs/conf",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	tests := map[string]struct {
		key               string
		expectedConfigMap interface{}
	}{
		"TestGetFirstLevel": {
			key: "api",
			expectedConfigMap: map[string]interface{}{
				"host": "127.0.0.1",
				"port": 8080,
			},
		},
		"TestGetNested": {
			key: "services.foo",
			expectedConfigMap: map[string]interface{}{
				"string": "foo",
				"int":    42,
				"float":  9.81,
				"bool":   false,
			},
		},
		"TestNonExistentConfig": {
			key:               "services.cux",
			expectedConfigMap: nil,
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			configMap := config.Get(testCase.key)
			if !reflect.DeepEqual(configMap, testCase.expectedConfigMap) {
				t.Errorf("Expected %+v, got %+v", testCase.expectedConfigMap, configMap)
			}
		})
	}
}

func TestEnvVariablesReplacement(t *testing.T) {
	os.Setenv("HOST", "0.0.0.0")
	os.Setenv("FLOAT", "10.5")
	config, err := cfg.Load(&cfg.Params{
		Path: "./confs/conf",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	tests := map[string]struct {
		key               string
		expectedConfigMap interface{}
	}{
		"TestReplaceString": {
			key: "api",
			expectedConfigMap: map[string]interface{}{
				"host": "0.0.0.0",
				"port": 8080,
			},
		},
		"TestFloatEnv": {
			key: "services.foo",
			expectedConfigMap: map[string]interface{}{
				"string": "foo",
				"int":    42,
				"float":  10.5,
				"bool":   false,
			},
		},
		"TestSingleMissingDeepKey": {
			key:               "services.foo.string.foo",
			expectedConfigMap: nil,
		},
		"TestSingleSingleKey": {
			key:               "services.foo.string",
			expectedConfigMap: "foo",
		},
		"TestNonExistentConfig": {
			key:               "services.cux",
			expectedConfigMap: nil,
		},
	}
	fmt.Println(config.AllSettings())
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			configMap := config.Get(testCase.key)
			if !reflect.DeepEqual(configMap, testCase.expectedConfigMap) {
				t.Errorf("Expected %+v, got %+v", testCase.expectedConfigMap, configMap)
			}
		})
	}
}
