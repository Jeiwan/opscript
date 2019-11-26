package gui

import (
	"github.com/jroimartin/gocui"
)

func (g *GUI) cursorDown(c *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if len(g.codeLines) == 0 {
			return nil
		}

		cx, cy := v.Cursor()
		nextLine := g.codeLines.next(cy)
		lastLine := g.codeLines.last()

		if nextLine.lineIdx > lastLine.lineIdx {
			return nil
		}

		if err := v.SetCursor(cx, nextLine.lineIdx); err != nil {
			ox, _ := v.Origin()
			if err := v.SetOrigin(ox, g.codeLines.first().lineIdx); err != nil {
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
		if len(g.codeLines) == 0 {
			return nil
		}

		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		prevLine := g.codeLines.previous(cy)
		firstLine := g.codeLines.first()

		if prevLine.lineIdx < firstLine.lineIdx {
			return nil
		}

		if err := v.SetCursor(cx, prevLine.lineIdx); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, firstLine.lineIdx); err != nil {
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
