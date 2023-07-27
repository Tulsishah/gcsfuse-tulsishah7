// Copyright 2023 Google Inc. All Rights Reserved.
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

// Provides integration tests when --rename-dir-limit flag is set.

package retry_test

//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"net/http"
//	"net/http/httputil"
//	"net/url"
//	"os"
//	"strings"
//	"testing"
//
//	"cloud.google.com/go/storage"
//	"google.golang.org/api/option"
//	htransport "google.golang.org/api/transport/http"
//)
//
//type resources struct {
//	bucket       *storage.BucketAttrs
//	object       *storage.ObjectAttrs
//	notification *storage.Notification
//	hmacKey      *storage.HMACKey
//}
//
//type emulatorTest struct {
//	*testing.T
//	name          string
//	id            string // ID to pass as a header in the test execution
//	resources     resources
//	host          *url.URL // set the path when using; path is not guaranteed between calls
//	wrappedClient *storage.Client
//}
//
//// retryTestRoundTripper sends the retry test ID to the emulator with each request
//type retryTestRoundTripper struct {
//	*testing.T
//	rt     http.RoundTripper
//	testID string
//}
//
//func (wt *retryTestRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
//	r.Header.Set("x-retry-test-id", wt.testID)
//
//	requestDump, err := httputil.DumpRequest(r, false)
//	if err != nil {
//		wt.Logf("error creating request dump: %v", err)
//	}
//
//	resp, err := wt.rt.RoundTrip(r)
//	if err != nil {
//		wt.Logf("roundtrip error (may be expected): %v\nrequest: %s", err, requestDump)
//	}
//	return resp, err
//}
//
//func wrappedClient(t *testing.T, testID string) (*storage.Client, error) {
//	ctx := context.Background()
//	base := http.DefaultTransport
//
//	trans, err := htransport.NewTransport(ctx, base, option.WithoutAuthentication(), option.WithUserAgent("custom-user-agent"))
//	if err != nil {
//		return nil, fmt.Errorf("failed to create http client: %v", err)
//	}
//
//	c := http.Client{Transport: trans}
//
//	// Add RoundTripper to the created HTTP client
//	wrappedTrans := &retryTestRoundTripper{rt: c.Transport, testID: testID, T: t}
//	c.Transport = wrappedTrans
//
//	// Supply this client to storage.NewClient
//	// STORAGE_EMULATOR_HOST takes care of setting the correct endpoint
//	client, err := storage.NewClient(ctx, option.WithHTTPClient(&c))
//	return client, err
//}
//
//func (et *emulatorTest) create(instructions map[string][]string) {
//	c := http.DefaultClient
//	data := struct {
//		Instructions map[string][]string `json:"instructions"`
//	}{
//		Instructions: instructions,
//	}
//
//	buf := new(bytes.Buffer)
//	if err := json.NewEncoder(buf).Encode(data); err != nil {
//		et.Fatalf("encoding request: %v", err)
//	}
//
//	et.host.Path = "retry_test"
//	resp, err := c.Post(et.host.String(), "application/json", buf)
//	if err != nil || resp.StatusCode != 200 {
//		et.Fatalf("creating retry test: err: %v, resp: %+v", err, resp)
//	}
//	defer func() {
//		closeErr := resp.Body.Close()
//		if err == nil {
//			err = closeErr
//		}
//	}()
//	testRes := struct {
//		TestID string `json:"id"`
//	}{}
//	if err := json.NewDecoder(resp.Body).Decode(&testRes); err != nil {
//		et.Fatalf("decoding test ID: %v", err)
//	}
//
//	et.id = testRes.TestID
//
//	// Create wrapped client which will send emulator instructions
//	et.host.Path = ""
//	client, err := wrappedClient(et.T, et.id)
//	if err != nil {
//		et.Fatalf("creating wrapped client: %v", err)
//	}
//	et.wrappedClient = client
//}
//
//// Deletes a retry test resource
//func (et *emulatorTest) delete() {
//	et.host.Path = strings.Join([]string{"retry_test", et.id}, "/")
//	c := http.DefaultClient
//	req, err := http.NewRequest("DELETE", et.host.String(), nil)
//	if err != nil {
//		et.Errorf("creating request: %v", err)
//	}
//	resp, err := c.Do(req)
//	if err != nil || resp.StatusCode != 200 {
//		et.Errorf("deleting test: err: %v, resp: %+v", err, resp)
//	}
//}
//
//func TestRetryConformance(t *testing.T) {
//	host := os.Getenv("STORAGE_EMULATOR_HOST")
//	if host == "" {
//		t.Skip("This test must use the testbench emulator; set STORAGE_EMULATOR_HOST to run.")
//	}
//	endpoint, err := url.Parse(host)
//	if err != nil {
//		t.Fatalf("error parsing emulator host (make sure it includes the scheme such as http://host): %v", err)
//	}
//
//	ctx := context.Background()
//
//	// Create non-wrapped client to use for setup steps.
//	client, err := storage.NewClient(ctx)
//
//	subtest := &emulatorTest{T: t, name: testName, host: endpoint}
//	subtest.create(map[string][]string{
//		method.Name: instructions.Instructions,
//	})
//
//}
