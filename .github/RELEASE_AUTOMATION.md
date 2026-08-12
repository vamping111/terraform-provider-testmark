# Release automation

The test provider repository owns both coordinated releases and the atomic
landing of provider changes that depend on an unreleased SDK change.

## Dependency invariant

The provider `develop` branch must always use an immutable SDK dependency: a
released tag or a Go pseudo-version for one exact commit. It must never use the
moving SDK `develop` branch.

The SDK keeps its original `module github.com/aws/aws-sdk-go` declaration so
existing imports remain valid. The right-hand side of the provider `replace`
directive points to the Rockit Cloud fork at an immutable tag or Go
pseudo-version.

The sandbox SDK repository currently has no `develop` branch, so the test
orchestrator reads its default `codex/release-automation` branch and records an
exact SHA before testing. The SDK required-check workflow temporarily includes
that branch as a push and pull-request target. Production configuration must
point these references back to `develop`.

SDK verification deliberately runs lint, vet, compilation, and unit tests without
regenerating every upstream AWS client. Generated SDK files are expected to be
committed by the feature pull request and are still checked for an unchanged
worktree.

## SDK-linked provider pull requests

Developers do not edit or commit provider `go.mod` or `go.sum` for an SDK-linked
change. The supported flow is:

1. Open the SDK and provider pull requests from forks. Apply `NeedSDK` to the
   provider pull request. Use equal branch names, or add
   `aws-sdk-go=<SDK branch>` to the provider pull request description, for its
   early fork-to-fork CI. Also add `aws-sdk-go-pr=#<SDK PR number>`; the landing
   workflow requires this exact audited link and rejects a mismatched input.
2. Review both pull requests and let their required checks pass.
   Because the landing App bypasses GitHub's last-push rule, the provider PR
   needs two distinct approvals on its current head before a real landing.
3. Merge the SDK pull request into the protected central SDK release branch,
   without creating an SDK release.
4. Leave `NeedSDK` on the provider pull request and leave its `go.mod` and
   `go.sum` untouched. The deliberately failing `sdk-linked-policy/head` check
   prevents an accidental normal GitHub merge.
5. In the provider repository, run **Land SDK-linked provider PR** with the two
   pull request numbers. First use `dry_run=true`. For the real run use
   `dry_run=false`, then approve the `provider-pr-land` environment.
6. The workflow generates one immutable merge candidate against the exact
   central SDK commit, tests that same candidate commit, and fast-forwards
   `develop` to that exact tested commit. GitHub sees the original provider pull request head as the
   second parent and marks that pull request merged.
7. Run the coordinated provider release later. It replaces an accumulated SDK
   pseudo-version with the new SDK release tag.

If several SDK pull requests are merged before landing a provider pull request,
the workflow tests and pins the exact current SDK release-branch head, so all
accumulated changes are covered. If provider `develop`, the provider pull request
head, or the SDK release branch changes between testing and writing, the run
fails without changing any repository. Rerun it; developers do not need to
rebase only because of a dependency bump.

The automatic path rejects provider pull requests that change dependency,
workflow, build, release, tool, script, symlink, submodule, or binary files.
Those uncommon changes use the normal reviewed merge path. This restriction is
intentional: the landing App bypasses the branch ruleset, so its workflow is
limited to ordinary provider source, documentation, examples, and changelog
fragments.

This prevents an SDK-only merge from making provider `develop` unstable. The
release workflow also tests the current provider against the SDK release
candidate before it creates an SDK tag or GitHub release.

### Landing security boundary

The landing workflow has four separate trust stages:

1. `test` authorizes immutable PR, branch, and check revisions. It uses no App
   token, environment, OIDC permission, or write permission and saves no Go
   cache.
2. `prepare` starts on a clean runner, creates the merge, runs only Go module
   resolution, and hands off a one-day thin Git bundle plus a SHA-256-bound
   manifest. It does not execute repository code.
3. `verify-candidate` has no token or permissions. It downloads, materializes,
   and runs the full provider tests on exactly that immutable Git tree, with
   fresh Go caches.
