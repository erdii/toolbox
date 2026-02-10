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
	"fmt"
	"os"

	"github.com/erdii/toolbox/pkg/container/image"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:")
		fmt.Printf("%s <input archive> <output folder>\n", os.Args[0])
		os.Exit(1)
	}

	tarball := os.Args[1]
	dest := os.Args[2]

	if fi, err := os.Stat(tarball); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("image tarball does not exist. aborting!")
			os.Exit(2)
		} else {
			panic(err)
		}
	} else if !fi.Mode().IsRegular() {
		fmt.Println("image tarball is not a regular file. aborting!")
		os.Exit(3)
	}

	if _, err := os.Stat(dest); err == nil {
		fmt.Println("output folder already exists. aborting!")
		os.Exit(4)
	} else if !os.IsNotExist(err) {
		panic(err)
	}

	if err := image.ExtractImage(tarball, dest); err != nil {
		panic(err)
	}
}
