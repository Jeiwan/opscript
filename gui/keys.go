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

		maxY-- // one-line padding at the bottom

		if nextLine.lineIdx == cy+oy {
			return nil
		}

		if nextLine.lineIdx > maxY && cy == maxY {
			oy++
			if err := v.SetOrigin(ox, oy); err != nil {
				return err
			}
		}

		if err := v.SetCursor(cx, nextLine.lineIdx-oy); err != nil {
			if err := v.SetOrigin(ox, maxY); err != nil {
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
		cx, cy := v.Cursor()

		prevLine := g.codeLines.previous(oy + cy)
		firstLine := g.codeLines.first()

		if prevLine.lineIdx == 0 {
			if err := v.SetOrigin(cx, 0); err != nil {
				return err
			}
			if err := v.SetCursor(cx, firstLine.lineIdx); err != nil {
				return err
			}
			return nil
		}

		if prevLine.lineIdx < oy {
			if err := v.SetOrigin(cx, oy-(oy-prevLine.lineIdx)); err != nil {
				return err
			}
			prevLine.lineIdx = oy
		}

		if err := v.SetCursor(cx, prevLine.lineIdx-oy); err != nil {
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
