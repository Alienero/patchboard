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
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/Alienero/patchboard/util"
)

type Route struct {
	routes *util.Map
	rp     *httputil.ReverseProxy
	sync.RWMutex
}

func NewRoute() *Route {
	r := &Route{
		RWMutex: sync.RWMutex{},
		routes:  util.NewMap(),
	}
	r.rp = &httputil.ReverseProxy{Director: r.director}
	return r
}

func (r *Route) AddRoute(host string, ip string) {
	r.routes.Add(host, ip)
}

func (r *Route) DelRoute(host string) {
	r.routes.Del(host)
}

func (r *Route) GetRoute(host string) string {
	if v := r.routes.Get(host); v == nil {
		return ""
	} else if vt, ok := v.(string); ok {
		return vt
	} else {
		return ""
	}
}

func (r *Route) ListenAndServe(addr string) {
	http.ListenAndServe(addr, r)
}

func (r *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if ip := r.GetRoute(req.Host); ip == "" {
		http.Error(w, "Host not found"+req.Host, 403)
		return
	}
	r.rp.ServeHTTP(w, req)
}

func (r *Route) director(req *http.Request) {
	addr := r.GetRoute(req.Host)
	target := url.URL{Scheme: "http", Host: addr, Path: "/"}
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
