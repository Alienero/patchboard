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

package dnspod

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type Response struct {
	Status `json:"status"`
	Record `json:"record"`
}

type Status struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Create_at string `json:"create_at"`
}

type Record struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type DnsPodError struct {
	message string
}

func newDnspodError(message string) *DnsPodError {
	return &DnsPodError{message: message}
}

func (err *DnsPodError) Error() string {
	return "dnspod error:" + err.message
}

type Dnspod struct {
	client *http.Client
	vpool  *sync.Pool
}

func New(name, psw, domain_id string) *Dnspod {
	return &Dnspod{
		client: new(http.Client),
		vpool: &sync.Pool{
			New: func() interface{} {
				v := make(url.Values)
				v.Add("login_email", name)
				v.Add("login_password", psw)
				v.Add("format", "json")
				v.Add("domain_id", domain_id)
				return v
			},
		},
	}
}

func (d *Dnspod) CreateA(domain, ip string) (string, error) {
	post_url := "https://dnsapi.cn/Record.Create"
	v := d.vpool.Get().(url.Values)
	defer d.vpool.Put(v)
	v.Add("sub_domain", domain)
	v.Add("record_type", "A")
	v.Add("record_line", "默认")
	v.Add("value", ip)
	resp, err := d.client.PostForm(post_url, v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if dnsresp, err := d.unmarshal(resp.Body); err != nil {
		return "", err
	} else {
		return dnsresp.Record.Id, nil
	}
}

func (d *Dnspod) DelA(id string) error {
	post_url := "https://dnsapi.cn/Record.Remove"
	v := d.vpool.Get().(url.Values)
	defer d.vpool.Put(v)
	v.Add("record_id", id)
	resp, err := d.client.PostForm(post_url, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = d.unmarshal(resp.Body)
	return err
}

func (d *Dnspod) unmarshal(reader io.Reader) (*Response, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	dnsresp := new(Response)
	if err = json.Unmarshal(data, dnsresp); err != nil {
		return nil, err
	} else {
		if dnsresp.Code != "1" {
			return nil, newDnspodError(dnsresp.Message)
		} else {
			return dnsresp, nil
		}
	}
}
