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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/magiconair/properties"
	"gopkg.in/yaml.v2"

	consulapi "github.com/hashicorp/consul/api"
)

type ConfigProperties map[string]string

func main() {
	// Load configuration data
	if err := loadConfigurationFile("./configuration.json"); err != nil {
		fmt.Println(err.Error())
		return
	}

	consulClient, err := getConsulCient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kv := consulClient.KV()

	if configuration.IsReset {
		removeStoredConfig(kv)
		loadConfigFromPath(kv)
	} else if !isConfigInitialized(kv) {
		loadConfigFromPath(kv)
	}
}

func loadConfigurationFile(path string) error {
	// Read the configuration file
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Decode the configuration as JSON
	err = json.Unmarshal(contents, &configuration)
	if err != nil {
		return err
	}

	return nil
}

func getConsulCient() (*consulapi.Client, error) {
	consulUrl := configuration.ConsulProtocol + "://" + configuration.ConsulHost + ":" + strconv.Itoa(configuration.ConsulPort)

	// Check the connection to Consul
	fails := 0
	for fails < configuration.FailLimit {
		resp, err := http.Get(consulUrl + CONSUL_STATUS_PATH)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Second * time.Duration(configuration.FailWaittime))
			fails++
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			break
		}
	}
	if fails >= configuration.FailLimit {
		return nil, errors.New("Cannot get connection to Consul")
	}

	// Connect to the Consul Agent
	config := consulapi.DefaultConfig()
	config.Address = consulUrl

	return consulapi.NewClient(config)
}

func removeStoredConfig(kv *consulapi.KV) {
	_, err := kv.DeleteTree(configuration.GlobalPrefix, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func isConfigInitialized(kv *consulapi.KV) bool {
	keys, _, err := kv.Keys(configuration.GlobalPrefix, "", nil)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if len(keys) > 0 {
		fmt.Printf("%s exists! The configuration data has been initialized.\n", configuration.GlobalPrefix)
		return true
	}
	fmt.Printf("%s doesn't exist! Start importing configuration data.\n", configuration.GlobalPrefix)
	return false
}

func loadConfigFromPath(kv *consulapi.KV) {
	err := filepath.Walk(configuration.ConfigPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories & unacceptable property extension
		if info.IsDir() || !isAcceptablePropertyExtensions(info.Name()) {
			return nil
		}

		dir, file := filepath.Split(path)
		configPath, err := filepath.Rel(".", configuration.ConfigPath)
		if err != nil {
			return err
		}

		dir = strings.TrimPrefix(dir, configPath+"/")
		fmt.Println("found config file:", file, "in context", dir)

		props, err := readPropertyFile(path)
		if err != nil {
			return err
		}

		prefix := configuration.GlobalPrefix + "/" + dir
		for k := range props {
			p := &consulapi.KVPair{Key: prefix + k, Value: []byte(props[k])}
			if _, err := kv.Put(p, nil); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func isAcceptablePropertyExtensions(file string) bool {
	for _, v := range configuration.AcceptablePropertyExtensions {
		if v == filepath.Ext(file) {
			return true
		}
	}
	return false
}

func readPropertyFile(filePath string) (ConfigProperties, error) {
	if isYamlExtensions(filePath) {
		// Read .yaml/.yml file
		return readYamlFile(filePath)
	} else {
		// Read .properties file
		return readPropertiesFile(filePath)
	}
}

func isYamlExtensions(file string) bool {
	for _, v := range configuration.YamlExtensions {
		if v == filepath.Ext(file) {
			return true
		}
	}
	return false
}

func readYamlFile(filePath string) (ConfigProperties, error) {
	configProps := ConfigProperties{}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	if err = yaml.Unmarshal(contents, &body); err != nil {
		return nil, err
	}

	for key, value := range body {
		configProps[key] = fmt.Sprint(value)
	}

	return configProps, nil
}

func readPropertiesFile(filePath string) (ConfigProperties, error) {
	configProps := ConfigProperties{}

	props, err := properties.LoadFile(filePath, properties.UTF8)
	if err != nil {
		return nil, err
	}
	configProps = props.Map()

	return configProps, nil
}
