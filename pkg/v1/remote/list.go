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

package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aaron-prindle/go-containerregistry/pkg/authn"
	"github.com/aaron-prindle/go-containerregistry/pkg/name"
	"github.com/aaron-prindle/go-containerregistry/pkg/v1/remote/transport"
)

type Tags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// TODO(jonjohnsonjr): return []name.Tag?
func List(repo name.Repository, auth authn.Authenticator, t http.RoundTripper) ([]string, error) {
	scopes := []string{repo.Scope(transport.PullScope)}
	tr, err := transport.New(repo.Registry, auth, t, scopes)
	if err != nil {
		return nil, err
	}

	uri := url.URL{
		Scheme: repo.Registry.Scheme(),
		Host:   repo.Registry.RegistryStr(),
		Path:   fmt.Sprintf("/v2/%s/tags/list", repo.RepositoryStr()),
	}

	client := http.Client{Transport: tr}
	resp, err := client.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := CheckError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	tags := Tags{}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	return tags.Tags, nil
}
