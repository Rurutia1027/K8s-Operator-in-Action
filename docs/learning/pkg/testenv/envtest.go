// Package testenv provides shared envtest bootstrap for learning issue tests.
package testenv

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	computev1 "github.com/shkatara/ec2Operator/api/v1"
)

// Environment holds envtest resources for controller learning tests.
type Environment struct {
	Ctx       context.Context
	Cancel    context.CancelFunc
	TestEnv   *envtest.Environment
	Config    *rest.Config
	K8sClient client.Client
}

// Setup boots envtest with the project CRD. Call Teardown in t.Cleanup.
func Setup(t *testing.T) *Environment {
	t.Helper()
	logf.SetLogger(zap.New(zap.WriteTo(os.Stderr), zap.UseDevMode(true)))

	g := NewWithT(t)
	ctx, cancel := context.WithCancel(context.Background())

	g.Expect(computev1.AddToScheme(scheme.Scheme)).To(Succeed())

	repoRoot := repoRoot(t)
	env := &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join(repoRoot, "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}
	if dir := firstEnvTestBinaryDir(repoRoot); dir != "" {
		env.BinaryAssetsDirectory = dir
	}

	cfg, err := env.Start()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cfg).NotTo(BeNil())

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	g.Expect(err).NotTo(HaveOccurred())

	te := &Environment{
		Ctx:       ctx,
		Cancel:    cancel,
		TestEnv:   env,
		Config:    cfg,
		K8sClient: k8sClient,
	}
	t.Cleanup(func() {
		cancel()
		g.Expect(env.Stop()).To(Succeed())
	})
	return te
}

func repoRoot(t *testing.T) string {
	t.Helper()
	if root := os.Getenv("REPO_ROOT"); root != "" {
		return root
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// Walk up to find go.mod
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repo root (go.mod)")
		}
		dir = parent
	}
}

func firstEnvTestBinaryDir(repoRoot string) string {
	basePath := filepath.Join(repoRoot, "bin", "k8s")
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return filepath.Join(basePath, entry.Name())
		}
	}
	return ""
}
