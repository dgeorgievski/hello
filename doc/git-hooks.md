# Git local hooks

Document the use of git hooks written in NuShell scripting language, which offers capabilities unparalled in Bash and other Unix Shells.

All of the git hook scripts used in this project are stored under `scripts/githooks` directory.

## Prepare commit message

This hook, [scripts/githooks/prepare-commit-msg](../scripts/githooks/prepare-commit-msg), verifies that the subject of the commit message complies with the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard.

The commit will be rejected if the subject of the message fails the test.

The script is utilizing the [is_conventional_commit_hook](../scripts/agit/core.nu) NuShell command.


