package main

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func newEngine(tx *wire.MsgTx, output []byte) (*txscript.Engine, error) {
	e, err := txscript.NewEngine(
		output,
		tx,
		0,
		txscript.ScriptVerifyWitness+txscript.ScriptBip16+txscript.ScriptVerifyCleanStack+txscript.ScriptVerifyMinimalData,
		nil,
		nil,
		-1,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func newWitnessEngine(sigScript []byte, input [][]byte, output []byte) (*txscript.Engine, error) {
	tx := &wire.MsgTx{
		Version: 1,
		TxIn: []*wire.TxIn{
			{
				SignatureScript: sigScript,
				Witness:         input,
			},
		},
	}

	e, err := txscript.NewEngine(
		output,
		tx,
		0,
		txscript.ScriptVerifyWitness+txscript.ScriptBip16+txscript.ScriptVerifyCleanStack+txscript.ScriptVerifyMinimalData,
		nil,
		nil,
		-1,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}
