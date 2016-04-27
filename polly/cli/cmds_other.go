package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/emccode/polly/util"
)

func (c *CLI) initOtherCmdsAndFlags() {
	c.initOtherCmds()
	c.initOtherFlags()
}

func (c *CLI) initOtherCmds() {
	c.versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			util.PrintVersion(os.Stdout)
		},
	}
	c.c.AddCommand(c.versionCmd)

	c.envCmd = &cobra.Command{
		Use:   "env",
		Short: "Print the Polly environment",
		Run: func(cmd *cobra.Command, args []string) {
			evs := c.p.Config.EnvVars()
			for _, ev := range evs {
				fmt.Println(ev)
			}
		},
	}
	c.c.AddCommand(c.envCmd)

	c.installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install Polly",
		Run: func(cmd *cobra.Command, args []string) {
			install()
		},
	}
	c.c.AddCommand(c.installCmd)

	c.uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Polly",
		Run: func(cmd *cobra.Command, args []string) {
			pkgManager, _ := cmd.Flags().GetBool("package")
			uninstall(pkgManager)
		},
	}
	c.c.AddCommand(c.uninstallCmd)
}

func (c *CLI) initOtherFlags() {
	cobra.HelpFlagShorthand = "?"
	cobra.HelpFlagUsageFormatString = "Help for %s"

	c.c.PersistentFlags().StringVarP(&c.cfgFile, "config", "c", "",
		"The path to a custom Polly configuration file")
	c.c.PersistentFlags().BoolP(
		"verbose", "v", false, "Print verbose help information")

	// add the flag sets
	for _, fs := range c.p.Config.FlagSets() {
		c.c.PersistentFlags().AddFlagSet(fs)
	}

	c.uninstallCmd.Flags().Bool("package", false,
		"A flag indicating a package manager is performing the uninstallation")
}
