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

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch [OPTIONS] SOURCE-PATH",
	Short: "Watch for FS Changes and Built Out the files",
	Long: `Currently Not Implmented, The plan is to create a utility like build monitor but inside of rome`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("Nothing to see hear yet")
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&destination,"destination", "d", "", "Where should the built files be put")
	watchCmd.Flags().StringVarP(&version, "version", "v", "","What Version is being built")
	watchCmd.Flags().StringVarP(&flavor, "flavor", "f", "ent","What Flavor of SugarCRM to build")
	watchCmd.Flags().BoolVar(&clean, "clean", false, "Remove Existing Build Before Building")

	watchCmd.MarkFlagRequired("version")
	watchCmd.MarkFlagRequired("flavor")
	watchCmd.MarkFlagRequired("destination")
}
