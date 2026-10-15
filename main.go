package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// properly handle exit (caused freezing before)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}
		return event
	})

	form, err := ServerConfigForm(app)
	if err != nil {
		panic(err)
	}

	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func ServerConfigForm(app *tview.Application) (*tview.Form, error) {
	form := tview.NewForm()

	minecraftManifest, err := FetchMinecraftManifest()
	if err != nil {
		return form, err
	}

	javaManifest, err := FetchJavaManifest()
	if err != nil {
		return form, err
	}

	allJavaVersions := javaManifest.GetStrList()

	initialMcVersionDetails, err := minecraftManifest.FetchVersionDetails("LATEST")
	if err != nil {
		return form, err
	}
	initialJavaVersions := initialMcVersionDetails.GetJavaVersions(allJavaVersions)
	if err != nil {
		return form, err
	}

	javaVersionDropdown := tview.NewDropDown().SetLabel("Java Version").
		SetOptions(initialJavaVersions, nil)
	if len(initialJavaVersions) > 0 {
		javaVersionDropdown.SetCurrentOption(0)
	}

	serverTypeDropdown := tview.NewDropDown().SetLabel("Server Type").
		SetOptions([]string{"Paper", "Forge", "Neoforge", "Fabric", "Quilt", "Spigot", "Vanilla"}, nil).
		SetCurrentOption(0)

	selectedServerVersion := "LATEST"
	showAllJavaVersions := false

	initializing := true
	var refreshJavaVersions func()
	serverVersionDropdown := tview.NewDropDown().SetLabel("Version").
		SetOptions(minecraftManifest.GetStrList(), func(text string, _ int) {
			if initializing {
				return
			}
			selectedServerVersion = text
			refreshJavaVersions()
		}).
		SetCurrentOption(0)
	initializing = false

	showAllJavaVersionsCheckbox := tview.NewCheckbox().SetLabel("Show all Java Versions")
	requestID := 0
	refreshJavaVersions = func() {
		requestID++
		currentRequestID := requestID

		if showAllJavaVersions {
			javaVersionDropdown.SetLabel("Java Version").
				SetOptions(allJavaVersions, nil)
			if len(allJavaVersions) > 0 {
				javaVersionDropdown.SetCurrentOption(0)
			}
			return
		}

		version := selectedServerVersion
		javaVersionDropdown.SetLabel("Java Version (loading)")
		go func() {
			versionDetails, err := minecraftManifest.FetchVersionDetails(version)
			var versions []string
			if err == nil {
				versions = versionDetails.GetJavaVersions(allJavaVersions)
			}

			app.QueueUpdateDraw(func() {
				if currentRequestID != requestID {
					return
				}
				if err != nil {
					javaVersionDropdown.SetLabel("Java Version (load failed)").
						SetOptions(nil, nil)
					return
				}

				javaVersionDropdown.SetLabel("Java Version").
					SetOptions(versions, nil)
				if len(versions) > 0 {
					javaVersionDropdown.SetCurrentOption(0)
				}
			})
		}()
	}

	showAllJavaVersionsCheckbox.SetChangedFunc(func(checked bool) {
		showAllJavaVersions = checked
		refreshJavaVersions()
	})

	form.AddFormItem(serverTypeDropdown).
		AddFormItem(serverVersionDropdown).
		AddFormItem(showAllJavaVersionsCheckbox).
		AddFormItem(javaVersionDropdown)

	form.SetBorder(true).SetTitle("Create new server").SetTitleAlign(tview.AlignCenter)

	return form, nil
}
