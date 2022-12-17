package strategy_test

import (
	"testing"

	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/internal/strategy"

	"github.com/stretchr/testify/assert"
)

func TestDetermineBumpStrategy(t *testing.T) {
	tests := map[string]struct {
		SourceBranch    string
		DestBranch      string
		Bump            string
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
		"source branch doc, dest branch develop and auto bump": {
			SourceBranch:    "doc/some",
			DestBranch:      "develop",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "",
		},
		"source branch misc, dest branch develop and auto bump": {
			SourceBranch:    "misc/some",
			DestBranch:      "develop",
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
			gf := strategy.New(strategy.Configuration{
				Bump:              test.Bump,
				BranchingModel:    "git-flow",
				DevelopBranchName: "develop",
				MainBranchName:    "master",
				PatchPattern:      regex.MustCompile(`(?i)^bugfix/.+`),
				MinorPattern:      regex.MustCompile(`(?i)^feature/.+`),
				MajorPattern:      regex.MustCompile(`(?i)^major/.+`),
				BuildPattern:      regex.MustCompile(`(?i)^(doc(s)?|misc)/.+`),
			})

			method, version := gf.DetermineBumpStrategy(test.SourceBranch, test.DestBranch)

			assert.Equal(t, test.ExpectedMethod, method)
			assert.Equal(t, test.ExpectedVersion, version)
		})
	}
}
