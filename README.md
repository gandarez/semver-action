# Semantic Versioning Action

This action calculates the next version relying on semantic versioning. It uses the branching name and branching model to calculate the next version.

## Strategies

If `auto` bump, it will try to extract the closest tag and calculate the next semantic version. If not, it will bump respecting the value passed.

### Branching Models

It supports the following branching models:

- [Gitflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)
- [Trunk Based Development](https://trunkbaseddevelopment.com/)

### Branch Names

These are the the default prefixes when `auto` bump:

- `(?i)^(.+:)?(bugfix/.+)` or `(?i)^(.+:)?(hotfix/.+)` - `patch`
- `(?i)^(.+:)?(feature/.+)` - `minor`
- `(?i)^(.+:)?(release/.+)` - `major`
- `(?i)^(.+:)?((doc(s)?|misc)/.+)` - `build`

### Scenarios

#### Gitflow

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

- Source branch is prefixed with `doc/` or `misc/` and dest branch is `develop` - Increments build version.

    ```text
    v0.1.0 results in v0.1.0-pre.1
    v1.5.3-pre.2 results in v1.5.3-pre.3
    ```

- Source branch is prefixed with `hotfix/` and dest branch is `master` - Increments patch version.

    ```text
    v0.1.0 results in v0.1.1-pre.1
    ```

- Source branch is `develop` and dest branch is `master` - Takes the closest tag and finalize it.

    ```text
    v1.5.3-pre.2 results in v1.5.3
    ```

#### Trunk-based

- Not a valid source branch prefix - Increments build version.

    ```text
    v0.1.0 results in v0.1.0+1
    ```

- Source branch is prefixed with `bugfix/` and dest branch is `master` - Increments patch version.

    ```text
    v0.1.0 results in v0.1.1
    ```

- Source branch is prefixed with `feature/` and dest branch is `master` - Increments minor version.

    ```text
    v0.1.0 results in v0.2.0
    ```

- Source branch is prefixed with `release/` and dest branch is `master` - Increments major version.

    ```text
    v0.1.0 results in v1.0.0
    ```

- Source branch is prefixed with `doc/` or `misc/` and dest branch is `master` - Increments build version.

    ```text
    v0.1.0 results in v0.1.0+1
    ```

## Github Environment Variables

Here are the environment variables it takes from Github Actions so far:

- `GITHUB_SHA`

## Example usage

### Gitflow basic

Uses `auto` bump strategy to calculate the next semantic version.

```yaml
- id: semver-tag
  uses: gandarez/semver-action@master
  with:
    branching_model: "git-flow"
- name: "Created tag"
  run: echo "tag ${{ steps.semver-tag.outputs.semver_tag }}"
```

### Gitflow custom

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
    hotfix_regex: "^hotfix/.+"
    exclude_regex: "^ignore/.+"
- name: "Created tag"
  run: echo "tag ${{ steps.semver-tag.outputs.semver_tag }}"
```

### Trunk-based basic

Uses `auto` bump strategy to calculate the next semantic version.

```yaml
- id: semver-tag
  uses: gandarez/semver-action@master
  with:
    branching_model: "trunk-based"
- name: "Created tag"
  run: echo "tag ${{ steps.semver-tag.outputs.semver_tag }}"
```

### Trunk-based custom

```yaml
- id: semver-tag
  uses: gandarez/semver-action@master
  with:
    branching_model: "trunk-based"
    prefix: "ver"
    main_branch_name: "main"
    patch_regex: "^fix/.+"
    minor_regex: "^feat/.+"
    major_regex: "^major/.+"
    build_regex: "^build/.+"
    exclude_regex: "^ignore/.+"
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
| patch_regex | false | Patch regex to match branch name for patch increment. | (?i)^(.+:)?(bugfix/.+) |
| minor_regex | false | Minor regex to match branch name for minor increment. | (?i)^(.+:)?(feature/.+) |
| major_regex | false | Major regex to match branch name for major increment. | (?i)^(.+:)?(release/.+) |
| build_regex | false | Build regex to match branch name for build increment. | (?i)^(.+:)?((doc(s)?|misc)/.+) |
| hotfix_regex | false | Hotfix regex to match branch name for patch increment. | (?i)^(.+:)?(hotfix/.+) |
| exclude_regex | false | Regex to exclude branches from semantic versioning. | |
| repo_dir | false | The repository path. | current dir |
| debug | false | Enable debug mode. | false |

## Outputs

| parameter     | description |
| ---           | --- |
| semver_tag    | The calculdated semantic version. |
| is_prerelease | True if calculated tag is pre-release. For trunk-based model it is always `false`. |
| previous_tag  | The tag used to calculate next semantic version. |
| ancestor_tag  | The ancestor tag based on specific pattern. For trunk-based model it is always empty .|
