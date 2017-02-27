// Copyright © 2017 Jon Whitcraft
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"os"

	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// self-updateCmd represents the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update Rome if a new version exists",
	Long:  `This will allow Rome to update it's self like copmoser or other new fangled tools do`,
	Run: func(cmd *cobra.Command, args []string) {
		var updater = &selfupdate.Updater{
			CurrentVersion: Version,
			ApiURL:         "http://h2ik.co/",
			BinURL:         "http://h2ik.co/",
			DiffURL:        "",
			Dir:            ".rome_update_check/",
			ForceCheck:     true,
			CmdName:        "rome", // app name
		}

		fmt.Println("Running Self Update")
		err := updater.BackgroundRun()
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		fmt.Println("Self Update Finished")
	},
}

func init() {
	RootCmd.AddCommand(selfUpdateCmd)
}
