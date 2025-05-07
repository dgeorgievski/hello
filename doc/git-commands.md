# Git Commands

## Git Squash

It squashes the last N commits into one. Commits messages are combined into one.

```nushell
# display pretty log history
> g ll
49c5b03 (HEAD -> feat/squash-cmd) docs(githook): Reference in main README to docs/git-hook.md added
c8506e4 test(git): draft git nushell commands
e4f0f9a (origin/main, origin/HEAD, main) docs(githook): Document the use of prepare-commit-hook
c6f7631 (tag: v0.0.4) ci(version): Bumped version to 0.0.4
...

# import agit module
> use scripts/agit

# squash last 2 commits
# This will open the configured editor to review commit messages.
> agit git-squash-last 2
Commit message is compliant with Conventional Commits
[feat/squash-cmd 60df5c8] docs(githook): Reference in main README to docs/git-hook.md added
 8 files changed, 369 insertions(+), 3 deletions(-)
 create mode 100644 scripts/agit/common.nu
 create mode 100644 scripts/agit/complete.nu
 create mode 100644 scripts/agit/git-flow.nu
 create mode 100644 scripts/agit/stat.nu
```
