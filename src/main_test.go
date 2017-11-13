/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/
package main
 
import (
	"testing"
	"os"

	consulapi "github.com/hashicorp/consul/api"	
)
 
type tearDown func(t *testing.T)

var (
	consulClient *consulapi.Client
	kv *consulapi.KV
)

func setUp(t *testing.T) tearDown {
    configuration.ConfigPath = "./config"
    configuration.GlobalPrefix = "config"
    configuration.ConsulProtocol = "http"
    configuration.ConsulHost = "localhost"
    configuration.ConsulPort = 8500
    configuration.IsReset = false
    configuration.FailLimit = 2
    configuration.FailWaittime = 0
    configuration.AcceptablePropertyExtensions = []string{".yaml", ".yml", ".properties"}
    configuration.YamlExtensions = []string{".yaml", ".yml"}
    configuration.LoggingFile = "edgex-core-config-seed.log"

	/*
	//@TODO: need a consul mock
	c, err := getConsulCient()
	if err != nil {
		t.Error(err.Error())
	}
	consulClient = c
	kv = consulClient.KV()
	*/
	return func(t *testing.T) {
		/*
		_, err := kv.DeleteTree(configuration.GlobalPrefix, nil)
		if err != nil {
			t.Error(err.Error())
		}
		*/
	}
}

func TestLoadConfigurationFile(t *testing.T) {
	tmpFile, err := os.Create("test.json")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer func() {
		os.Remove(tmpFile.Name())
	}()

	_, err = tmpFile.Write([]byte("{\"key\":\"value\"}"))
	if err != nil {
		t.Error(err.Error())
		return
	}

	var paths = []struct {
		path	string
		err		string
	}{
		{
			"./test.json",
			"",
		},
		{
			"./invalid.file",
			"open ./invalid.file: no such file or directory",
		},
	}

	for _, p := range paths {
		err := loadConfigurationFile(p.path)
		if err != nil && err.Error() != p.err {
			t.Error("Expected error :" + p.err + ", Actual error :" + err.Error())
		}
	}
}
/* @TODO
func TestGetConsulCient(t *testing.T) {
	// - work well
	// - invalid url
	// - fail limit exceeded
}

func TestIsConfigInitialized(t *testing.T) {
	// - "Config should not be initialized by default"
	// - "Config should be initialized after load"
	
}

func TestLoadConfigFromPath(t *testing.T) {
	// - work well (get with valid key after loading)
	// - invalid path
	// - get with invalid key after loading
}

func TestIsAcceptablePropertyExtensions(t *testing.T) {
	// - work well (.properies)
	// - not work (.json)
}

func TestReadPropertiesFile(t *testing.T) {
	// - invalid path
}

func TestIsYamlExtensions(t *testing.T) {
	// - work well (.yaml = true)
	// - not work (.properties = true)
}
*/
