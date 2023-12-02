package cmd

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/moviesbudget"
	"github.com/spf13/cobra"
)

var downloadBudgetCmd = &cobra.Command{
	Use:   "download-budget",
	Short: "Downloads budgets for movies.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := moviesbudget.DownloadData("data/original/movie_budget.csv")
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		return nil
	},
}

var processBudgetCmd = &cobra.Command{
	Use:   "process-budget",
	Short: "Processes budgets for movies.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := moviesbudget.ProcessData("data/original/movie_budget.csv", "data/processed/edx.csv")
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		return nil
	},
}
