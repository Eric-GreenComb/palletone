package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/tokenengine"

	"github.com/Eric-GreenComb/palletone/bean"
)

// GenRawTransaction 创建RawTransaction
func GenRawTransaction(utxos bean.Utxos, from, to string, pay, amount uint64) (*modules.Transaction, error) {

	// 根据支付金额，获取花费的utxo数组和找零
	_takenUtxo, _change, err := SelectUtxoGreedy(utxos, pay)
	if err != nil {
		return nil, err
	}

	_mtx, err := GenMessageByUxto(_takenUtxo, from, to, amount, _change)
	if err != nil {
		return nil, err
	}

	return _mtx, nil
}

// GenRawTransactionEnCoding 创建RawTransaction Rlp encoding
func GenRawTransactionEnCoding(mtx *modules.Transaction) (string, error) {

	// 由于message里只有payment，直接对tx进行rlp
	_hexRlp, err := GenMtxRlpHex(mtx)
	if err != nil {
		return "", err
	}

	return _hexRlp, nil
}

// GenMtxRlpHex 构造交易编码
func GenMtxRlpHex(mtx *modules.Transaction) (string, error) {

	_bytetxjson, err := json.Marshal(mtx)
	if err != nil {
		return "", err
	}
	fmt.Println(string(_bytetxjson))
	mtxbt, err := rlp.EncodeToBytes(_bytetxjson)
	if err != nil {
		return "", err
	}
	mtxHex := "0x" + hex.EncodeToString(mtxbt)
	return mtxHex, nil

}

// GenRawTransactionEx 创建RawTransaction
func GenRawTransactionEx(utxos bean.Utxos, from, to string, pay, amount uint64) (*modules.Transaction, error) {

	// 根据支付金额，获取花费的utxo数组和找零
	_takenUtxo, _change, err := SelectUtxoGreedy(utxos, pay)
	if err != nil {
		return nil, err
	}

	_mtx, err := GenMessageByUxto(_takenUtxo, from, to, amount, _change)
	if err != nil {
		return nil, err
	}
	return _mtx, nil
}

// GenSignHash GenSignHash
func GenSignHash(from common.Address, mtx *modules.Transaction) (string, string, []bean.TxHash) {

	var _hashList []bean.TxHash
	var _hash bean.TxHash
	var _buffer bytes.Buffer

	var _ppscript []byte
	_ppscript = tokenengine.GenerateLockScript(from)

	// var _ppscript []byte
	// _ppscript, err := base64.StdEncoding.DecodeString("tmZvA5sAu5B2mCRIBfF4Pc59gEjO8PH6SRjMX59xncU=")
	// if err != nil {
	// 	fmt.Println("base64 decoding error", err.Error())
	// }

	_mtxtmp := mtx
	for _msgindex, _msg := range _mtxtmp.TxMessages {
		_payload, _ok := _msg.Payload.(*modules.PaymentPayload)
		if !_ok {
			continue
		}
		for _inputindex, _input := range _payload.Inputs {
			hashforsign, err := tokenengine.CalcSignatureHash(_mtxtmp, tokenengine.SigHashAll, _msgindex, _inputindex, _ppscript)
			if err != nil {
				fmt.Println("sign error", err.Error())
				continue
			}

			_buffer.WriteString(_input.PreviousOutPoint.TxHash.String())
			_buffer.WriteString(",")

			_hash.Hash = "0x" + hex.EncodeToString(hashforsign)

			_hashList = append(_hashList, _hash)
		}
	}

	_buf := make([]byte, _buffer.Len()-1)
	_buffer.Read(_buf)

	return base64.StdEncoding.EncodeToString(_ppscript), string(_buf), _hashList
}

