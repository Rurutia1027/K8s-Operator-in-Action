// Reference: study this file, then implement in api/v1/ec2instance_types.go
package reference

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ec2InstanceSpec defines the desired state (what the user wants).
type Ec2InstanceSpec struct {
	InstanceType      string            `json:"instanceType"`
	AMIId             string            `json:"amiId"`
	Region            string            `json:"region"`
	AvailabilityZone  string            `json:"availabilityZone,omitempty"`
	KeyPair           string            `json:"keyPair,omitempty"`
	SecurityGroups    []string          `json:"securityGroups,omitempty"`
	Subnet            string            `json:"subnet,omitempty"`
	UserData          string            `json:"userData,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
	Storage           StorageConfig     `json:"storage,omitempty"`
	AssociatePublicIP bool              `json:"associatePublicIP,omitempty"`
}

// Ec2InstanceStatus defines the observed state (what the operator reports).
type Ec2InstanceStatus struct {
	InstanceID string `json:"instanceId,omitempty"`
	State      string `json:"state,omitempty"`
	PublicIP   string `json:"publicIP,omitempty"`
	PrivateIP  string `json:"privateIP,omitempty"`
	PublicDNS  string `json:"publicDNS,omitempty"`
	PrivateDNS string `json:"privateDNS,omitempty"`
	LaunchTime *metav1.Time `json:"launchTime,omitempty"`
}

type StorageConfig struct {
	RootVolume        VolumeConfig   `json:"rootVolume"`
	AdditionalVolumes []VolumeConfig `json:"additionalVolumes,omitempty"`
}

type VolumeConfig struct {
	Size       int32  `json:"size"`
	Type       string `json:"type,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	Encrypted  bool   `json:"encrypted,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Ec2Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Ec2InstanceSpec   `json:"spec,omitempty"`
	Status            Ec2InstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type Ec2InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ec2Instance `json:"items"`
}

type CreatedInstanceInfo struct {
	InstanceID string `json:"instanceId"`
	PublicIP   string `json:"publicIP"`
	PrivateIP  string `json:"privateIP"`
	PublicDNS  string `json:"publicDNS"`
	PrivateDNS string `json:"privateDNS"`
	State      string `json:"state"`
}
