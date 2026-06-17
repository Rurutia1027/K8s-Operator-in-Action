/*
Copyright 2026.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controller

import (
	"encoding/base64"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
)

func TestBuildRunInstancesInput(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: testInstanceType,
			AMIId:        "ami-123",
			KeyPair:      "my-key",
			Subnet:       "subnet-abc",
		},
	}

	in := BuildRunInstancesInput(cr)
	g.Expect(*in.ImageId).To(Equal("ami-123"))
	g.Expect(string(in.InstanceType)).To(Equal(testInstanceType))
	g.Expect(*in.KeyName).To(Equal("my-key"))
	g.Expect(*in.SubnetId).To(Equal("subnet-abc"))
}

func TestBuildRunInstancesInput_SecurityGroups(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Spec: computev1.Ec2InstanceSpec{
			SecurityGroups: []string{"sg-1", "sg-2"},
			AMIId:          testAMIX,
			InstanceType:   testInstanceType,
		},
	}
	in := BuildRunInstancesInput(cr)
	g.Expect(in.SecurityGroupIds).To(Equal([]string{"sg-1", "sg-2"}))
}
func TestBuildRunInstancesInput_Tags(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Spec: computev1.Ec2InstanceSpec{
			AMIId:        testAMIX,
			InstanceType: testInstanceType,
			Tags:         map[string]string{"Name": "web", "env": "prod"},
		},
	}
	in := BuildRunInstancesInput(cr)
	g.Expect(in.TagSpecifications).To(HaveLen(1))
	g.Expect(in.TagSpecifications[0].Tags).To(HaveLen(2))
}
func TestBuildRunInstancesInput_UserDataBase64(t *testing.T) {
	g := NewWithT(t)
	script := "#!/bin/bash\necho hello"
	cr := &computev1.Ec2Instance{
		Spec: computev1.Ec2InstanceSpec{
			AMIId:        testAMIX,
			InstanceType: testInstanceType,
			UserData:     script,
		},
	}
	in := BuildRunInstancesInput(cr)
	decoded, err := base64.StdEncoding.DecodeString(*in.UserData)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(decoded)).To(Equal(script))
}
func TestBuildRunInstancesInput_Storage(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: "storage"},
		Spec: computev1.Ec2InstanceSpec{
			AMIId:        testAMIX,
			InstanceType: testInstanceType,
			Storage: computev1.StorageConfig{
				RootVolume: computev1.VolumeConfig{Size: 30, Type: "gp3", Encrypted: true},
				AdditionalVolumes: []computev1.VolumeConfig{
					{Size: 100, Type: "gp3", DeviceName: "/dev/sdf", Encrypted: true},
				},
			},
		},
	}
	in := BuildRunInstancesInput(cr)
	g.Expect(in.BlockDeviceMappings).To(HaveLen(2))
	g.Expect(*in.BlockDeviceMappings[0].Ebs.VolumeSize).To(Equal(int32(30)))
	g.Expect(*in.BlockDeviceMappings[1].DeviceName).To(Equal("/dev/sdf"))
}
