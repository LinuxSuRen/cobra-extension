package version

import (
	"fmt"
	"github.com/google/go-github/v29/github"
	"net/http"
)

// PrintOption is the version option
type PrintOption struct {
	Changelog  bool
	ShowLatest bool

	Org  string
	Repo string
}

// CustomDownloadFunc is the function interface for custom download URL
type CustomDownloadFunc func(string) string

// SelfUpgradeOption is the option for self upgrade command
type SelfUpgradeOption struct {
	ShowProgress       bool
	Privilege          bool
	Org                string
	Repo               string
	Name               string
	CustomDownloadFunc CustomDownloadFunc
	PathSeparate       string
	Thread             int

	GitHubClient *github.Client
	RoundTripper http.RoundTripper
}

var (
	version string
	commit  string
	date    string
)

// GetVersion returns the version
func GetVersion() string {
	return version
}

// SetVersion is only for the test purpose
func SetVersion(ver string) {
	version = ver
}

// GetCommit returns the commit id
func GetCommit() string {
	return commit
}

// GetDate returns the build date time
func GetDate() string {
	return date
}

// GetCombinedVersion returns the version and commit id
func GetCombinedVersion() string {
	return fmt.Sprintf("jcli; %s; %s", GetVersion(), GetCommit())
}
