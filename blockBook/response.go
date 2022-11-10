package blockBook

type StatusResponse struct {
	BlockBook StatusBlockBook `json:"blockbook"`
	Backend   StatusBackend   `json:"backend"`
}

type StatusBlockBook struct {
	Coin          string `json:"coin"`
	Host          string `json:"host"`
	InSync        bool   `json:"inSync"`
	InSyncMempool bool   `json:"inSyncMempool"`
}

type StatusBackend struct {
	Chain    string `json:"chain"`
	Blocks   int64  `json:"blocks"`
	Warnings string `json:"warnings,omitempty"`
}

type AddressResponse struct {
	Page               int32          `json:"page"`
	TotalPages         int32          `json:"totalPages"`
	ItemsOnPage        int32          `json:"itemsOnPage"`
	Address            string         `json:"address"`
	Balance            string         `json:"balance"`
	TotalReceived      string         `json:"totalReceived"`
	TotalSent          string         `json:"totalSent"`
	UnconfirmedBalance string         `json:"unconfirmedBalance"`
	UnconfirmedTxs     int32          `json:"unconfirmedTxs"`
	Txs                int32          `json:"txs"`
	TxIds              []string       `json:"txids"`
	Nonce              string         `json:"nonce"`
	NonTokenTxs        int32          `json:"nonTokenTxs"`
	Tokens             []AddressToken `json:"tokens"`
}

type AddressToken struct {
	TokenType string `json:"type"`
	Name      string `json:"name"`
	Contract  string `json:"contract"`
	Transfers int32  `json:"transfers"`
	Symbol    string `json:"symbol"`
	SubUnit   int    `json:"decimals"`
	Balance   string `json:"balance"`
}

// ========================================== //

type TransactionResponse struct {
	TxId             string                    `json:"txid"`
	VIn              []ResponseTxVIn           `json:"vin"`
	VOut             []ResponseTxVOut          `json:"vout"`
	BlockHash        string                    `json:"blockHash"`
	BlockHeight      int64                     `json:"blockHeight"`
	Confirmations    int64                     `json:"confirmations"`
	BlockTime        int64                     `json:"blockTime"`
	Value            string                    `json:"value"`
	Fees             string                    `json:"fees"`
	TokenTransfers   []*TokenTransfers         `json:"tokenTransfers,omitempty"`
	EthereumSpecific *ResponseEthereumSpecific `json:"ethereumSpecific"`
}

type ResponseEthereumSpecific struct {
	Status   int    `json:"status"`
	Nonce    int64  `json:"nonce"`
	GasLimit int64  `json:"gasLimit"`
	GasUsed  int64  `json:"gasUsed"`
	GasPrice string `json:"gasPrice"`
}

type TokenTransfers struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	To      string `json:"to"`
	Token   string `json:"token"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	SubUnit int    `json:"decimals"`
	Value   string `json:"value"`
}

type ResponseTxVIn struct {
	N         int64    `json:"n"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
}

type ResponseTxVOut struct {
	Value     string   `json:"value"`
	N         int64    `json:"n"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
}

// ==========================================

// ========================================== //

type BlockResponseTxVIn struct {
	N         int64    `json:"n"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
}

type BlockResponseTxVOut struct {
	Value     string   `json:"value"`
	N         int64    `json:"n"`
	Addresses []string `json:"addresses"`
	IsAddress bool     `json:"isAddress"`
}

type BlockResponseEthereumSpecific struct {
	Status   int    `json:"status"`
	Nonce    int64  `json:"nonce"`
	GasLimit int64  `json:"gasLimit"`
	GasUsed  int64  `json:"gasUsed"`
	GasPrice string `json:"gasPrice"`
}

type BlockResponseTx struct {
	TxId             string                         `json:"txid"`
	VIn              []BlockResponseTxVIn           `json:"vin"`
	VOut             []BlockResponseTxVOut          `json:"vout"`
	BlockHash        string                         `json:"blockHash"`
	BlockHeight      int64                          `json:"blockHeight"`
	Confirmations    int64                          `json:"confirmations"`
	BlockTime        int64                          `json:"blockTime"`
	Value            string                         `json:"value"`
	Fees             string                         `json:"fees"`
	TokenTransfers   []*TokenTransfers              `json:"tokenTransfers"`
	EthereumSpecific *BlockResponseEthereumSpecific `json:"ethereumSpecific"`
}

type BlockResponse struct {
	Page              int64             `json:"page"`
	TotalPages        int64             `json:"totalPages"`
	ItemsOnPage       int64             `json:"itemsOnPage"`
	Hash              string            `json:"hash"`
	PreviousBlockHash string            `json:"previousBlockHash"`
	NextBlockHash     string            `json:"nextBlockHash"`
	Height            int64             `json:"height"`
	Confirmations     int64             `json:"confirmations"`
	Size              int64             `json:"size"`
	Time              int64             `json:"time"`
	Version           int64             `json:"version"`
	MerkleRoot        string            `json:"merkleRoot"`
	Nonce             string            `json:"nonce"`
	Bits              string            `json:"bits"`
	Difficulty        string            `json:"difficulty"`
	TxCount           int64             `json:"txCount"`
	Txs               []BlockResponseTx `json:"txs"`
}

type BlockIndexResponse struct {
	BlockHash string `json:"blockHash"`
}

// ==========================================

// ========================================== //

type BlockBookWsNewBlockResponseBody struct {
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
}

type WsBlockResponse struct {
	Id   string                          `json:"id"`
	Data BlockBookWsNewBlockResponseBody `json:"data"`
}

// ==========================================

// ========================================== //

type WsResponse struct {
	Id   string              `json:"id"`
	Data TransactionResponse `json:"data"`
}

type WsSubscribeNewTransaction struct {
	Id     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// ==========================================

// ========================================== //

type BroadcastTransactionResponse struct {
	TxId string `json:"result"`
}

type Utxo struct {
	Txid          string `json:"txid"`
	Vout          uint32 `json:"vout"`
	Value         string `json:"value"`
	Height        int64  `json:"height"`
	Confirmations int64  `json:"confirmations"`
}
