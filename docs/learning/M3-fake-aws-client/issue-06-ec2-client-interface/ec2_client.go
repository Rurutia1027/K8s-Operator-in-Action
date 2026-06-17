// Package issue06 defines EC2Client and a fake for local testing (Issue #6).
package issue06

import (
	"context"
	"fmt"
	"sync"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

// InstanceDetails is a cloud-agnostic view of an EC2 instance.
type InstanceDetails struct {
	InstanceID string
	State      string
	PublicIP   string
	PrivateIP  string
	PublicDNS  string
	PrivateDNS string
}

// EC2Client abstracts AWS EC2 operations for the reconciler.
type EC2Client interface {
	RunInstance(ctx context.Context, instance *computev1.Ec2Instance) (*computev1.CreatedInstanceInfo, error)
	TerminateInstance(ctx context.Context, instanceID, region string) error
	DescribeInstance(ctx context.Context, instanceID, region string) (exists bool, details *InstanceDetails, err error)
}

// FakeEC2Client is an in-memory EC2 for tests and local dev (USE_FAKE_EC2=true).
type FakeEC2Client struct {
	mu        sync.RWMutex
	instances map[string]*InstanceDetails
	counter   int
}

func NewFakeEC2Client() *FakeEC2Client {
	return &FakeEC2Client{
		instances: make(map[string]*InstanceDetails),
		counter:   0,
	}
}

func (f *FakeEC2Client) RunInstance(_ context.Context, instance *computev1.Ec2Instance) (*computev1.CreatedInstanceInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.counter++
	id := fmt.Sprintf("i-fake%03d", f.counter)
	details := &InstanceDetails{
		InstanceID: id,
		State:      "running",
		PublicIP:   "203.0.113.10",
		PrivateIP:  "10.0.0.10",
		PublicDNS:  id + ".example.com",
		PrivateDNS: "ip-10-0-0-10.internal",
	}
	f.instances[id] = details

	_ = instance // region/spec validated in Issue #11
	return &computev1.CreatedInstanceInfo{
		InstanceID: id,
		State:      details.State,
		PublicIP:   details.PublicIP,
		PrivateIP:  details.PrivateIP,
		PublicDNS:  details.PublicDNS,
		PrivateDNS: details.PrivateDNS,
	}, nil
}

func (f *FakeEC2Client) TerminateInstance(_ context.Context, instanceID, _ string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.instances[instanceID]; !ok {
		return fmt.Errorf("instance %s not found", instanceID)
	}
	delete(f.instances, instanceID)
	return nil
}

func (f *FakeEC2Client) DescribeInstance(_ context.Context, instanceID, _ string) (bool, *InstanceDetails, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	d, ok := f.instances[instanceID]
	if !ok {
		return false, nil, nil
	}
	copy := *d
	return true, &copy, nil
}

// SetInstanceState changes state for drift-detection tests (Issue #10).
func (f *FakeEC2Client) SetInstanceState(instanceID, state string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if d, ok := f.instances[instanceID]; ok {
		d.State = state
	}
}

// DeleteInstance removes instance without terminate (simulate external deletion).
func (f *FakeEC2Client) DeleteInstance(instanceID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.instances, instanceID)
}
