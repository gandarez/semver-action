package generate

import (
	"fmt"

	"github.com/gandarez/semver-action/internal/strategy"
	"github.com/gandarez/semver-action/pkg/git"

	"github.com/apex/log"
	"github.com/blang/semver/v4"
)

const initialTag = "0.0.0"

// Result contains the result of Run().
type Result struct {
	PreviousTag  string
	AncestorTag  string
	SemverTag    string
	IsPrerelease bool
}

// Run generates a semantic version using the commit sha.
func Run() (Result, error) {
	params, err := LoadParams()
	if err != nil {
		return Result{}, fmt.Errorf("failed to load parameters: %s", err)
	}

	if params.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logs enabled\n")
	}

	log.Debug(params.String())

	gc := git.New(params.RepoDir)

	return Tag(params, gc)
}

// Tag returns the calculated semantic version.
func Tag(params Params, gc git.Git) (Result, error) {
	err := gc.MakeSafe()
	if err != nil {
		return Result{}, fmt.Errorf("failed to make safe: %s", err)
	}

	if !gc.IsRepo() {
		return Result{}, fmt.Errorf("current folder is not a git repository")
	}

	dest, err := gc.CurrentBranch()
	if err != nil {
		return Result{}, fmt.Errorf("failed to extract dest branch from commit: %s", err)
	}

	log.Debugf("dest branch: %q\n", dest)

	source, err := gc.SourceBranch(params.CommitSha)
	if err != nil {
		return Result{}, fmt.Errorf("failed to extract source branch from commit: %s", err)
	}

	log.Debugf("source branch: %q\n", source)

	branchingStrategy, err := strategy.New(strategy.Configuration{
		Bump:              params.Bump,
		BranchingModel:    params.BranchingModel,
		MainBranchName:    params.MainBranchName,
		DevelopBranchName: params.DevelopBranchName,
		PatchPattern:      params.PatchPattern,
		MinorPattern:      params.MinorPattern,
		MajorPattern:      params.MajorPattern,
		BuildPattern:      params.BuildPattern,
		HotfixPattern:     params.HotfixPattern,
		ExcludePattern:    params.ExcludePattern,
	})
	if err != nil {
		return Result{}, fmt.Errorf("failed to decide branching strategy: %s", err)
	}

	log.Debugf("using branching strategy: %q\n", branchingStrategy.Name())

	method, version := branchingStrategy.DetermineBumpStrategy(source, dest)

	log.Debugf("method: %q, version: %q", method, version)

	if method == "" && version == "" {
		log.Info("no version bump required")

		return Result{}, nil
	}

	latestTag := gc.LatestTag()

	var tag *semver.Version

	if latestTag == "" {
		tag, _ = semver.New(initialTag)
	} else {
		parsed, err := semver.ParseTolerant(latestTag)
		if err != nil {
			return Result{}, fmt.Errorf("failed to parse tag %q or not valid semantic version: %s", latestTag, err)
		}
		tag = &parsed
	}

	previousTag := params.Prefix + tag.String()

	if params.BaseVersion != nil {
		tag = params.BaseVersion
	}

	result, err := branchingStrategy.Tag(strategy.TagParams{
		DestBranch:   dest,
		Method:       method,
		Prefix:       params.Prefix,
		PrereleaseID: params.PrereleaseID,
		Tag:          tag,
		Version:      version,
	}, gc)
	if err != nil {
		return Result{}, fmt.Errorf("failed to tag: %s", err)
	}

	log.Debugf("result: %+v\n", result)

	return Result{
		PreviousTag:  previousTag,
		AncestorTag:  result.AncestorTag,
		SemverTag:    result.SemverTag,
		IsPrerelease: result.IsPrerelease,
	}, nil
}
