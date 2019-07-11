package handler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gin-gonic/gin"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/dag/modules"

	"github.com/Eric-GreenComb/palletone/bean"
)

// GetRawTx GetRawTx
func GetRawTx(c *gin.Context) {

	var _params bean.TxParams

	c.BindJSON(&_params)

	var _ret bean.TxReturn
	_ret.Status = "false"

	// 付给to的金额
	_amountTx, err := strconv.ParseUint(_params.Amount, 10, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// 交易费
	_fee, err := strconv.ParseUint(_params.Fee, 10, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// from需要支付的金额，包括to的金额+交易费
	_pay := _amountTx + _fee

	_mtx, err := GenRawTransaction(_params.Utxos, _params.SendAddr, _params.RecvAddr, _pay, _amountTx)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	c.JSON(http.StatusOK, _mtx)
}

// GetRawTxEncoding GetRawTxEncoding
func GetRawTxEncoding(c *gin.Context) {

	var _params bean.TxParams

	c.BindJSON(&_params)

	var _ret bean.TxReturn
	_ret.Status = "false"

	// 付给to的金额
	_amountTx, err := strconv.ParseUint(_params.Amount, 10, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// 交易费
	_fee, err := strconv.ParseUint(_params.Fee, 10, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// from需要支付的金额，包括to的金额+交易费
	_pay := _amountTx + _fee

	_mtx, err := GenRawTransaction(_params.Utxos, _params.SendAddr, _params.RecvAddr, _pay, _amountTx)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_hex, err := GenRawTransactionEnCoding(_mtx)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	c.JSON(http.StatusOK, _hex)
}

// GetTxHash GetTxHash
func GetTxHash(c *gin.Context) {

	var _params bean.TxParams

	c.BindJSON(&_params)

	var _ret bean.TxReturn
	_ret.Status = "false"

	// 付给to的金额
	_amountTx, err := strconv.ParseUint(_params.Amount, 10, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// 交易费
	_fee, err := strconv.ParseUint(_params.Fee, 10, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// from需要支付的金额，包括to的金额+交易费
	_pay := _amountTx + _fee

	_mtx, err := GenRawTransactionEx(_params.Utxos, _params.SendAddr, _params.RecvAddr, _pay, _amountTx)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_addr, err := common.StringToAddress(_params.SendAddr)
	if err != nil {
		c.JSON(http.StatusOK, _ret)
		return
	}

	_signatureScript, _tx, _hashList := GenSignHash(_addr, _mtx)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_ret.Status = "true"
	_ret.SignatureScript = _signatureScript
	_ret.Tx = _tx
	_ret.HashList = _hashList

	c.JSON(http.StatusOK, _ret)
}

// GetRawTxDecoding GetRawTxDecoding
func GetRawTxDecoding(c *gin.Context) {

	var _params bean.TxParams

	c.BindJSON(&_params)

	newTx := &modules.Transaction{}

	_bytes, err := hex.DecodeString(_params.Rlp)
	if err != nil {
		fmt.Println("hex decoding error:", err.Error())
	}

	err = rlp.DecodeBytes(_bytes, newTx)
	if err != nil {
		fmt.Println("rlp decoding error:", err.Error())
	}

	c.JSON(http.StatusOK, newTx)
}
