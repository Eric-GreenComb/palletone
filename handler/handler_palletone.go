package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Eric-GreenComb/palletone/bean"
)

// GetRawTx GetRawTx
func GetRawTx(c *gin.Context) {

	var _params bean.TxParams

	c.BindJSON(&_params)

	var _ret bean.TxReturn
	_ret.Status = "false"

	_amountTx, err := strconv.ParseFloat(_params.Amount, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_fee, err := strconv.ParseFloat(_params.Fee, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// 付给to的金额
	_amount := int64(_amountTx * 100000000)
	if _amount <= 0 {
		_ret.Message = "付款金额小于零"
		c.JSON(http.StatusOK, _ret)
		return
	}

	// from需要支付的金额，包括to的金额+交易费
	_pay := int64((_amountTx + _fee) * 100000000)
	if _pay <= 0 {
		_ret.Message = "支付金额小于零"
		c.JSON(http.StatusOK, _ret)
		return
	}

	_mtx, err := GenRawTransaction(_params.Utxos, _params.SendAddr, _params.RecvAddr, uint64(_pay), uint64(_amount))
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

	_amountTx, err := strconv.ParseFloat(_params.Amount, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_fee, err := strconv.ParseFloat(_params.Fee, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// 付给to的金额
	_amount := int64(_amountTx * 100000000)
	if _amount <= 0 {
		_ret.Message = "付款金额小于零"
		c.JSON(http.StatusOK, _ret)
		return
	}

	// from需要支付的金额，包括to的金额+交易费
	_pay := int64((_amountTx + _fee) * 100000000)
	if _pay <= 0 {
		_ret.Message = "支付金额小于零"
		c.JSON(http.StatusOK, _ret)
		return
	}

	_mtx, err := GenRawTransaction(_params.Utxos, _params.SendAddr, _params.RecvAddr, uint64(_pay), uint64(_amount))
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

	_amountTx, err := strconv.ParseFloat(_params.Amount, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_fee, err := strconv.ParseFloat(_params.Fee, 64)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	// 付给to的金额
	_amount := int64(_amountTx * 100000000)
	if _amount <= 0 {
		_ret.Message = "付款金额小于零"
		c.JSON(http.StatusOK, _ret)
		return
	}

	// from需要支付的金额，包括to的金额+交易费
	_pay := int64((_amountTx + _fee) * 100000000)
	if _pay <= 0 {
		_ret.Message = "支付金额小于零"
		c.JSON(http.StatusOK, _ret)
		return
	}

	_mtx, _taken, err := GenRawTransactionEx(_params.Utxos, _params.SendAddr, _params.RecvAddr, uint64(_pay), uint64(_amount))
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_tx, _hashList := GenSignHash(_mtx, _taken)
	if err != nil {
		_ret.Message = err.Error()
		c.JSON(http.StatusOK, _ret)
		return
	}

	_ret.Status = "true"
	_ret.Tx = _tx
	_ret.HashList = _hashList

	c.JSON(http.StatusOK, _ret)
}
