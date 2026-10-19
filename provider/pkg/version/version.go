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

// Name is the provider's plugin name, and the product token of its
// User-Agent.
const Name = "pulumi-nvidia-aicr"

// devVersion is reported when the binary was built without the release
// ldflags (`go build ./...`, `go test`, a local `make provider` with no
// VERSION). Keeping the fallback here means the plugin version and the
// User-Agent can never disagree about which build is running.
const devVersion = "0.0.1-dev"

// Version is set at build time via ldflags.
var Version string

// Semver returns the build's version, falling back to devVersion for a
// build without the release ldflags.
func Semver() string {
	if Version == "" {
		return devVersion
	}
	return Version
}

// UserAgent identifies this provider, and the build of it, on the API calls
// the provider makes itself — currently the Kubernetes API server (see
// aicr.sweepValidatorJobs). The registry asks every package to send one so
// that traffic originating from Pulumi can be attributed as such instead of
// landing in a vendor's logs under a client library's default agent.
//
// Scope, so the next reader does not assume more coverage than exists:
// recipe component deployment goes out through the Pulumi Kubernetes
// provider in its own process, and the AICR SDK builds its own clients for
// the snapshot/validator harness. Neither exposes a user-agent hook today
// (the SDK has no WithUserAgent; tagging that traffic needs an upstream
// change), so neither carries this string.
func UserAgent() string {
	return Name + "/" + Semver()
}
