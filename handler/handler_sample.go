package handler

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/ptnjson"
	// "github.com/palletone/go-palletone/ptnjson/walletjson"
	"github.com/palletone/go-palletone/tokenengine"

	"github.com/Eric-GreenComb/palletone/bean"
)

// CreateRawTransaction CreateRawTransaction
func CreateRawTransaction(utxos bean.Utxos, from, to string, amount uint64) error {
	amounts := []ptnjson.AddressAmt{}
	amounts = append(amounts, ptnjson.AddressAmt{to, ptnjson.Dao2Ptn(amount)})
	if len(amounts) == 0 {
		return fmt.Errorf("amounts is invalid")
	}

	_takenUtxo, _change, err := SelectUtxoGreedy(utxos, amount)
	if err != nil {
		return err
	}

	var inputs []ptnjson.TransactionInput
	var input ptnjson.TransactionInput
	for _, u := range _takenUtxo {
		input.Txid = u.TxID
		input.MessageIndex = u.MessageIndex
		input.Vout = u.OutIndex
		inputs = append(inputs, input)
	}

	if _change > 0 {
		amounts = append(amounts, ptnjson.AddressAmt{from, ptnjson.Dao2Ptn(_change)})
	}

	var LockTime int64
	LockTime = 0
	arg := ptnjson.NewCreateRawTransactionCmd(inputs, amounts, &LockTime)

	fmt.Println(arg)

	result, err := WalletCreateTransaction(arg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(result)

	return nil
}

// const 常量
const (
	MaxTxInSequenceNum uint32 = 0xffffffff
)

// WalletCreateTransaction WalletCreateTransaction
func WalletCreateTransaction(c *ptnjson.CreateRawTransactionCmd) (string, error) {

	// Validate the locktime, if given.
	if c.LockTime != nil &&
		(*c.LockTime < 0 || *c.LockTime > int64(MaxTxInSequenceNum)) {
		return "", errors.New("Locktime out of range")
	}
	// Add all transaction inputs to a new transaction after performing
	// some validity checks.
	//先构造PaymentPayload结构，再组装成Transaction结构
	pload := new(modules.PaymentPayload)
	// var inputjson []walletjson.InputJson
	for _, input := range c.Inputs {
		txHash := common.HexToHash(input.Txid)
		// inputjson = append(inputjson, walletjson.InputJson{TxHash: input.Txid, MessageIndex: input.MessageIndex, OutIndex: input.Vout, HashForSign: "", Signature: ""})
		prevOut := modules.NewOutPoint(txHash, input.MessageIndex, input.Vout)
		txInput := modules.NewTxIn(prevOut, []byte{})
		pload.AddTxIn(txInput)
	}
	// var OutputJSON []walletjson.OutputJson
	// Add all transaction outputs to the transaction after performing
	//	// some validity checks.
	//	//only support mainnet
	//	var params *chaincfg.Params
	var ppscript []byte
	for _, addramt := range c.Amounts {
		encodedAddr := addramt.Address
		ptnAmt := addramt.Amount
		amount := ptnjson.Ptn2Dao(ptnAmt)
		// Ensure amount is in the valid range for monetary amounts.
		if amount <= 0 /*|| amount > ptnjson.MaxDao*/ {
			return "", errors.New("Invalid amount")
		}
		addr, err := common.StringToAddress(encodedAddr)
		if err != nil {
			return "", errors.New("Invalid address or key")
		}
		switch addr.GetType() {
		case common.PublicKeyHash:
		case common.ScriptHash:
		case common.ContractHash:
			//case *ptnjson.AddressPubKeyHash:
			//case *ptnjson.AddressScriptHash:
		default:
			return "", &ptnjson.RPCError{
				Code:    ptnjson.ErrRPCInvalidAddressOrKey,
				Message: "Invalid address or key",
			}
		}
		// Create a new script which pays to the provided address.
		pkScript := tokenengine.GenerateLockScript(addr)
		ppscript = pkScript
		// Convert the amount to satoshi.
		dao := ptnjson.Ptn2Dao(ptnAmt)
		if err != nil {
			context := "Failed to convert amount"
			return "", errors.New(context)
		}
		_assetID := modules.PTNCOIN
		txOut := modules.NewTxOut(uint64(dao), pkScript, _assetID.ToAsset())
		pload.AddTxOut(txOut)
		// OutputJSON = append(OutputJSON, walletjson.OutputJson{Amount: uint64(dao), Asset: _assetID.String(), ToAddress: addr.String()})
	}
	//	// Set the Locktime, if given.
	if c.LockTime != nil {
		pload.LockTime = uint32(*c.LockTime)
	}
	//	// Return the serialized and hex-encoded transaction.  Note that this
	//	// is intentionally not directly returning because the first return
	//	// value is a string and it would result in returning an empty string to
	//	// the client instead of nothing (nil) in the case of an error.

	mtx := &modules.Transaction{
		TxMessages: make([]*modules.Message, 0),
	}
	mtx.TxMessages = append(mtx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, pload))
	//mtx.TxHash = mtx.Hash()
	//sign mtx
	mtxtmp := mtx
	for msgindex, msg := range mtxtmp.TxMessages {
		payload, ok := msg.Payload.(*modules.PaymentPayload)
		if ok == false {
			continue
		}
		for inputindex := range payload.Inputs {
			hashforsign, err := tokenengine.CalcSignatureHash(mtxtmp, tokenengine.SigHashAll, msgindex, inputindex, ppscript)
			if err != nil {
				return "", err
			}
			payloadtmp := mtx.TxMessages[msgindex].Payload.(*modules.PaymentPayload)
			payloadtmp.Inputs[inputindex].SignatureScript = hashforsign
		}
	}

	bytetxjson, err := json.Marshal(mtx)
	if err != nil {
		return "", err
	}
	mtxbt, err := rlp.EncodeToBytes(bytetxjson)
	if err != nil {
		return "", err
	}
	//log.Debugf("payload input outpoint:%s", pload.Input[0].PreviousOutPoint.TxHash.String())
	mtxHex := hex.EncodeToString(mtxbt)
	return mtxHex, nil
	//return string(bytetxjson), nil
}
