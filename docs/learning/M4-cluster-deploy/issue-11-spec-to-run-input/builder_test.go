package issue11

import (
	"encoding/base64"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

func TestBuildRunInstancesInput_Basics(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: "t3.micro",
			AMIId:        "ami-123",
			KeyPair:      "my-key",
			Subnet:       "subnet-abc",
		},
	}
	in := BuildRunInstancesInput(cr)
	g.Expect(*in.ImageId).To(Equal("ami-123"))
	g.Expect(string(in.InstanceType)).To(Equal("t3.micro"))
	g.Expect(*in.KeyName).To(Equal("my-key"))
	g.Expect(*in.SubnetId).To(Equal("subnet-abc"))
}

func TestBuildRunInstancesInput_SecurityGroups(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Spec: computev1.Ec2InstanceSpec{
			SecurityGroups: []string{"sg-1", "sg-2"},
			AMIId:          "ami-x",
			InstanceType:   "t3.micro",
		},
	}
	in := BuildRunInstancesInput(cr)
	g.Expect(in.SecurityGroupIds).To(Equal([]string{"sg-1", "sg-2"}))
}

func TestBuildRunInstancesInput_Tags(t *testing.T) {
	g := NewWithT(t)
	cr := &computev1.Ec2Instance{
		Spec: computev1.Ec2InstanceSpec{
			AMIId:        "ami-x",
			InstanceType: "t3.micro",
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
			AMIId:        "ami-x",
			InstanceType: "t3.micro",
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
			AMIId:        "ami-x",
			InstanceType: "t3.micro",
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
