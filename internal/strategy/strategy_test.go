package strategy_test

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func initGitClientMock(
	t *testing.T,
	latestTag,
	ancestorTag,
	currentBranch,
	sourceBranch,
	expectedCommitHash string) *gitClientMock {
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
	m.CurrentBranchFnInvoked++
	return m.CurrentBranchFn()
}

func (m *gitClientMock) MakeSafe() error {
	m.MakeSafeFnInvoked++
	return m.MakeSafeFn()
}

func (m *gitClientMock) IsRepo() bool {
	m.IsRepoFnInvoked++
	return m.IsRepoFn()
}

func (m *gitClientMock) LatestTag() string {
	m.LatestTagFnInvoked++
	return m.LatestTagFn()
}

func (m *gitClientMock) AncestorTag(include, exclude, branch string) string {
	m.AncestorTagFnInvoked++
	return m.AncestorTagFn(include, exclude, branch)
}

func (m *gitClientMock) SourceBranch(commitHash string) (string, error) {
	m.SourceBranchFnInvoked++
	return m.SourceBranchFn(commitHash)
}

func newSemVerPtr(t *testing.T, s string) *semver.Version {
	version, err := semver.New(s)
	require.NoError(t, err)

	return version
}
