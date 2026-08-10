# Release automation

The provider repository owns the coordinated SDK and provider release. The
release workflow is intentionally the only workflow that can publish both
components.

## Dependency invariant

The provider `develop` branch must always use an immutable SDK dependency: a
released tag or a Go pseudo-version for one exact commit. It must never use the
moving SDK `develop` branch.

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
choose **Run workflow**, enter the version from the first `CHANGELOG.md` entry,
and keep **release_sdk** set to `auto`. Approve the `provider-release`
environment deployment when the preflight summary is correct. Keep
**dry_run** enabled for the first validation run; disable it only for a real
release.

The workflow tests both repositories, optionally releases accumulated SDK
changes, pins the provider dependency, publishes the signed provider release,
waits for the public Terraform Registry, and updates the private registry.

## Recovery

The **Building the provider for the public registry** workflow is manual and is
only for republishing an existing provider tag. The private registry workflow
is also retained as a manual recovery operation. Neither workflow is part of a
normal release.

## Required repository settings

- Install the release GitHub App on both repositories with `Contents: read and
  write` and `Metadata: read-only`.
- Add the App to the `develop` branch ruleset bypass list with `Always allow`.
- Protect `v*` tags and allow only the App and emergency administrators to
  create or update them.
- Require one pull request approval and the `required` status check on
  `develop`.
- Do not require pull request branches to be up to date; the push check on
  `develop` provides the integration safety net without repeated rebases.
- Create the `provider-release` environment, add required reviewers, prevent
  self-review when possible, and restrict deployment branches to `develop`.
- Store `GPG_PRIVATE_KEY` and `s3_config` as environment secrets of
  `provider-release`. `APP_ID` and `APP_PEM` currently remain repository
  secrets because the trusted `push: develop` changelog workflow also needs
  them; fork pull request workflows cannot access repository secrets.
- Create an `internal-docs` environment containing its own `s3_config` secret
  for the documentation publication workflow.
- Restrict Actions to selected actions and require actions to be pinned to a
  full commit SHA. All actions in these workflows are SHA-pinned.
