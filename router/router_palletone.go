package router

import (
	"github.com/Eric-GreenComb/palletone/handler"
	"github.com/gin-gonic/gin"
)

// SetupPalletoneRouter SetupPalletoneRouter
func SetupPalletoneRouter(g *gin.Engine) {
	rpalletone := g.Group("/")
	{
		// 根据参数，生成交易结构
		rpalletone.POST("tx/raw", handler.GetRawTx)

		// 根据参数，生成交易结构，并进行rlp编码
		rpalletone.POST("tx/encoding", handler.GetRawTxEncoding)

		rpalletone.POST("gettxhash", handler.GetTxHash)
	}
}
