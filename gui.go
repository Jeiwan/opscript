package main

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

type gui struct {
	cui      *gocui.Gui
	debugger *debugger.Debugger
}

func newGui(debugger *debugger.Debugger) (*gui, error) {
	c, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	g := &gui{
		cui:      c,
		debugger: debugger,
	}

	c.SetManagerFunc(g.layout)

	if err := c.SetKeybinding(viewScript, gocui.KeyArrowUp, gocui.ModNone, g.cursorUp); err != nil {
		return nil, err
	}

	if err := c.SetKeybinding(viewScript, gocui.KeyArrowDown, gocui.ModNone, g.cursorDown); err != nil {
		return nil, err
	}

	if err := c.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return nil, fmt.Errorf("setKey q: %+v", err)
	}

	return g, nil
}

func (g gui) Stop() {
	g.cui.Close()
}

func (g gui) Start() error {
	if err := g.cui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (g gui) layout(c *gocui.Gui) error {
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
			fmt.Fprintf(v, "%s\n", s.Disasm)
		}
	}

	if v, err := c.SetView(viewStack, int(0.5*float64(maxX)), 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		step := g.debugger.CurrentStep()
		for _, s := range step.Stack {
			fmt.Fprintf(v, "%x\n", s)
		}
	}

	return nil
}

func (g gui) updateStack() error {
	v, err := g.cui.View(viewStack)
	if err != nil {
		return err
	}

	v.Clear()

	step := g.debugger.CurrentStep()
	for _, s := range step.Stack {
		fmt.Fprintf(v, "%x\n", s)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
