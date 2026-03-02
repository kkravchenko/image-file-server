package handlers

import (
	"is/internal/download"

	"github.com/gin-gonic/gin"
)

func AddImage(imagePath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.PostForm("url")
		shop := ctx.PostForm("shop")
		entity := ctx.PostForm("entity")

		req := download.Request{
			Url: url, ImagePath: imagePath, EntityType: entity, Shop: shop,
		}

		imagePath := req.Download()
		ctx.JSON(200, imagePath)
	}
}
