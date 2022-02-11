package main

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(waas completion bash)

# To load completions for each session, execute once:
Linux:
  $ waas completion bash > /etc/bash_completion.d/waas
MacOS:
  $ waas completion bash > /usr/local/etc/bash_completion.d/waas

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ waas completion zsh > "${fpath[1]}/_waas"

# You will need to start a new shell for this setup to take effect.

Fish:

$ waas completion fish | source

# To load completions for each session, execute once:
$ waas completion fish > ~/.config/fish/completions/waas.fish

Powershell:

PS> waas completion powershell | Out-String | Invoke-Expression

# To load completions for every new session, run:
PS> waas completion powershell > waas.ps1
# and source this file from your powershell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	})
}
