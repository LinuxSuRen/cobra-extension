package pkg

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"path"
)

// completionOptions is the option of completion command
type completionOptions struct {
	shellType string
	auto      bool
}

// ShellTypes contains all types of shell
var ShellTypes = []string{
	"zsh", "bash", "powerShell",
}

// NewCompletionCmd creates the completion command
func NewCompletionCmd(rootCmd *cobra.Command) (cmd *cobra.Command) {
	rootName := rootCmd.Name()
	opt := completionOptions{}

	cmd = &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: fmt.Sprintf(`Generate shell completion scripts
Normally you don't need to do more extra work to have this feature if you've installed %s by brew`, rootName),
		Example: fmt.Sprintf(`# Installing bash completion on macOS using homebrew
## If running Bash 3.2 included with macOS
brew install bash-completion
## or, if running Bash 4.1+
brew install bash-completion@2
## If %[1]s is installed via homebrew, this should start working immediately.
## If you've installed via other means, you may need add the completion to your completion directory
%[1]s completion > $(brew --prefix)/etc/bash_completion.d/%[1]s
## If you get trouble, please visit https://github.com/jenkins-zh/jenkins-cli/issues/83

# Installing bash completion on Linux
## If bash-completion is not installed on Linux, please install the 'bash-completion' package
## via your distribution's package manager.
## Load the %[1]s completion code for bash into the current shell
source <(%[1]s completion bash)
## Write bash completion code to a file and source if from .bash_profile
mkdir -p ~/.config/%[1]s/ && %[1]s completion bash > ~/.config/%[1]s/completion.bash.inc
printf "
# %[1]s shell completion
source '$HOME/.config/%[1]s/completion.bash.inc'
" >> $HOME/.bash_profile
source $HOME/.bash_profile

# In order to have good experience on zsh completion, ohmyzsh is a good choice.
# Please install ohmyzsh by the following command
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
# Get more details about onmyzsh from https://github.com/ohmyzsh/ohmyzsh

# Load the %[1]s completion code for zsh[1] into the current shell
source <(%[1]s completion --type zsh)
# Set the %[1]s completion code for zsh[1] to autoload on startup
%[1]s completion --type zsh > "${fpath[1]}/_%[1]s"`, rootName),
		RunE: opt.runE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.shellType, "type", "", "bash",
		fmt.Sprintf("Generate different types of shell which are %v", ShellTypes))
	flags.BoolVarP(&opt.auto, "auto", "", true,
		"Indicate if install the completion script automatically, print the script if it's false")

	err := cmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) (
		i []string, directive cobra.ShellCompDirective) {
		return ShellTypes, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		cmd.PrintErrf("register flag type for sub-command doc failed %#v\n", err)
	}
	return
}

func (o *completionOptions) runE(cmd *cobra.Command, _ []string) (err error) {
	rootCmd := cmd.Root()
	writer := cmd.OutOrStdout()
	buf := bytes.NewBuffer([]byte{})
	if o.auto {
		writer = buf
	}

	shellType := o.shellType
	switch shellType {
	case "zsh":
		err = rootCmd.GenZshCompletion(writer)
	case "powerShell":
		err = rootCmd.GenPowerShellCompletion(writer)
	case "bash":
		err = rootCmd.GenBashCompletion(writer)
	default:
		err = fmt.Errorf("unknown shell type %s", shellType)
	}

	if err == nil && o.auto {
		err = installCompletionScript(rootCmd.Use, buf.Bytes())
	}
	return
}

func installCompletionScript(name string, script []byte) (err error) {
	// create ~/.bash_completion if it does not exist or is empty
	if err = writeIfNotExist("~/.bash_completion", []byte(`
for bcfile in ~/.bash_completion.d/* ; do
  . $bcfile
done`), false); err != nil {
		return
	}

	// install script for the current command to '~/.bash_completion.d/cmd'
	if err = writeIfNotExist(fmt.Sprintf("~/.bash_completion.d/%s", name), script, true); err != nil {
		return
	}
	return
}

func writeIfNotExist(requestPath string, data []byte, force bool) (err error) {
	var targetPath string
	if targetPath, err = homedir.Expand(requestPath); err != nil {
		return
	}

	var ok bool
	if ok, err = pathExists(targetPath); !ok || err != nil || force {
		if err = os.MkdirAll(path.Dir(targetPath), os.FileMode(0755)); err != nil {
			return
		}
		err = os.WriteFile(targetPath, data, 0644)
	}
	return
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
