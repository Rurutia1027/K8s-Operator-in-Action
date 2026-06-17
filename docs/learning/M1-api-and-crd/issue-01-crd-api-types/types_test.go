package issue01

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

func TestEc2InstanceSpecStatusSeparation(t *testing.T) {
	g := NewWithT(t)

	spec := computev1.Ec2InstanceSpec{
		InstanceType: "t3.micro",
		AMIId:        "ami-placeholder",
		Region:       "eu-central-1",
	}
	status := computev1.Ec2InstanceStatus{
		InstanceID: "i-fake123",
		State:      "running",
	}

	cr := computev1.Ec2Instance{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		Spec:       spec,
		Status:     status,
	}

	g.Expect(cr.Spec.InstanceType).To(Equal("t3.micro"))
	g.Expect(cr.Status.InstanceID).To(Equal("i-fake123"))
	g.Expect(cr.Spec.AMIId).NotTo(Equal(cr.Status.InstanceID))
}

func TestEc2InstanceRoundTripJSON(t *testing.T) {
	g := NewWithT(t)

	original := &computev1.Ec2Instance{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "compute.cloud.com/v1",
			Kind:       "Ec2Instance",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "roundtrip",
			Namespace: "default",
		},
		Spec: computev1.Ec2InstanceSpec{
			InstanceType: "t3.micro",
			AMIId:        "ami-abc",
			Region:       "us-east-1",
			Tags:         map[string]string{"env": "learning"},
		},
	}

	data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(original)
	g.Expect(err).NotTo(HaveOccurred())

	restored := &computev1.Ec2Instance{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(data, restored)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(restored.Spec).To(Equal(original.Spec))
}

func TestSampleYAMLUnmarshals(t *testing.T) {
	g := NewWithT(t)

	path := filepath.Join("sample.yaml")
	raw, err := os.ReadFile(path)
	g.Expect(err).NotTo(HaveOccurred())

	var doc map[string]interface{}
	g.Expect(yaml.Unmarshal(raw, &doc)).To(Succeed())
	g.Expect(doc["kind"]).To(Equal("Ec2Instance"))

	cr := &computev1.Ec2Instance{}
	g.Expect(yaml.Unmarshal(raw, cr)).To(Succeed())
	g.Expect(cr.Spec.InstanceType).To(Equal("t3.micro"))
	g.Expect(cr.Spec.Region).To(Equal("eu-central-1"))
	g.Expect(cr.Spec.Tags["ManagedBy"]).To(Equal("ec2-operator"))
	g.Expect(cr.Status.InstanceID).To(BeEmpty(), "status should be empty in a new CR")
}

func TestCreatedInstanceInfoFields(t *testing.T) {
	g := NewWithT(t)

	info := computev1.CreatedInstanceInfo{
		InstanceID: "i-001",
		State:      "running",
		PublicIP:   "203.0.113.1",
	}
	g.Expect(info.InstanceID).NotTo(BeEmpty())
	g.Expect(info.State).To(Equal("running"))
}
