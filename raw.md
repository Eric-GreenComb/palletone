# 交易过程

## 从接口获取utxo列表，及交易信息

## 构建tx

1. 根据支付的数额pay（付款+手续费）和utxo根据贪心算法获取使用的utxo和找零数量

  ps:如果小于pay的utxo够用，则返回使用的utxo数组和找零，不够则返回大于pay的最小utxo和找零

2. 构建input

  ps:使用1返回的utxo数组，构建pload的TxIn数组

3. 构建output

  ps:先构建to账户的output，再构建找零的output，组成数组设置为pload的TxOut

4. 根据pload组成tx的Messages，就此生成mtx

5. 对每个选用的utxo进行签名，生成每组utxo的hash

  ps:tokenengine.CalcSignatureHash

