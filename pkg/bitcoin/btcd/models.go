package btcd

import (
	"encoding/json"
)

type BlockTemplate struct {
	Version           int    `json:"version"`
	VbRequired        uint32 `json:"vbrequired"`
	PreviousBlockHash string `json:"previousblockhash"`
	Transactions      []struct {
		Data    string `json:"data"`
		TxId    string `json:"txid"`
		Hash    string `json:"hash"`
		Depends []int  `json:"depends"`
		Fee     int    `json:"fee"`
		SigOps  int    `json:"sigops"`
		Weight  int    `json:"weight"`
	} `json:"transactions"`
	CoinbaseAux struct {
		Key string `json:"key"`
	}
	CoinbaseValue            int      `json:"coinbasevalue"`
	LongPollId               string   `json:"longpollid"`
	Target                   string   `json:"target"`
	MinTime                  int      `json:"mintime"`
	Mutable                  []string `json:"mutable"`
	NonceRange               string   `json:"noncerange"`
	SigOpLimit               int      `json:"sigoplimit"`
	SizeLimit                int      `json:"sizelimit"`
	WeightLimit              int      `json:"weightlimit"`
	CurTime                  int      `json:"curtime"`
	Bits                     string   `json:"bits"`
	Height                   int      `json:"height"`
	DefaultWitnessCommitment string   `json:"default_witness_commitment"`
}

func (bt *BlockTemplate) Marshal() string {
	btBytes, _ := json.Marshal(bt)
	return string(btBytes)
}

func (bt *BlockTemplate) Unmarshal(jsonStr string) {
	json.Unmarshal([]byte(jsonStr), bt)
}
