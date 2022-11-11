package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	bitcoinWallet "github.com/Amirilidan78/bitcoin-wallet"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const privateKeyHex = "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
const utxoHash = "c77fe6125260df1702382ff89716dd5873051ffb95872d6d9407132ef52c4e84"
const utxoAmount = 67000
const toAddressHex = "tb1qtvnf9xcnyw34qrxc0aufqr34el7l4fec4fnknp"

func createWallet() *bitcoinWallet.BitcoinWallet {
	w, _ := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, privateKeyHex)
	return w
}

func walletPrivateAndPublicKey(wallet *bitcoinWallet.BitcoinWallet) (*btcec.PrivateKey, *btcec.PublicKey) {
	b, _ := wallet.PrivateKeyBytes()
	return btcec.PrivKeyFromBytes(b)
}

func getFromAddress(wallet *bitcoinWallet.BitcoinWallet) btcutil.Address {
	address, _ := btcutil.DecodeAddress(wallet.Address, &chaincfg.TestNet3Params)
	return address
}

func getToAddress() btcutil.Address {
	address, _ := btcutil.DecodeAddress(toAddressHex, &chaincfg.TestNet3Params)
	return address
}

func getRawTransaction(redeemTx *wire.MsgTx) (string, error) {

	var tx bytes.Buffer
	err := redeemTx.Serialize(&tx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tx.Bytes()), nil
}

func signTransaction(redeemTx *wire.MsgTx, priv *btcec.PrivateKey, fromAddress btcutil.Address, amount int64, outPoint *wire.OutPoint, redeemTxOut *wire.TxOut) *wire.MsgTx {

	subscript := priv.PubKey().SerializeCompressed()
	a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*outPoint: redeemTxOut,
	})

	txSigHashes := txscript.NewTxSigHashes(redeemTx, a)
	wit, err := txscript.WitnessSignature(redeemTx, txSigHashes, 0, amount, subscript, txscript.SigHashAll, priv, true)
	if err != nil {
		panic(err)
	}

	redeemTx.TxIn[0].Witness = wit
	return redeemTx
}

func main() {

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	amount := int64(utxoAmount) - 50

	wallet := createWallet()
	priv, _ := walletPrivateAndPublicKey(wallet)

	fromAddress := getFromAddress(wallet)

	toAddress := getToAddress()
	toAddressByte, err := txscript.PayToAddrScript(toAddress)
	if err != nil {
		panic(err)
	}

	hash, err := chainhash.NewHashFromStr(utxoHash)
	if err != nil {
		panic(err)
	}

	outPoint := wire.NewOutPoint(hash, 0)
	txIn := wire.NewTxIn(outPoint, nil, [][]byte{})
	redeemTx.AddTxIn(txIn)

	redeemTxOut := wire.NewTxOut(amount, toAddressByte)
	redeemTx.AddTxOut(redeemTxOut)

	redeemTx = signTransaction(redeemTx, priv, fromAddress, amount, outPoint, redeemTxOut)

	finalRawTx, err := getRawTransaction(redeemTx)
	fmt.Println(finalRawTx, err)

}

//
//func generate() string {
//	wallet := bitcoinWallet.GenerateBitcoinWallet(enums.TEST_NODE)
//	fmt.Println(wallet.PrivateKey)
//	fmt.Println(wallet.PublicKey)
//	fmt.Println(wallet.Address)
//	return wallet.PrivateKey
//}
//
//func create(priv string) {
//	wallet, err := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, priv)
//	fmt.Println(err)
//	fmt.Println(wallet.Address)
//}
