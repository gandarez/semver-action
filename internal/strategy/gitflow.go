package strategy

import (
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/pkg/git"

	"github.com/blang/semver/v4"
)

// GitFlow implements the git-flow strategy.
type GitFlow struct {
	bump              string
	developBranchName string
	mainBranchName    string
	patchPattern      regex.Regex
	minorPattern      regex.Regex
	majorPattern      regex.Regex
	buildPattern      regex.Regex
	hotfixPattern     regex.Regex
	excludePattern    regex.Regex
}

// DetermineBumpStrategy determines the strategy for semver to bump product version.
func (g *GitFlow) DetermineBumpStrategy(sourceBranch, destBranch string) (string, string) {
	// if source branch is excluded, do not bump
	if g.excludePattern != nil && g.excludePattern.MatchString(sourceBranch) {
		return "", ""
	}

	// if bump is not auto, return it
	if g.bump != "auto" {
		return g.bump, ""
	}

	// bugfix into develop branch
	if g.patchPattern.MatchString(sourceBranch) && destBranch == g.developBranchName {
		return "build", "patch"
	}

	// feature into develop
	if g.minorPattern.MatchString(sourceBranch) && destBranch == g.developBranchName {
		return "build", "minor"
	}

	// major into develop
	if g.majorPattern.MatchString(sourceBranch) && destBranch == g.developBranchName {
		return "build", "major"
	}

	// build into develop branch
	if g.buildPattern.MatchString(sourceBranch) && destBranch == g.developBranchName {
		return "build", ""
	}

	// hotfix into main branch
	if g.hotfixPattern.MatchString(sourceBranch) && destBranch == g.mainBranchName {
		return "hotfix", ""
	}

	// develop branch into main branch
	if sourceBranch == g.developBranchName && destBranch == g.mainBranchName {
		return "final", ""
	}

	return "build", ""
}

// Tag implements the Strategy interface.
func (g *GitFlow) Tag(params TagParams, gc git.Git) (Result, error) {
	var (
		finalTag       string
		includePattern string
		excludePattern string
		isPrerelease   bool
	)

	if (params.Version == "major" && params.Method == "build") || params.Method == "major" {
		log.Debug("incrementing major")

		if err := params.Tag.IncrementMajor(); err != nil {
			return Result{}, fmt.Errorf("failed to increment major version: %s", err)
		}
	}

	if (params.Version == "minor" && params.Method == "build") || params.Method == "minor" {
		log.Debug("incrementing minor")

		if err := params.Tag.IncrementMinor(); err != nil {
			return Result{}, fmt.Errorf("failed to increment minor version: %s", err)
		}
	}

	if (params.Version == "patch" && params.Method == "build") || params.Method == "patch" || params.Method == "hotfix" {
		log.Debug("incrementing patch")

		if err := params.Tag.IncrementPatch(); err != nil {
			return Result{}, fmt.Errorf("failed to increment patch version: %s", err)
		}
	}

	switch params.Method {
	case "build":
		{
			isPrerelease = true
			includePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)

			buildNumber, _ := semver.NewPRVersion("0")

			if len(params.Tag.Pre) > 1 && params.Version == "" {
				buildNumber = params.Tag.Pre[1]
			}

			params.Tag.Pre = nil

			preVersion, err := semver.NewPRVersion(params.PrereleaseID)
			if err != nil {
				return Result{}, fmt.Errorf("failed to create new pre-release version: %s", err)
			}

			params.Tag.Pre = append(params.Tag.Pre, preVersion)

			buildVersion, err := semver.NewPRVersion(strconv.Itoa(int(buildNumber.VersionNum + 1)))
			if err != nil {
				return Result{}, fmt.Errorf("failed to create new build version: %s", err)
			}

			params.Tag.Pre = append(params.Tag.Pre, buildVersion)

			finalTag = params.Prefix + params.Tag.String()
		}
	case "major", "minor", "patch":
		if len(params.Tag.Pre) > 0 {
			isPrerelease = true
			includePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)
		} else {
			includePattern = fmt.Sprintf("%s[0-9]*", params.Prefix)
			excludePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)
		}

		finalTag = params.Prefix + params.Tag.String()
	default:
		includePattern = fmt.Sprintf("%s[0-9]*", params.Prefix)
		excludePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)
		finalTag = params.Prefix + params.Tag.FinalizeVersion()
	}

	return Result{
		AncestorTag:  gc.AncestorTag(includePattern, excludePattern, params.DestBranch),
		SemverTag:    finalTag,
		IsPrerelease: isPrerelease,
	}, nil
}

// Name returns the name of the strategy.
func (GitFlow) Name() string {
	return "git-flow"
}
