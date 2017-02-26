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

	"strings"

	"golang.org/x/net/context"

	"github.com/fatih/color"
	pb "github.com/jwhitcraft/rome/aqueduct"
	"github.com/jwhitcraft/rome/build"
	"github.com/rjeczalik/notify"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:     "watch",
	Example: "rome watch -v 7.9.0.0 -f ent -d /tmp/sugar /path/to/mango/git/checkout",
	Args:    validSourceArg,
	Short:   "Watch the file system for changes and built any files that change",
	Long:    `Watch for file changes, and then build them as they happen.`,
	PreRunE: buildPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		// Make the channel buffered to ensure no event is dropped. Notify will drop
		// an event if the receiver is not able to keep up the sending pace.
		c := make(chan notify.EventInfo, 1)

		// Set up a watchpoint listening for events within a directory tree rooted
		// at current working directory. Dispatch remove events to c.
		if err := notify.Watch(source+"/...", c, notify.Create, notify.Write, notify.Rename); err != nil {
			log.Fatal(err)
		}
		defer notify.Stop(c)

		fmt.Printf("%v %v %v\n",
			color.GreenString("Starting Build Watcher, press"),
			color.RedString("ctrl+c"),
			color.GreenString("to exit"))

		// keep the looping open
		for {
			// todo figure out symlinks in here
			file := <-c
			// silly jetbrains and how it saves files
			if !strings.Contains(file.Path(), "___jb_") &&
				!isExcluded(strings.Replace(file.Path(), source, "", -1), flavor) {

				switch file.Event() {
				case notify.Create:
					fallthrough
				case notify.Rename:
					fallthrough
				case notify.Write:
					fileChanged(build.CreateFile(file.Path(), convertToTargetPath(file.Path())))
				default:
					log.Printf("%s is not handled yet, moving along", file.Event().String())
				}
			}

		}
	},
}

func fileChanged(file iFile) {
	// this is a bit complex, but it works, should look at cleaning it up
	if cleanCache {
		if conduit != nil {
			conduit.CleanCache(context.Background(), &pb.CleanCacheRequest{})
		} else {
			build.CleanCache(destination, cleanCacheItems)
		}
	}
	if conduit != nil {
		file.SendToAqueduct(conduit)
	} else {
		file.Process(flavor, version)
	}
	log.Printf("%v %s",
		color.GreenString("[Built]"),
		file.GetTarget())
}

func init() {
	RootCmd.AddCommand(watchCmd)

	addBuildCommands(watchCmd)

	watchCmd.MarkFlagRequired("version")
	watchCmd.MarkFlagRequired("flavor")
	watchCmd.MarkFlagRequired("destination")
}
