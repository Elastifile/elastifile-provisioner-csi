#!/bin/bash

kubectl create -f csidriver.yaml --validate=false 
kubectl create -f csinodeinfo.yaml --validate=false 
