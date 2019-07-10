package bean

import ()

// Params Params
type Params struct {
	Key   string `form:"key" json:"key"`     //
	Value string `form:"value" json:"value"` //
}

// TxParams TxParams
type TxParams struct {
	SendAddr string `form:"sendAddr" json:"sendAddr"` //
	RecvAddr string `form:"recvAddr" json:"recvAddr"` //
	Amount   string `form:"amount" json:"amount"`     //
	Fee      string `form:"fee" json:"fee"`           //
	Utxos    Utxos  `form:"utxo" json:"utxo"`         //
	Rlp      string `form:"rlp" json:"rlp"`           //
}

// TxUtxo TxUtxo
type TxUtxo struct {
	TxID               string `json:"txid"`                 //
	MessageIndex       uint32 `json:"message_index"`        //
	OutIndex           uint32 `json:"out_index"`            //
	Amount             uint64 `json:"amount"`               //
	Asset              string `json:"asset"`                //
	PkScriptHex        string `json:"pk_script_hex"`        //
	PkScriptString     string `json:"pk_script_string"`     //
	CreateTime         string `json:"create_time"`          //
	LockTime           uint   `json:"lock_time"`            //
	FlagStatus         string `json:"flag_status"`          //
	CoinDays           int64  `json:"coin_days"`            //
	AmountWithInterest int64  `json:"amount_with_interest"` //
}

// Utxos Utxos
type Utxos []TxUtxo

func (c Utxos) Len() int {
	return len(c)
}
func (c Utxos) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Utxos) Less(i, j int) bool {
	return c[i].Amount < c[j].Amount
}

// FindMin FindMin
func FindMin(utxos []TxUtxo) TxUtxo {
	amout := utxos[0].Amount
	minUtxo := utxos[0]
	for _, utxo := range utxos {
		if utxo.Amount < amout {
			minUtxo = utxo
			amout = minUtxo.Amount
		}
	}
	return minUtxo
}

// TxReturn TxReturn
type TxReturn struct {
	Status          string   `json:"status"`           // status true false
	Message         string   `json:"message"`          // 如有失败原因，请写失败原因，成功请空着
	SignatureScript string   `json:"signature_script"` // from地址的 signature_script
	Tx              string   `json:"tx"`               // 待签名的tx字符串”//这个是标准的rlp编码后的字符串，跟例子上的数据一致
	HashList        []TxHash `json:"hashlist"`         // 待签名数据列表，按照input的顺序排列
}

// TxHash TxHash
type TxHash struct {
	Hash string `json:"hash"` //
}
