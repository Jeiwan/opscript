package debugger

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/sirupsen/logrus"
)

// Debugger ...
type Debugger struct {
	CurrentPos uint16
	Engine     *txscript.Engine
	Steps      []Step
}

// New ...
func New(en *txscript.Engine) (*Debugger, error) {
	var steps []Step

Loop:
	for {
		var step Step

		d, err := en.DisasmPC()
		if err != nil {
			return nil, err
		}
		logrus.Debugf("EXECUTING: %s", d)

		done, err := en.Step()
		if err != nil {
			return nil, err
		}

		step.Disasm = d
		step.Stack = en.GetStack()

		steps = append(steps, step)

		if done {
			logrus.Debugf("FINISHED")
			break Loop
		}
	}

	return &Debugger{
		CurrentPos: 0,
		Engine:     en,
		Steps:      steps,
	}, nil
}

// Next ...
func (d *Debugger) Next() bool {
	if int(d.CurrentPos+1) >= len(d.Steps) {
		return false
	}

	d.CurrentPos++

	return true
}

// Previous ...
func (d *Debugger) Previous() bool {
	if int(d.CurrentPos-1) < 0 {
		return false
	}

	d.CurrentPos--

	return true
}

// CurrentStep ...
func (d *Debugger) CurrentStep() *Step {
	return &d.Steps[d.CurrentPos]
}
