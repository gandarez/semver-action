# Semantic Versioning Action

This action calculates the next version relying on semantic versioning.

## Strategies

If `auto` bump, it will try to extract the closest tag and calculate the next semantic version. If not, it will bump respecting the value passed.

### Branching Models

It supports the following branching models:

- [Gitflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)
- [Trunk Based Development](https://trunkbaseddevelopment.com/)

### Branch Names

These are the the default prefixes when `auto` bump:

- `^bugfix/.+` - `patch`
- `^feature/.+` - `minor`
- `^release/.+` - `major`
- `^(doc(s)?|misc)/.+` - `build`

### Scenarios

#### Gitglow and auto bump

- Not a valid source branch prefix - Increments build version.

    ```text
        v0.1.0 results in v0.1.0-pre.1
        v1.5.3-pre.2 results in v1.5.3-pre.3
    ```

- Source branch is prefixed with `bugfix/` and dest branch is `develop` - Increments patch version.

    ```text
    v0.1.0 results in v0.1.1-pre.1
    v1.5.3-pre.2 results in v1.5.4-pre.1
    ```

- Source branch is prefixed with `feature/` and dest branch is `develop` - Increments minor version.

    ```text
    v0.1.0 results in v0.2.0-pre.1
    v1.5.3-pre.2 results in v1.6.0-pre.1
    ```

- Source branch is prefixed with `release/` and dest branch is `develop` - Increments major version.

    ```text
    v0.1.0 results in v1.0.0-pre.1
    v1.5.3-pre.2 results in v2.0.0-pre.1
    ```

- Source branch is `develop` and dest branch is `master` - Takes the closest tag and finalize it.

    ```text
    v1.5.3-pre.2 results in v1.5.3
    ```

## Github Environment Variables

Here are the environment variables it takes from Github Actions so far:

- `GITHUB_SHA`

## Example usage

### Gitflow

Uses `auto` bump strategy to calculate the next semantic version.

```yaml
- id: semver-tag
  uses: gandarez/semver-action@master
  with:
    branching_model: "git-flow"
- name: "Created tag"
  run: echo "tag ${{ steps.semver-tag.outputs.semver_tag }}"
```

### Custom

```yaml
- id: semver-tag
  uses: gandarez/semver-action@master
  with:
    branching_model: "git-flow"
    prefix: "ver"
    prerelease_id: "alpha"
    main_branch_name: "main"
    develop_branch_name: "dev"
    patch_regex: "^fix/.+"
    minor_regex: "^feat/.+"
    major_regex: "^major/.+"
    build_regex: "^build/.+"
- name: "Created tag"
  run: echo "tag ${{ steps.semver-tag.outputs.semver_tag }}"
```

## Inputs

| parameter | required | description | default |
| --- | --- | --- | --- |
| bump | false | Bump strategy for semantic versioning. Can be `auto`, `major`, `minor`, `patch`. | auto |
| base_version | false | Version to use as base for the generation, skips version bumps. | |
| prefix | false | Prefix used to prepend the final version.| v |
| prerelease_id | false | Text representing the prerelease identifier. | pre |
| main_branch_name | false | The main branch name. | master |
| develop_branch_name | false | The develop branch name. | develop |
| patch_regex | false | Regex to match patch branches. | ^bugfix/.+ |
| minor_regex | false | Regex to match minor branches. | ^feature/.+ |
| major_regex | false | Regex to match major branches. | ^release/.+ |
| build_regex | false | Regex to match build branches. | ^(doc(s)?|misc)/.+ |
| repo_dir | false | The repository path. | current dir |
| debug | false | Enables debug mode. | false |

## Outpus

| parameter     | description                                      |
| ---           | ---                                              |
| semver_tag    | The calculdated semantic version.                |
| is_prerelease | True if calculated tag is pre-release.           |
| previous_tag  | The tag used to calculate next semantic version. |
| ancestor_tag  | The ancestor tag based on specific pattern.      |
