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

// Struct used to pase the JSON configuration file
type ConfigurationStruct struct{
    ConfigPath string
    GlobalPrefix string
    ConsulProtocol string
    ConsulHost string
    ConsulPort int
    IsReset bool
    FailLimit int
    FailWaittime int
    AcceptablePropertyExtensions []string
    YamlExtensions []string
    LoggingFile string
}

// Configuration data for the config-seed service
var configuration ConfigurationStruct = ConfigurationStruct{}	// Needs to be initialized before used

var (
	CONSUL_STATUS_PATH string = "/v1/agent/self"
)