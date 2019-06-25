
签名函数

tokenengine.CalcSignatureHash(mtx, tokenengine.SigHashAll, int(_utxo.MessageIndex), int(_utxo.OutIndex), nil)

//为钱包计算要签名某个Input对应的Hash
func CalcSignatureHash(tx *modules.Transaction, hashType uint32, msgIdx, inputIdx int, lockOrRedeemScript []byte) ([]byte, error) {
	acc := &account{}
	return txscript.CalcSignatureHash(lockOrRedeemScript, txscript.SigHashType(hashType), tx, msgIdx, inputIdx, acc)
}

此处为txscript.CalcSignatureHash函数

// calcSignatureHash will, given a script and hash type for the current script
// engine instance, calculate the signature hash to be used for signing and
// verification.
func calcSignatureHash(script []parsedOpcode, hashType SigHashType, tx *modules.Transaction, msgIdx, idx int, crypto ICrypto) []byte {
	pay := tx.TxMessages[msgIdx].Payload.(*modules.PaymentPayload)
	// The SigHashSingle signature type signs only the corresponding input
	// and output (the output with the same index number as the input).
	//
	// Since transactions can have more inputs than outputs, this means it
	// is improper to use SigHashSingle on input indices that don't have a
	// corresponding output.
	//
	// A bug in the original Satoshi client implementation means specifying
	// an index that is out of range results in a signature hash of 1 (as a
	// uint256 little endian).  The original intent appeared to be to
	// indicate failure, but unfortunately, it was never checked and thus is
	// treated as the actual signature hash.  This buggy behavior is now
	// part of the consensus and a hard fork would be required to fix it.
	//
	// Due to this, care must be taken by software that creates transactions
	// which make use of SigHashSingle because it can lead to an extremely
	// dangerous situation where the invalid inputs will end up signing a
	// hash of 1.  This in turn presents an opportunity for attackers to
	// cleverly construct transactions which can steal those coins provided
	// they can reuse signatures.
	if hashType&sigHashMask == SigHashSingle && idx >= len(pay.Outputs) {
		var hash common.Hash
		hash[0] = 0x01
		return hash[:]
	}

	// Remove all instances of OP_CODESEPARATOR from the script.
	script = removeOpcode(script, OP_CODESEPARATOR)

	// Make a deep copy of the transaction, zeroing out the script for all
	// inputs that are not currently being processed.

	txCopy := tx.Clone()
	payCopy := txCopy.TxMessages[msgIdx].Payload.(*modules.PaymentPayload)
	requestIndex := tx.GetRequestMsgIndex()
	isInResult := false
	if msgIdx > requestIndex && requestIndex != -1 {
		isInResult = true
	}

	for mIdx, mCopy := range txCopy.TxMessages {

====只处理APP_PAYMENT，其他保持原值，并且SignatureScript=nil

		if mCopy.App == modules.APP_PAYMENT {
			pay := txCopy.TxMessages[mIdx].Payload.(*modules.PaymentPayload)
			if isInResult && mIdx < requestIndex {
				continue // 对于请求部分的Payment，不做任何处理
			}
			for i := range pay.Inputs {
				//Devin: for contract payout, remove all lockscript
				if i == idx && mIdx == msgIdx && hashType != SigHashRaw {
					// UnparseScript cannot fail here because removeOpcode
					// above only returns a valid script.
					sigScript, _ := unparseScript(script)
					pay.Inputs[idx].SignatureScript = sigScript
				} else {
					pay.Inputs[i].SignatureScript = nil
				}
			}
		}
	}

====签名类型为SigHashAll，不处理

	switch hashType & sigHashMask {
	case SigHashNone:
		payCopy.Outputs = payCopy.Outputs[0:0] // Empty slice.
		//for i := range payCopy.Inputs {
		//	if i != idx {
		//		payCopy.Inputs[i].Sequence = 0
		//	}
		//}

	case SigHashSingle:
		// Resize output array to up to and including requested index.
		payCopy.Outputs = payCopy.Outputs[:idx+1]

		// All but current output get zeroed out.
		for i := 0; i < idx; i++ {
			payCopy.Outputs[i].Value = 0
			payCopy.Outputs[i].PkScript = nil
		}

		// Sequence on all other inputs is 0, too.
		//for i := range payCopy.TxIn {
		//	if i != idx {
		//		payCopy.TxIn[i].Sequence = 0
		//	}
		//}

	default:
		// Consensus treats undefined hashtypes like normal SigHashAll
		// for purposes of hash generation.
		fallthrough
	case SigHashOld:
		fallthrough
	case SigHashAll:
		// Nothing special here.
	}
	if hashType&SigHashAnyOneCanPay != 0 {
		payCopy.Inputs = payCopy.Inputs[idx : idx+1]
		idx = 0
	}

	// The final hash is the double sha256 of both the serialized modified
	// transaction and the hash type (encoded as a 4-byte little-endian
	// value) appended.
	//var wbuf bytes.Buffer
	//payCopy.Serialize(&wbuf)
	//binary.Write(&wbuf, binary.LittleEndian, hashType)
	//return wire.DoubleSha256(wbuf.Bytes())


====此部分有效
====rlp为以太的rlp编码规则
====Hash为crypto.Keccak256(msg)

	data, err := rlp.EncodeToBytes(&txCopy)
	if err != nil {
		log.Error("Rlp encode tx error:" + err.Error())
	}
	hash, _ := crypto.Hash(data)
	return hash




}