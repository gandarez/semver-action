package generate_test

import (
	"os"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/gandarez/semver-action/cmd/generate"

	"github.com/alecthomas/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadParams_Prefix(t *testing.T) {
	os.Setenv("INPUT_PREFIX", "ver")
	defer os.Unsetenv("INPUT_PREFIX")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "ver", params.Prefix)
}

func TestLoadParams_Prefix_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "v", params.Prefix)
}

func TestLoadParams_PrereleaseID(t *testing.T) {
	os.Setenv("INPUT_PRERELEASE_ID", "alpha")
	defer os.Unsetenv("INPUT_PRERELEASE_ID")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "alpha", params.PrereleaseID)
}

func TestLoadParams_PrereleaseID_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "pre", params.PrereleaseID)
}

func TestLoadParams_MainBranchName(t *testing.T) {
	os.Setenv("INPUT_MAIN_BRANCH_NAME", "main")
	defer os.Unsetenv("INPUT_MAIN_BRANCH_NAME")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "main", params.MainBranchName)
}

func TestLoadParams_MainBranchName_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "master", params.MainBranchName)
}

func TestLoadParams_DevelopBranchName(t *testing.T) {
	os.Setenv("INPUT_DEVELOP_BRANCH_NAME", "dev")
	defer os.Unsetenv("INPUT_DEVELOP_BRANCH_NAME")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "dev", params.DevelopBranchName)
}

func TestLoadParams_DevelopBranchName_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "develop", params.DevelopBranchName)
}

func TestLoadParams_PatchPattern(t *testing.T) {
	os.Setenv("INPUT_PATCH_PATTERN", "^fix/.+")
	defer os.Unsetenv("INPUT_PATCH_PATTERN")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "^fix/.+", params.PatchPattern.String())
}

func TestLoadParams_PatchPattern_Invalid(t *testing.T) {
	os.Setenv("INPUT_PATCH_PATTERN", "[")
	defer os.Unsetenv("INPUT_PATCH_PATTERN")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_PatchPattern_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "(?i)^(.+:)?(bugfix/.+)", params.PatchPattern.String())
}

func TestLoadParams_MinorPattern(t *testing.T) {
	os.Setenv("INPUT_MINOR_PATTERN", "^feat/.+")
	defer os.Unsetenv("INPUT_MINOR_PATTERN")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "^feat/.+", params.MinorPattern.String())
}

func TestLoadParams_MinorPattern_Invalid(t *testing.T) {
	os.Setenv("INPUT_MINOR_PATTERN", "[")
	defer os.Unsetenv("INPUT_MINOR_PATTERN")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_MinorPattern_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "(?i)^(.+:)?(feature/.+)", params.MinorPattern.String())
}

func TestLoadParams_MajorPattern(t *testing.T) {
	os.Setenv("INPUT_MAJOR_PATTERN", "^major/.+")
	defer os.Unsetenv("INPUT_MAJOR_PATTERN")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "^major/.+", params.MajorPattern.String())
}

func TestLoadParams_MajorPattern_Invalid(t *testing.T) {
	os.Setenv("INPUT_MAJOR_PATTERN", "[")
	defer os.Unsetenv("INPUT_MAJOR_PATTERN")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_MajorPattern_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "(?i)^(.+:)?(release/.+)", params.MajorPattern.String())
}

func TestLoadParams_BuildPattern(t *testing.T) {
	os.Setenv("INPUT_BUILD_PATTERN", "^build/.+")
	defer os.Unsetenv("INPUT_BUILD_PATTERN")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "^build/.+", params.BuildPattern.String())
}

func TestLoadParams_BuildPattern_Invalid(t *testing.T) {
	os.Setenv("INPUT_BUILD_PATTERN", "[")
	defer os.Unsetenv("INPUT_BUILD_PATTERN")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_BuildPattern_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "(?i)^(.+:)?((doc(s)?|misc)/.+)", params.BuildPattern.String())
}

func TestLoadParams_HotfixPattern(t *testing.T) {
	os.Setenv("INPUT_HOTFIX_PATTERN", "^hotfix/.+")
	defer os.Unsetenv("INPUT_HOTFIX_PATTERN")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "^hotfix/.+", params.HotfixPattern.String())
}

func TestLoadParams_HotfixPattern_Invalid(t *testing.T) {
	os.Setenv("INPUT_HOTFIX_PATTERN", "[")
	defer os.Unsetenv("INPUT_HOTFIX_PATTERN")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_HotfixPattern_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "(?i)^(.+:)?(hotfix/.+)", params.HotfixPattern.String())
}

func TestLoadParams_ExcludePattern(t *testing.T) {
	os.Setenv("INPUT_EXCLUDE_PATTERN", "^ignore/.+")
	defer os.Unsetenv("INPUT_EXCLUDE_PATTERN")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "^ignore/.+", params.ExcludePattern.String())
}

func TestLoadParams_ExcludePattern_Invalid(t *testing.T) {
	os.Setenv("INPUT_EXCLUDE_PATTERN", "[")
	defer os.Unsetenv("INPUT_EXCLUDE_PATTERN")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_ExcludePattern_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Nil(t, params.ExcludePattern)
}

func TestLoadParams_CommitSha(t *testing.T) {
	os.Setenv("GITHUB_SHA", "2f08f7b455ec64741d135216d19d7e0c4dd46458")
	defer os.Unsetenv("GITHUB_SHA")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "2f08f7b455ec64741d135216d19d7e0c4dd46458", params.CommitSha)
}

