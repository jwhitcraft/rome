package build


import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"regexp"
	"io/ioutil"
	"path"

	"github.com/jwhitcraft/rome/utils"
)

var (
	ProcessibleExtensions = []string{
		".php", ".json", ".js",
	}
	Flavors = map[string][]string{
		"pro": {"pro"},
		"corp": {"pro", "corp"},
		"ent": {"pro", "ent"},
		"ult": {"pro", "ent", "ult"},
	}

	License = map[string][]string {
		"lic": {"sub"},
	}

	TagRegex = regexp.MustCompile("//[[:space:]]*(BEGIN|END|FILE|ELSE)[[:space:]]*SUGARCRM[[:space:]]*(.*) ONLY")

	VarRegex = regexp.MustCompile( "@_SUGAR_(FLAV|VERSION)")
)

func ProcessBuildTag(tag string, flavors []string ) bool {
	// first things first, check for &&
	tags := strings.Split(tag, "&&")
	ok := true
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		tagSep := getTagSperator(tag)
		tagKey := getTagKey(tag, tagSep)
		tagVal := getTagValue(tag, tagSep)
		var testValue []string
		var tagOk bool
		// default the tag to be allowed, only change it something else is off
		switch tagKey {
		case "flav":
			testValue = flavors
		case "lic":
			testValue = License[tagKey]
		}
		if tagSep == "!=" {
			tagOk = notContains(testValue, tagVal)
		} else {
			tagOk = contains(testValue, tagVal)
		}
		ok = ok && tagOk
	}

	return ok
}

func getTagSperator(tag string) string {
	if strings.Contains(tag, "!=") {
		return "!="
	}

	return "="
}

func BuildFile(srcPath string, destPath string, buildFlavor string, buildVersion string) bool {
	var useLine bool = true
	var shouldProcess bool = false
	var canProcess bool = false

	var skippedLines utils.Counter

	// lets make sure the that folder exists
	var destFolder string = path.Dir(destPath)
	var fileExt string = path.Ext(destPath)
	// var fileName string = path.Base(destPath)
	os.MkdirAll(destFolder, 0775)

	// regardless, if the file is in the node_modules folder
	// don't try and process it
	if !strings.Contains(destFolder, "node_modules") {
		canProcess = contains(ProcessibleExtensions, fileExt)
	}

	// first load the whole file to check for the build tags
	fileBytes, err := ioutil.ReadFile(srcPath)
	fileString := string(fileBytes)
	if canProcess && TagRegex.MatchString(fileString) {
		shouldProcess = true
		// check to see if it's a type of FILE
		matches := TagRegex.FindStringSubmatch(fileString)
		if matches[1] == "FILE" {
			tagOk := ProcessBuildTag(matches[2], Flavors[buildFlavor])
			//fmt.Printf("// File Tag Found for flavor: %s and building %s, should build file: %t\n", tagFlav, buildFlavor, tagOk)
			if tagOk == false {
				return false
			}
		}
	}

	// do the variable replacement
	if canProcess && VarRegex.MatchString(fileString) {
		fileString = strings.Replace(fileString, "@_SUGAR_VERSION", buildVersion, -1)
		fileString = strings.Replace(fileString, "@_SUGAR_FLAV", buildFlavor, -1)
	}


	if err != nil {
		fmt.Printf("pre-preocess error: %v\n",err)
		return false
	}

	fw, err := os.Create(destPath)
	defer fw.Close()

	if shouldProcess {
		f := strings.NewReader(fileString)
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
					//fmt.Printf("// Begin Tag Found for flavor: %s and building %s, should use lines: %t\n", tagFlav, buildFlavor, tagOk)
					useLine = ProcessBuildTag(matches[2], Flavors[buildFlavor])
					if !useLine {
						skippedLines.Increment()
					}
				case "END":
					//fmt.Printf("// Skipped %d lines\n", skippedLines.get())
					skippedLines.Reset()
					useLine = true
				}
			} else if useLine {
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

func getTagValue(eval string, splitOn string) string {
	splitFlav := strings.Split(eval, splitOn)
	if len(splitFlav) == 1 {
		return strings.ToLower(splitFlav[0])
	}

	return strings.ToLower(splitFlav[1])
}

func getTagKey(eval string, splitOn string) string {
	splitFlav := strings.Split(eval, splitOn)

	if len(splitFlav) == 1 {
		return "flav"
	}

	return strings.ToLower(splitFlav[0])
}

func notContains(slice []string, item string) bool {
	return !contains(slice, item)
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}