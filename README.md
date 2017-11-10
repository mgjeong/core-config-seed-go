# go-core-config-seed
This repository is for initializing the Configuration Management micro service.
It loads the default configuration from property or YAML files, and push values to the Consul Key/Value store.

## Configuration Guidelines ##

The configuration of this tool is located in res/configuration.json.
There are several properties in it, and here are the default values and explanation:

\#The root path of the configuration files which would be loaded by this tool
ConfigPath=./config

\#The global prefix namespace which will be created on the Consul Key/Value store
GlobalPrefix=config

\#The communication protocol of the Consul server
ConsulProtocol=http

\#The hostname of the Consul server
ConsulHost=localhost

\#The communication port number of the Consul server
ConsulPort=8500

\#If isReset=true, it will remove all the original values under the globalPrefix and import the configuration data
\#If isReset=false, it will check the globalPrefix exists or not, and it only imports configuration data when the globalPrefix doesn't exist. 
IsReset=false

\#The number for retry to connect to the Consul server when connection fails
FailLimit=30

\#The seconds how long to wait for the next retry
FailWaittime=3

## Configuration File Structure ##

In /config folder, there are some sample files for testing.
The structure of the keys on the Consul server will be the same as the folders of the configPath, and the folder name must be the same as the microservice id registered on the Consul server.

For example, the files under /config/edgex-core-data folder will be loaded and create /{global_prefix}/edgex-core-data/{property_name} on the Consul server.
In addition, "edgex-core-data" is the micro service id of Core Data micro service.

However, you can use different profile name to categorize the usage on the same microservice. For instance,
"/config/edgex-core-data" contains the default configuration of Core Data Microservice.
"/config/edgex-core-data,dev" contains the specific configuration for development time, and "dev" is the profile name.
"/config/edgex-core-data,test" contains the specific configuration for test time, and "test" is the profile name.

## Docker
### Build
- `docker build -t go-core-config-seed .`
### Run
- `docker run -p 8400:8400 -p 8500:8500 -p 8600:8600 --name="edgex-config-seed" --hostname="edgex-core-config-seed" go-core-config-seed`

## TODO ## 
- Apply go-logging-client
- Apply gofmt
- TC with gomock
- Docker Compose file