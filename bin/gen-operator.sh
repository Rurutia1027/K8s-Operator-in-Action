#!/bin/sh 

export PATH="$HOME/go/bin:$PATH"

source ~/.zshrc

# gen operator skeleton 
kubebuilder init --domain cloud.com --repo github.com/Rurutia1027/K8s-Operator-in-Action

# gen api --> controller skeleton 
kubebuilder create api --group compute --version v1 --kind Ec2Instance