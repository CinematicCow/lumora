package cli

import (
	"fmt"
	"os"

	"github.com/CinematicCow/lumora/internal/config"
	"github.com/spf13/cobra"
)

var (
	dataDir string
	rootCmd = &cobra.Command{
		Use:   "lumora",
		Short: "A simple key-value database",
		Long:  `A key-value store CLI that allows you to interact with your data.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "init" || cmd.Use == "config" {
				return nil
			}

			cfg, err := config.InitConfig()
			if err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
			dbName, _ := cmd.Flags().GetString("db-name")
			if dbName == "" {
				dbName = cfg.DefaultDB
			}
			if dbName == "" {
				return fmt.Errorf("no database specified. Please use --db-name or set a default database using 'lumora config set-default'")
			}

			dbPath, exists := cfg.GetDBPath(dbName)
			if !exists {
				return fmt.Errorf("database %s not found. Please initialize it using 'lumora init <db-name>'", dbName)
			}

			dataDir = dbPath
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
}
