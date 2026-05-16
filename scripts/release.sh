#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/release.sh vX.Y.Z [--push] [--github-release]

Runs the release preflight, updates the documented Docker image version,
commits the bump, creates an annotated tag, and optionally pushes/creates
the GitHub Release.

Options:
  --push             Push main and the release tag to origin.
  --github-release   Create a GitHub Release with generated notes after push.
USAGE
}

version="${1:-}"
if [[ -z "$version" || "$version" == "-h" || "$version" == "--help" ]]; then
  usage
  exit 0
fi
shift || true

push_release=0
github_release=0
for arg in "$@"; do
  case "$arg" in
    --push) push_release=1 ;;
    --github-release) github_release=1 ;;
    *)
      echo "Unknown option: $arg" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Version must match vX.Y.Z, got: $version" >&2
  exit 2
fi

if [[ "$github_release" -eq 1 && "$push_release" -ne 1 ]]; then
  echo "--github-release requires --push so the tag exists on GitHub first." >&2
  exit 2
fi

branch="$(git branch --show-current)"
if [[ "$branch" != "main" ]]; then
  echo "Release must run from main; current branch is $branch." >&2
  exit 1
fi

if [[ -n "$(git status --short)" ]]; then
  echo "Working tree must be clean before release." >&2
  git status --short >&2
  exit 1
fi

if git rev-parse -q --verify "refs/tags/$version" >/dev/null; then
  echo "Local tag already exists: $version" >&2
  exit 1
fi

if git ls-remote --exit-code --tags origin "refs/tags/$version" >/dev/null 2>&1; then
  echo "Remote tag already exists: $version" >&2
  exit 1
fi

plain_version="${version#v}"

echo "==> Updating documented Docker image tag to $plain_version"
perl -0pi -e 's#ghcr\.io/nczz/phishguard:[0-9]+\.[0-9]+\.[0-9]+#ghcr.io/nczz/phishguard:'"$plain_version"'#g' README.md

if [[ -z "$(git diff -- README.md)" ]]; then
  echo "README.md did not change; check whether the version is already current." >&2
  exit 1
fi

echo "==> Running backend tests"
(cd backend && go test ./...)

echo "==> Running frontend lint"
(cd frontend && npm run lint)

echo "==> Running frontend build"
(cd frontend && npm run build)

echo "==> Validating production compose config"
env \
  DB_ROOT_PASS=root \
  DB_PASS=db \
  JWT_SECRET=jwt \
  ENCRYPT_KEY=00000000000000000000000000000000 \
  ADMIN_PASSWORD=admin \
  TRACKER_BASE_URL=http://tracker.local \
  docker compose -f docker-compose.prod.yml config --quiet

git add README.md
git commit -m "Bump release version to $version"
git tag -a "$version" -m "Release $version"

if [[ "$push_release" -eq 1 ]]; then
  echo "==> Pushing main and $version"
  git push origin main
  git push origin "$version"
fi

if [[ "$github_release" -eq 1 ]]; then
  repo="${GH_REPO:-$(git remote get-url origin | sed -E 's#.*github.com[:/]([^/]+/[^/.]+)(\.git)?#\1#')}"
  echo "==> Creating GitHub Release for $version in $repo"
  gh api "repos/$repo/releases" \
    -X POST \
    -f tag_name="$version" \
    -f name="$version" \
    -F generate_release_notes=true >/dev/null
fi

echo "Release prepared: $version"
