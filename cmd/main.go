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

func main() {

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	amount := int64(67000)
	utxoHash := "7229ce852799857ea1834bc33d6e5214a4b6c7734933121fdbb87fd127f5e240"
	wallet, _ := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, privateKeyHex)

	b, _ := wallet.PrivateKeyBytes()
	priv, publicKey := btcec.PrivKeyFromBytes(b)

	address, err := btcutil.DecodeAddress(wallet.Address, &chaincfg.TestNet3Params)
	toAddress, err := btcutil.DecodeAddress("tb1qtvnf9xcnyw34qrxc0aufqr34el7l4fec4fnknp", &chaincfg.TestNet3Params)
	fmt.Println(priv.PubKey().SerializeCompressed(), "NNNNN")

	toAddressByte, err := txscript.PayToAddrScript(toAddress)
	fmt.Println(address, err)
	bldr := txscript.NewScriptBuilder()
	bldr.AddOp(txscript.OP_1)
	bldr.AddData(priv.PubKey().SerializeCompressed())
	sigScript, err := bldr.Script()
	fmt.Println(sigScript, err)

	hash, err := chainhash.NewHashFromStr(utxoHash)
	fmt.Println(hash, err)

	witness := [][]byte{
		publicKey.SerializeCompressed(),
	}

	outPoint := wire.NewOutPoint(hash, 10)
	txIn := wire.NewTxIn(outPoint, nil, witness)
	redeemTx.AddTxIn(txIn)

	redeemTxOut := wire.NewTxOut(amount, toAddressByte)
	redeemTx.AddTxOut(redeemTxOut)

	pkScripts, err := txscript.ComputePkScript(nil, witness)
	fmt.Println(err)

	finalRawTx, err := SignTx(priv, pkScripts.Script(), redeemTx)
	fmt.Println("==========")
	fmt.Println(finalRawTx, err)

}

func SignTx(privKey *btcec.PrivateKey, pkScript []byte, redeemTx *wire.MsgTx) (string, error) {

	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}

func generate() string {
	wallet := bitcoinWallet.GenerateBitcoinWallet(enums.TEST_NODE)
	fmt.Println(wallet.PrivateKey)
	fmt.Println(wallet.PublicKey)
	fmt.Println(wallet.Address)
	return wallet.PrivateKey
}

func create(priv string) {
	wallet, err := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, priv)
	fmt.Println(err)
	fmt.Println(wallet.Address)
}
