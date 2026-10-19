// Copyright 2026, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aicr

import (
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/rest"

	"github.com/pulumi/pulumi-nvidia-aicr/provider/pkg/version"
)

// A minimal, never-dialed kubeconfig: BuildKubeClient only parses the file
// when given an explicit path, so no API server is involved.
const fixtureKubeconfig = `apiVersion: v1
kind: Config
clusters:
  - name: fixture
    cluster:
      server: https://127.0.0.1:6443
contexts:
  - name: fixture
    context:
      cluster: fixture
      user: fixture
current-context: fixture
users:
  - name: fixture
    user:
      token: fixture-token
`

func TestKubeClientConfigSetsProviderUserAgent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(path, []byte(fixtureKubeconfig), 0o600); err != nil {
		t.Fatalf("writing fixture kubeconfig: %v", err)
	}

	config, err := kubeClientConfig(path)
	if err != nil {
		t.Fatalf("kubeClientConfig: %v", err)
	}

	if got, want := config.UserAgent, version.UserAgent(); got != want {
		t.Errorf("UserAgent = %q, want %q", got, want)
	}
	// The failure this guards against is silent: dropping the assignment
	// leaves a perfectly working client that reports itself as client-go,
	// which is exactly the unattributable traffic the registry asks
	// packages to avoid.
	if config.UserAgent == rest.DefaultKubernetesUserAgent() {
		t.Errorf("UserAgent is client-go's default (%q); the provider's agent was not applied",
			config.UserAgent)
	}
}

// The kubeconfig the SDK resolved must otherwise be untouched — the copy
// exists only to carry the agent.
func TestKubeClientConfigPreservesServer(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(path, []byte(fixtureKubeconfig), 0o600); err != nil {
		t.Fatalf("writing fixture kubeconfig: %v", err)
	}

	config, err := kubeClientConfig(path)
	if err != nil {
		t.Fatalf("kubeClientConfig: %v", err)
	}
	if got, want := config.Host, "https://127.0.0.1:6443"; got != want {
		t.Errorf("Host = %q, want %q", got, want)
	}
	if config.BearerToken != "fixture-token" {
		t.Errorf("BearerToken was not carried through from the kubeconfig")
	}
}
