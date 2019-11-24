package gui

import (
	"github.com/jroimartin/gocui"
)

func (g *GUI) cursorDown(c *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy+1 >= len(g.debugger.Steps) {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}

		g.debugger.Next()

		if err := g.updateStack(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GUI) cursorUp(c *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy-1 < 0 {
			return nil
		}

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}

		g.debugger.Previous()

		if err := g.updateStack(); err != nil {
			return err
		}
	}
	return nil
}
