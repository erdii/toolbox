/*
Copyright (C) 2023 erdii <me@erdii.engineering>

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
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/erdii/toolbox/pkg/container/registry"
)

var rootCmd = &cobra.Command{
	Use:   "hash2tags",
	Short: "Lookup docker repo tags for a specific image hash.",
	Long:  `Lookup docker repo tags for a specific image hash.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		src := args[0]
		// split repo and hash
		repoHash := strings.SplitN(src, "@", 2)
		repo := repoHash[0]
		hash := repoHash[1]

		fmt.Printf("Image: `%s@%s`.\n", repo, hash)
		fmt.Println("Looking up tags...")

		tags, err := registry.FindTagsForImageHash(ctx, repo, hash)
		if err != nil {
			return err
		}

		fmt.Printf("Found these tags matching hash `%s`:\n", hash)
		for _, tag := range tags {
			fmt.Printf("%s\n", tag)
		}

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
