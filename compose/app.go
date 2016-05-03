package compose

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"

	"fmt"
	"github.com/ayufan/docker-composer/helpers"
	"io/ioutil"
)

var AppsDirectory string

const (
	StatusRunning      string = "running"
	StatusNotRunning          = "not running"
	StatusPartial             = "partial"
	StatusError               = "error"
	StatusNoContainers        = "no containers"
)

type App struct {
	Name string
}

func (a *App) log(method ...string) *logrus.Entry {
	entry := logrus.WithField("app-name", a.Name)
	if len(method) > 0 {
		entry = entry.WithField("method", method[0])
	}
	return entry
}

func (a *App) Init() error {
	cmd := helpers.Git("init", a.Path())
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (a *App) UpdateHooks() error {
	hooksDir := a.Path(".git", "hooks")
	err := os.RemoveAll(hooksDir)
	if err != nil {
		return fmt.Errorf("RemoveHooks: %v", err)
	}
	err = os.MkdirAll(hooksDir, 0700)
	if err != nil {
		return fmt.Errorf("MkdirHooks: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(hooksDir, "push-to-checkout"), []byte(pushToCheckout), 0700)
	if err != nil {
		return fmt.Errorf("WriteHooks: %v", err)
	}
	return nil
}

func (a *App) UpdateConfig() error {
	cmd := helpers.Git("config", "receive.denyCurrentBranch", "updateInstead")
	cmd.Dir = a.Path()
	return cmd.Run()
}

func (a *App) Path(elem ...string) string {
	return filepath.Join(append([]string{AppsDirectory, a.Name}, elem...)...)
}

func (a *App) Compose(action string, args ...string) (err error) {
	cmd := helpers.Compose(action, a.Path(), args...)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (a *App) Git(args ...string) (err error) {
	cmd := helpers.Git(args...)
	cmd.Dir = a.Path()
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (a *App) Changed() (bool, error) {
	cmd := helpers.Git("status", "-s")
	cmd.Dir = a.Path()
	data, err := cmd.Output()
	if err != nil {
		return false, err
	}
	if len(data) == 0 {
		return false, nil
	}
	return true, nil
}

func (a *App) Commit(message string) error {
	return a.Git("commit", "-m", message)
}

func (a *App) Tag() (err error) {
	return a.Git("tag", "-f", "latest")
}

func (a *App) Reset(name string) (err error) {
	return a.Git("reset", "--hard", name)
}

func (a *App) Revert() (err error) {
	return a.Reset("latest")
}

func (a *App) Validate() (err error) {
	cmd := helpers.Compose("config", a.Path(), "-q")
	return cmd.Run()
}

func (a *App) Deploy(args ...string) (err error) {
	cmd := helpers.Compose("up", a.Path(), "--build", "--remove-orphans", "-d")
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return
	}
	return
}

func (a *App) Reference() (reference string, err error) {
	cmd := helpers.Git("symbolic-ref", "-q", "--short", "HEAD")
	cmd.Dir = a.Path()
	data, err := cmd.Output()
	if err != nil {
		return
	}
	reference = strings.TrimSpace(string(data))
	return
}

func (a *App) Containers() (containers []string, err error) {
	cmd := helpers.Compose("ps", a.Path(), "-q")
	data, err := cmd.Output()
	if err != nil {
		return
	}
	return helpers.Lines(data), nil
}

func (a *App) ContainerNames() (names []string, err error) {
	containers, err := a.Containers()
	if err != nil || len(containers) == 0 {
		return
	}

	cmd := helpers.Docker("inspect", "--format", "{{.Name}}")
	cmd.Args = append(cmd.Args, containers...)
	data, err := cmd.Output()
	if err != nil {
		return
	}

	return helpers.Lines(data), nil
}

func (a *App) Services() (services []string, err error) {
	cmd := helpers.Compose("config", a.Path(), "--services")
	data, err := cmd.Output()
	if err != nil {
		return
	}
	return helpers.Lines(data), nil
}

func (a *App) Status() (status string, err error) {
	containers, err := a.Containers()
	if err != nil {
		return StatusNoContainers, err
	}
	if len(containers) == 0 {
		return StatusNotRunning, nil
	}
	cmd := helpers.Docker("inspect", "--format", "{{.State.Running}}")
	cmd.Args = append(cmd.Args, containers...)
	data, err := cmd.Output()
	if err != nil {
		return StatusError, err
	}
	running := false
	notrunning := false
	for _, status := range helpers.Lines(data) {
		if status == "true" {
			running = true
		} else if status == "false" {
			notrunning = true
		}
	}
	if running && notrunning {
		return StatusPartial, err
	} else if running {
		return StatusRunning, err
	} else if notrunning {
		return StatusNotRunning, err
	} else {
		return StatusError, err
	}
}

func (a *App) IsValid() bool {
	if strings.HasPrefix(a.Name, ".") {
		a.log("IsValid").Debugln("has dot in name")
		return false
	}
	return true
}

func (a *App) Exists() bool {
	path := a.Path()
	fi, err := os.Stat(path)
	if err != nil {
		a.log("IsValid").WithError(err).Debugln("Stat")
		return false
	}
	if !fi.IsDir() {
		a.log("IsValid").WithError(err).Debugln("Not a dir")
		return false
	}
	return true
}

func (a *App) Match(filters ...string) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		matched, err := filepath.Match(filter, a.Name)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
