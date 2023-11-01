// Package common is the support for HPSA
// Copyright 2021 The ChromiumOS Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
// init HPSA classes
package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"go.chromium.org/tast/core/errors"
)

// Welcome  set up json
type Welcome struct {
	Welcome []welcomeInfo
}
type welcomeInfo struct {
	Name  string
	Class string
	NTH   int
}

// ProfileJSON  set up json for profile info
type ProfileJSON struct {
	ProfileJSON []profileInfo `json:"profile"`
}
type profileInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Dashboard set up json
type Dashboard struct {
	Dashboard []dashboardInfo `json:"dashboard"`
}
type dashboardInfo struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	NTH   int    `json:"nth"`
}

// GetJSON to get the json from the path
func GetJSON(name, path string) (string, int, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return "", 0, errors.Wrapf(err, "can not read json from : %q ", path)
		// s.Fatal("Can not read json from hpsa.json error: ", err)

	}
	var welcome Welcome
	json.Unmarshal(jsonData, &welcome)
	for _, val := range welcome.Welcome {
		if val.Name == name {
			return val.Class, val.NTH, nil
		}
	}
	return "", 0, errors.Wrapf(err, "can not find json data : %q ", name)

}

// GetProfileJSON to get the json from the path
func GetProfileJSON(id, path string) (string, string, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", errors.Wrapf(err, "can not read json from : %q ", path)
		// s.Fatal("Can not read json from hpsa.json error: ", err)

	}
	var profile ProfileJSON
	json.Unmarshal(jsonData, &profile)
	for _, val := range profile.ProfileJSON {
		if val.ID == id {
			fmt.Sprintf("Find the target json %v, %v", val.Username, val.Password)
			return val.Username, val.Password, nil
		}
	}
	return "", "", errors.Wrapf(err, "can not find json data : %q ", profile)

}

// GetJSONDashboard to get the json from the path
func GetJSONDashboard(name, path string) (string, int, error) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return "", 0, errors.Wrapf(err, "can not read json from : %q ", path)
		// s.Fatal("Can not read json from hpsa.json error: ", err)

	}
	var dashboard Dashboard
	json.Unmarshal(jsonData, &dashboard)
	for _, val := range dashboard.Dashboard {
		if val.Name == name {
			return val.Class, val.NTH, nil
		}
	}
	return "", 0, errors.Wrapf(err, "can not find json data : %q ", name)

}
