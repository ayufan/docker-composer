package helpers

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func GitEditor() (editor string, err error) {
	editor = os.Getenv("GIT_EDITOR")
	term := os.Getenv("TERM")
	isDumb := term == "" || term == "dumb"

	if editor == "" && !isDumb {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" && isDumb {
		return "", errors.New("No GIT_EDITOR defined")
	}
	if editor == "" {
		editor = "vi"
	}
	return
}

func EditFile(name string) (err error) {
	editor, err := GitEditor()
	if err != nil {
		return
	}

	logrus.Infoln("Editing", filepath.Base(name), "...")
	cmd := Command(editor, name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
