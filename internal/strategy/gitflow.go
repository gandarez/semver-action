package strategy

import (
	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/pkg/git"
)

// GitFlow implements the git-flow strategy.
type GitFlow struct {
	bump              string
	DevelopBranchName string
	MainBranchName    string
	patchPattern      regex.Regex
	minorPattern      regex.Regex
	majorPattern      regex.Regex
	buildPattern      regex.Regex
}

// DetermineBumpStrategy determines the strategy for semver to bump product version.
func (g *GitFlow) DetermineBumpStrategy(sourceBranch, destBranch string) (string, string) {
	if g.bump != "auto" {
		return g.bump, ""
	}

	// bugfix into develop branch
	if g.patchPattern.MatchString(sourceBranch) && destBranch == g.DevelopBranchName {
		return "build", "patch"
	}

	// feature into develop
	if g.minorPattern.MatchString(sourceBranch) && destBranch == g.DevelopBranchName {
		return "build", "minor"
	}

	// major into develop
	if g.majorPattern.MatchString(sourceBranch) && destBranch == g.DevelopBranchName {
		return "build", "major"
	}

	// build into develop branch
	if g.buildPattern.MatchString(sourceBranch) && destBranch == g.DevelopBranchName {
		return "build", ""
	}

	// develop branch into main branch
	if sourceBranch == g.DevelopBranchName && destBranch == g.MainBranchName {
		return "final", ""
	}

	return "build", ""
}

// Tag implements the Strategy interface.
func (g *GitFlow) Tag(config Configuration, gc git.Client) (Result, error) {
	panic("not implemented")
}
