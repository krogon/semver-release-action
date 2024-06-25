package git

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/K-Phoen/semver-release-action/internal/pkg/action"
	"github.com/blang/semver/v4"
	"github.com/google/go-github/v45/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func LatestTagCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "latest-tag [REPOSITORY] [GH_TOKEN] [TAG_FORMAT]",
		Args: cobra.ExactArgs(3),
		Run:  executeLatestTag,
	}
}

func createRegexFromTagFormat(tagFormat string) string {
	tagFormatRegex := strings.ReplaceAll(tagFormat, "%major%", "(\\d+)")
	tagFormatRegex = strings.ReplaceAll(tagFormatRegex, "%minor%", "(\\d+)")
	tagFormatRegex = strings.ReplaceAll(tagFormatRegex, "%patch%", "(\\d+)")
	tagFormatRegex = "^" + tagFormatRegex

	return tagFormatRegex
}

func executeLatestTag(cmd *cobra.Command, args []string) {
	repository := args[0]
	githubToken := args[1]
	tagFormat := args[2]

	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, tokenSource))

	parts := strings.Split(repository, "/")
	owner := parts[0]
	repo := parts[1]

	refs, response, err := client.Git.ListMatchingRefs(ctx, owner, repo, &github.ReferenceListOptions{
		Ref: "tags",
	})
	if response != nil && response.StatusCode == http.StatusNotFound {
		cmd.Print("v0.0.0")

		return
	}
	action.AssertNoError(cmd, err, "could not list git refs: %s", err)

	latest := filterRemoteTags(refs, tagFormat)
	cmd.Printf("v%s", latest)
}

func filterRemoteTags(refs []*github.Reference, tagFormat string) semver.Version {
	latest := semver.MustParse("0.0.0")
	tagFormatRegex := createRegexFromTagFormat(tagFormat)

	for _, ref := range refs {
		versionStr := strings.Replace(*ref.Ref, "refs/tags/", "", 1)

		formatValid, _ := regexp.MatchString(tagFormatRegex, versionStr)
		if !formatValid {
			continue
		}

		re := regexp.MustCompile(tagFormatRegex)
		versionArr := re.FindStringSubmatch(versionStr)

		major, _ := strconv.ParseUint(versionArr[1], 10, 64)
		minor, _ := strconv.ParseUint(versionArr[2], 10, 64)
		patch, _ := strconv.ParseUint(versionArr[3], 10, 64)

		version := semver.Version{
			Major: major,
			Minor: minor,
			Patch: patch,
		}

		if version.GT(latest) {
			latest = version
		}
	}

	return latest
}
