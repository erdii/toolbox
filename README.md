# erdii's toolbox

This is an assortment of niche use-case tools that I wrote for myself.

If you want to use one of the tools for yourself you can install a tool directly from the main branch like this: `go install github.com/erdii/toolbox/cmd/<toolname>@main`.

## tools

### imagehash2tags

[`imagehash2tags`](./cmd/imagehash2tags) resolves a docker image reference with a hash to all the tags in the repository pointing to it. Beware: Due to how the registry/v2 api work, this will be very slow on repos with a lot of images.

As an example, running `imagehash2tags quay.io/app-sre/package-operator-manager@sha256:bc252a6bbb96bd98b505b3f7d66e934580ab8a0f0cbb3b93a1f8af9f52452a03` will output:
```
Image: `quay.io/app-sre/package-operator-manager@sha256:bc252a6bbb96bd98b505b3f7d66e934580ab8a0f0cbb3b93a1f8af9f52452a03`.
Looking up tags...
Found these tags matching hash `sha256:bc252a6bbb96bd98b505b3f7d66e934580ab8a0f0cbb3b93a1f8af9f52452a03`:
978fabf
imagehash2tags quay.io/app-sre/package-operator-manager@sha256:bc252a6bbb96bd98b505b3f7d66e934580ab8a0f0cbb3b93a1f8af9f52452a03
```

### keep

[`keep`](./cmd/keep) repeatedly runs and restarts a long-running command.
An example that I use very often: `keep kubectl logs -f -n package-operator-system deployments/package-operator-manager`

### extract-image

[`extract-image`](./cmd/extract-image) extracts the combined filesystem of all layers in a docker image tarball exported by `docker save` or `podman save`.

### workday-excel-to-ical

I don't want to talk about it 😅
