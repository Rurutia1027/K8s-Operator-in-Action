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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	computev1 "github.com/Rurutia1027/K8s-Operator-in-Action/api/v1"
)

// BuildRunInstancesInput converts a CR spec into AWS SDK input (no API call).
func BuildRunInstancesInput(cr *computev1.Ec2Instance) *ec2.RunInstancesInput {
	spec := cr.Spec
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(spec.AMIId),
		InstanceType: ec2types.InstanceType(spec.InstanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	}

	if spec.KeyPair != "" {
		input.KeyName = aws.String(spec.KeyPair)
	}
	if spec.Subnet != "" {
		input.SubnetId = aws.String(spec.Subnet)
	}
	if len(spec.SecurityGroups) > 0 {
		input.SecurityGroupIds = spec.SecurityGroups
	}
	if spec.UserData != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(spec.UserData))
		input.UserData = aws.String(encoded)
	}
	if len(spec.Tags) > 0 {
		tags := make([]ec2types.Tag, 0, len(spec.Tags))
		for k, v := range spec.Tags {
			tags = append(tags, ec2types.Tag{Key: aws.String(k), Value: aws.String(v)})
		}
		input.TagSpecifications = []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeInstance,
				Tags:         tags,
			},
		}
	}

	if spec.Storage.RootVolume.Size > 0 {
		volType := ec2types.VolumeTypeGp3
		if spec.Storage.RootVolume.Type != "" {
			volType = ec2types.VolumeType(spec.Storage.RootVolume.Type)
		}
		input.BlockDeviceMappings = append(input.BlockDeviceMappings, ec2types.BlockDeviceMapping{
			DeviceName: aws.String("/dev/xvda"),
			Ebs: &ec2types.EbsBlockDevice{
				VolumeSize:          aws.Int32(spec.Storage.RootVolume.Size),
				VolumeType:          volType,
				Encrypted:           aws.Bool(spec.Storage.RootVolume.Encrypted),
				DeleteOnTermination: aws.Bool(true),
			},
		})
	}

	for _, vol := range spec.Storage.AdditionalVolumes {
		volType := ec2types.VolumeTypeGp3
		if vol.Type != "" {
			volType = ec2types.VolumeType(vol.Type)
		}
		device := vol.DeviceName
		if device == "" {
			device = "/dev/sdf"
		}
		input.BlockDeviceMappings = append(input.BlockDeviceMappings, ec2types.BlockDeviceMapping{
			DeviceName: aws.String(device),
			Ebs: &ec2types.EbsBlockDevice{
				VolumeSize:          aws.Int32(vol.Size),
				VolumeType:          volType,
				Encrypted:           aws.Bool(vol.Encrypted),
				DeleteOnTermination: aws.Bool(false),
			},
		})
	}

	return input
}
