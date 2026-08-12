# Release automation

The test provider repository owns the coordinated SDK and provider release. The
release workflow is intentionally the only workflow that can publish both
components.

## Dependency invariant

The provider `develop` branch must always use an immutable SDK dependency: a
released tag or a Go pseudo-version for one exact commit. It must never use the
moving SDK `develop` branch.

This test branch currently pins
`github.com/vamping111/testmark-sdk-go@d43f5c5d7036` through its immutable Go
pseudo-version. The SDK keeps its original `module github.com/aws/aws-sdk-go`
declaration so existing imports remain valid.

The sandbox SDK repository currently has no `develop` branch, so the test
orchestrator reads its default `codex/release-automation` branch and records an
exact SHA before testing. The SDK required-check workflow temporarily includes
that branch as a push and pull-request target. Production configuration must
point these references back to `develop`.

SDK verification deliberately runs lint, vet, compilation, and unit tests without
regenerating every upstream AWS client. Generated SDK files are expected to be
committed by the feature pull request and are still checked for an unchanged
worktree.

For a feature that changes both repositories:

1. Test the provider pull request against the SDK pull request with `NeedSDK`.
2. Merge the SDK pull request without publishing an SDK release.
3. Replace the provider dependency with a pseudo-version for the merged SDK
   commit and commit `go.mod` and `go.sum` to the provider pull request.
4. Merge the provider pull request only after its required check passes.
5. Run the coordinated provider release. It replaces the pseudo-version with
   the new SDK release tag.

This prevents an SDK-only merge from making provider `develop` unstable. The
release workflow also tests the current provider against the SDK release
candidate before it creates an SDK tag or GitHub release.

## Normal release

From the provider repository, open **Actions**, select **Release provider**,
choose **Run workflow**, leave the provider version empty so it is read from the
first `CHANGELOG.md` entry, and keep **release_sdk** set to `auto`. For the first SDK release, use
`release_sdk=always` and explicitly provide an SDK tag such as
`v1.44.10-TESTMARK1`. Approve the `provider-release`
environment deployment when the preflight summary is correct. Keep
**dry_run** enabled for the first validation run; disable it only for a real
test release. Keep **publish_registries** disabled: test releases are created
only in `vamping111/testmark-sdk-go` and
`vamping111/terraform-provider-testmark`.

The workflow tests both repositories, optionally releases accumulated SDK
changes, pins the provider dependency, publishes the signed provider release,
creates the GitHub release, and adds the next minor `Unreleased` changelog
section in a follow-up commit delivered by the same atomic push. The release tag
still points to the clean release commit. Registry waiting and the
private-registry update run only when **publish_registries** is explicitly
enabled.

## Recovery

The **Building the provider for the public registry** workflow is manual and is
only for republishing an existing provider tag. The private registry workflow
is also retained as a manual recovery operation. Neither workflow is part of a
normal release.

## Required repository settings

- Install the release GitHub App on `vamping111/terraform-provider-testmark`
  and `vamping111/testmark-sdk-go` with `Contents: read and
  write` and `Metadata: read-only`.
- Add the App to the `develop` branch ruleset bypass list with `Always allow`.
- Protect `v*` tags and allow only the App and emergency administrators to
  create or update them.
- Require one pull request approval and the `required` status check on
  `develop`.
- Do not require pull request branches to be up to date; the push check on
  `develop` provides the integration safety net without repeated rebases.
- Create the `provider-release` environment, add required reviewers, prevent
  self-review when another release operator is available, and restrict
  deployment branches to `develop`.
- Store a test-only `GPG_PRIVATE_KEY` as an environment secret of
  `provider-release`. `APP_ID` and `APP_PEM` currently remain repository
  secrets because the trusted `push: develop` changelog workflow also needs
  them; fork pull request workflows cannot access repository secrets.
- For the initial sandbox release, do **not** create `s3_config`, do not set
  `ENABLE_TEST_REGISTRY_WRITE`, and do not set `ENABLE_TEST_DOCS_PUBLISH`.
  Consequently neither private-registry nor documentation publication can run.
- A later private-registry test additionally requires dedicated non-production
  values for `TEST_S3_REGISTRY_URL`, `TEST_S3_BUCKET_NAME`, and
  `TEST_RELEASE_MIRROR_URL`, plus an environment `s3_config` secret. The script
  rejects the known production endpoint, bucket, and release mirror.
- A later documentation test requires a dedicated non-production
  `TEST_S3_DOCS_BUCKET_NAME`, an `internal-docs` environment with its own
  `s3_config`, and `ENABLE_TEST_DOCS_PUBLISH=true`. The production docs bucket
  is rejected by the sandbox script.
- Restrict Actions to selected actions and require actions to be pinned to a
  full commit SHA. All actions in these workflows are SHA-pinned.

## Sandbox setup links

1. Open [GitHub App installations](https://github.com/settings/installations),
   select the release App, choose **Only select repositories**, add both
   `terraform-provider-testmark` and `testmark-sdk-go`, and save.
2. In the provider [Actions secrets](https://github.com/vamping111/terraform-provider-testmark/settings/secrets/actions),
   add repository secrets `APP_ID` (the numeric GitHub App ID) and `APP_PEM`
   (the complete private-key PEM). Do not add these secrets to the SDK repo.
3. In the provider [environments](https://github.com/vamping111/terraform-provider-testmark/settings/environments),
   create `provider-release`, restrict deployments to `develop`, add one
   required reviewer, and store a test-only `GPG_PRIVATE_KEY` environment
   secret. For a single-person sandbox, leave self-review prevention disabled;
   production should use a second release operator.
4. In provider [Actions settings](https://github.com/vamping111/terraform-provider-testmark/settings/actions),
   set the default `GITHUB_TOKEN` workflow permission to read-only. The
   workflows mint short-lived, repository-scoped App tokens for writes.
5. In provider [rulesets](https://github.com/vamping111/terraform-provider-testmark/settings/rules),
   protect `develop`, require one approval and the `required` check, disable
   force pushes/deletions, and add the release App to the bypass list with
   **Always allow**. Do not require pull-request branches to be up to date.
6. In SDK [rulesets](https://github.com/vamping111/testmark-sdk-go/settings/rules),
   apply the equivalent rule to the sandbox release branch
   `codex/release-automation`; protect `v*` tags in both repositories and allow
   tag creation only through the App bypass (plus emergency administrators).
7. Leave `ENABLE_TEST_REGISTRY_WRITE` and `ENABLE_TEST_DOCS_PUBLISH` absent in
   provider [Actions variables](https://github.com/vamping111/terraform-provider-testmark/settings/variables/actions).
   Do not configure any production S3 credential in the sandbox.

After the provider workflow is merged into `develop`, first run **Release
provider** with `release_sdk=always`,
`sdk_version=v1.44.10-TESTMARK1`, `dry_run=true`, and
`publish_registries=false`. Leave `provider_version` empty. A successful dry
run creates no commits, tags, releases, or registry writes. Repeat with
`dry_run=false` and approve `provider-release` to perform the sandbox release.
