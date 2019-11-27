package gui

import (
	"github.com/jroimartin/gocui"
)

func (g *GUI) cursorDown(c *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if len(g.codeLines) == 0 {
			return nil
		}

		ox, oy := v.Origin()
		_, maxY := v.Size()
		cx, cy := v.Cursor()
		nextLine := g.codeLines.next(cy + oy)
		lastLine := g.codeLines.last()

		maxY-- // one-line padding at the bottom

		if nextLine.lineIdx > lastLine.lineIdx {
			return nil
		}

		if nextLine.lineIdx > maxY {
			if err := v.SetOrigin(cx, oy+1); err != nil {
				return err
			}
			nextLine.lineIdx = maxY
		}

		if err := v.SetCursor(cx, nextLine.lineIdx); err != nil {
			if err := v.SetOrigin(ox, g.codeLines.first().lineIdx); err != nil {
				return err
			}
		}

		g.debugger.Next()

		if err := g.updateStack(); err != nil {
			return err
		}

		if err := g.updateSpec(); err != nil {
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
		_, maxY := v.Size()
		cx, cy := v.Cursor()

		prevLine := g.codeLines.previous(oy + cy)
		firstLine := g.codeLines.first()

		maxY-- // one-line padding at the bottom

		if prevLine.lineIdx < firstLine.lineIdx {
			return nil
		}

		if prevLine.lineIdx >= maxY {
			if err := v.SetOrigin(cx, oy-1); err != nil {
				return err
			}
			prevLine.lineIdx = maxY
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

		if err := g.updateSpec(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GUI) showDebugView(c *gocui.Gui, v *gocui.View) error {
	if _, err := g.cui.SetViewOnTop(viewDebug); err != nil {
		return err
	}

	return nil
}
