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
	"os"
	"strings"
	"strconv"
	// "log"
	"fmt"
	"time"
	"path/filepath"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"

	// "bitbucket.org/clientcto/go-logging-client"
	// "github.com/bullapse/bootloader"
	// "github.com/natefinch/lumberjack"

	consulapi "github.com/hashicorp/consul/api"	

	"github.com/magiconair/properties"
	"gopkg.in/yaml.v2"

)
// var logger *log.Logger

type ConfigProperties map[string]string

func main(){
	// Load configuration data
	if err := readConfigurationFile("./configuration.json"); err != nil {
		fmt.Println(err.Error())
		return
	}

	// logger = log.New(&lumberjack.Logger{
	// 		Filename:   configuration.LoggingFile,
	// 		MaxSize:    500,	// megabytes
	// 		MaxBackups: 3,
	// 		MaxAge:     28, 	//days
	// }, "", log.Ldate | log.Ltime | log.Lshortfile)
	// bootloader.RunWithLog("banner.txt", logger)

	consulClient, err := getConsulCient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kv := consulClient.KV()

	if configuration.IsReset {
		_, err := kv.DeleteTree(configuration.GlobalPrefix, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else if !isConfigInitialized(kv) {
		loadConfigFromPath(kv)
	}
}

func readConfigurationFile(path string) error {
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
		resp, err := http.Get(consulUrl + CONSUL_STATUS_PATH) //@TODO: not sure if this method is proper
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
	if fails == configuration.FailLimit {
		return nil, errors.New("Cannot get connection to Consul")
	}

	// Connect to the Consul Agent
	config := consulapi.DefaultConfig()
	config.Address = consulUrl

	return consulapi.NewClient(config)
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

		// skip directories
		if info.IsDir() {
			return nil
		}
		// verify file extension
		if !isAcceptablePropertyExtensions(info.Name()) {
			return nil
		}
		
		dir, file := filepath.Split(path)
		fmt.Println("found config file:", file, "in context", strings.TrimPrefix(dir, configuration.GlobalPrefix + "/"))

		props, err := readPropertiesFile(path)
		if err != nil {
			return err
		}
	
		for k := range props {
			p := &consulapi.KVPair{Key: dir + k, Value: []byte(props[k])}
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

func readPropertiesFile(filePath string) (ConfigProperties, error) {
	configProps := ConfigProperties{}

	if isYamlExtensions(filePath) {
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
	} else {
		props, err := properties.LoadFile(filePath, properties.UTF8)
		if err != nil {
			return nil, err
		}
		configProps = props.Map()
	}

	return configProps, nil
}

func isYamlExtensions(file string) bool {
	for _, v := range configuration.YamlExtensions {
		if v == filepath.Ext(file) {
			return true
		}
	}
	return false
}