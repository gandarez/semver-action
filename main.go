package main

import (
	"os"
	"strconv"

	"github.com/gandarez/semver-action/cmd/generate"
	"github.com/gandarez/semver-action/pkg/actions"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func main() {
	log.SetHandler(cli.Default)

	result, err := generate.Run()
	if err != nil {
		log.Fatalf("failed to generate semver version: %s\n", err)
	}

	outputFilepath := os.Getenv("GITHUB_OUTPUT")

	// Print previous tag.
	log.Infof("PREVIOUS_TAG: %s", result.PreviousTag)

	if err := actions.SetOutput(outputFilepath, "PREVIOUS_TAG", result.PreviousTag); err != nil {
		log.Fatalf("%s\n", err)
	}

	// Print ancestor tag.
	log.Infof("ANCESTOR_TAG: %s", result.AncestorTag)

	if err := actions.SetOutput(outputFilepath, "ANCESTOR_TAG", result.AncestorTag); err != nil {
		log.Fatalf("%s\n", err)
	}

	// Print calculated semver tag.
	log.Infof("SEMVER_TAG: %s", result.SemverTag)

	if err := actions.SetOutput(outputFilepath, "SEMVER_TAG", result.SemverTag); err != nil {
		log.Fatalf("%s\n", err)
	}

	// Print is prerelease.
	log.Infof("IS_PRERELEASE: %v", result.IsPrerelease)

	if err := actions.SetOutput(outputFilepath, "IS_PRERELEASE", strconv.FormatBool(result.IsPrerelease)); err != nil {
		log.Fatalf("%s\n", err)
	}
}
