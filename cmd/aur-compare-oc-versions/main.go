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

package main

import (
	"log"

	"github.com/erdii/toolbox/pkg/utils"
)

func main() {
	releaseVersion, err := utils.GetLatestAvailableV4Version()
	dieOnErr(err)
	log.Println("Latest release version:", releaseVersion)

	packageVersion, err := utils.GetPackageVersion()
	dieOnErr(err)
	log.Println("Current package version:", packageVersion)

	if releaseVersion == packageVersion {
		log.Println("Nothing to do.")
		return
	}

	log.Println("Versions do not match!")
}

func dieOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
