package build


import (
	"fmt"
	"os"
	"bytes"
	"bufio"
	"strings"
	"regexp"
	"io/ioutil"
	"path"

	"github.com/jwhitcraft/rome/utils"
)

var (
	ProcessibleExtensions = []string{
		"php", "json", "js",
	}
	Flavors = map[string][]string{
		"pro": {"pro"},
		"corp": {"pro", "corp"},
		"ent": {"pro", "corp", "ent"},
		"ult": {"pro", "corp", "ent", "ult"},
	}

	TagRegex = regexp.MustCompile("//[[:space:]]*(BEGIN|END|FILE|ELSE)[[:space:]]*SUGARCRM[[:space:]]*(.*) ONLY")

	VarRegex = regexp.MustCompile( "@_SUGAR_(FLAV|VERSION)")
)

func BuildFile(srcPath string, destPath string, buildFlavor string, buildVersion string) bool {
	var useLine bool = true
	var shouldProcess bool = false

	var skippedLines utils.Counter

	// lets make sure the that folder exists
	var destFolder string = path.Dir(destPath)
	var fileExt string = path.Ext(destPath)
	// var fileName string = path.Base(destPath)
	os.MkdirAll(destFolder, 0775)

	var canProcess bool = contains(ProcessibleExtensions, fileExt)

	// regardless, if the file is in the node_modules folder
	// don't try and process it
	if strings.Contains(destFolder, "node_modules") {
		canProcess = false
	}

	// first load the whole file to check for the build tags
	fileBytes, err := ioutil.ReadFile(srcPath)
	fileString := string(fileBytes)
	if canProcess && TagRegex.MatchString(fileString) {
		shouldProcess = true
		// check to see if it's a type of FILE
		matches := TagRegex.FindStringSubmatch(fileString)
		if matches[1] == "FILE" {
			tagFlav := getTagFlavor(matches[2])
			tagOk := contains(Flavors[buildFlavor], tagFlav)
			//fmt.Printf("// File Tag Found for flavor: %s and building %s, should build file: %t\n", tagFlav, buildFlavor, tagOk)
			if tagOk == false {
				return false
			}
		}
	}

	// do the variable replacement
	if canProcess && VarRegex.MatchString(fileString) {
		matches := VarRegex.FindStringSubmatch(fileString)
		switch matches[1] {
		case "VERSION":
			fileString = strings.Replace(fileString, "@_SUGAR_VERSION", buildVersion, -1)
		case "FLAV":
			fileString = strings.Replace(fileString, "@_SUGAR_FLAV", buildFlavor, -1)
		}
	}


	if err != nil {
		fmt.Printf("pre-preocess error: %v\n",err)
		return false
	}

	fw, err := os.Create(destPath)
	defer fw.Close()

	if shouldProcess {
		f := bytes.NewReader(fileBytes)
		if err != nil {
			fmt.Printf("error opening file: %v\n",err)
			os.Exit(1)
		}
		writer := bufio.NewWriter(fw)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			val := scanner.Text()

			if TagRegex.MatchString(val) {
				// get the matches
				matches := TagRegex.FindStringSubmatch(val)

				switch matches[1] {
				case "BEGIN":
					tagFlav := getTagFlavor(matches[2])
					tagOk := contains(Flavors[buildFlavor], tagFlav)
					//fmt.Printf("// Begin Tag Found for flavor: %s and building %s, should use lines: %t\n", tagFlav, buildFlavor, tagOk)
					useLine = tagOk
					if tagOk == false {
						skippedLines.Increment()
					}
				case "END":
					//fmt.Printf("// Skipped %d lines\n", skippedLines.get())
					skippedLines.Reset()
					useLine = true
				}
			} else if useLine {
				fmt.Println(val) // Println will add back the final '\n'
				fmt.Fprintln(writer, val)
			} else {
				skippedLines.Increment()
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		} else {
			// write the file to the disk
			writer.Flush()
		}
	} else {
		fw.WriteString(fileString)
	}

	return true
}

func getTagFlavor(eval string) string {
	splitFlav := strings.Split(eval, "=")
	if len(splitFlav) == 1 {
		return strings.ToLower(splitFlav[0])
	}

	return strings.ToLower(splitFlav[1])
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}