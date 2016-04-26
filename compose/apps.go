package compose

import (
	"errors"
	"fmt"
	"os"
)

func Apps(filters ...string) (apps []*App, err error) {
	dir, err := os.Open(AppsDirectory)
	if err != nil {
		return
	}
	defer dir.Close()

	fis, err := dir.Readdir(0)
	if err != nil {
		return
	}

	for _, fi := range fis {
		app := App{Name: fi.Name()}
		if !app.IsValid() || !app.Exists() {
			continue
		}
		if !app.Match(filters...) {
			continue
		}
		apps = append(apps, &app)
	}
	return
}

func Application(name ...string) (app *App, err error) {
	if len(name) == 0 {
		return nil, errors.New("specify application name")
	}
	if len(name) > 1 {
		return nil, errors.New("specify only one application name")
	}
	app = &App{Name: name[0]}
	if !app.IsValid() {
		return nil, fmt.Errorf("%v is invalid", app.Name)
	}
	return
}

func ExistingApplication(name ...string) (app *App, err error) {
	app, err = Application(name...)
	if err != nil {
		return
	}
	if !app.Exists() {
		return nil, fmt.Errorf("%v doesn't exist", app.Name)
	}
	return
}
