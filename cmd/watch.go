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
	"os"
	"path"
	"strings"
	"time"

	"github.com/cortesi/moddwatch"
	"github.com/fatih/color"
	pb "github.com/jwhitcraft/rome/aqueduct"
	"github.com/jwhitcraft/rome/build"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
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

		fmt.Printf("%v %v %v\n",
			color.GreenString("Starting Build Watcher, press"),
			color.RedString("ctrl+c"),
			color.GreenString("to exit"))

		c := make(chan *moddwatch.Mod, 1)

		paths := []string{
			source + "/...",
		}

		watch, err := moddwatch.Watch(paths, 300*time.Millisecond, c)
		defer watch.Stop()
		if err != nil {
			os.Exit(500)
		}

		for {
			files := <-c
			for a := 0; a < len(files.Added); a++ {
				file := files.Added[a]
				if isValidFile(file) {
					fileChanged(build.CreateFile(file, convertToTargetPath(file)), false)
				}
			}
			for c := 0; c < len(files.Changed); c++ {
				file := files.Changed[c]
				if isValidFile(file) {
					fileChanged(build.CreateFile(file, convertToTargetPath(file)), false)
				}
			}
			for d := 0; d < len(files.Deleted); d++ {
				file := files.Deleted[d]
				if isValidFile(file) {
					fileChanged(build.CreateFile(file, convertToTargetPath(file)), true)
				}
			}
		}
	},
}

func isValidFile(file string) bool {
	var fileExt string = path.Ext(file)
	return !strings.Contains(file, "___jb_") &&
		fileExt != ".swp" &&
		!isExcluded(strings.Replace(file, source, "", -1), flavor)
}

func fileChanged(file iFile, isRemove bool) {
	tag := color.GreenString("[Built]")
	if isRemove {
		tag = color.RedString("[Removed]")
	}
	if conduit != nil {
		if cleanCache {
			conduit.CleanCache(context.Background(), &pb.CleanCacheRequest{})
		}
		file.SendToAqueduct(conduit)
	} else {
		if cleanCache {
			build.CleanCache(destination, cleanCacheItems)
		}
		file.Process(flavor, version)
	}

	log.Printf("%v %s", tag, file.GetTarget())
}

func init() {
	RootCmd.AddCommand(watchCmd)

	addBuildCommands(watchCmd)

	watchCmd.MarkFlagRequired("version")
	watchCmd.MarkFlagRequired("flavor")
	watchCmd.MarkFlagRequired("destination")
}
