package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	dataDir string
	rootCmd = &cobra.Command{
		Use:   "lumora",
		Short: "kvdb to store yo mom's fat!",
		Long:  `A key-value store CLI that allows you to interact with your mom's stored fat!.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use != "init" && dataDir == "" {
				return fmt.Errorf("data directory not specified. Please use --data-dir or LUMORA_DATA_DIR environment variable")
			}
			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", "", "Data directory for the database")
	rootCmd.MarkPersistentFlagRequired("data-dir")
}
