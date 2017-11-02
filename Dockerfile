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
# Docker image for EdgeX Foundry Config Seed 
FROM golang:1.7.5-alpine AS build-env

RUN mkdir -p /go/src \
 && mkdir -p /go/bin \
 && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN apk add --update git
RUN go get github.com/hashicorp/consul/api
RUN go get github.com/magiconair/properties
RUN go get gopkg.in/yaml.v2

RUN mkdir -p $GOPATH/src/go-core-config-seed
WORKDIR src/go-core-config-seed
COPY /src/. .
COPY /res/. .

RUN go build .


# Consul Docker image for EdgeX Foundry
FROM consul:0.7.3

# environment variables
ENV APP_DIR=/edgex/go-core-config-seed
ENV APP=go-core-config-seed
ENV WAIT_FOR_A_WHILE=10
ENV CONSUL_ARGS="-server -client=0.0.0.0 -bootstrap -ui"

#set the working directory
WORKDIR $APP_DIR

#copy Go App and default config files to the image
COPY --from=build-env /go/src/go-core-config-seed/ .

COPY launch-consul-config.sh .
COPY ./config ./config

#call the wrapper to launch consul and main app
CMD $APP_DIR/launch-consul-config.sh
