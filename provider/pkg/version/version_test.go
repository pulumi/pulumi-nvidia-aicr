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

package version

import "testing"

func TestSemverFallsBackWhenUnset(t *testing.T) {
	saveVersion(t)
	Version = ""
	if got := Semver(); got != devVersion {
		t.Fatalf("Semver() = %q, want %q", got, devVersion)
	}
}

func TestSemverUsesLinkerValue(t *testing.T) {
	saveVersion(t)
	Version = "0.3.3"
	if got := Semver(); got != "0.3.3" {
		t.Fatalf("Semver() = %q, want %q", got, "0.3.3")
	}
}

// The registry asks for "<product>/<version>" (e.g. pulumi-nvidia-aicr/1.2.3)
// so a vendor can attribute the traffic; assert the exact shape rather than
// just "non-empty", since a silently malformed agent still sends fine.
func TestUserAgentCarriesNameAndVersion(t *testing.T) {
	saveVersion(t)
	Version = "0.3.3"
	if got, want := UserAgent(), "pulumi-nvidia-aicr/0.3.3"; got != want {
		t.Fatalf("UserAgent() = %q, want %q", got, want)
	}
}

func TestUserAgentIsWellFormedInDevBuilds(t *testing.T) {
	saveVersion(t)
	Version = ""
	if got, want := UserAgent(), "pulumi-nvidia-aicr/0.0.1-dev"; got != want {
		t.Fatalf("UserAgent() = %q, want %q", got, want)
	}
}

// saveVersion snapshots Version and restores it after the test. It must be
// called (not deferred) before the test mutates Version: a deferred call
// would snapshot the already-mutated value.
func saveVersion(t *testing.T) {
	t.Helper()
	saved := Version
	t.Cleanup(func() { Version = saved })
}
