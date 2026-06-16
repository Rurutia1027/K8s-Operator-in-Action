#!/bin/sh 

kind create cluster --name operator-dev 

kind get clusters 

exit 0 
# 1) Create the project skeleton 

kubebuilder init --domain cloud.com --repo github.com/Rurutia1027/K8s-Operator-in-Action

# 2) Add your custom resource + controller stub 