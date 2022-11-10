package blockBook

import (
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/Amirilidan78/bitcoin-wallet/httpClient"
	"strconv"
)

type HttpBlockBook interface {
	get(path string, res interface{}) error
	getHost() string
	GetStatus() (StatusResponse, error)
	GetBlockIndex() (BlockIndexResponse, error)
	GetBlock(hash string) (BlockResponse, error)
	GetAddress(address string) (AddressResponse, error)
	GetTransaction(txId string) (TransactionResponse, error)
	BroadcastTransaction(hex string) (BroadcastTransactionResponse, error)
	GetAddressUTXO(address string) ([]Utxo, error)
}

type httpBlockBook struct {
	node enums.Node
	hc   httpClient.HttpClient
}

func (b *httpBlockBook) getHost() string {
	return b.node.Http
}

func (b *httpBlockBook) get(path string, res interface{}) error {

	host := b.getHost()

	url := host + path

	err := b.hc.SimpleGet(url, res)

	return err
}

func (b *httpBlockBook) GetBlockIndex() (BlockIndexResponse, error) {

	res := BlockIndexResponse{}

	path := BlockIndexPath

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (b *httpBlockBook) BroadcastTransaction(hex string) (BroadcastTransactionResponse, error) {

	res := BroadcastTransactionResponse{}

	path := BroadcastPath + hex

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (b *httpBlockBook) GetBlock(hash string) (BlockResponse, error) {

	res := BlockResponse{}

	path := BlockPath + hash

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	if res.TotalPages > 1 {
		for i := 2; i <= int(res.TotalPages); i++ {

			res2, errSecond := b.getBlockPage(hash, strconv.Itoa(i))

			if errSecond != nil {
				return res, err
			}

			res.Txs = append(res.Txs, res2.Txs...)
		}
	}

	return res, nil
}

func (b *httpBlockBook) getBlockPage(hash string, page string) (BlockResponse, error) {

	res := BlockResponse{}

	path := BlockPath + hash + "?page=" + page

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (b *httpBlockBook) GetStatus() (StatusResponse, error) {

	res := StatusResponse{}

	path := StatusPath

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (b *httpBlockBook) GetAddress(address string) (AddressResponse, error) {

	res := AddressResponse{}

	path := AddressPath + address

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (b *httpBlockBook) GetTransaction(txId string) (TransactionResponse, error) {

	res := TransactionResponse{}

	path := TXPath + txId

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (b *httpBlockBook) GetAddressUTXO(address string) ([]Utxo, error) {

	res := make([]Utxo, 0)

	path := UTXOPath + address

	err := b.get(path, &res)

	if err != nil {
		return res, err
	}

	return res, nil
}

func NewHttpBlockBookService(node enums.Node) HttpBlockBook {
	return &httpBlockBook{node, httpClient.NewHttpClient()}
}