func TestLoadParams_CommitSha_Invalid(t *testing.T) {
	os.Setenv("GITHUB_SHA", "any")
	defer os.Unsetenv("GITHUB_SHA")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_RepoDir(t *testing.T) {
	os.Setenv("INPUT_REPO_DIR", "/var/tmp/project")
	defer os.Unsetenv("INPUT_REPO_DIR")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "/var/tmp/project", params.RepoDir)
}

func TestLoadParams_RepoDir_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, ".", params.RepoDir)
}

func TestLoadParams_Bump(t *testing.T) {
	tests := map[string]string{
		"auto":  "auto",
		"major": "major",
		"minor": "minor",
		"patch": "patch",
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv("INPUT_BUMP", value)
			defer os.Unsetenv("INPUT_BUMP")

			params, err := generate.LoadParams()
			require.NoError(t, err)

			assert.Equal(t, value, params.Bump)
		})
	}
}

func TestLoadParams_Bump_Invalid(t *testing.T) {
	os.Setenv("INPUT_BUMP", "invalid")
	defer os.Unsetenv("INPUT_BUMP")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_BranchingModel(t *testing.T) {
	tests := map[string]string{
		"git flow":    "git-flow",
		"trunk based": "trunk-based",
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv("INPUT_BRANCHING_MODEL", value)
			defer os.Unsetenv("INPUT_BRANCHING_MODEL")

			params, err := generate.LoadParams()
			require.NoError(t, err)

			assert.Equal(t, value, params.BranchingModel)
		})
	}
}

func TestLoadParams_BranchingModel_Invalid(t *testing.T) {
	os.Setenv("INPUT_BRANCHING_MODEL", "invalid")
	defer os.Unsetenv("INPUT_BRANCHING_MODEL")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_BaseVersion(t *testing.T) {
	os.Setenv("INPUT_BASE_VERSION", "1.2.3")
	defer os.Unsetenv("INPUT_BASE_VERSION")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	var expected = semver.MustParse("1.2.3")

	assert.True(t, expected.EQ(*params.BaseVersion))
}

func TestLoadParams_BaseVersion_Invalid(t *testing.T) {
	os.Setenv("INPUT_BASE_VERSION", "invalid")
	defer os.Unsetenv("INPUT_BASE_VERSION")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_Debug(t *testing.T) {
	os.Setenv("INPUT_DEBUG", "true")
	defer os.Unsetenv("INPUT_DEBUG")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.True(t, params.Debug)
}

func TestLoadParams_Debug_Invalid(t *testing.T) {
	os.Setenv("INPUT_DEBUG", "invalid")
	defer os.Unsetenv("INPUT_DEBUG")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_String(t *testing.T) {
	os.Setenv("INPUT_BUMP", "auto")
	os.Setenv("INPUT_BASE_VERSION", "1.2.3")
	os.Setenv("INPUT_PREFIX", "r")
	os.Setenv("INPUT_PRERELEASE_ID", "alpha")
	os.Setenv("INPUT_MAIN_BRANCH_NAME", "main")
	os.Setenv("INPUT_DEVELOP_BRANCH_NAME", "dev")
	os.Setenv("INPUT_REPO_DIR", "/var/tmp/project")
	os.Setenv("GITHUB_SHA", "2f08f7b455ec64741d135216d19d7e0c4dd46458")
	os.Setenv("INPUT_PATCH_PATTERN", "^bugfix/.+")
	os.Setenv("INPUT_MINOR_PATTERN", "^feat/.+")
	os.Setenv("INPUT_MAJOR_PATTERN", "^major/.+")
	os.Setenv("INPUT_BUILD_PATTERN", "^build/.+")
	os.Setenv("INPUT_HOTFIX_PATTERN", "^hotfix/.+")
	os.Setenv("INPUT_EXCLUDE_PATTERN", "^ignore/.+")
	os.Setenv("INPUT_DEBUG", "true")

	defer func() {
		os.Unsetenv("INPUT_BUMP")
		os.Unsetenv("INPUT_BASE_VERSION")
		os.Unsetenv("INPUT_PREFIX")
		os.Unsetenv("INPUT_PRERELEASE_ID")
		os.Unsetenv("INPUT_MAIN_BRANCH_NAME")
		os.Unsetenv("INPUT_DEVELOP_BRANCH_NAME")
		os.Unsetenv("INPUT_REPO_DIR")
		os.Unsetenv("GITHUB_SHA")
		os.Unsetenv("INPUT_PATCH_PATTERN")
		os.Unsetenv("INPUT_MINOR_PATTERN")
		os.Unsetenv("INPUT_MAJOR_PATTERN")
		os.Unsetenv("INPUT_BUILD_PATTERN")
		os.Unsetenv("INPUT_HOTFIX_PATTERN")
		os.Unsetenv("INPUT_EXCLUDE_PATTERN")
		os.Unsetenv("INPUT_DEBUG")
	}()

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, `commit sha: "2f08f7b455ec64741d135216d19d7e0c4dd46458",`+
		` bump: "auto",`+
		` base version: "1.2.3",`+
		` prefix: "r",`+
		` prerelease id: "alpha",`+
		` main branch name: "main",`+
		` develop branch name: "dev",`+
		` patch pattern: "^bugfix/.+",`+
		` minor pattern: "^feat/.+",`+
		` major pattern: "^major/.+",`+
		` build pattern: "^build/.+",`+
		` hotfix pattern "^hotfix/.+",`+
		` exclude pattern: "^ignore/.+",`+
		` repo dir: "/var/tmp/project",`+
		` debug: true`,
		params.String())
}
