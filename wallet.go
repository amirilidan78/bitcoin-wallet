package bitcoinWallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type BitcoinWallet struct {
	Node       enums.Node
	Address    string
	PrivateKey string
	PublicKey  string
}

// generating

func GenerateBitcoinWallet(node enums.Node) *BitcoinWallet {

	privateKey, _ := generatePrivateKey()
	privateKeyHex := convertPrivateKeyToHex(privateKey)

	publicKey, _ := getPublicKeyFromPrivateKey(privateKey)
	publicKeyHex := convertPublicKeyToHex(publicKey)

	address, _ := getAddressFromPublicKey(node, publicKey)

	return &BitcoinWallet{
		Node:       node,
		Address:    address,
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
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

	address, err := getAddressFromPublicKey(node, publicKey)
	if err != nil {
		return nil, err
	}

	return &BitcoinWallet{
		Node:       node,
		Address:    address,
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
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

func getAddressFromPublicKey(node enums.Node, publicKey *ecdsa.PublicKey) (string, error) {

	privateKeyBytes := crypto.FromECDSAPub(publicKey)

	pkHash := btcutil.Hash160(privateKeyBytes)

	chainConfig := &chaincfg.MainNetParams
	if node.Test {
		chainConfig = &chaincfg.TestNet3Params
	}

	addr, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, chainConfig)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return addr.EncodeAddress(), nil
}
