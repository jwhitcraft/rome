package main

import (
	"path/filepath"
	"os"
	"strings"
	"sync"
	"fmt"
	"time"
	"path"
)

type File string
type Link struct {
	Link string
	Target string
}

func fileWorker(files <-chan File, quit <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case file, ok := <-files:
			if !ok {
				return
			}
			shortPath := strings.Replace(string(file), "/Users/jwhitcraft/Projects/Mango/sugarcrm", "", -1)
			destination := "/Users/jwhitcraft/test_build" + shortPath
			BuildFile(string(file), destination, "ent")
		case <-quit:
			return
		}
	}
}

func linkWorker(links <- chan Link, quit <- chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case link, ok := <-links:
			if !ok {
				return
			}
			shortPath := strings.Replace(string(link.Link), "/Users/jwhitcraft/Projects/Mango/sugarcrm", "", -1)
			destination := "/Users/jwhitcraft/test_build" + shortPath
			os.MkdirAll(path.Dir(destination), 0775)
			os.Symlink(link.Target, destination)
		case <-quit:
			return
		}
	}
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf(" in %.3f seconds\n", elapsed.Seconds())
}

func main() {
	fmt.Println("Starting Rome...")
	defer timeTrack(time.Now())
	var builtFiles counter
	files := make(chan File, 4096)
	links := make(chan Link, 2048)
	quit := make(chan bool)
	var wg sync.WaitGroup
	var linkWg sync.WaitGroup

	// spawn 5 workers
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go fileWorker(files, quit, &wg)
	}

	for i := 0; i < 5; i++ {
		linkWg.Add(1)
		go linkWorker(links, quit, &linkWg)
	}

	filepath.Walk("/Users/jwhitcraft/Projects/Mango/sugarcrm", func(path string, f os.FileInfo, err error) error {
		// ignore the node_modules dir in the root, but lead sidecar
		if f.Name() == "node_modules" && strings.Contains(path, "sugarcrm/node_modules") {
			return filepath.SkipDir
		}
		if !f.IsDir() {
			builtFiles.increment()
			// handle symlinks differently than normal files
			if f.Mode()&os.ModeSymlink != 0 {
				originFile, _ := os.Readlink(path)
				links <- Link{Link: path, Target: originFile}
			} else {
				files <- File(path)
			}
		}
		return nil
	})

	// end of tasks. the workers should quit afterwards
	close(files)
	close(links)
	// use "close(quit)", if you do not want to wait for the remaining tasks

	// wait for all workers to shut down properly
	wg.Wait()
	linkWg.Wait()

	fmt.Printf("Built %d files", builtFiles.get())
}