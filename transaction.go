package bitcoinWallet

//func getUTXOs(node enums.Node, addressHex string) (string, int64, string, error) {
//
//
//	btcutil.WIF{
//		PrivKey:        nil,
//		CompressPubKey: false,
//	}
//	address, err := btcutil.DecodeAddress(addressHex, &chaincfg.TestNet3Params)
//	if err != nil {
//		return "", 0, "", err
//	}
//
//	utxos, err := blockBook.NewHttpBlockBookService(node).GetAddressUTXO(addressHex)
//	if err != nil {
//		return "", 0, "", err
//	}
//
//	var previousTxid string = "16688d2946c3e029ca91ce730109994c2bcafb859d580a6f7c820fb2bb5b6afc"
//	var balance int64 = 62000
//	var pubKeyScript string = "76a91455d5e92958a8b06b4ff15cd2dd3d254f375e98db88ac"
//	return previousTxid, balance, pubKeyScript, nil
//}
//
//func createTransaction(node enums.Node, fromAddressHex string, privateKeyHex string, toAddressHex string, amountInSatoshi int64) (*wire.MsgTx, error) {
//
//	msg := wire.NewMsgTx(wire.TxVersion)
//
//	wire.NewTxIn()
//
//	return wire.NewMsgTx(wire.TxVersion), nil
//}
