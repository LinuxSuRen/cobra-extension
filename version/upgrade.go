package version

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/linuxsuren/cobra-extension/common"
	gh "github.com/linuxsuren/cobra-extension/github"
	httpdownloader "github.com/linuxsuren/http-downloader/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

// NewSelfUpgradeCmd create a command for self upgrade
func NewSelfUpgradeCmd(org, repo, name string, customDownloadFunc CustomDownloadFunc) (cmd *cobra.Command) {
	opt := &SelfUpgradeOption{
		Org:                org,
		Repo:               repo,
		Name:               name,
		CustomDownloadFunc: customDownloadFunc,
	}

	cmd = &cobra.Command{
		Use:   "upgrade",
		Short: fmt.Sprintf("Upgrade %s itself", name),
		Long: fmt.Sprintf(`Upgrade %s itself
You can use any exists version to upgrade %s itself. If there's no argument given, it will upgrade to the latest release.
You can upgrade to the latest developing version, please use it like: %s version upgrade dev'`, name, name, name),
		RunE: opt.RunE,
	}
	opt.addFlags(cmd.Flags())
	return
}

func (o *SelfUpgradeOption) addFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&o.ShowProgress, "show-progress", "", true,
		fmt.Sprintf("If you want to show the progress of download %s", o.Name))
	flags.BoolVarP(&o.Privilege, "privilege", "", true,
		fmt.Sprintf("Try to take the privilege from system if there's no write permission on %s", o.Name))
	flags.IntVarP(&o.Thread, "thread", "t", 0,
		"Download the target binary file in multi-thread mode. It only works when its value is bigger than 1")
}

// RunE is the main point of current command
func (o *SelfUpgradeOption) RunE(cmd *cobra.Command, args []string) (err error) {
	var version string
	if len(args) > 0 {
		version = args[0]
	}

	// copy binary file into system path
	var targetPath string
	if targetPath, err = exec.LookPath(o.Name); err != nil {
		err = fmt.Errorf("cannot find %s from system path, error: %v", o.Name, err)
		return
	}

	var f *os.File
	if f, err = os.OpenFile(targetPath, os.O_WRONLY, 0666); os.IsPermission(err) {
		if !o.Privilege {
			return
		}

		var sudo string
		if sudo, err = exec.LookPath("sudo"); err == nil {
			sudoArgs := []string{"sudo", o.Name, "version", "upgrade", "--privilege=false"}
			sudoArgs = append(sudoArgs, args...)

			env := os.Environ()
			err = syscall.Exec(sudo, sudoArgs, env)
		}
		return
	}
	defer func() {
		_ = f.Close()
	}()

	currentVersion := GetVersion()
	err = o.Download(cmd, version, currentVersion, targetPath)
	return
}

// Download downloads the binary file from GitHub release
// Org, Repo, Name is necessary
func (o *SelfUpgradeOption) Download(log common.Printer, version, currentVersion, targetPath string) (err error) {
	// try to understand the version from user input
	switch version {
	case "dev":
		version = "master"
	case "":
		o.GitHubClient = github.NewClient(nil)
		ghClient := &gh.GitHubReleaseClient{
			Client: o.GitHubClient,
			Org:    o.Org,
			Repo:   o.Repo,
		}
		if asset, assetErr := ghClient.GetLatestJCLIAsset(); assetErr == nil && asset != nil {
			version = asset.TagName
		} else {
			err = fmt.Errorf("cannot get the latest version, error: %s", assetErr)
			return
		}
	}

	// version review
	if currentVersion == version {
		log.Printf("no need to upgrade %s\n", o.Name)
		return
	}
	log.Println(fmt.Sprintf("prepare to upgrade to %s", version))

	// download the tar file of target file
	tmpDir := os.TempDir()
	output := fmt.Sprintf("%s/%s.tar.gz", tmpDir, o.Name)

	if o.PathSeparate == "" {
		o.PathSeparate = "-"
	}

	var fileURL string
	if o.CustomDownloadFunc == nil {
		fileURL = fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s%s%s%s%s.tar.gz",
			o.Org, o.Repo, version, o.Name, o.PathSeparate, runtime.GOOS, o.PathSeparate, runtime.GOARCH)
	} else {
		fileURL = o.CustomDownloadFunc(version)

		// make sure we count the download action
		go func() {
			o.downloadCount(version, runtime.GOOS)
		}()
	}
	log.Println("start to download from", fileURL)

	defer func() {
		_ = os.RemoveAll(output)
	}()

	if o.Thread > 1 {
		if err = httpdownloader.DownloadFileWithMultipleThread(fileURL, output, o.Thread, o.ShowProgress); err != nil {
			err = fmt.Errorf("cannot download %s from %s, error: %v", o.Name, fileURL, err)
			return
		}
	} else {
		// keep this exists, it can avoid error due to the new feature
		downloader := httpdownloader.HTTPDownloader{
			RoundTripper:   o.RoundTripper,
			TargetFilePath: output,
			URL:            fileURL,
			ShowProgress:   o.ShowProgress,
		}
		if err = downloader.DownloadFile(); err != nil {
			err = fmt.Errorf("cannot download %s from %s, error: %v", o.Name, fileURL, err)
			return
		}
	}

	if err = o.extractFiles(output); err == nil {
		err = o.overWriteBinary(fmt.Sprintf("%s/%s", filepath.Dir(output), o.Name), targetPath)
	} else {
		err = fmt.Errorf("cannot extract %s from tar file, error: %v", o.Name, err)
	}
	return
}

func (o *SelfUpgradeOption) overWriteBinary(sourceFile, targetPath string) (err error) {
	switch runtime.GOOS {
	case "linux":
		var cp string
		if cp, err = exec.LookPath("cp"); err == nil {
			err = syscall.Exec(cp, []string{"cp", sourceFile, targetPath}, os.Environ())
		}
	default:
		sourceF, _ := os.Open(sourceFile)
		targetF, _ := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0664)
		if _, err = io.Copy(targetF, sourceF); err != nil {
			err = fmt.Errorf("cannot copy %s from %s to %v, error: %v", o.Name, sourceFile, targetPath, err)
		}
	}
	return
}

func (o *SelfUpgradeOption) downloadCount(version string, arch string) {
	countURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s-%s-%s.tar.gz",
		o.Org, o.Repo, version, o.Name, runtime.GOOS, runtime.GOARCH)

	if tempDir, err := ioutil.TempDir(".", "download-count"); err == nil {
		tempFile := tempDir + fmt.Sprintf("/%s.tar.gz", o.Name)
		defer func() {
			_ = os.RemoveAll(tempDir)
		}()

		downloader := httpdownloader.HTTPDownloader{
			RoundTripper:   o.RoundTripper,
			TargetFilePath: tempFile,
			URL:            countURL,
		}
		// we don't care about the result, just for counting
		_ = downloader.DownloadFile()
	}
}

func (o *SelfUpgradeOption) extractFiles(tarFile string) (err error) {
	var f *os.File
	var gzf *gzip.Reader
	if f, err = os.Open(tarFile); err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	if gzf, err = gzip.NewReader(f); err != nil {
		return
	}

	tarReader := tar.NewReader(gzf)
	var header *tar.Header
	for {
		if header, err = tarReader.Next(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if name != o.Name {
				continue
			}
			var targetFile *os.File
			if targetFile, err = os.OpenFile(fmt.Sprintf("%s/%s", filepath.Dir(tarFile), name),
				os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)); err != nil {
				break
			}
			if _, err = io.Copy(targetFile, tarReader); err != nil {
				break
			}
			_ = targetFile.Close()
		}
	}
	return
}
