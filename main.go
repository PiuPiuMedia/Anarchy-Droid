package main

import (
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"anarchy-droid/logger"

	"github.com/getsentry/sentry-go"
)

const AppName = "Anarchy-Droid"
const AppVersion = "0.9.0"
var BuildDate string	// Build date injected during build
// use: go run -ldflags "-X main.BuildDate=`date +%Y-%m-%d`" .

var a fyne.App
var w fyne.Window
var active_screen string

func main() {
	// Propagate AppVersion and BuildDate to the logger package
	logger.Consent = true
	logger.AppName = AppName
	logger.AppVersion = AppVersion
	logger.BuildDate = BuildDate

	err := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: "https://26d9d7416f0e45ac8ab1733c8d691f1d@o1013551.ingest.sentry.io/5978898",
		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		// Environment: "",
		Release:     AppName + "@" + AppVersion,
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: false,
	})
	if err != nil {
		logger.Log("sentry.Init: " + err.Error())
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: logger.Sessionmap["id"]})
	})
	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(5 * time.Second)

	a = app.New()
	icon, err := fyne.LoadResourceFromPath("icon.png")
	if err != nil {
		logger.LogError("Error loading icon.png:", err)
	}
	a.SetIcon(icon)

	w = a.NewWindow(AppName)
	w.SetMaster()

	w.Resize(fyne.NewSize(562, 226))
	w.SetFixedSize(true)

	active_screen = "initScreen"
	w.SetContent(initScreen())

	// Set working directory to a subdir named like the app
	_, err = os.Stat(AppName)
	if os.IsNotExist(err) {
		err = os.Mkdir(AppName, 0755)
	    if err != nil {
	        logger.LogError("Error setting working directory:", err)
	    }
	}
	os.Chdir(AppName)

	go func() {
		go logger.Report(map[string]string{"progress":"Setup App"})
		_, err := initApp()
		// logger.Log(success)
		if err != nil {
			logger.LogError("Setup failed:", err)
		}
	}()

	w.ShowAndRun()
}