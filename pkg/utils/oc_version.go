/*
Copyright (C) 2021 erdii <me@erdii.engineering>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"
)

const latestV4ReleaseTxtURL = "https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/stable/release.txt"

var (
	releaseVersionRegex        = regexp.MustCompile(`Version:\s{0,}(4\.\d{1,}\.\d{1,})`)
	packageVersionRegex        = regexp.MustCompile(`pkgver=\s{0,}(4\.\d{1,}\.\d{1,})`)
	ErrCouldNotExtractVersion  = errors.New("could not extract version string from release file")
	ErrCouldNotValidateVersion = errors.New("could not validate semver string from release file")
)

// GetLatestAvailableV4Version fetches the latest release.txt and extracts/validates the version string
func GetLatestAvailableV4Version() (string, error) {
	resp, err := http.Get(latestV4ReleaseTxtURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := releaseVersionRegex.FindStringSubmatch(string(body))

	// the 2nd match should be the version number
	if len(matches) < 2 {
		return "", ErrCouldNotExtractVersion
	}

	version := strings.Trim(matches[1], " \n\t\r")

	if !semver.IsValid(fmt.Sprintf("v%s", version)) {
		return "", ErrCouldNotValidateVersion
	}

	return version, nil
}

func GetPackageVersion() (string, error) {
	body, err := ioutil.ReadFile("./PKGBUILD")
	if err != nil {
		return "", err
	}

	matches := packageVersionRegex.FindStringSubmatch(string(body))

	// the 2nd match should be the version number
	if len(matches) < 2 {
		return "", ErrCouldNotExtractVersion
	}

	version := strings.Trim(matches[1], " \n\t\r")

	if !semver.IsValid(fmt.Sprintf("v%s", version)) {
		return "", ErrCouldNotValidateVersion
	}

	return version, nil
}
