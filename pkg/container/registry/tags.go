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
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"context"
)

// FindTagsForImageHash queries all tags in a repo and compares their manifest hashes against the given hash. Uses standard docker authentication (ie ~/.docker/config.json).
func FindTagsForImageHash(ctx context.Context, image, hash string) ([]string, error) {
	matchingTags := []string{}

	withDefaultKeychainAuth := remote.WithAuthFromKeychain(authn.DefaultKeychain)

	repo, err := name.NewRepository(image)
	if err != nil {
		return nil, err
	}

	tags, err := listTags(ctx, repo, withDefaultKeychainAuth)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		manifestHash, err := getManifestHash(ctx, fmt.Sprintf("%s:%s", image, tag), withDefaultKeychainAuth)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(manifestHash, hash) {
			matchingTags = append(matchingTags, tag)
		}
	}

	return matchingTags, nil
}

// Fetch manifest for image `reference`. Reference must be an image ref with a tag or hash.
func getManifestHash(ctx context.Context, reference string, remoteOpts ...remote.Option) (string, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return "", nil
	}

	desc, err := remote.Get(ref, remoteOpts...)
	if err != nil {
		return "", nil
	}

	return desc.Digest.String(), nil
}

// Fetch and list all tags in the `image` repository.
func listTags(ctx context.Context, repo name.Repository, remoteOpts ...remote.Option) ([]string, error) {
	tagNames := []string{}

	puller, err := remote.NewPuller(remoteOpts...)
	if err != nil {
		return nil, err
	}

	lister, err := puller.Lister(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("reading tags for %s: %w", repo, err)
	}

	for lister.HasNext() {
		tags, err := lister.Next(ctx)
		if err != nil {
			return nil, err
		}
		tagNames = append(tagNames, tags.Tags...)
	}

	return tagNames, nil
}
