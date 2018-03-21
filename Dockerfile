###############################################################################
# Copyright 2017 Samsung Electronics All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
###############################################################################

# Consul Docker image for EdgeX Foundry
FROM consul:0.7.3

# environment variables
ENV APP_DIR=/edgex/core-config-seed-go
ENV APP=core-config-seed-go
ENV WAIT_FOR_A_WHILE=5
ENV CONSUL_ARGS="-server -client=0.0.0.0 -bootstrap -ui"

# copy files
COPY $APP launch-consul-config.sh $APP_DIR/
COPY ./res $APP_DIR/res
COPY ./config $APP_DIR/config

# set the working directory
WORKDIR $APP_DIR

# call the wrapper to launch consul and config-seed application
CMD ["sh", "launch-consul-config.sh"]
