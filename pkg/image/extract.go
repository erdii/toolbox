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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type imageManifest struct {
	Layers []string `json:"Layers"`
}

func ExtractImage(tarballPath, dst string) error {
	tmpdir, err := os.MkdirTemp("", "extract-image.*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	f, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer f.Close()

	te := newTarExtractor()
	if err := te.extract(f, tmpdir); err != nil {
		return err
	}

	manifestBytes, err := os.ReadFile(filepath.Join(tmpdir, "manifest.json"))
	if err != nil {
		return err
	}
	manifests := &[]imageManifest{}
	if err := json.Unmarshal(manifestBytes, manifests); err != nil {
		return err
	}
	if len(*manifests) != 1 {
		return fmt.Errorf("image did not contain exactly one manifest. contained %d instead", len(*manifests))
	}

	for _, layer := range (*manifests)[0].Layers {
		fmt.Println("extracting layer", layer)

		f, err := os.Open(filepath.Join(tmpdir, layer))
		if err != nil {
			return err
		}
		defer f.Close()
		te := newTarExtractor()
		if err := te.extract(f, dst); err != nil {
			return err
		}
	}

	return nil
}
