package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/b4b4r07/go-colon"
)

func Filter(selecter, text string) ([]string, error) {
	var (
		selectedLines []string
		buf           bytes.Buffer
		err           error
	)
	if text == "" {
		return selectedLines, errors.New("No input")
	}
	if selecter == "" {
		return selectedLines, errors.New("no selectcmd specified")
	}
	err = runFilter(selecter, strings.NewReader(text), &buf)
	if err != nil {
		return selectedLines, err
	}
	if buf.Len() == 0 {
		return selectedLines, errors.New("no lines selected")
	}
	for _, line := range strings.Split(buf.String(), "\n") {
		if line == "" {
			continue
		}
		selectedLines = append(selectedLines, line)
	}
	return selectedLines, nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}

func runFilter(command string, r io.Reader, w io.Writer) error {
	command = expandPath(command)
	result, err := colon.Parse(command)
	if err != nil {
		return err
	}
	first, err := result.Executable().First()
	if err != nil {
		return err
	}
	command = first.Item
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}

func Run(command string, args ...string) error {
	if command == "" {
		return errors.New("command not found")
	}
	command += " " + strings.Join(args, " ")
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
