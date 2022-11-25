package main

//func awm() {
//
//	privKey := "cS9Zef6XdN3jHTFJFSsyJAtmDgCCdnygyVUJsLoyB8neuwhidUNJ"
//	spendAddrStr := "tb1qppv790u4dz48ctnk3p7ss7fmspckagp3wrfyp0"
//	destAddrStr := "tb1q9dkhf8vxlvujxjmnslxsv97nseg9pjmxqsku3v"
//	chain := &chaincfg.TestNet3Params
//	txHash := "4de74b4af672742f331eb3e2712ab54eca9a49d7e34d132b5dd39b98186057d7"
//	position := 1
//
//	txAll := int64(1200000)
//	txAmount := int64(8700)
//	fee := int64(10000)
//	back := txAll - txAmount - fee
//
//	/**
//	Output of c996d2c6cc0794a8dda2c421b211d5f7bda3a2d3e6c1d555a2a2718c2df4696c
//	is at position/index 1
//
//	The value below was gotten from electrum using `electrumt gettransaction $txHash`
//	"outputs": [
//			{
//				"scriptpubkey": "00142b6d749d86fb39234b7387cd0617d3865050cb66",
//				"address": "tb1q9dkhf8vxlvujxjmnslxsv97nseg9pjmxqsku3v",
//				"value_sats": 2000
//			},
//			{
//				"scriptpubkey": "0014160976523eb1c345b9d91025f7e2b98c95b84669",
//				"address": "tb1qzcyhv537k8p5twwezqjl0c4e3j2ms3nfe3kdgt",
//				"value_sats": 7800
//			}
//		]
//	**/
//
//	spendAddr, err := btcutil.DecodeAddress(spendAddrStr, chain)
//	if err != nil {
//		log.Println("DecodeAddress spendAddr err", err)
//		return
//	}
//
//	destAddr, err := btcutil.DecodeAddress(destAddrStr, chain)
//	if err != nil {
//		log.Println("DecodeAddress destAddrStr err", err)
//		return
//	}
//
//	wif, err := btcutil.DecodeWIF(privKey)
//	if err != nil {
//		log.Println("wif err", err)
//		return
//	}
//
//	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
//	if err != nil {
//		log.Println("spendAddr PayToAddrScript err", err)
//		return
//	}
//
//	addrPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed()), chain)
//	witnessProgram, err := txscript.PayToAddrScript(addrPubKey)
//	log.Println(hex.EncodeToString(spenderAddrByte), hex.EncodeToString(witnessProgram))
//	// either witnessProgram or spenderAddrByte works
//
//	destAddrByte, err := txscript.PayToAddrScript(destAddr)
//	if err != nil {
//		log.Println("destAddr PayToAddrScript err", err)
//		return
//	}
//
//	utxoHash, err := chainhash.NewHashFromStr(txHash)
//	if err != nil {
//		log.Println("NewHashFromStr err", err)
//		return
//	}
//
//	outPoint := wire.NewOutPoint(utxoHash, uint32(position))
//
//	redeemTx := wire.NewMsgTx(2)
//
//	txIn := wire.NewTxIn(outPoint, nil, [][]byte{})
//	txIn.Sequence = txIn.Sequence - 2
//	redeemTx.AddTxIn(txIn)
//
//	redeemTxOut0 := wire.NewTxOut(txAmount, destAddrByte)
//	redeemTxOut1 := wire.NewTxOut(back, spenderAddrByte)
//
//	redeemTx.AddTxOut(redeemTxOut0)
//	redeemTx.AddTxOut(redeemTxOut1)
//	redeemTx.LockTime = 2407372
//
//	if err != nil {
//		log.Println("DecodeString pkScript err", err)
//		return
//	}
//
//	a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
//		*outPoint: {},
//	})
//	sigHashes := txscript.NewTxSigHashes(redeemTx, a)
//
//	signature, err := txscript.WitnessSignature(redeemTx, sigHashes, 0, txAll, spenderAddrByte, txscript.SigHashAll, wif.PrivKey, true)
//	if err != nil {
//		log.Println("WitnessSignature err", err)
//		return
//	}
//	redeemTx.TxIn[0].Witness = signature
//
//	var signedTx bytes.Buffer
//	redeemTx.Serialize(&signedTx)
//
//	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
//	log.Println(hexSignedTx)
//
//}
