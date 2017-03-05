package build

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"time"

	"bytes"

	pb "github.com/jwhitcraft/rome/aqueduct"
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
	IdRegex  = regexp.MustCompile(`\$Id(.*)\$`)
	VarRegex = regexp.MustCompile("@_SUGAR_(FLAV|VERSION|BUILD)")

	// reusable variables in the loops
	fileByte               = []byte("File")
	variableSugarVersion   = []byte("@_SUGAR_VERSION")
	variableSugarFlav      = []byte("@_SUGAR_FLAV")
	variableSugarBuildNum  = []byte("@_SUGAR_BUILD_NUMBER")
	variableSugarBuildTime = []byte("@_SUGAR_BUILD_TIME")

	// golang has a weird date formatting thing
	// see: https://gobyexample.com/time-formatting-parsing
	timeOfBuild = []byte(time.Now().Format("2006-01-02 15:04pm"))
)

type File struct {
	Path         string
	Target       string
	fileContents []byte
	removed      bool
}

// CreateFile, Create the a File Struct and return it
func CreateFile(path, target string) *File {
	return &File{Path: path, Target: target}
}

func RemoveFile(path, target string) *File {
	f := CreateFile(path, target)
	f.removed = true

	return f
}

func CreateRemoteFile(target string, contents []byte) *File {
	return &File{Target: target, fileContents: contents}
}

func (f *File) SendToAqueduct(aqueduct pb.AqueductClient) (*pb.FileResponse, error) {

	f.readFile()
	fr := &pb.FileRequest{
		Path:     f.Path,
		Target:   f.Target,
		Contents: f.fileContents,
	}

	if f.removed == true {
		return aqueduct.DeleteFile(context.Background(), fr)
	} else {
		return aqueduct.BuildFile(context.Background(), fr)
	}
}

func (f *File) GetSource() string {
	return f.Path
}

func (f *File) Process(flavor string, version string, buildNumber string) error {
	// todo: return errors from processFile
	f.processFile(flavor, version, buildNumber)

	return nil
}

func (f *File) Delete() error {
	err := os.Remove(f.Target)
	if err != nil {
		return err
	}
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

func (f *File) processFile(buildFlavor string, buildVersion string, buildNumber string) bool {
	var useLine bool = true
	var shouldProcess bool = false
	var canProcess bool = false

	// lets make sure the that folder exists
	var destFolder string = path.Dir(f.Target)
	var fileExt string = path.Ext(f.Target)
	// var fileName string = path.Base(destPath)
	os.MkdirAll(destFolder, 0775)

	// regardless, if the file is in the node_modules folder
	// don't try and process it
	if !strings.Contains(destFolder, "node_modules") {
		canProcess = contains(processableExtensions, []byte(fileExt))
	}

	// first load the whole file to check for the build tags
	err := f.readFile()
	if err != nil {
		// todo return the error instead of false
		//return err;
		return false
	}
	fileBytes := f.fileContents
	//fileString := string(f.fileContents)
	if canProcess && TagRegex.Match(fileBytes) {
		shouldProcess = true
		// check to see if it's a type of FILE
		matches := TagRegex.FindSubmatch(fileBytes)
		if bytes.Equal(matches[1], fileByte) {
			tagOk := processBuildTag(matches[2], Flavors[buildFlavor])
			if tagOk == false {
				// todo return nil here as no file should be built
				return false
			}
		}
	}

	// do the variable replacement
	if canProcess && VarRegex.Match(fileBytes) {
		fileBytes = bytes.Replace(fileBytes, variableSugarVersion, []byte(buildVersion), -1)
		fileBytes = bytes.Replace(fileBytes, variableSugarFlav, []byte(buildFlavor), -1)
		fileBytes = bytes.Replace(fileBytes, variableSugarBuildNum, []byte(buildNumber), -1)
		fileBytes = bytes.Replace(fileBytes, variableSugarBuildTime, timeOfBuild, -1)
	}
	fw, err := os.Create(f.Target)
	defer fw.Close()

	if shouldProcess {
		f := bytes.NewReader(fileBytes)
		writer := bufio.NewWriter(fw)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			val := scanner.Bytes()

			if TagRegex.Match(val) {
				// get the matches
				matches := TagRegex.FindSubmatch(val)

				switch string(matches[1]) {
				case "BEGIN":
					useLine = processBuildTag(matches[2], Flavors[buildFlavor])
				case "END":
					useLine = true
				}
			} else if IdRegex.Match(val) {
				fmt.Fprintln(writer, "")
			} else if useLine {
				fmt.Fprintln(writer, string(val))
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		} else {
			// write the file to the disk
			writer.Flush()
		}
	} else {
		fw.WriteString(string(fileBytes))
	}

	return true
}
