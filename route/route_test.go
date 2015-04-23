// Copyright © 2015 Alienero. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package route

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"testing"
)

func TestRoute(t *testing.T) {
	check1 := "ts1"
	check2 := "ts2"
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(check1))
	}))
	defer ts1.Close()
	u1, err := url.Parse(ts1.URL)
	if err != nil {
		t.Error(err)
	}

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(check2))
	}))
	defer ts1.Close()
	u2, err := url.Parse(ts2.URL)
	if err != nil {
		t.Error(err)
	}

	route := NewRoute()
	route.AddRoute(check1, u1.Host)
	route.AddRoute(check2, u2.Host)

	listen_addr := "127.0.0.1:8800"
	go route.ListenAndServe(listen_addr)
	runtime.Gosched()

	r1, err := http.NewRequest("GET", u1.Scheme+"://"+listen_addr+u1.Path, nil)
	if err != nil {
		t.Error(err)
	}
	r1.Header.Set("Host", check1)
	r1.Host = check1
	r2, err := http.NewRequest("GET", u2.Scheme+"://"+listen_addr+u2.Path, nil)
	if err != nil {
		t.Error(err)
	}
	r2.Header.Set("Host", check2)
	r2.Host = check2
	client := new(http.Client)
	resp, err := client.Do(r1)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(data) != check1 {
		t.Error("not", check1, string(data))
	}

	resp, err = client.Do(r2)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(data) != check2 {
		t.Error("not", check2, string(data))
	}
}
