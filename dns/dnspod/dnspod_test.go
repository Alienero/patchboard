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
	"fmt"
	"testing"
)

var (
	user_name = "example@mail.com"
	user_psw  = "user_password"
	domainID  = "domainID"
	recordID  string
)

func TestCteateA(t *testing.T) {
	dp := New(user_name, user_psw, domainID)
	if id, err := dp.CreateA("yim.so", "1.1.1.1"); err != nil {
		t.Error(err)
	} else {
		recordID = id
		fmt.Println("A ID is:", id)
	}
}

func TestDelA(t *testing.T) {
	dp := New(user_name, user_psw, domainID)
	if err := dp.DelA(recordID); err != nil {
		t.Error(err)
	}
}

// func Benchmark_interface2p(b *testing.B) {
// 	t := make([]byte, 10)
// 	for i := 0; i < b.N; i++ { //use b.N for looping
// 		func(v interface{}) {
// 			println(v.([]byte))
// 		}(t)
// 	}
// }

// func Benchmark_interface2v(b *testing.B) {
// 	t := make([]byte, 10)
// 	for i := 0; i < b.N; i++ { //use b.N for looping
// 		func(v interface{}) {
// 			println(v.(*[]byte))
// 		}(&t)
// 	}
// }
