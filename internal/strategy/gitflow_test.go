package strategy_test

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/internal/strategy"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineBumpStrategy_Gitflow(t *testing.T) {
	tests := map[string]struct {
		SourceBranch    string
		DestBranch      string
		Bump            string
		ExcludePattern  regex.Regex
		ExpectedMethod  string
		ExpectedVersion string
	}{
		"source branch bugfix, dest branch develop and auto bump": {
			SourceBranch:    "bugfix/some",
			DestBranch:      "develop",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "patch",
		},
		"source branch feature, dest branch develop and auto bump": {
			SourceBranch:    "feature/some",
			DestBranch:      "develop",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "minor",
		},
		"source branch major, dest branch develop and auto bump": {
			SourceBranch:    "major/some",
			DestBranch:      "develop",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "major",
		},
		"source branch build, dest branch develop and auto bump": {
			SourceBranch:    "doc/some",
			DestBranch:      "develop",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "",
		},
		"source branch hotfix, dest branch master and auto bump": {
			SourceBranch:   "hotfix/some",
			DestBranch:     "master",
			Bump:           "auto",
			ExpectedMethod: "hotfix",
		},
		"source branch ignore": {
			SourceBranch:    "ignore/some",
			ExcludePattern:  regex.MustCompile(`(?i)^ignore/.+`),
			ExpectedMethod:  "",
			ExpectedVersion: "",
		},
		"source branch ignore, no exclude pattern": {
			SourceBranch:    "ignore/some",
			ExpectedMethod:  "",
			ExpectedVersion: "",
		},
		"source branch ignore, no exclude pattern and auto bump": {
			SourceBranch:    "ignore/some",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "",
		},
		"source branch develop, dest branch master and auto bump": {
			SourceBranch:   "develop",
			DestBranch:     "master",
			Bump:           "auto",
			ExpectedMethod: "final",
		},
		"not a valid source branch prefix and auto bump": {
			SourceBranch:   "some-branch",
			Bump:           "auto",
			ExpectedMethod: "build",
		},
		"patch bump": {
			Bump:           "patch",
			ExpectedMethod: "patch",
		},
		"minor bump": {
			Bump:           "minor",
			ExpectedMethod: "minor",
		},
		"major bump": {
			Bump:           "major",
			ExpectedMethod: "major",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			branchingStrategy, err := strategy.New(strategy.Configuration{
				Bump:              test.Bump,
				BranchingModel:    "git-flow",
				DevelopBranchName: "develop",
				MainBranchName:    "master",
				PatchPattern:      regex.MustCompile(`(?i)^bugfix/.+`),
				MinorPattern:      regex.MustCompile(`(?i)^feature/.+`),
				MajorPattern:      regex.MustCompile(`(?i)^major/.+`),
				BuildPattern:      regex.MustCompile(`(?i)^(doc(s)?|misc)/.+`),
				HotfixPattern:     regex.MustCompile(`(?i)^hotfix/.+`),
				ExcludePattern:    test.ExcludePattern,
			})
			require.NoError(t, err)

			method, version := branchingStrategy.DetermineBumpStrategy(test.SourceBranch, test.DestBranch)

			assert.Equal(t, test.ExpectedMethod, method)
			assert.Equal(t, test.ExpectedVersion, version)
		})
	}
}

func TestTag_Gitflow(t *testing.T) {
	tests := map[string]struct {
		Method      string
		AncestorTag string
		Tag         *semver.Version
		Version     string
		Expected    strategy.Result
	}{
		"build": {
			Method:      "build",
			Tag:         newSemVerPtr(t, "1.2.3-alpha.2"),
			AncestorTag: "v1.2.3-alpha.1",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.3-alpha.1",
				SemverTag:    "v1.2.3-alpha.3",
				IsPrerelease: true,
			},
		},
		"major with pre release tag": {
			Method:      "build",
			Version:     "major",
			Tag:         newSemVerPtr(t, "1.2.3-alpha.1"),
			AncestorTag: "v1.2.2-alpha.1",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.2-alpha.1",
				SemverTag:    "v2.0.0-alpha.1",
				IsPrerelease: true,
			},
		},
		"major without pre release tag": {
			Method:      "major",
			Tag:         newSemVerPtr(t, "1.2.3"),
			AncestorTag: "v1.2.2",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.2",
				SemverTag:    "v2.0.0",
				IsPrerelease: false,
			},
		},
		"minor with pre release tag": {
			Method:      "build",
			Version:     "minor",
			Tag:         newSemVerPtr(t, "1.2.3-alpha.1"),
			AncestorTag: "v1.2.2-alpha.1",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.2-alpha.1",
				SemverTag:    "v1.3.0-alpha.1",
				IsPrerelease: true,
			},
		},
		"minor without pre release tag": {
			Method:      "minor",
			Tag:         newSemVerPtr(t, "1.2.3"),
			AncestorTag: "v1.2.2",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.2",
				SemverTag:    "v1.3.0",
				IsPrerelease: false,
			},
		},
		"patch with pre release tag": {
			Method:      "build",
			Version:     "patch",
			Tag:         newSemVerPtr(t, "1.2.3-alpha.1"),
			AncestorTag: "v1.2.2-alpha.1",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.2-alpha.1",
				SemverTag:    "v1.2.4-alpha.1",
				IsPrerelease: true,
			},
		},
		"patch without pre release tag": {
			Method:      "patch",
			Tag:         newSemVerPtr(t, "1.2.3"),
			AncestorTag: "v1.2.2",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.2",
				SemverTag:    "v1.2.4",
				IsPrerelease: false,
			},
		},
		"final version": {
			Method:      "final",
			Tag:         newSemVerPtr(t, "1.2.3-alpha.1"),
			AncestorTag: "v1.2.2-alpha.1",
			Expected: strategy.Result{
				AncestorTag:  "v1.2.2-alpha.1",
				SemverTag:    "v1.2.3",
				IsPrerelease: false,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gf := strategy.GitFlow{}

			gc := initGitClientMock(
				t, "", test.AncestorTag, "", "", "",
			)

			result, err := gf.Tag(strategy.TagParams{
				DestBranch:   "not-used",
				Prefix:       "v",
				PrereleaseID: "alpha",
				Method:       test.Method,
				Tag:          test.Tag,
				Version:      test.Version,
			}, gc)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, result)
		})
	}
}
