package strategy_test

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/internal/strategy"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineBumpStrategy_TrunkBased(t *testing.T) {
	tests := map[string]struct {
		SourceBranch   string
		DestBranch     string
		Bump           string
		ExcludePattern regex.Regex
		ExpectedMethod string
	}{
		"source branch bugfix, dest branch master and auto bump": {
			SourceBranch:   "bugfix/some",
			DestBranch:     "master",
			Bump:           "auto",
			ExpectedMethod: "patch",
		},
		"source branch feature, dest branch master and auto bump": {
			SourceBranch:   "feature/some",
			DestBranch:     "master",
			Bump:           "auto",
			ExpectedMethod: "minor",
		},
		"source branch major, dest branch master and auto bump": {
			SourceBranch:   "major/some",
			DestBranch:     "master",
			Bump:           "auto",
			ExpectedMethod: "major",
		},
		"source branch build, dest branch master and auto bump": {
			SourceBranch:   "doc/some",
			DestBranch:     "master",
			Bump:           "auto",
			ExpectedMethod: "build",
		},
		"source branch ignore": {
			SourceBranch:   "ignore/some",
			ExcludePattern: regex.MustCompile(`(?i)^ignore/.+`),
			ExpectedMethod: "",
		},
		"source branch ignore, no exclude pattern": {
			SourceBranch:   "ignore/some",
			ExpectedMethod: "",
		},
		"source branch ignore, no exclude pattern and auto bump": {
			SourceBranch:   "ignore/some",
			Bump:           "auto",
			ExpectedMethod: "build",
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
				Bump:           test.Bump,
				BranchingModel: "trunk-based",
				MainBranchName: "master",
				PatchPattern:   regex.MustCompile(`(?i)^bugfix/.+`),
				MinorPattern:   regex.MustCompile(`(?i)^feature/.+`),
				MajorPattern:   regex.MustCompile(`(?i)^major/.+`),
				BuildPattern:   regex.MustCompile(`(?i)^(doc(s)?|misc)/.+`),
				ExcludePattern: test.ExcludePattern,
			})
			require.NoError(t, err)

			method, version := branchingStrategy.DetermineBumpStrategy(test.SourceBranch, test.DestBranch)

			assert.Equal(t, test.ExpectedMethod, method)
			assert.Empty(t, version)
		})
	}
}

func TestTag_Trunkbased(t *testing.T) {
	tests := map[string]struct {
		Method   string
		Tag      *semver.Version
		Version  string
		Expected strategy.Result
	}{
		"build": {
			Method: "build",
			Tag:    newSemVerPtr(t, "1.2.3"),
			Expected: strategy.Result{
				SemverTag:    "v1.2.3+1",
				IsPrerelease: false,
			},
		},
		"build with previous tag containing build": {
			Method: "build",
			Tag:    newSemVerPtr(t, "1.2.3+1"),
			Expected: strategy.Result{
				SemverTag:    "v1.2.3+2",
				IsPrerelease: false,
			},
		},
		"major": {
			Method: "major",
			Tag:    newSemVerPtr(t, "1.2.3"),
			Expected: strategy.Result{
				SemverTag:    "v2.0.0",
				IsPrerelease: false,
			},
		},
		"minor": {
			Method: "minor",
			Tag:    newSemVerPtr(t, "1.2.3"),
			Expected: strategy.Result{
				SemverTag:    "v1.3.0",
				IsPrerelease: false,
			},
		},
		"patch": {
			Method: "patch",
			Tag:    newSemVerPtr(t, "1.2.3"),
			Expected: strategy.Result{
				SemverTag:    "v1.2.4",
				IsPrerelease: false,
			},
		},
		"default": {
			Method: "not-in-use",
			Tag:    newSemVerPtr(t, "1.2.3"),
			Expected: strategy.Result{
				SemverTag:    "v1.2.3",
				IsPrerelease: false,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tb := strategy.TrunkBased{}

			gc := initGitClientMock(
				t, "", "", "", "", "",
			)

			result, err := tb.Tag(strategy.TagParams{
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
