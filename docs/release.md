# Release Process

PhishGuard releases are tag-driven. Pushing a `vX.Y.Z` tag triggers
`.github/workflows/docker-publish.yml`, which builds and publishes the
multi-architecture Docker image to GHCR with these tags:

- `ghcr.io/nczz/phishguard:X.Y.Z`
- `ghcr.io/nczz/phishguard:X.Y`
- `ghcr.io/nczz/phishguard:X`
- `ghcr.io/nczz/phishguard:latest`

## Preconditions

- Work from `main`.
- Working tree is clean.
- The target tag does not exist locally or on `origin`.
- GitHub CLI is authenticated if creating a GitHub Release.
- Docker is available for compose config validation.

## Standard Release

```bash
scripts/release.sh vX.Y.Z --push --github-release
```

The script performs the required local gate before it pushes anything:

1. Updates the README Docker image example to `X.Y.Z`.
2. Runs `go test ./...` in `backend`.
3. Runs `npm run lint` in `frontend`.
4. Runs `npm run build` in `frontend`.
5. Validates `docker-compose.prod.yml` with required environment values.
6. Commits the version bump.
7. Creates an annotated `vX.Y.Z` tag.
8. Pushes `main` and the tag.
9. Creates a GitHub Release with generated notes.

## Post-Release Checks

After pushing the tag, confirm the Docker workflow:

```bash
gh run list --workflow docker-publish.yml --limit 3
gh run watch <run-id>
```

Then confirm the package tags in GHCR before announcing the release.
