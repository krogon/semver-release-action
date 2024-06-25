
# Semver Release Github Action ![badge](https://github.com/krogon/semver-release-action/workflows/CI/badge.svg)

## Fork information

Based on https://github.com/K-Phoen/semver-release-action

Additions:
* default_increment - default increment bump when label is not found, default `"skip"`
* tag_prefix - allows to manage separate tags in monorepo, default `""`

## Description

Automatically create [SemVer](https://semver.org/) compliant releases based on
PR labels.

Assuming that a PR is tagged with a "*semver-compliant*" label (*patch*, *minor* or *major*),
then this action can create a tag and a GitHub release when it is merged.

**Note:** to determine the base tag for the increment, this action will try to
find the most recent tag complying to [SemVer](https://semver.org/). No
additional setup is required.

## Inputs

### `release_branch`

**Required** Branch to tag. Default `"master"`.

### `release_strategy`

**Required** Release strategy. Default `"release"` (`release`: creates a GitHub
release ; `tag`: creates a lightweight tag ; `none`: computes the next
[SemVer](https://semver.org/) version but does not create a release or tag).

### `default_increment`

**Required** Default increment (skip/patch/minor/major). Default `"skip"`.

### `tag_format`

**Optional** Format used to create tags. Default `"v%major%.%minor%.%patch%"`.

### `tag_prefixt`

**Optional** Tag prefix for monorepo. Default `""`.

### `tag`

**Optional** Tag to use. If left undefined, it will be computed using the tags
already present in the repository.

## Outputs

### `tag`

The newly created tag.

## Example usage

```yaml
# .github/workflows/release.yml
name: Release

on:
  pull_request:
    types: [closed]

jobs:
  build:
    runs-on: ubuntu-latest

    if: github.event.pull_request.merged
    
    steps:
      - name: Tag
        uses: krogon/semver-release-action@master
        with:
          release_branch: master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

```

## License

This library is under the [MIT](LICENSE.md) license.
