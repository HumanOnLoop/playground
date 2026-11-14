#!/usr/bin/env bash
# Creates or updates the ruleset that protects main (ruleset/main.json): the "ci" check is required, and main can be
# neither force-pushed nor deleted. Needs a token with administration rights on the repo (gh auth as an org owner).
set -euo pipefail
repo=${1:-HumanOnLoop/playground}
dir=$(cd "$(dirname "$0")" && pwd)
id=$(gh api "repos/$repo/rulesets" --jq '.[] | select(.name == "main") | .id')
if [ -n "$id" ]; then
  gh api -X PUT "repos/$repo/rulesets/$id" --input "$dir/main.json" --jq '"updated ruleset \(.id)"'
else
  gh api -X POST "repos/$repo/rulesets" --input "$dir/main.json" --jq '"created ruleset \(.id)"'
fi
