package cmd

import (
	"log"

	gp2app "github.com/j4ng5y/gphotos_downloader/pkg/gphotos_downloader"

	"github.com/spf13/cobra"
)

// Run is a function to run the gp2app cli.
func Run() {
	var (
		rootCmd = &cobra.Command{
			Use:   "gphotos-downloader",
			Short: "gphotos-downloader: Google Photos media download tool",
			Args:  cobra.NoArgs,
			Run: func(*cobra.Command, []string) {
				if err := gp2app.Run(); err != nil {
					log.Fatal(err)
				}
			},
		}
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
