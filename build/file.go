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
	"path/filepath"
)

var (
	ProcessableExtensions = []string{
		".php", ".json", ".js", ".html", ".tpl", ".css", ".hbs",
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

	TagRegex = regexp.MustCompile("(?i)//[[:space:]]*(BEGIN|END|FILE|ELSE)[[:space:]]*SUGARCRM[[:space:]]*(.*) ONLY")

	IdRegex = regexp.MustCompile(`\$Id(.*)\$`)

	VarRegex = regexp.MustCompile( "@_SUGAR_(FLAV|VERSION)")
)

type File struct {
	SourcePath string
	DestinationPath string
	linkTarget string
}

func (f *File) Process(flavor string, version string) bool {
	if f.linkTarget == "" {
		return processFile(f.SourcePath, f.DestinationPath, flavor, version)
	} else {
		return f.link()
	}
}

func (f *File) link() bool {
	// this will create the symlink
	os.Symlink(f.linkTarget, f.DestinationPath)
	return true
}

func (f *File) SetDestination(source string, destination string) {
	shortPath := strings.Replace(string(f.SourcePath), source, "", -1)
	f.DestinationPath = filepath.Join(destination, shortPath)
	os.MkdirAll(path.Dir(f.DestinationPath), 0775)
}

func (f *File) GetDestination() string {
	return f.DestinationPath
}

func CreateFile(path string) *File {
	return &File{SourcePath: path}
}

func CreateSymLink(path, target string) *File {
	return &File{SourcePath: path, linkTarget: target}
}

func processBuildTag(tag string, flavors []string ) bool {
	// first things first, check for &&
	tags := strings.Split(tag, "&&")
	ok := true
	for _, tag := range tags {

		var tagOk bool
		if strings.Contains(tag, "||") {
			// split on the ||
			var orOk bool
			orTags := strings.Split(tag, "||")
			for _, orTag := range orTags {
				orOk = orOk || getTagBooleanValue(orTag, flavors)
			}
			tagOk = orOk
		} else {
			tagOk = getTagBooleanValue(tag, flavors)
		}

		ok = ok && tagOk
	}

	return ok
}

func getTagBooleanValue(tag string, flavors []string) bool {
	tag = strings.TrimSpace(tag)
	tagSep := getTagSperator(tag)
	tagKey, tagVal := splitTag(tag, tagSep)

	var testValue []string
	// default the tag to be allowed, only change it something else is off
	switch tagKey {
	case "flav":
		testValue = flavors
	case "lic":
		testValue = License[tagKey]
	case "dep":
		testValue = []string{"os"}
	}
	if tagSep == "!=" {
		return notContains(testValue, tagVal)
	}

	return contains(testValue, tagVal)
}

func getTagSperator(tag string) string {
	if strings.Contains(tag, "!=") {
		return "!="
	}

	return "="
}

func processFile(srcPath string, destPath string, buildFlavor string, buildVersion string) bool {
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
		canProcess = contains(ProcessableExtensions, fileExt)
	}

	// first load the whole file to check for the build tags
	fileBytes, err := ioutil.ReadFile(srcPath)
	fileString := string(fileBytes)
	if canProcess && TagRegex.MatchString(fileString) {
		shouldProcess = true
		// check to see if it's a type of FILE
		matches := TagRegex.FindStringSubmatch(fileString)
		if matches[1] == "FILE" {
			tagOk := processBuildTag(matches[2], Flavors[buildFlavor])
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
					useLine = processBuildTag(matches[2], Flavors[buildFlavor])
					if !useLine {
						skippedLines.Increment()
					}
				case "END":
					skippedLines.Reset()
					useLine = true
				}
			} else if IdRegex.MatchString(val) {
				fmt.Fprintln(writer, "")
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

func splitTag(eval, splitOn string) (key, val string) {
	splitVal := strings.Split(eval, splitOn)

	// if the value only has item in it after split, assume it's the value not the key
	if len(splitVal) == 1 {
		key = "flav"
		val = splitVal[0]
	} else {
		key = splitVal[0]
		val = splitVal[1]
	}

	return strings.TrimSpace(strings.ToLower(key)), strings.TrimSpace(strings.ToLower(val))
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