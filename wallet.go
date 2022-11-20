package bitcoinWallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/Amirilidan78/bitcoin-wallet/blockBook"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
)

type BitcoinWallet struct {
	Node       enums.Node
	Address    string
	PrivateKey string
	PublicKey  string
	blockBook  blockBook.HttpBlockBook
}

// generating

func GenerateBitcoinWallet(node enums.Node) *BitcoinWallet {

	privateKey, _ := generatePrivateKey()
	privateKeyHex := convertPrivateKeyToHex(privateKey)

	publicKey, _ := getPublicKeyFromPrivateKey(privateKey)
	publicKeyHex := convertPublicKeyToHex(publicKey)

	address, _ := getAddressFromPrivateKey(node, privateKey)

	return &BitcoinWallet{
		Node:       node,
		Address:    address,
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
		blockBook:  blockBook.NewHttpBlockBookService(node),
	}
}

func CreateBitcoinWallet(node enums.Node, privateKeyHex string) (*BitcoinWallet, error) {

	privateKey, err := privateKeyFromHex(privateKeyHex)
	if err != nil {
		return nil, err
	}

	publicKey, err := getPublicKeyFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	publicKeyHex := convertPublicKeyToHex(publicKey)

	address, err := getAddressFromPrivateKey(node, privateKey)
	if err != nil {
		return nil, err
	}

	return &BitcoinWallet{
		Node:       node,
		Address:    address,
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
		blockBook:  blockBook.NewHttpBlockBookService(node),
	}, nil
}

// struct functions

func (bw *BitcoinWallet) PrivateKeyRCDSA() (*ecdsa.PrivateKey, error) {
	return privateKeyFromHex(bw.PrivateKey)
}

func (bw *BitcoinWallet) PrivateKeyBytes() ([]byte, error) {

	priv, err := bw.PrivateKeyRCDSA()
	if err != nil {
		return []byte{}, err
	}

	return crypto.FromECDSA(priv), nil
}

// private key

func generatePrivateKey() (*ecdsa.PrivateKey, error) {

	return crypto.GenerateKey()
}

func convertPrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {

	privateKeyBytes := crypto.FromECDSA(privateKey)

	return hexutil.Encode(privateKeyBytes)[2:]
}

func privateKeyFromHex(hex string) (*ecdsa.PrivateKey, error) {

	return crypto.HexToECDSA(hex)
}

// public key

func getPublicKeyFromPrivateKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error in getting public key")
	}

	return publicKeyECDSA, nil
}

func convertPublicKeyToHex(publicKey *ecdsa.PublicKey) string {

	privateKeyBytes := crypto.FromECDSAPub(publicKey)

	return hexutil.Encode(privateKeyBytes)[2:]
}

// address

func getAddressFromPrivateKey(node enums.Node, privateKey *ecdsa.PrivateKey) (string, error) {

	chainConfig := &chaincfg.MainNetParams
	if node.Test {
		chainConfig = &chaincfg.TestNet3Params
	}

	temp := crypto.FromECDSA(privateKey)
	priv, _ := btcec.PrivKeyFromBytes(temp)
	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(priv.PubKey().SerializeCompressed()), chainConfig)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return addr.EncodeAddress(), nil
}

// balance

func (bw *BitcoinWallet) Balance() (int64, error) {

	res, err := bw.blockBook.GetAddress(bw.Address)
	if err != nil {
		return 0, err
	}

	balance, err := strconv.Atoi(res.Balance)
	if err != nil {
		return 0, err
	}

	return int64(balance), nil
}

// transactions

func (bw *BitcoinWallet) UTXOs() ([]blockBook.Utxo, error) {

	utxos, err := bw.blockBook.GetAddressUTXO(bw.Address)
	if err != nil {
		return utxos, err
	}

	var res []blockBook.Utxo

	for _, utxo := range utxos {
		if utxo.Confirmations > 2 {
			res = append(res, utxo)
		}
	}

	return res, nil
}

func (bw *BitcoinWallet) TxIds() ([]string, error) {

	var txIds []string

	res, err := bw.blockBook.GetAddress(bw.Address)
	if err != nil {
		return txIds, err
	}

	return res.TxIds, nil
}
