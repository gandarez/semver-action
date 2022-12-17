package strategy

import (
	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/pkg/git"
)

type (
	// Strategy defines the interface for a strategy.
	Strategy interface {
		DetermineBumpStrategy(sourceBranch, destBranch string) (string, string)
		Tag(config Configuration, gc git.Client) (Result, error)
	}

	// Configuration contains the strategy configuration.
	Configuration struct {
		Bump              string
		BranchingModel    string
		MainBranchName    string
		DevelopBranchName string
		PatchPattern      regex.Regex
		MinorPattern      regex.Regex
		MajorPattern      regex.Regex
		BuildPattern      regex.Regex
	}

	// Result contains the result of strategy execution.
	Result struct {
		PreviousTag  string
		AncestorTag  string
		SemverTag    string
		IsPrerelease bool
	}
)

// New returns a new strategy.
func New(config Configuration) Strategy {
	switch config.BranchingModel {
	case "git-flow":
		return &GitFlow{
			bump:              config.Bump,
			DevelopBranchName: config.DevelopBranchName,
			MainBranchName:    config.MainBranchName,
			patchPattern:      config.PatchPattern,
			minorPattern:      config.MinorPattern,
			majorPattern:      config.MajorPattern,
			buildPattern:      config.BuildPattern,
		}
	case "trunk-based":
		return &TrunkBased{}
	default:
		return &GitFlow{}
	}
}
