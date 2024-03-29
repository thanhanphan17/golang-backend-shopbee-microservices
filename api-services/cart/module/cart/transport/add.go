package carttransport

import (
	"net/http"
	"shopbee/common"
	"shopbee/component/appctx"
	cartbiz "shopbee/module/cart/biz"
	cartmodel "shopbee/module/cart/model"
	cartstorage "shopbee/module/cart/storage"

	"github.com/gin-gonic/gin"
)

func AddToCart(appctx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appctx.GetMainDBConnection()
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		type CartCreate struct {
			ProductId       string `json:"product_id"`
			ProductQuantity int    `json:"quantity"`
		}

		var cartCreate CartCreate

		if err := c.ShouldBind(&cartCreate); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		uid, err := common.FromBase58(cartCreate.ProductId)

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data cartmodel.CartCreate

		data.UserId = requester.GetUserId()
		data.ProductId = int(uid.GetLocalID())
		data.ProductQuantity = cartCreate.ProductQuantity

		store := cartstorage.NewSQLStore(db)
		biz := cartbiz.NewAddToCartBiz(store, requester)

		if err := biz.AddToCart(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSucessResponse(true))
	}
}
