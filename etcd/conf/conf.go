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

package conf

type conf struct {
	Admin string

	// Default configs.
	Mysql mysql
	// etcd configs.
	Etcd etcd
	// Dns configs.
	Dns dns
	// Docker config.
	Docker docker
	// Controller config.
	control Control
}

type mysql struct {
	Space int
}

type etcd struct {
	Addr string

	// Etcd workspace.
	Dir string

	// Keys.
	MySqlDir string
	Route    string
	Git      string
}

type dns struct {
}

type docker struct {
}

type control struct {
}
