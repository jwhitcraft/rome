package build

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	pb "github.com/jwhitcraft/rome/aqueduct"
	"github.com/jwhitcraft/rome/utils"
	"golang.org/x/net/context"
)

var (
	processableExtensions = []string{
		".php", ".json", ".js", ".html", ".tpl", ".css", ".hbs",
	}
	Flavors = map[string][]string{
		"pro":  {"pro"},
		"corp": {"pro", "corp"},
		"ent":  {"pro", "ent"},
		"ult":  {"pro", "ent", "ult"},
	}

	License = map[string][]string{
		"lic": {"sub"},
	}

	TagRegex = regexp.MustCompile("(?i)//[[:space:]]*(BEGIN|END|FILE|ELSE)[[:space:]]*SUGARCRM[[:space:]]*(.*) ONLY")

	IdRegex = regexp.MustCompile(`\$Id(.*)\$`)

	VarRegex = regexp.MustCompile("@_SUGAR_(FLAV|VERSION)")
)

type File struct {
	Path         string
	Target       string
	fileContents []byte
}

// CreateFile, Create the a File Struct and return it
func CreateFile(path, target string) *File {
	return &File{Path: path, Target: target}
}

func CreateRemoteFile(target string, contents []byte) *File {
	return &File{Target: target, fileContents: contents}
}

func (f *File) SendToAqueduct(aqueduct pb.AqueductClient) (*pb.FileResponse, error) {
	f.readFile()
	return aqueduct.BuildFile(context.Background(), &pb.FileRequest{
		Path:     f.Path,
		Target:   f.Target,
		Contents: f.fileContents,
	})
}

func (f *File) GetSource() string {
	return f.Path
}

func (f *File) Process(flavor string, version string) error {
	// todo: return errors from processFile
	f.processFile(flavor, version)

	return nil
}

func (f *File) GetTarget() string {
	return f.Target
}

func (f *File) readFile() error {
	var err error

	// prevent multiple ReadFile calls
	if f.fileContents == nil {
		f.fileContents, err = ioutil.ReadFile(f.Path)
	}

	return err
}

func (f *File) processFile(buildFlavor string, buildVersion string) bool {
	var useLine bool = true
	var shouldProcess bool = false
	var canProcess bool = false

	var skippedLines utils.Counter

	// lets make sure the that folder exists
	var destFolder string = path.Dir(f.Target)
	var fileExt string = path.Ext(f.Target)
	// var fileName string = path.Base(destPath)
	os.MkdirAll(destFolder, 0775)

	// regardless, if the file is in the node_modules folder
	// don't try and process it
	if !strings.Contains(destFolder, "node_modules") {
		canProcess = contains(processableExtensions, fileExt)
	}

	// first load the whole file to check for the build tags
	err := f.readFile()
	if err != nil {
		// todo return the error instead of false
		//return err;
		return false
	}
	fileString := string(f.fileContents)
	if canProcess && TagRegex.MatchString(fileString) {
		shouldProcess = true
		// check to see if it's a type of FILE
		matches := TagRegex.FindStringSubmatch(fileString)
		if matches[1] == "FILE" {
			tagOk := processBuildTag(matches[2], Flavors[buildFlavor])
			if tagOk == false {
				// todo return nil here as no file should be built
				return false
			}
		}
	}

	// do the variable replacement
	if canProcess && VarRegex.MatchString(fileString) {
		fileString = strings.Replace(fileString, "@_SUGAR_VERSION", buildVersion, -1)
		fileString = strings.Replace(fileString, "@_SUGAR_FLAV", buildFlavor, -1)
	}
	fw, err := os.Create(f.Target)
	defer fw.Close()

	if shouldProcess {
		f := strings.NewReader(fileString)
		if err != nil {
			fmt.Printf("error opening file: %v\n", err)
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
