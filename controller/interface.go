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

package controller

type control struct {
}

type user struct {
	id string
}

// Register a new user.
func (ctrl *control) register(mail, psw string) error {
	return nil
}

// User login.
func (ctrl *control) login(id, psw string) (bool, error) {
	return false, nil
}

// User recharge.

// Cost.
func (ctrl *control) cost(appid string) {
}

// Add a Image.
func (ctrl *control) addImage(path string) (imgID string, err error) {
}

// Start App.
func (ctrl *control) startApp(mem, cpu, space, nodeNum int) (appid []string, err error) {
}

// Delete app.
func (ctrl *control) delApp(appid string) error {

}

// DataBase.
// Create a new database.

// Scaling a database.
