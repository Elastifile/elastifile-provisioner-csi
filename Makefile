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

.PHONY: all binary image deployrunner push clean

#REGISTRY ?= hub.docker.com
#PLUGIN_IMAGE_NAME = $(REGISTRY)/elastifileio/ecfs-provisioner-csi
PLUGIN_IMAGE_NAME = elastifileio/ecfs-provisioner-csi
PLUGIN_TAG ?= dev

PROJECT_ROOT=$(CURDIR)
TEMP_DIR = $(PROJECT_ROOT)/_output
PLUGIN_DOCKER_DIR = $(PROJECT_ROOT)/images/plugin
PLUGIN_BINARY = ecfsplugin

VENDOR_DIR="$(PROJECT_ROOT)/src/vendor"
GOPATH = "$(PROJECT_ROOT):$(VENDOR_DIR)"

$(info Setting GOPATH to $(GOPATH))
$(info Elastifile CSI plugin image: $(PLUGIN_IMAGE_NAME) tag $(PLUGIN_TAG))

DEPLOYRUNNER_IMAGE_NAME = elastifileio/ecfs-provisioner-csi-deployrunner
DEPLOYRUNNER_DOCKER_DIR = $(PROJECT_ROOT)/images/deployrunner
DEPLOYRUNNER_COPY_DIR = $(DEPLOYRUNNER_DOCKER_DIR)/deploy

GKEDEPLOY_IMAGE_NAME = elastifileio/ecfs-provisioner-csi-gke-deploy
GKEDEPLOY_DOCKER_DIR = $(PROJECT_ROOT)/images/gke-deploy
GKEDEPLOY_COPY_DIR = $(GKEDEPLOY_DOCKER_DIR)/deploy

# Compile, create image and push it
all: image push

# Compile the plugin binary
binary:
	if [ ! -d $(VENDOR_DIR) ]; then pushd $(VENDOR_DIR)/..; dep ensure; popd; fi
	GOPATH=$(GOPATH) CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  $(TEMP_DIR)/$(PLUGIN_BINARY) $(PROJECT_ROOT)/src/ecfs

# Create docker image
image: binary
	cp $(TEMP_DIR)/$(PLUGIN_BINARY) $(PLUGIN_DOCKER_DIR)
	docker build -t $(PLUGIN_IMAGE_NAME):$(PLUGIN_TAG) $(PLUGIN_DOCKER_DIR)

# Push image to docker registry
push:
	docker push $(PLUGIN_IMAGE_NAME):$(PLUGIN_TAG)

deployrunner:
	mkdir -p $(DEPLOYRUNNER_COPY_DIR)
	cp -r deploy/* $(DEPLOYRUNNER_COPY_DIR)
	# kubectl version installed on the host running make is used in the resulting image
	cp -f $(shell which kubectl) $(DEPLOYRUNNER_DOCKER_DIR)/
	docker build -t $(DEPLOYRUNNER_IMAGE_NAME):$(PLUGIN_TAG) $(DEPLOYRUNNER_DOCKER_DIR)
	docker push $(DEPLOYRUNNER_IMAGE_NAME):$(PLUGIN_TAG)

gkedeploy:
	mkdir -p $(GKEDEPLOY_COPY_DIR)
	cp -r deploy/* $(GKEDEPLOY_COPY_DIR)
	cp -r gke-deploy/*.sh $(GKEDEPLOY_COPY_DIR)
	# kubectl version installed on the host running make is used in the resulting image
	cp -f $(shell which kubectl) $(GKEDEPLOY_DOCKER_DIR)/
	docker build -t $(GKEDEPLOY_IMAGE_NAME):$(PLUGIN_TAG) $(GKEDEPLOY_DOCKER_DIR)
	docker push $(GKEDEPLOY_IMAGE_NAME):$(PLUGIN_TAG)

# Remove previous build's artifacts
clean:
	go clean -r -x
	rm -f $(TEMP_DIR)/$(PLUGIN_BINARY)
	rm -f $(PLUGIN_DOCKER_DIR)/$(PLUGIN_BINARY)

	rm -rf $(DEPLOYRUNNER_COPY_DIR)/*
	rm -f $(DEPLOYRUNNER_DOCKER_DIR)/kubectl

	rm -rf $(GKEDEPLOY_COPY_DIR)/*
	rm -f $(GKEDEPLOY_DOCKER_DIR)/kubectl
