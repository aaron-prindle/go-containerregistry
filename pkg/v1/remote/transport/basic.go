// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"net/http"

	"github.com/aaron-prindle/go-containerregistry/pkg/authn"
)

type basicTransport struct {
	inner  http.RoundTripper
	auth   authn.Authenticator
	target string
}

var _ http.RoundTripper = (*basicTransport)(nil)

// RoundTrip implements http.RoundTripper
func (bt *basicTransport) RoundTrip(in *http.Request) (*http.Response, error) {
	hdr, err := bt.auth.Authorization()
	if err != nil {
		return nil, err
	}

	// http.Client handles redirects at a layer above the http.RoundTripper
	// abstraction, so to avoid forwarding Authorization headers to places
	// we are redirected, only set it when the authorization header matches
	// the host with which we are interacting.
	// In case of redirect http.Client can use an empty Host, check URL too.
	if in.Host == bt.target || in.URL.Host == bt.target {
		in.Header.Set("Authorization", hdr)
	}
	in.Header.Set("User-Agent", transportName)
	return bt.inner.RoundTrip(in)
}
