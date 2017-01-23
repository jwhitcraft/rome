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
	"log"

	"github.com/spf13/cobra"
	"github.com/rjeczalik/notify"
	"os"
	"github.com/jwhitcraft/rome/build"
	"strings"
	"path/filepath"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Example: "rome watch -v 7.9.0.0 -f ent -d /tmp/sugar /path/to/mango/git/checkout",
	Short: "Watch for FS Changes and Built Out the files",
	Long: `Currently Not Implmented, The plan is to create a utility like build monitor but inside of rome`,
	PreRun: func(cmd *cobra.Command, args[]string) {
		// in the preRun, make sure that the source and destination exists
		source = args[0]

		destExists, err := exists(destination)
		if err != nil || !destExists {
			fmt.Printf("Destination Path (%s) does not exists, Creating Now\n", destination)
			os.MkdirAll(destination, 0775)
			// since we had to create the destination dir, set clean to false
			clean = false
		}

		sourceExists, err := exists(source)
		if err != nil || !sourceExists {
			fmt.Printf("\n\nSource Path (%s) does not exists!!\n\n", source)
			os.Exit(401)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		// Make the channel buffered to ensure no event is dropped. Notify will drop
		// an event if the receiver is not able to keep up the sending pace.
		c := make(chan notify.EventInfo, 1)

		// Set up a watchpoint listening for events within a directory tree rooted
		// at current working directory. Dispatch remove events to c.
		if err := notify.Watch(source + "/...", c, notify.All); err != nil {
			log.Fatal(err)
		}
		defer notify.Stop(c)

		fmt.Println("Starting Build Watcher, press ctrl+c to exit")

		// keep the looping open
		for {
			file := <-c
			switch file.Event() {
			case notify.Create:
				fallthrough
			case notify.Write:
				log.Printf("Building %s", file.Path())
				shortPath := strings.Replace(file.Path(), source, "", -1)
				build.BuildFile(file.Path(), filepath.Join(destination, shortPath), flavor, version)
			default:
				log.Printf("%s is not handled yet, moving along", file.Event().String())
			}

		}
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&destination,"destination", "d", "", "Where should the built files be put")
	watchCmd.Flags().StringVarP(&version, "version", "v", "","What Version is being built")
	watchCmd.Flags().StringVarP(&flavor, "flavor", "f", "ent","What Flavor of SugarCRM to build")

	watchCmd.MarkFlagRequired("version")
	watchCmd.MarkFlagRequired("flavor")
	watchCmd.MarkFlagRequired("destination")
}
