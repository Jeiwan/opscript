package gui

import (
	"fmt"

	"github.com/Jeiwan/scriptdbg/debugger"
	"github.com/jroimartin/gocui"
)

const (
	keyQ       = 'q'
	viewScript = "script"
	viewStack  = "stack"
)

// GUI ...
type GUI struct {
	cui      *gocui.Gui
	debugger *debugger.Debugger
}

// New returns a new GUI.
func New(debugger *debugger.Debugger) (*GUI, error) {
	c, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	g := &GUI{
		cui:      c,
		debugger: debugger,
	}

	c.SetManagerFunc(g.layout)
	g.setKeybindings(c)

	return g, nil
}

// Stop ...
func (g GUI) Stop() {
	g.cui.Close()
}

// Start ...
func (g GUI) Start() error {
	if err := g.cui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (g GUI) layout(c *gocui.Gui) error {
	maxX, maxY := c.Size()

	if v, err := c.SetView(viewScript, 0, 0, int(0.5*float64(maxX))-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return fmt.Errorf("setView 'script': %+v", err)
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		cx, cy := v.Cursor()
		v.SetCursor(cx, cy)
		c.SetCurrentView(viewScript)

		for _, s := range g.debugger.Steps {
			fmt.Fprintf(v, "%s\n", formatOpcode(s.Disasm))
		}
	}

	if _, err := c.SetView(viewStack, int(0.5*float64(maxX)), 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := g.updateStack(); err != nil {
			return err
		}
	}

	return nil
}

func (g GUI) setKeybindings(c *gocui.Gui) error {
	if err := c.SetKeybinding(viewScript, gocui.KeyArrowUp, gocui.ModNone, g.cursorUp); err != nil {
		return err
	}

	if err := c.SetKeybinding(viewScript, gocui.KeyArrowDown, gocui.ModNone, g.cursorDown); err != nil {
		return err
	}

	if err := c.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func (g GUI) updateStack() error {
	v, err := g.cui.View(viewStack)
	if err != nil {
		return err
	}

	v.Clear()

	step := g.debugger.CurrentStep()
	for i := range step.Stack {
		fmt.Fprintf(v, "%x\n", step.Stack[len(step.Stack)-i-1])
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
