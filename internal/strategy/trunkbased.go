package strategy

import (
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/pkg/git"

	"github.com/blang/semver/v4"
)

// TrunkBased implements the trunk-based strategy.
type TrunkBased struct {
	bump           string
	branchName     string
	patchPattern   regex.Regex
	minorPattern   regex.Regex
	majorPattern   regex.Regex
	buildPattern   regex.Regex
	excludePattern regex.Regex
}

// DetermineBumpStrategy determines the strategy for semver to bump product version.
func (t *TrunkBased) DetermineBumpStrategy(sourceBranch, destBranch string) (string, string) {
	// if source branch is excluded, do not bump
	if t.excludePattern != nil && t.excludePattern.MatchString(sourceBranch) {
		return "", ""
	}

	// if bump is not auto, return it
	if t.bump != "auto" {
		return t.bump, ""
	}

	// bugfix into main branch
	if t.patchPattern.MatchString(sourceBranch) && destBranch == t.branchName {
		return "patch", ""
	}

	// feature into main branch
	if t.minorPattern.MatchString(sourceBranch) && destBranch == t.branchName {
		return "minor", ""
	}

	// major into main branch
	if t.majorPattern.MatchString(sourceBranch) && destBranch == t.branchName {
		return "major", ""
	}

	// build into main branch
	if t.buildPattern.MatchString(sourceBranch) && destBranch == t.branchName {
		return "build", ""
	}

	return "build", ""
}

// Tag implements the Strategy interface.
func (t *TrunkBased) Tag(params TagParams, gc git.Git) (Result, error) {
	var finalTag string

	switch params.Method {
	case "build":
		{
			buildNumberStr, _ := semver.NewBuildVersion("0")

			if len(params.Tag.Build) > 0 && params.Version == "" {
				buildNumberStr = params.Tag.Build[len(params.Tag.Build)-1]
			}

			buildNumber, _ := strconv.Atoi(buildNumberStr)
			buildNumber++

			params.Tag.Build = []string{strconv.Itoa(buildNumber)}

			finalTag = params.Prefix + params.Tag.String()
		}
	case "major":
		{
			log.Debug("incrementing major")

			if err := params.Tag.IncrementMajor(); err != nil {
				return Result{}, fmt.Errorf("failed to increment major version: %s", err)
			}

			finalTag = params.Prefix + params.Tag.FinalizeVersion()
		}
	case "minor":
		{
			log.Debug("incrementing minor")

			if err := params.Tag.IncrementMinor(); err != nil {
				return Result{}, fmt.Errorf("failed to increment minor version: %s", err)
			}

			finalTag = params.Prefix + params.Tag.FinalizeVersion()
		}
	case "patch":
		{
			log.Debug("incrementing patch")

			if err := params.Tag.IncrementPatch(); err != nil {
				return Result{}, fmt.Errorf("failed to increment patch version: %s", err)
			}

			finalTag = params.Prefix + params.Tag.FinalizeVersion()
		}
	default:
		finalTag = params.Prefix + params.Tag.FinalizeVersion()
	}

	return Result{
		AncestorTag:  "",
		SemverTag:    finalTag,
		IsPrerelease: false,
	}, nil
}

// Name returns the name of the strategy.
func (TrunkBased) Name() string {
	return "trunk-based"
}