// GenMessageByUxto GenMessageByUxto
func GenMessageByUxto(utxos bean.Utxos, from, to string, amount, change uint64) (*modules.Transaction, error) {

	// 构造PaymentPayload
	pload := new(modules.PaymentPayload)

	pload.LockTime = 0

	_from, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}

	var _ppscript []byte
	// _ppscript, err := base64.StdEncoding.DecodeString("tmZvA5sAu5B2mCRIBfF4Pc59gEjO8PH6SRjMX59xncU=")
	// if err != nil {
	// 	fmt.Println("base64 decoding error", err.Error())
	// }
	_ppscript = tokenengine.GenerateLockScript(_from)

	// 构造Input,Output
	for _, _utxo := range utxos {

		// 构造 payload input
		txHash := common.HexToHash(_utxo.TxID)
		prevOut := modules.NewOutPoint(txHash, _utxo.MessageIndex, _utxo.OutIndex)
		txInput := modules.NewTxIn(prevOut, _ppscript)
		pload.AddTxIn(txInput)

	}

	// 构造 payload output，默认PTNCOIN
	_toOut, err := GenOutput(to, amount)
	if err != nil {
		return nil, errors.New("构建To output error:" + err.Error())
	}
	pload.AddTxOut(_toOut)

	// 如果找零大于0，则构造output，并加入payload output数组
	if change > 0 {
		_changeOut, err := GenOutput(from, change)
		if err != nil {
			return nil, errors.New("构建找零 output error:" + err.Error())
		}
		pload.AddTxOut(_changeOut)
	}

	mtx := &modules.Transaction{
		TxMessages: make([]*modules.Message, 0),
	}
	mtx.TxMessages = append(mtx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, pload))

	return mtx, nil
}

// GenOutput 构建output
func GenOutput(addr string, amount uint64) (*modules.Output, error) {

	_addr, err := common.StringToAddress(addr)
	if err != nil {
		return nil, errors.New("Invalid address or key")
	}

	switch _addr.GetType() {
	case common.PublicKeyHash:
	case common.ScriptHash:
	case common.ContractHash:
	default:
		return nil, errors.New("Invalid address or key")
	}

	// Create a new script which pays to the provided address.
	_pkScript := tokenengine.GenerateLockScript(_addr)
	_assetID := modules.PTNCOIN

	_output := modules.NewTxOut(amount, _pkScript, _assetID.ToAsset())

	return _output, nil

}

// SelectUtxoGreedy 贪吃算法获取utxos和找零
func SelectUtxoGreedy(utxos bean.Utxos, amount uint64) (bean.Utxos, uint64, error) {
	var greaters bean.Utxos
	var lessers bean.Utxos
	var takenLutxo bean.Utxos
	var takenGutxo bean.Utxos
	var accum uint64
	var change uint64
	logPickedAmt := ""
	accum = 0
	// 根据需要支付的金额分组utxo，分为大于等于和小于两个组
	for _, utxo := range utxos {
		if utxo.Amount >= amount {
			greaters = append(greaters, utxo)
		}
		if utxo.Amount < amount {
			lessers = append(lessers, utxo)
		}
	}

	// 判断小于的组是否够支付
	if len(lessers) > 0 {
		// 排序小于的组，从小到大累计，看是否够支付，够的话返回可以支付的utxo数组，和找零
		sort.Sort(bean.Utxos(lessers))
		for _, utxo := range lessers {
			accum += utxo.Amount
			logPickedAmt += fmt.Sprintf("%d,", utxo.Amount)
			takenLutxo = append(takenLutxo, utxo)
			if accum >= amount {
				change = accum - amount
				return takenLutxo, change, nil
			}
		}
	}

	// 如果不够支付，则返回错误
	fmt.Println(accum, amount)
	if accum < amount && len(greaters) == 0 {
		return nil, 0, errors.New("Amount Not Enough to pay")
	}

	// 找到最小的可支付utxo，返回数组（一个值）和找零
	var minGreater bean.TxUtxo
	minGreater = bean.FindMin(greaters)
	change = minGreater.Amount - amount
	logPickedAmt = fmt.Sprintf("%d,", minGreater.Amount)
	takenGutxo = append(takenGutxo, minGreater)

	return takenGutxo, change, nil
}
