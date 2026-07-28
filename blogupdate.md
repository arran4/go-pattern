# Updates to CI Workflow for Go Pattern

This document summarizes the changes applied to the CI workflow based on the blog post [Simplified GitHub CI](https://arran4.github.io/blog/post/2026/011-simplified-github-ci/).

## Fix: "git tag was not made against commit" for manual releases
When triggering a manual release via `workflow_dispatch` on the main branch, `goreleaser` previously failed because the current HEAD commit lacked a git tag matching the latest tag in the repository. Goreleaser expects the tag being built to exactly match the checked-out commit.
*   **Resolution:** Added a `Calculate and Create Tag` step to the `goreleaser` job. This step automatically bumps the semantic version based on `inputs.mode` (e.g., `release-major`), or uses `inputs.release_version_override`. It then creates and pushes the tag to the repository *before* `goreleaser` runs. We then set `GORELEASER_CURRENT_TAG` to this new tag to guarantee `goreleaser` targets the correct, newly-created tag for the current commit.
*   **Snapshot Mode:** For test and rc releases (`release-test`, `release-rc`, `release-alpha`), we still pass the `--snapshot` flag to Goreleaser, as these do not require a permanent tag.

## Fix: Direct GitHub Releases Publishing
Previously, the workflow separated artifact building and draft publishing into two jobs (`goreleaser` and `publish-draft`). However, the `goreleaser` job was not actually uploading artifacts to GitHub Actions storage, causing downstream jobs to fail or act improperly.
*   **Resolution:** We removed the `publish-draft` and `promote-release` jobs entirely to eliminate GitHub Action artifact storage quotas and costs.
*   **Native Publishing:** `goreleaser` inherently supports uploading directly to GitHub Releases. By simply providing it with `GITHUB_TOKEN`, it will automatically create the release and upload the binaries without intermediate artifact storage.

## Ensuring Tag Context
We ensured that the `actions/checkout@v4` step in the `goreleaser` job properly uses `fetch-depth: 0` alongside `fetch-tags: true`. Fetching tags is critical for Goreleaser to establish correct changelogs and validate that the build commit matches the tag accurately.
