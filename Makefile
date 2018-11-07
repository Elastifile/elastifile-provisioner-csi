# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: all plugin

#REGISTRY=hub.docker.com
#IMAGE_NAME = $(REGISTRY)/elastifileio/ecfs-provisioner-csi
IMAGE_NAME = elastifileio/ecfs-provisioner-csi
PLUGIN_TAG ?= next

TEMP_DIR=_output
DOCKER_DIR=deploy/docker
PLUGIN_BINARY=ecfsplugin

$(info ecfs image settings: $(IMAGE_NAME) version $(PLUGIN_TAG))

# Compile, create image and push it
all: image push

# Compile the plugin binary
binary:
	if [ ! -d ./vendor ]; then dep ensure; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  $(TEMP_DIR)/$(PLUGIN_BINARY) ./ecfs

# Create docker image
image: binary
	cp $(TEMP_DIR)/$(PLUGIN_BINARY) $(DOCKER_DIR)
	docker build -t $(IMAGE_NAME):$(PLUGIN_TAG) $(DOCKER_DIR)

# Push image to docker registry
push:
	docker push $(IMAGE_NAME):$(PLUGIN_TAG)

# Remove previous build's artifacts
clean:
	go clean -r -x
	rm -f $(DOCKER_DIR)/$(PLUGIN_BINARY)
	rm -f $(TEMP_DIR)/$(PLUGIN_BINARY)

