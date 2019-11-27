package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Jeiwan/opscript/debugger"
	"github.com/Jeiwan/opscript/spec"
	"github.com/jroimartin/gocui"
)

const (
	keyQ       = 'q'
	viewScript = "script"
	viewSpec   = "spec"
	viewStack  = "stack"
)

// GUI ...
type GUI struct {
	codeLines codeLines
	cui       *gocui.Gui
	debugger  *debugger.Debugger
	spec      spec.Script
}

// New returns a new GUI.
func New(debugger *debugger.Debugger, spec spec.Script) (*GUI, error) {
	c, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	g := &GUI{
		codeLines: []codeLine{},
		cui:       c,
		debugger:  debugger,
		spec:      spec,
	}

	g.populateCodeLines()
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

func (g *GUI) layout(c *gocui.Gui) error {
	maxX, maxY := c.Size()

	if v, err := c.SetView(viewScript, 0, 0, int(0.5*float64(maxX))-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return fmt.Errorf("setView 'script': %+v", err)
		}

		v.Title = "Script"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		g.renderCodeLines(v)

		cx, _ := v.Cursor()
		v.SetCursor(cx, g.codeLines.first().lineIdx)
		c.SetCurrentView(viewScript)

	}

	if v, err := c.SetView(viewStack, int(0.5*float64(maxX)), 0, maxX-1, int(0.7*float64(maxY))-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Stack"

		if err := g.updateStack(); err != nil {
			return err
		}
	}

	if v, err := c.SetView(viewSpec, int(0.5*float64(maxX)), int(0.7*float64(maxY)), maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Spec"
		v.Wrap = true

		if err := g.updateSpec(); err != nil {
			return err
		}
	}

	return nil
}

func (g *GUI) populateCodeLines() {
	g.codeLines = nil
	curLine := 0

	var hasSigScript bool
	var hasPkScript bool

	for _, s := range g.debugger.Steps {
		if isFirstScriptLine(s.Disasm) {
			var line codeLine
			line.isSeparator = true
			line.lineIdx = curLine

			if isSignatureScript(s.Disasm) {
				hasSigScript = true
				line.text = "        Signature Script\n"
				curLine++

			} else if isPubkeyScript(s.Disasm) {
				hasPkScript = true
				line.text = "\n        Pubkey Script\n"
				curLine++
				if hasSigScript {
					curLine++
				}

			} else if isWitnessScript(s.Disasm) {
				line.text = "\n        Witness Script\n"
				curLine++
				if hasPkScript || hasSigScript {
					curLine++
				}
			}

			g.codeLines = append(g.codeLines, line)
		}

		var line codeLine
		line.lineIdx = curLine
		line.text = fmt.Sprintln(formatDisasm(s.Disasm))

		g.codeLines = append(g.codeLines, line)
		curLine++
	}
}

func (g *GUI) renderCodeLines(v *gocui.View) {
	for _, cl := range g.codeLines {
		fmt.Fprint(v, cl.text)
	}
}

func (g GUI) setKeybindings(c *gocui.Gui) error {
	if err := c.SetKeybinding(viewScript, gocui.KeyArrowUp, gocui.ModNone, g.cursorUp); err != nil {
		return err
	}

	if err := c.SetKeybinding(viewScript, gocui.KeyArrowDown, gocui.ModNone, g.cursorDown); err != nil {
		return err
	}

	if err := c.SetKeybinding("", keyQ, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func (g GUI) updateSpec() error {
	v, err := g.cui.View(viewSpec)
	if err != nil {
		return err
	}

	v.Clear()

	step := g.debugger.CurrentStep()
	opcodeRegex := regexp.MustCompile(`OP_[\w_]+`)
	opcode := opcodeRegex.FindString(step.Disasm)
	if opcode == "" || strings.HasPrefix(opcode, "OP_DATA") {
		return nil
	}

	spec := g.spec[opcode]
	if spec.Word == "" {
		fmt.Fprintf(v, " Missing spec for %s.", opcode)
		return nil
	}

	fmt.Fprintf(v, " %s (%s)\n\n", spec.Word, spec.Opcode)
	fmt.Fprintf(v, " Input:  %s\n", spec.Input)
	fmt.Fprintf(v, " Output: %s\n", spec.Output)
	fmt.Fprintf(v, " \n %s", spec.Short)

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
		fmt.Fprintf(v, " %x\n", step.Stack[len(step.Stack)-i-1])
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
