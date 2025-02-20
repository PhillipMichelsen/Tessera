package controllers

import (
	"AlgorithmicTraderDistributed/internal/api"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUIController struct {
	instance api.InstanceAPIExternal
    app      *tview.Application
    layout   *tview.Flex
    modules  *tview.TextView
}

func NewTUIController(instance api.InstanceAPIExternal) *TUIController {
    controller := &TUIController{
        instance: instance,
        app:      tview.NewApplication(),
    }

    controller.initializeUI()
    return controller
}

func (t *TUIController) initializeUI() {
    // Create the modules view
    t.modules = tview.NewTextView()
    t.modules.SetDynamicColors(true)
    t.modules.SetRegions(true)
    t.modules.SetScrollable(true)
    t.modules.SetTitle("Running Modules")
    t.modules.SetBorder(true)

    buttons := tview.NewFlex().SetDirection(tview.FlexRow)
    shutdownBtn := tview.NewButton("Shutdown").
        SetSelectedFunc(func() {
            t.instance.Shutdown()
        })

    buttons.AddItem(shutdownBtn, 1, 0, false)

    // Create main layout
    t.layout = tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(t.modules, 0, 1, false).
        AddItem(buttons, 3, 0, true)

    // Set up keyboard shortcuts
    t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyEsc {
            t.instance.Shutdown()
            return nil
        }
        return event
    })
}

func (t *TUIController) updateModules() {
    for {
        modules := t.instance.GetModules()
        var content string
        content = "[yellow]Active Modules:[white]\n"
        
        for _, moduleUUID := range modules {
            status := t.instance.GetModuleStatus(moduleUUID)
            content += fmt.Sprintf("• Module %s (Status: %s)\n", moduleUUID, status)
        }

        t.app.QueueUpdateDraw(func() {
            t.modules.SetText(content)
        })
        
        time.Sleep(1 * time.Second)
    }
}

func (t *TUIController) Start() {
    go t.updateModules()
    
    if err := t.app.SetRoot(t.layout, true).EnableMouse(true).Run(); err != nil {
        panic(err)
    }
}

func (t *TUIController) Stop() {
    t.app.Stop()
}