/*
Copyright © 2023 Pax Newman pax.newman@gmail.com
*/
package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func validateDimension(d string) bool {
	// FIXME this will currently validate sizes too large (e.g. 900% or more px than window size?)
	m, _ := regexp.MatchString("^((([0-9]{1,3}%))|([0-9]+(|(px))))$", d)
	return m
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "img",
	Short: "display an inline image",
	Long:  `img provides the ability to display an inline image in the iterm2 terminal`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MatchAll(
		cobra.ExactArgs(1),
		func(cmd *cobra.Command, args []string) error {
			// Check if the file exists
			_, err := os.Stat(args[0])
			return err
		},
	),
	Run: func(cmd *cobra.Command, args []string) {
		fname := args[0]
		content, _ := os.ReadFile(fname)

		encoding := base64.StdEncoding.EncodeToString(content)

		sequence := ""
		sequence += "\u001B]1337;File=inline=1"

		flags := cmd.Flags()

		if w, _ := flags.GetString("width"); w != "" && validateDimension(w) {
			sequence += fmt.Sprintf(";width=%s", w)
		}

		if h, _ := flags.GetString("height"); h != "" && validateDimension(h) {
			sequence += fmt.Sprintf(";height=%s", h)
		}

		if o, _ := flags.GetInt("offset"); o >= 0 {
			sequence = strings.Repeat(" ", o) + sequence
		}

		fmt.Printf("%s;preserveAspectRatio=0:%s\a\n", sequence, encoding)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("width", "auto", "Sets the width of the image (N: N Char Cells, Npx: N pixels, N%: N% of the sesions width")
	rootCmd.Flags().String("height", "auto", "Sets the height of the image (N: N Char Cells, Npx: N pixels, N%: N% of the sesions height")
	rootCmd.Flags().Int("offset", 0, "Adds a left offset to the displayed image")
}
