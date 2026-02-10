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
package image

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func newTarExtractor() *tarExtractor {
	return &tarExtractor{
		madeDirs: map[string]struct{}{},
	}
}

type tarExtractor struct {
	madeDirs map[string]struct{}
}

func (te *tarExtractor) extract(archiveReader io.Reader, dst string) error {
	tr := tar.NewReader(archiveReader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := te.ensureDir(target); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := te.ensureDir(filepath.Dir(target)); err != nil {
				return err
			}
			if err := te.writeFile(target, tr, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := te.ensureDir(filepath.Dir(target)); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("unhandled typeflag: %s - %c\n", target, header.Typeflag))
		}
	}
}

func (te *tarExtractor) writeFile(name string, tr *tar.Reader, perm fs.FileMode) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, tr); err != nil {
		return err
	}
	return nil
}

func (te *tarExtractor) ensureDir(name string) error {
	if _, ok := te.madeDirs[name]; ok {
		return nil
	}

	if _, err := os.Stat(name); os.IsNotExist(err) {
		if err := os.MkdirAll(name, 0755); err != nil {
			return err
		}
	} else {
		return err
	}

	te.madeDirs[name] = struct{}{}
	return nil
}
