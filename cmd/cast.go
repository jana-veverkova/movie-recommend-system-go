package cmd

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/moviescast"
	"github.com/spf13/cobra"
)

var downloadCastCmd = &cobra.Command{
	Use:   "download-cast",
	Short: "Downloads casts for movies.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := moviescast.DownloadData("data/processed/edx.csv", "data/processed/cast.csv")
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		return nil
	},
}
