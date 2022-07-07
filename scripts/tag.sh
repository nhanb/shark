#!/usr/bin/env sh
set -euf

# Usage: ./scripts/tag.sh v0.0.0
# which will open an annotated tag editor prefilled with a shortlog comparing
# against the previous tag.

PREVIOUS_TAG=$(git describe --abbrev=0 HEAD^1)
SHORTLOG=$(git shortlog "$PREVIOUS_TAG..HEAD")

printf "CHANGEME\n\n%s" "$SHORTLOG" | git tag "$1" --file -
git tag "$1" -f -a
