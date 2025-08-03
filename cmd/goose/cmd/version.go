package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GoOse",
	Long:  `Print the version number of GoOse`,
	Run: func(cmd *cobra.Command, args []string) {
		version := getVersion()
		fmt.Printf("GoOse v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersion() string {
	// Try to read from VERSION file in project root
	versionFile := filepath.Join("..", "..", "..", "VERSION")
	if content, err := ioutil.ReadFile(versionFile); err == nil {
		return strings.TrimSpace(string(content))
	}

	// Fallback version
	return "2.1.25"
}
