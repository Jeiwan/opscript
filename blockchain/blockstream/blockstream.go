package blockstream

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

const baseURL = "https://blockstream.info/api/"

// Blockstream is a client for https://blockstream.info
type Blockstream struct {
}

// New returns a new Blockstream client.
func New() *Blockstream {
	return &Blockstream{}
}

// GetTransaction returns a transaction by its hash.
func (b Blockstream) GetTransaction(txHash string) (*wire.MsgTx, error) {
	url := fmt.Sprintf("%s/tx/%s", baseURL, txHash)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get a transaction: %s", b)
	}

	var tx transaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return nil, err
	}

	msgTx := wire.NewMsgTx(tx.Version)

	for _, vin := range tx.Vin {
		voutHash, err := chainhash.NewHashFromStr(vin.Txid)
		if err != nil {
			return nil, err
		}

		sigScript, err := hex.DecodeString(vin.Scriptsig)
		if err != nil {
			return nil, err
		}

		var witness [][]byte
		for _, w := range vin.Witness {
			ws, err := hex.DecodeString(w)
			if err != nil {
				return nil, err
			}

			witness = append(witness, ws)
		}

		msgTx.AddTxIn(
			wire.NewTxIn(
				wire.NewOutPoint(voutHash, vin.Vout),
				sigScript,
				witness,
			),
		)
	}

	for _, vout := range tx.Vout {
		pkScript, err := hex.DecodeString(vout.Scriptpubkey)
		if err != nil {
			return nil, err
		}

		msgTx.AddTxOut(
			wire.NewTxOut(
				vout.Value,
				pkScript,
			),
		)
	}

	if msgTx.TxHash().String() != tx.Txid {
		return nil, fmt.Errorf("transaction hash doesn't match")
	}

	return msgTx, nil
}
