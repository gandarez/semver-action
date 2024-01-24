package strategy

import (
	"errors"

	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/pkg/git"

	"github.com/blang/semver/v4"
)

type (
	// Strategy defines the interface for a strategy.
	Strategy interface {
		DetermineBumpStrategy(sourceBranch, destBranch string) (string, string)
		Tag(params TagParams, gc git.Git) (Result, error)
		Name() string
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
		HotfixPattern     regex.Regex
		ExcludePattern    regex.Regex
	}

	// TagParams contains the parameters for Tag().
	TagParams struct {
		DestBranch   string
		Method       string
		Prefix       string
		PrereleaseID string
		Tag          *semver.Version
		Version      string
	}

	// Result contains the result of strategy execution.
	Result struct {
		AncestorTag  string
		SemverTag    string
		IsPrerelease bool
	}
)

// New returns a new strategy.
func New(config Configuration) (Strategy, error) {
	switch config.BranchingModel {
	case "git-flow":
		return &GitFlow{
			bump:              config.Bump,
			developBranchName: config.DevelopBranchName,
			mainBranchName:    config.MainBranchName,
			patchPattern:      config.PatchPattern,
			minorPattern:      config.MinorPattern,
			majorPattern:      config.MajorPattern,
			buildPattern:      config.BuildPattern,
			hotfixPattern:     config.HotfixPattern,
			excludePattern:    config.ExcludePattern,
		}, nil
	case "trunk-based":
		return &TrunkBased{
			bump:           config.Bump,
			branchName:     config.MainBranchName,
			patchPattern:   config.PatchPattern,
			minorPattern:   config.MinorPattern,
			majorPattern:   config.MajorPattern,
			buildPattern:   config.BuildPattern,
			excludePattern: config.ExcludePattern,
		}, nil
	default:
		return nil, errors.New("invalid branching model")
	}
}
