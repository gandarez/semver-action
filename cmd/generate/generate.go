package generate

import (
	"fmt"
	"strconv"

	"github.com/gandarez/semver-action/internal/strategy"
	"github.com/gandarez/semver-action/pkg/git"

	"github.com/apex/log"
	"github.com/blang/semver/v4"
)

const tagDefault = "0.0.0"

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

	branchingStrategy := strategy.New(strategy.Configuration{
		Bump:              params.Bump,
		BranchingModel:    params.BranchingModel,
		MainBranchName:    params.MainBranchName,
		DevelopBranchName: params.DevelopBranchName,
		PatchPattern:      params.PatchPattern,
		MinorPattern:      params.MinorPattern,
		MajorPattern:      params.MajorPattern,
		BuildPattern:      params.BuildPattern,
	})

	method, version := branchingStrategy.DetermineBumpStrategy(source, dest)

	log.Debugf("method: %q, version: %q", method, version)

	latestTag := gc.LatestTag()

	var tag *semver.Version

	if latestTag == "" {
		tag, _ = semver.New(tagDefault)
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

	if (version == "major" && method == "build") || method == "major" {
		log.Debug("incrementing major")
		if err := tag.IncrementMajor(); err != nil {
			return Result{}, fmt.Errorf("failed to increment major version: %s", err)
		}
	}

	if (version == "minor" && method == "build") || method == "minor" {
		log.Debug("incrementing minor")
		if err := tag.IncrementMinor(); err != nil {
			return Result{}, fmt.Errorf("failed to increment minor version: %s", err)
		}
	}

	if (version == "patch" && method == "build") || method == "patch" || method == "hotfix" {
		log.Debug("incrementing patch")
		if err := tag.IncrementPatch(); err != nil {
			return Result{}, fmt.Errorf("failed to increment patch version: %s", err)
		}
	}

	var (
		finalTag       string
		ancestorTag    string
		includePattern string
		excludePattern string
		isPrerelease   bool
	)

	switch method {
	case "build":
		{
			isPrerelease = true
			includePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)

			buildNumber, _ := semver.NewPRVersion("0")

			if len(tag.Pre) > 1 && version == "" {
				buildNumber = tag.Pre[1]
			}

			tag.Pre = nil

			preVersion, err := semver.NewPRVersion(params.PrereleaseID)
			if err != nil {
				return Result{}, fmt.Errorf("failed to create new pre-release version: %s", err)
			}

			tag.Pre = append(tag.Pre, preVersion)

			buildVersion, err := semver.NewPRVersion(strconv.Itoa(int(buildNumber.VersionNum + 1)))
			if err != nil {
				return Result{}, fmt.Errorf("failed to create new build version: %s", err)
			}

			tag.Pre = append(tag.Pre, buildVersion)

			finalTag = params.Prefix + tag.String()
		}
	case "major", "minor", "patch":
		if len(tag.Pre) > 0 {
			isPrerelease = true
			includePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)
		} else {
			includePattern = fmt.Sprintf("%s[0-9]*", params.Prefix)
			excludePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)
		}

		finalTag = params.Prefix + tag.String()
	default:
		includePattern = fmt.Sprintf("%s[0-9]*", params.Prefix)
		excludePattern = fmt.Sprintf("%s[0-9]*-%s*", params.Prefix, params.PrereleaseID)
		finalTag = params.Prefix + tag.FinalizeVersion()
	}

	ancestorTag = gc.AncestorTag(includePattern, excludePattern, dest)

	return Result{
		PreviousTag:  previousTag,
		AncestorTag:  ancestorTag,
		SemverTag:    finalTag,
		IsPrerelease: isPrerelease,
	}, nil
}
