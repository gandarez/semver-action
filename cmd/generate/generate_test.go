package generate_test

import (
	"errors"
	"testing"

	"github.com/gandarez/semver-action/cmd/generate"
	"github.com/gandarez/semver-action/internal/regex"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	tests := map[string]struct {
		CurrentBranch string
		LatestTag     string
		AncestorTag   string
		SourceBranch  string
		Params        func() generate.Params
		Result        generate.Result
	}{
		"no previous tag": {
			CurrentBranch: "develop",
			SourceBranch:  "release/semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v0.0.0",
				AncestorTag:  "",
				SemverTag:    "v1.0.0-pre.1",
				IsPrerelease: true,
			},
		},
		"first non-development tag": {
			CurrentBranch: "master",
			LatestTag:     "1.0.0-pre.1",
			AncestorTag:   "e63c125b",
			SourceBranch:  "develop",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v1.0.0-pre.1",
				AncestorTag:  "e63c125b",
				SemverTag:    "v1.0.0",
				IsPrerelease: false,
			},
		},
		"feature branch into develop": {
			CurrentBranch: "develop",
			LatestTag:     "0.2.1",
			SourceBranch:  "feature/semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v0.2.1",
				SemverTag:    "v0.3.0-pre.1",
				IsPrerelease: true,
			},
		},
		"doc branch into develop": {
			CurrentBranch: "develop",
			LatestTag:     "0.2.1-pre.1",
			SourceBranch:  "doc/semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v0.2.1-pre.1",
				SemverTag:    "v0.2.1-pre.2",
				IsPrerelease: true,
			},
		},
		"misc branch into develop": {
			CurrentBranch: "develop",
			LatestTag:     "0.2.1-pre.1",
			SourceBranch:  "misc/semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v0.2.1-pre.1",
				SemverTag:    "v0.2.1-pre.2",
				IsPrerelease: true,
			},
		},
		"exclude branch": {
			CurrentBranch: "develop",
			LatestTag:     "0.2.1-pre.1",
			SourceBranch:  "ignore/semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				p.ExcludePattern = regex.MustCompile(`(?i)^ignore/.+`)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "",
				AncestorTag:  "",
				SemverTag:    "",
				IsPrerelease: false,
			},
		},
		"merge develop into master": {
			CurrentBranch: "master",
			LatestTag:     "1.4.17-pre.1",
			SourceBranch:  "develop",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v1.4.17-pre.1",
				SemverTag:    "v1.4.17",
				IsPrerelease: false,
			},
		},
		"merge develop into master with previous matching tag": {
			CurrentBranch: "master",
			LatestTag:     "1.4.17-pre.1",
			AncestorTag:   "v1.4.16",
			SourceBranch:  "develop",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v1.4.17-pre.1",
				AncestorTag:  "v1.4.16",
				SemverTag:    "v1.4.17",
				IsPrerelease: false,
			},
		},
		"base version set": {
			CurrentBranch: "develop",
			LatestTag:     "2.6.19",
			SourceBranch:  "feature/semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				p.BaseVersion = newSemVerPtr(t, "4.2.0")

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19",
				SemverTag:    "v4.3.0-pre.1",
				IsPrerelease: true,
			},
		},
		"invalid branch name": {
			CurrentBranch: "develop",
			LatestTag:     "2.6.19-pre.1",
			SourceBranch:  "semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19-pre.1",
				SemverTag:    "v2.6.19-pre.2",
				IsPrerelease: true,
			},
		},
		"force bump major": {
			CurrentBranch: "develop",
			LatestTag:     "2.6.19-pre.1",
			SourceBranch:  "semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				p.Bump = "major"

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19-pre.1",
				SemverTag:    "v3.0.0-pre.1",
				IsPrerelease: true,
			},
		},
		"force bump minor": {
			CurrentBranch: "develop",
			LatestTag:     "2.6.19-pre.1",
			SourceBranch:  "semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				p.Bump = "minor"

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19-pre.1",
				SemverTag:    "v2.7.0-pre.1",
				IsPrerelease: true,
			},
		},
		"force bump patch": {
			CurrentBranch: "develop",
			LatestTag:     "2.6.19-pre.1",
			SourceBranch:  "semver-initial",
			Params: func() generate.Params {
				p, err := generate.LoadParams()
				require.NoError(t, err)

				p.Bump = "patch"

				return p
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19-pre.1",
				SemverTag:    "v2.6.20-pre.1",
				IsPrerelease: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := test.Params()

			gc := initGitClientMock(
				t,
				test.LatestTag,
				test.AncestorTag,
				test.CurrentBranch,
				test.SourceBranch,
				p.CommitSha,
			)

			result, err := generate.Tag(p, gc)
			require.NoError(t, err)

			assert.Equal(t, test.Result, result)
		})
	}
}

