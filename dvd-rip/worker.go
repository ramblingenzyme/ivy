package dvdrip

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MountPath = "/mnt/dvdfs"

type DVDModel struct {
	Title     string
	MainTrack string
}

type Session struct {
	mountPath string
	id        string
	logOffset int64
}

func NewSession(mountPath string) (*Session, error) {
	data, err := os.ReadFile(filepath.Join(mountPath, "clone"))
	if err != nil {
		return nil, fmt.Errorf("read clone: %w", err)
	}
	id := strings.TrimSpace(string(data))
	return &Session{mountPath: mountPath, id: id}, nil
}

func (s *Session) dir() string {
	return filepath.Join(s.mountPath, s.id)
}

func (s *Session) ctl(cmd string) error {
	f, err := os.OpenFile(filepath.Join(s.dir(), "ctl"), os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, cmd)
	return err
}

// tailLog polls the log file from s.logOffset, writing lines to w.
// Returns nil on "done" or an error if the operation reported failure.
func (s *Session) tailLog(w io.Writer) error {
	f, err := os.Open(filepath.Join(s.dir(), "log"))
	if err != nil {
		return err
	}

	var pending []byte
	buf := make([]byte, 4096)
	for {
		n, readErr := f.ReadAt(buf, s.logOffset)
		if n > 0 {
			s.logOffset += int64(n)
			pending = append(pending, buf[:n]...)
			for {
				idx := bytes.IndexByte(pending, '\n')
				if idx < 0 {
					break
				}
				line := string(pending[:idx])
				pending = pending[idx+1:]
				if line == "done" {
					return nil
				}
				if strings.HasPrefix(line, "error: ") {
					return errors.New(strings.TrimPrefix(line, "error: "))
				}
				fmt.Fprintln(w, line)
			}
		}
		if readErr == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if readErr != nil {
			return readErr
		}
	}
}

func (s *Session) Info(device string, w io.Writer) (DVDModel, error) {
	if err := s.ctl("device " + device); err != nil {
		return DVDModel{}, err
	}
	if err := s.ctl("info"); err != nil {
		return DVDModel{}, err
	}
	if err := s.tailLog(w); err != nil {
		return DVDModel{}, err
	}
	data, err := os.ReadFile(filepath.Join(s.dir(), "info"))
	if err != nil {
		return DVDModel{}, fmt.Errorf("read info: %w", err)
	}
	return parseInfo(data)
}

func parseInfo(data []byte) (DVDModel, error) {
	var title, mainTrack string
	for _, line := range bytes.Split(data, []byte("\n")) {
		f := strings.Fields(string(line))
		switch {
		case len(f) >= 2 && f[0] == "name":
			title = strings.Join(f[1:], " ")
		case len(f) >= 2 && f[0] == "title" && f[len(f)-1] == "main":
			mainTrack = f[1]
		}
	}
	if title == "" {
		return DVDModel{}, errors.New("couldn't find title in dvd info")
	}
	if mainTrack == "" {
		return DVDModel{}, errors.New("couldn't find main track in dvd info")
	}
	return DVDModel{Title: title, MainTrack: mainTrack}, nil
}

func (s *Session) Backup(outputDir string, w io.Writer) error {
	if err := s.ctl("output " + outputDir); err != nil {
		return err
	}
	if err := s.ctl("backup"); err != nil {
		return err
	}
	return s.tailLog(w)
}

func (s *Session) MergeMKV(w io.Writer) error {
	if err := s.ctl("mkv"); err != nil {
		return err
	}
	return s.tailLog(w)
}

func (s *Session) Close() {
	s.ctl("close")
}
