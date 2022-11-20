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
	"log"
)

const privateKeyHex = "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
const utxoHash = "e6160c52401949139688623ce33a6290eed43d8d564d6e16c38006c4dc28f4a8"
const utxoAmount = 57821
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

	addrPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(priv.PubKey().SerializeCompressed()), &chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}
	subscript, err := txscript.PayToAddrScript(addrPubKey)
	if err != nil {
		panic(err)
	}

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

func ccc() {
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	amount := int64(utxoAmount)

	wallet := createWallet()
	priv, _ := walletPrivateAndPublicKey(wallet)

	fromAddress := getFromAddress(wallet)
	fmt.Println(priv)
	fmt.Println("========")

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
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	redeemTxOut := wire.NewTxOut(amount, toAddressByte)
	redeemTx.AddTxOut(redeemTxOut)

	redeemTx = signTransaction(redeemTx, priv, fromAddress, amount, outPoint, redeemTxOut)

	fmt.Println(*redeemTx.TxOut[0])
	fmt.Println(*redeemTx.TxIn[0])
	finalRawTx, err := getRawTransaction(redeemTx)
	fmt.Println(finalRawTx, err)
}

func main() {

	createSignedP2wkhTx()

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

func createSignedP2wkhTx() {

	spendAddrStr := "tb1qppv790u4dz48ctnk3p7ss7fmspckagp3wrfyp0"
	destAddrStr := "tb1q9dkhf8vxlvujxjmnslxsv97nseg9pjmxqsku3v"
	chain := &chaincfg.TestNet3Params
	txHash := "3a65dbebef06449f551c13541c9e2d6e2334aaaf133165705e3630745dbedef0"
	position := 1
	txAll := int64(1258855)
	txAmount := int64(8700)
	back := int64(1250000)

	wallet := createWallet()
	priv, _ := walletPrivateAndPublicKey(wallet)

	spendAddr, err := btcutil.DecodeAddress(spendAddrStr, chain)
	if err != nil {
		log.Println("DecodeAddress spendAddr err", err)
		return
	}

	destAddr, err := btcutil.DecodeAddress(destAddrStr, chain)
	if err != nil {
		log.Println("DecodeAddress destAddrStr err", err)
		return
	}

	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
	if err != nil {
		log.Println("spendAddr PayToAddrScript err", err)
		return
	}

	addrPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(priv.PubKey().SerializeCompressed()), chain)
	fmt.Println("===========")
	fmt.Println(addrPubKey)
	witnessProgram, err := txscript.PayToAddrScript(addrPubKey)
	log.Println(hex.EncodeToString(spenderAddrByte), hex.EncodeToString(witnessProgram))
	// either witnessProgram or spenderAddrByte works

	destAddrByte, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		log.Println("destAddr PayToAddrScript err", err)
		return
	}

	utxoHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Println("NewHashFromStr err", err)
		return
	}

	outPoint := wire.NewOutPoint(utxoHash, uint32(position))

	redeemTx := wire.NewMsgTx(2)

	txIn := wire.NewTxIn(outPoint, nil, [][]byte{})
	txIn.Sequence = txIn.Sequence - 2
	redeemTx.AddTxIn(txIn)

	redeemTxOut0 := wire.NewTxOut(txAmount, destAddrByte)
	redeemTxOut1 := wire.NewTxOut(back, spenderAddrByte)

	redeemTx.AddTxOut(redeemTxOut0)
	redeemTx.AddTxOut(redeemTxOut1)
	redeemTx.LockTime = 2407372

	if err != nil {
		log.Println("DecodeString pkScript err", err)
		return
	}

	a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*outPoint: {},
	})
	sigHashes := txscript.NewTxSigHashes(redeemTx, a)

	signature, err := txscript.WitnessSignature(redeemTx, sigHashes, 0, txAll, spenderAddrByte, txscript.SigHashAll, priv, true)
	if err != nil {
		log.Println("WitnessSignature err", err)
		return
	}
	redeemTx.TxIn[0].Witness = signature

	var signedTx bytes.Buffer
	err = redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	log.Println(hexSignedTx)

}
