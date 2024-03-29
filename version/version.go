package version

import (
	"fmt"
	"github.com/google/go-github/v29/github"
	gh "github.com/linuxsuren/cobra-extension/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

// NewVersionCmd create a command for version
func NewVersionCmd(org, repo, name string, customDownloadFunc CustomDownloadFunc) (cmd *cobra.Command) {
	opt := &PrintOption{
		Org:  org,
		Repo: repo,
	}

	cmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   fmt.Sprintf("Print the version of %s", name),
		Long:    fmt.Sprintf("Print the version of %s", name),
		RunE:    opt.RunE,
	}

	flags := cmd.Flags()
	opt.addFlags(flags)

	cmd.AddCommand(NewSelfUpgradeCmd(org, repo, name, customDownloadFunc))
	return
}

func (o *PrintOption) addFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&o.Changelog, "changelog", "c", false,
		"Output the changelog")
	flags.BoolVarP(&o.ShowLatest, "show-latest", "s", false,
		"Output the latest version")
}

// RunE is the main point of current command
func (o *PrintOption) RunE(cmd *cobra.Command, _ []string) (err error) {
	version := GetVersion()
	cmd.Printf("Version: %s\n", version)
	cmd.Printf("Last Commit: %s\n", GetCommit())
	cmd.Printf("Build Date: %s\n", GetDate())
	cmd.Printf("https://github.com/%s/%s", o.Org, o.Repo)

	if strings.HasPrefix(version, "dev-") {
		version = strings.ReplaceAll(version, "dev-", "")
	}

	ghClient := &gh.ReleaseClient{
		Client: github.NewClient(nil),
		Org:    o.Org,
		Repo:   o.Repo,
	}
	var asset *gh.ReleaseAsset
	if o.Changelog && !o.ShowLatest {
		// only print the changelog of current version
		if asset, err = ghClient.GetJCLIAsset(version); err == nil && asset != nil {
			cmd.Println("Changelog:")
			cmd.Println(asset.Body)
		}
	}

	if o.ShowLatest {
		if asset, err = ghClient.GetLatestReleaseAsset(o.Org, o.Repo); err == nil && asset != nil {
			cmd.Println("The latest version", asset.TagName)
			if o.Changelog {
				cmd.Println("Changelog:")
				cmd.Println(asset.Body)
			}
		}
	}
	return
}
