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

package registry

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// FindTagsForImageHash queries all tags in a repo and compares their manifest hashes against the given hash. Uses standard docker authentication (ie ~/.docker/config.json).
func FindTagsForImageHash(ctx context.Context, image, hash string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		withDefaultKeychainAuth := remote.WithAuthFromKeychain(authn.DefaultKeychain)

		repo, err := name.NewRepository(image)
		if err != nil {
			yield("", err)
			return
		}

		for tag, err := range listTags(ctx, repo, withDefaultKeychainAuth) {
			if err != nil {
				yield("", err)
				return
			}

			manifestHash, err := getManifestHash(ctx, fmt.Sprintf("%s:%s", image, tag), withDefaultKeychainAuth)
			if err != nil {
				yield("", err)
				return
			}
			if strings.HasPrefix(manifestHash, hash) {
				if !yield(tag, nil) {
					return
				}
			}
		}
	}
}

// Fetch manifest for image `reference`. Reference must be an image ref with a tag or hash.
func getManifestHash(ctx context.Context, reference string, remoteOpts ...remote.Option) (string, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return "", nil
	}

	opts := append([]remote.Option{remote.WithContext(ctx)}, remoteOpts...)

	desc, err := remote.Get(ref, opts...)
	if err != nil {
		return "", nil
	}

	return desc.Digest.String(), nil
}

// Fetch and list all tags in the `image` repository.
func listTags(ctx context.Context, repo name.Repository, remoteOpts ...remote.Option) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		puller, err := remote.NewPuller(remoteOpts...)
		if err != nil {
			yield("", err)
			return
		}

		lister, err := puller.Lister(ctx, repo)
		if err != nil {
			yield("", fmt.Errorf("reading tags for %s: %w", repo, err))
			return
		}

		for lister.HasNext() {
			tags, err := lister.Next(ctx)
			if err != nil {
				yield("", err)
				return
			}
			for _, tag := range tags.Tags {
				if !yield(tag, nil) {
					return
				}
			}
		}
	}
}