func TestTag_IsNotRepo(t *testing.T) {
	gc := &gitClientMock{
		MakeSafeFn: func() error {
			return nil
		},
		IsRepoFn: func() bool {
			return false
		},
	}

	_, err := generate.Tag(generate.Params{}, gc)
	require.Error(t, err)

	assert.EqualError(t, err, "current folder is not a git repository")
}

func TestTag_MakeSafeErr(t *testing.T) {
	gc := &gitClientMock{
		MakeSafeFn: func() error {
			return errors.New("error")
		},
	}

	_, err := generate.Tag(generate.Params{}, gc)
	require.Error(t, err)

	assert.EqualError(t, err, "failed to make safe: error")
}

type gitClientMock struct {
	CurrentBranchFn        func() (string, error)
	CurrentBranchFnInvoked int
	IsRepoFn               func() bool
	IsRepoFnInvoked        int
	MakeSafeFn             func() error
	MakeSafeFnInvoked      int
	LatestTagFn            func() string
	LatestTagFnInvoked     int
	AncestorTagFn          func(include, exclude, branch string) string
	AncestorTagFnInvoked   int
	SourceBranchFn         func(commitHash string) (string, error)
	SourceBranchFnInvoked  int
}

func initGitClientMock(t *testing.T, latestTag, ancestorTag, currentBranch, sourceBranch, expectedCommitHash string) *gitClientMock {
	return &gitClientMock{
		CurrentBranchFn: func() (string, error) {
			return currentBranch, nil
		},
		IsRepoFn: func() bool {
			return true
		},
		MakeSafeFn: func() error {
			return nil
		},
		LatestTagFn: func() string {
			return latestTag
		},
		AncestorTagFn: func(include, exclude, branch string) string {
			return ancestorTag
		},
		SourceBranchFn: func(commitHash string) (string, error) {
			assert.Equal(t, expectedCommitHash, commitHash)
			return sourceBranch, nil
		},
	}
}

func (m *gitClientMock) CurrentBranch() (string, error) {
	m.CurrentBranchFnInvoked += 1
	return m.CurrentBranchFn()
}
func (m *gitClientMock) IsRepo() bool {
	m.IsRepoFnInvoked += 1
	return m.IsRepoFn()
}

func (m *gitClientMock) MakeSafe() error {
	m.MakeSafeFnInvoked++
	return m.MakeSafeFn()
}

func (m *gitClientMock) LatestTag() string {
	m.LatestTagFnInvoked += 1
	return m.LatestTagFn()
}

func (m *gitClientMock) AncestorTag(include, exclude, branch string) string {
	m.AncestorTagFnInvoked += 1
	return m.AncestorTagFn(include, exclude, branch)
}

func (m *gitClientMock) SourceBranch(commitHash string) (string, error) {
	m.SourceBranchFnInvoked += 1
	return m.SourceBranchFn(commitHash)
}

func newSemVerPtr(t *testing.T, s string) *semver.Version {
	version, err := semver.New(s)
	require.NoError(t, err)

	return version
}