4. `land` starts only after the environment approval. It never checks out a
   worktree or executes provider code. It rechecks the approval, resolved review
   threads, two distinct current-head reviewers, exact required workflow/check identity, PR and branch SHAs, SDK
   ancestry, Git parents/tree, generated module delta, and artifact digest. Only
   then does it mint a provider-only App token and push the exact tested commit
   without force. A
   concurrent `develop` update makes that push fail rather than overwrite work.

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
- Also require the `sdk-linked-policy/head` status check on `develop`. It succeeds
  for ordinary pull requests and intentionally fails for `NeedSDK` pull
  requests, which must use the landing workflow and App bypass. The decision is
  sticky: after `NeedSDK` has appeared in a PR timeline, removing the label does
  not reopen the ordinary merge path.
- Do not require pull request branches to be up to date; the push check on
  `develop` provides the integration safety net without repeated rebases.
- Create the `provider-release` environment, add required reviewers, prevent
  self-review when another release operator is available, and restrict
  deployment branches to `develop`.
- Store a test-only `GPG_PRIVATE_KEY` as an environment secret of
  `provider-release`. `APP_ID` and `APP_PEM` currently remain repository
  secrets because the trusted `push: develop` changelog workflow also needs
  them; fork pull request workflows cannot access repository secrets.
- Create a separate `provider-pr-land` environment, restrict deployments to
  `develop`, and add a required reviewer. Store `LAND_APP_ID` and
  `LAND_APP_PEM` as environment secrets. In production these belong to a
  dedicated landing App installed only on the provider repository, with
  `Contents: read and write` and `Metadata: read-only`. It may bypass only the
  provider `develop` ruleset; it must not bypass `v*` tag rules and must not be
  installed on the SDK repository. The existing release App remains separate.
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
   protect `develop`, require one approval plus the `required` and
   `sdk-linked-policy/head` checks, disable force pushes/deletions, and add the
   release App to the bypass list with **Always allow**. Do not require
   pull-request branches to be up to date.
   Bootstrap this in order: first merge the pull request that adds
   `sdk-linked-policy.yml`; then reopen, relabel, or update a disposable pull
   request so the new check runs at least once; only then add
   `sdk-linked-policy/head` to the ruleset. Adding it before the workflow exists on
   `develop` would block the bootstrap pull request itself.
6. In SDK [rulesets](https://github.com/vamping111/testmark-sdk-go/settings/rules),
   apply the equivalent rule to the sandbox release branch
   `codex/release-automation`; protect `v*` tags in both repositories and allow
   tag creation only through the App bypass (plus emergency administrators).
7. Leave `ENABLE_TEST_REGISTRY_WRITE` and `ENABLE_TEST_DOCS_PUBLISH` absent in
   provider [Actions variables](https://github.com/vamping111/terraform-provider-testmark/settings/variables/actions).
   Do not configure any production S3 credential in the sandbox.
8. Create a separate landing GitHub App even for the sandbox. Install it only
   on `terraform-provider-testmark`, grant `Contents: read and write` and
   `Metadata: read-only`, and put it only in the provider `develop` bypass list.
   Do not install it on the SDK and do not put it in either tag ruleset.
9. In provider [environments](https://github.com/vamping111/terraform-provider-testmark/settings/environments),
   create `provider-pr-land`, restrict deployments to `develop`, and add a
   required reviewer. Store the dedicated landing App ID and complete PEM as
   environment secrets `LAND_APP_ID` and `LAND_APP_PEM`.

To test landing, create disposable SDK and provider pull requests, approve both,
apply `NeedSDK` to the provider pull request, and merge only the SDK pull request.
Add `aws-sdk-go-pr=#<SDK PR number>` to the provider PR description.
Run **Land SDK-linked provider PR** from `develop` with both pull request numbers
and `dry_run=true`. Confirm that no refs changed. Repeat with `dry_run=false` and
approve `provider-pr-land`; the provider pull request should be represented by
one new merge commit on `develop`, including the generated immutable SDK pin.

After the provider workflow is merged into `develop`, first run **Release
provider** with `release_sdk=always`,
`sdk_version=v1.44.10-TESTMARK1`, `dry_run=true`, and
`publish_registries=false`. Leave `provider_version` empty. A successful dry
run creates no commits, tags, releases, or registry writes. Repeat with
`dry_run=false` and approve `provider-release` to perform the sandbox release.
