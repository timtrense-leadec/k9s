// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package watch

import (
	"testing"

	"github.com/derailed/k9s/internal/client"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

func TestFactoryForResourceDropsNamespacePrefix(t *testing.T) {
	conn := newFakeConnection(t)
	f := NewFactory(conn)
	f.stopChan = make(chan struct{})
	t.Cleanup(func() {
		close(f.stopChan)
	})

	inf, err := f.ForResource("kube-system", client.NsGVR)
	require.NoError(t, err)
	require.NotNil(t, inf)

	if _, ok := f.factories[client.BlankNamespace]; !ok {
		t.Fatalf("expected factory entry for blank namespace")
	}
	if _, ok := f.factories["kube-system"]; ok {
		t.Fatalf("unexpected factory entry scoped to kube-system")
	}
}

type fakeConnection struct {
	dyn dynamic.Interface
}

func newFakeConnection(t *testing.T) *fakeConnection {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))

	return &fakeConnection{dyn: dynamicfake.NewSimpleDynamicClient(scheme)}
}

func (f *fakeConnection) CanI(string, *client.GVR, string, []string) (bool, error) {
	return true, nil
}

func (*fakeConnection) Config() *client.Config { return client.NewConfig(nil) }

func (*fakeConnection) ConnectionOK() bool { return true }

func (*fakeConnection) Dial() (kubernetes.Interface, error) { return nil, nil }

func (*fakeConnection) DialLogs() (kubernetes.Interface, error) { return nil, nil }

func (*fakeConnection) SwitchContext(string) error { return nil }

func (*fakeConnection) CachedDiscovery() (*disk.CachedDiscoveryClient, error) {
	return nil, nil
}

func (*fakeConnection) RestConfig() (*rest.Config, error) { return nil, nil }

func (*fakeConnection) MXDial() (*versioned.Clientset, error) { return nil, nil }

func (f *fakeConnection) DynDial() (dynamic.Interface, error) { return f.dyn, nil }

func (*fakeConnection) HasMetrics() bool { return false }

func (*fakeConnection) ValidNamespaceNames() (client.NamespaceNames, error) {
	return nil, nil
}

func (*fakeConnection) IsValidNamespace(string) bool { return true }

func (*fakeConnection) ServerVersion() (*version.Info, error) { return nil, nil }

func (*fakeConnection) CheckConnectivity() bool { return true }

func (*fakeConnection) ActiveContext() string { return "" }

func (*fakeConnection) ActiveNamespace() string { return client.DefaultNamespace }

func (*fakeConnection) IsActiveNamespace(string) bool { return false }
