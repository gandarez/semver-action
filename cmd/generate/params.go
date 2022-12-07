package generate

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/gandarez/semver-action/internal/regex"
	"github.com/gandarez/semver-action/pkg/actions"

	"github.com/blang/semver/v4"
)

var (
	//nolint
	commitShaRegex = regexp.MustCompile(`\b[0-9a-f]{5,40}\b`)
	// nolint
	validBumpStrategies = []string{"auto", "major", "minor", "patch"}
	// nolint
	validBranchingModels = []string{"git-flow", "trunk-based"}
)

// Params contains semver generate command parameters.
type Params struct {
	CommitSha         string
	RepoDir           string
	Bump              string
	BranchingModel    string
	BaseVersion       *semver.Version
	Prefix            string
	PrereleaseID      string
	MainBranchName    string
	DevelopBranchName string
	PatchPattern      regex.Regex
	MinorPattern      regex.Regex
	MajorPattern      regex.Regex
	BuildPattern      regex.Regex
	Debug             bool
}

// LoadParams loads semver generate config params.
func LoadParams() (Params, error) {
	var commitSha string

	if commitShaStr := os.Getenv("GITHUB_SHA"); commitShaStr != "" {
		if !commitShaRegex.MatchString(commitShaStr) {
			return Params{}, fmt.Errorf("invalid commit-sha format: %s", commitShaStr)
		}

		commitSha = commitShaStr
	}

	var repoDir string = "."

	if repoDirStr := actions.GetInput("repo_dir"); repoDirStr != "" {
		repoDir = repoDirStr
	}

	var bump string = "auto"

	if bumpStr := actions.GetInput("bump"); bumpStr != "" {
		if !stringInSlice(bumpStr, validBumpStrategies) {
			return Params{}, fmt.Errorf("invalid bump value: %s", bumpStr)
		}

		bump = bumpStr
	}

	var branchingModel string = "git-flow"

	if branchingModelStr := actions.GetInput("branching_model"); branchingModelStr != "" {
		if !stringInSlice(branchingModelStr, validBranchingModels) {
			return Params{}, fmt.Errorf("invalid branching model value: %s", branchingModelStr)
		}

		branchingModel = branchingModelStr
	}

	var patchPattern = regex.MustCompile("(?i)^bugfix/.+")

	if patchPatternStr := actions.GetInput("patch_pattern"); patchPatternStr != "" {
		compiled, err := regex.Compile(patchPatternStr)
		if err != nil {
			return Params{}, fmt.Errorf("invalid patch pattern value: %s", patchPatternStr)
		}

		patchPattern = compiled
	}

	var minorPattern = regex.MustCompile("(?i)^feature/.+")

	if minorPatternStr := actions.GetInput("minor_pattern"); minorPatternStr != "" {
		compiled, err := regex.Compile(minorPatternStr)
		if err != nil {
			return Params{}, fmt.Errorf("invalid minor pattern value: %s", minorPatternStr)
		}

		minorPattern = compiled
	}

	var majorPattern = regex.MustCompile("(?i)^release/.+")

	if majorPatternStr := actions.GetInput("major_pattern"); majorPatternStr != "" {
		compiled, err := regex.Compile(majorPatternStr)
		if err != nil {
			return Params{}, fmt.Errorf("invalid major pattern value: %s", majorPatternStr)
		}

		majorPattern = compiled
	}

	var buildPattern = regex.MustCompile("(?i)^(doc(s)?|misc)/.+")

	if buildPatternStr := actions.GetInput("build_pattern"); buildPatternStr != "" {
		compiled, err := regex.Compile(buildPatternStr)
		if err != nil {
			return Params{}, fmt.Errorf("invalid build pattern value: %s", buildPatternStr)
		}

		buildPattern = compiled
	}

	var debug bool

	if debugStr := actions.GetInput("debug"); debugStr != "" {
		parsed, err := strconv.ParseBool(debugStr)
		if err != nil {
			return Params{}, fmt.Errorf("invalid debug argument: %s", debugStr)
		}

		debug = parsed
	}

	var prefix string = "v"

	if prefixStr := actions.GetInput("prefix"); prefixStr != "" {
		prefix = prefixStr
	}

	var baseVersion *semver.Version

	if baseVersionStr := actions.GetInput("base_version"); baseVersionStr != "" {
		prefixRe := regexp.MustCompile(fmt.Sprintf("^%s", prefix))
		baseVersionStr = prefixRe.ReplaceAllLiteralString(baseVersionStr, "")

		parsed, err := semver.Parse(baseVersionStr)
		if err != nil {
			return Params{}, fmt.Errorf("invalid base_version format: %s", baseVersionStr)
		}

		baseVersion = &parsed
	}

	var mainBranchName string = "master"

	if mainBranchNameStr := actions.GetInput("main_branch_name"); mainBranchNameStr != "" {
		mainBranchName = mainBranchNameStr
	}

	var developBranchName string = "develop"

	if developBranchNameStr := actions.GetInput("develop_branch_name"); developBranchNameStr != "" {
		developBranchName = developBranchNameStr
	}

	var prereleaseID string = "pre"

	if prereleaseIDStr := actions.GetInput("prerelease_id"); prereleaseIDStr != "" {
		prereleaseID = prereleaseIDStr
	}

	return Params{
		CommitSha:         commitSha,
		RepoDir:           repoDir,
		Bump:              bump,
		BranchingModel:    branchingModel,
		BaseVersion:       baseVersion,
		Prefix:            prefix,
		PrereleaseID:      prereleaseID,
		MainBranchName:    mainBranchName,
		DevelopBranchName: developBranchName,
		PatchPattern:      patchPattern,
		MinorPattern:      minorPattern,
		MajorPattern:      majorPattern,
		BuildPattern:      buildPattern,
		Debug:             debug,
	}, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func (p Params) String() string {
	var baseVersion string
	if p.BaseVersion != nil {
		baseVersion = p.BaseVersion.String()
	}

	return fmt.Sprintf(
		"commit sha: %q, bump: %q, base version: %q, prefix: %q,"+
			" prerelease id: %q, main branch name: %q, develop branch name: %q,"+
			" patch pattern: %q, minor pattern: %q, major pattern: %q, build pattern: %q,"+
			" repo dir: %q, debug: %t",
		p.CommitSha,
		p.Bump,
		baseVersion,
		p.Prefix,
		p.PrereleaseID,
		p.MainBranchName,
		p.DevelopBranchName,
		p.PatchPattern.String(),
		p.MinorPattern.String(),
		p.MajorPattern.String(),
		p.BuildPattern.String(),
		p.RepoDir,
		p.Debug,
	)
}
