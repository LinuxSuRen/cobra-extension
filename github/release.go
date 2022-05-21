package github

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v29/github"
)

// ReleaseClient is the client of GitHub
type ReleaseClient struct {
	Client *github.Client
	Org    string
	Repo   string
}

// ReleaseAsset is the asset from GitHub release
type ReleaseAsset struct {
	TagName string
	Body    string
}

// Release represents a GitHub release
type Release struct {
	TagName string
	ID      int64
}

// Tag represents a tag of a git repository
type Tag struct {
	Name string
}

// Init inits the GitHub client
func (g *ReleaseClient) Init() {
	token := os.Getenv("GITHUB_TOKEN")
	var tc *http.Client
	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	g.Client = github.NewClient(tc)
}

// GetLatestReleaseAsset returns the latest release asset
func (g *ReleaseClient) GetLatestReleaseAsset(owner, repo string) (ra *ReleaseAsset, err error) {
	ctx := context.Background()

	var release *github.RepositoryRelease
	if release, _, err = g.Client.Repositories.GetLatestRelease(ctx, owner, repo); err == nil {
		ra = &ReleaseAsset{
			TagName: release.GetTagName(),
			Body:    release.GetBody(),
		}
	}
	return
}

// GetReleaseList returns a list of release
func (g *ReleaseClient) GetReleaseList(owner, repo string, count int) (list []Release, err error) {
	ctx := context.Background()

	var releaseList []*github.RepositoryRelease
	if releaseList, _, err = g.Client.Repositories.ListReleases(ctx, owner, repo, &github.ListOptions{PerPage: count}); err == nil {
		for i := range releaseList {
			release := releaseList[i]
			list = append(list, Release{
				TagName: *release.Name,
				ID:      *release.ID,
			})
		}
	}
	return
}

// GetTagList returns a list of tag
func (g *ReleaseClient) GetTagList(owner, repo string, count int) (list []Tag, err error) {
	ctx := context.Background()

	var tagList []*github.RepositoryTag
	if tagList, _, err = g.Client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{PerPage: count}); err == nil {
		for i := range tagList {
			tag := tagList[i]
			list = append(list, Tag{
				Name: *tag.Name,
			})
		}
	}
	return
}

// GetJCLIAsset returns the asset from a tag name
func (g *ReleaseClient) GetJCLIAsset(tagName string) (*ReleaseAsset, error) {
	return g.GetReleaseAssetByTagName(g.Org, g.Repo, tagName)
}

// GetReleaseAssetByTagName returns the release asset by tag name
func (g *ReleaseClient) GetReleaseAssetByTagName(owner, repo, tagName string) (ra *ReleaseAsset, err error) {
	ctx := context.Background()

	opt := &github.ListOptions{
		PerPage: 99999,
	}

	var releaseList []*github.RepositoryRelease
	if releaseList, _, err = g.Client.Repositories.ListReleases(ctx, owner, repo, opt); err == nil {
		for _, item := range releaseList {
			if item.GetTagName() == tagName {
				ra = &ReleaseAsset{
					TagName: item.GetTagName(),
					Body:    item.GetBody(),
				}
				break
			}
		}
	}
	return
}
