package strategy

import "github.com/gandarez/semver-action/pkg/git"

// TrunkBased implements the trunk-based strategy.
type TrunkBased struct{}

// DetermineBumpStrategy determines the strategy for semver to bump product version.
func (t *TrunkBased) DetermineBumpStrategy(sourceBranch, destBranch string) (string, string) {
	panic("not implemented")
}

// Tag implements the Strategy interface.
func (t *TrunkBased) Tag(config Configuration, gc git.Client) (Result, error) {
	panic("not implemented")
}
