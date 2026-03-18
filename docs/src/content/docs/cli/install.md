---
title: CLI — install
description: Install a workflow pack from another repository.
---

```bash
aynig install <repo> [ref] [subfolder]
```

`aynig install` clones another repository into a temporary directory and copies
its workflow files into the current repository.

By default it copies the source repository's `.dwp/` directory into your local
`.dwp/` directory.

## Behavior

- Requires a source repository argument.
- Accepts a GitHub shorthand such as `owner/name` and expands it to
  `https://github.com/owner/name.git`.
- Clones with `--depth 1`.
- Uses the provided ref as the clone branch when `[ref]` is given.
- Copies from `.dwp/` by default, or from `[subfolder]` when provided.
- Refuses to run when `.dwp/` has uncommitted changes in the current repository.
- Skips source `README.md` files while copying.
- Reports installed files and overwritten files.
- When `.dwp/COMMANDS.md` changes, attempts an automatic merge with `opencode`
  or `claude` when either tool is available.

## Arguments

### `<repo>`

Required. The source repository to clone.

Examples:

```bash
aynig install hacknlove/all-you-need-is-git
```

```bash
aynig install https://github.com/hacknlove/all-you-need-is-git.git
```

### `[ref]`

Optional. The branch, tag, or ref to clone.

Example:

```bash
aynig install hacknlove/all-you-need-is-git main
```

### `[subfolder]`

Optional. Copies files from a specific subfolder in the source repository
instead of the default `.dwp/` directory.

Example:

```bash
aynig install hacknlove/all-you-need-is-git main workflow-packs/reviewer
```

## Options

This command has no flags.

## Examples

Install the default `.dwp/` directory from a GitHub repository:

```bash
aynig install hacknlove/all-you-need-is-git
```

Install from a specific branch:

```bash
aynig install hacknlove/all-you-need-is-git ops-workflow-pack
```

Install from a specific subfolder:

```bash
aynig install hacknlove/all-you-need-is-git main workflow-packs/reviewer
```

Install from a full git URL:

```bash
aynig install https://github.com/hacknlove/all-you-need-is-git.git
```
