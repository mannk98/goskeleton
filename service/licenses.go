// Copyright © 2021 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Parts inspired by https://github.com/ryanuber/go-license

package service

import (
	"fmt"
	"github.com/spf13/cobra"
	"goske/models"
	"strings"
)

// Licenses contains all possible licenses a user can choose from.
var Licenses = make(map[string]models.License)

func init() {
	// Allows a user to not use a license.
	Licenses["none"] = models.License{"None", []string{"none", "false"}, "", ""}

	initApache2()
	initMit()
	initBsdClause3()
	initBsdClause2()
	initGpl2()
	initGpl3()
	initLgpl()
	initAgpl()
}

// getLicense returns license specified by user in flag or in config.
// If user didn't specify the license, it returns none
//
// TODO: Inspect project for existing license
func getLicense(userLicense, license_header, license_text string) models.License {
	// If explicitly flagged, use that.
	if userLicense != "" {
		return findLicense(userLicense)
	}
	if license_header != "" && license_text != "" {
		return models.License{Header: license_header,
			Text: license_text}
	}
	// If user didn't set any license, use none by default
	return Licenses["none"]
}

func copyrightLine(year, author string) string {
	return "Copyright © " + year + " " + author
}

// findLicense looks for License object of built-in licenses.
// If it didn't find license, then the app will be terminated and
// error will be printed.
func findLicense(name string) models.License {
	found := matchLicense(name)
	if found == "" {
		cobra.CheckErr(fmt.Errorf("unknown license: " + name))
	}
	return Licenses[found]
}

// matchLicense compares the given a license name
// to PossibleMatches of all built-in licenses.
// It returns blank string, if name is blank string or it didn't find
// then appropriate match to name.
func matchLicense(name string) string {
	if name == "" {
		return ""
	}

	for key, lic := range Licenses {
		for _, match := range lic.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}
