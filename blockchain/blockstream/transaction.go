package blockstream

type transaction struct {
	Txid     string `json:"txid"`
	Version  int32  `json:"version"`
	Locktime int    `json:"locktime"`
	Size     int    `json:"size"`
	Weight   int    `json:"weight"`
	Fee      int    `json:"fee"`
	Vin      []vin  `json:"vin"`
	Vout     []vout `json:"vout"`
	Status   status `json:"status"`
}

type vin struct {
	Txid                  string   `json:"txid"`
	Vout                  uint32   `json:"vout"`
	IsCoinbase            bool     `json:"is_coinbase"`
	Scriptsig             string   `json:"scriptsig"`
	ScriptsigAsm          string   `json:"scriptsig_asm"`
	InnerRedeemscriptAsm  string   `json:"inner_redeemscript_asm"`
	InnerWitnessscriptAsm string   `json:"inner_witnessscript_asm"`
	Sequence              int      `json:"sequence"`
	Witness               []string `json:"witness"`
	Prevout               vout     `json:"prevout"`
}

type vout struct {
	Scriptpubkey        string `json:"scriptpubkey"`
	ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
	ScriptpubkeyType    string `json:"scriptpubkey_type"`
	ScriptpubkeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

type status struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int    `json:"block_time"`
}
