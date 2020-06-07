package main

import (
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/ajaxray/geek-life/model"
	"github.com/ajaxray/geek-life/repository"
	repo "github.com/ajaxray/geek-life/repository/storm"
	"github.com/ajaxray/geek-life/util"
)

var (
	app                            *tview.Application
	newProject, newTask            *tview.InputField
	projectList, taskList          *tview.List
	projectPane, projectDetailPane *tview.Flex
	taskPane, taskDetailPane       *tview.Flex
	layout, contents               *tview.Flex

	statusBar *StatusBar

	db          *storm.DB
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository

	projects       []model.Project
	currentProject *model.Project
)

func main() {
	app = tview.NewApplication()

	db = util.ConnectStorm()
	defer db.Close()

	projectRepo = repo.NewProjectRepository(db)
	taskRepo = repo.NewTaskRepository(db)

	titleText := tview.NewTextView().SetText("[lime::b]Geek-life [::-]- Task Manager for geeks!").SetDynamicColors(true)
	cloudStatus := tview.NewTextView().SetText("[::d]Version: 0.0.3").SetTextAlign(tview.AlignRight).SetDynamicColors(true)

	titleBar := tview.NewFlex().
		AddItem(titleText, 0, 2, false).
		AddItem(cloudStatus, 0, 1, false)

	prepareProjectPane()
	prepareProjectDetail()
	prepareTaskPane()
	statusBar = makeStatusBar(app)
	prepareDetailPane()

	contents = tview.NewFlex().
		AddItem(projectPane, 25, 1, true).
		AddItem(taskPane, 0, 2, false)

	layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBar, 2, 1, false).
		AddItem(contents, 0, 2, true).
		AddItem(statusBar, 1, 1, false)

	setKeyboardShortcuts()

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func setKeyboardShortcuts() *tview.Application {
	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvt() {
			return event
		}

		// Handle based on current focus. Handlers may modify event
		switch {
		case projectPane.HasFocus():
			event = handleProjectPaneShortcuts(event)
		case taskPane.HasFocus():
			event = handleTaskPaneShortcuts(event)
		case taskDetailPane.HasFocus():
			event = handleDetailPaneShortcuts(event)
		}

		// Global shortcuts
		switch event.Rune() {
		case 'p':
			app.SetFocus(projectPane)
		case 't':
			app.SetFocus(taskPane)
		case 'f':
			// @TODO : Remove
			statusBar.showForSeconds(reflect.TypeOf(app.GetFocus()).String(), 5)
		}

		return event
	})
}
