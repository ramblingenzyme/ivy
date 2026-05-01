package dvdrip

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
)

type DVDModel struct {
	Title     string
	MainTrack string
}

func DVDInfo(device string, w io.Writer) (DVDModel, error) {
	fmt.Fprintf(w, "$ dvdbackup -i %s -I\n", device)
	dvdOutput, err := exec.Command("dvdbackup", "-i", device, "-I").Output()
	if err != nil {
		return DVDModel{}, err
	}

	re := regexp.MustCompile(`(?m)DVD with title "(?<title>.+)"$`)
	match := re.FindSubmatch(dvdOutput)
	if match == nil {
		return DVDModel{}, errors.New("couldn't find title")
	}
	title := string(match[re.SubexpIndex("title")])

	re = regexp.MustCompile(`(?m)containing the main feature is (?<track>\d+)$`)
	match = re.FindSubmatch(dvdOutput)
	if match == nil {
		return DVDModel{}, errors.New("couldn't find main track")
	}
	mainTrackInt, err := strconv.Atoi(string(match[re.SubexpIndex("track")]))
	if err != nil {
		return DVDModel{}, err
	}

	return DVDModel{
		Title:     title,
		MainTrack: fmt.Sprintf("%02d", mainTrackInt),
	}, nil
}

func BackupDVD(device, outputDir string, w io.Writer) error {
	fmt.Fprintf(w, "$ dvdbackup -i %s -o %s -F -p\n", device, outputDir)
	cmd := exec.Command("dvdbackup", "-i", device, "-o", outputDir, "-F", "-p")
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

// MergeMKV merges a VOB file to MKV. path should be <movie>/VIDEO_TS/VTS_<track>_0.VOB
func MergeMKV(path, output string, w io.Writer) error {
	fmt.Fprintf(w, "$ mkvmerge -o %s %s\n", output, path)
	cmd := exec.Command("mkvmerge", "-o", output, path)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}
