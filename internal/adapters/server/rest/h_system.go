package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
)

func (o *St) hSystemFilterUnusedFiles(c *gin.Context) {
	reqObj := make([]string, 0)
	if !dopHttps.BindJSON(c, &reqObj) {
		return
	}

	c.JSON(http.StatusOK, o.ucs.SystemFilterUnusedFiles(reqObj))
}

func (o *St) hSystemCronTick5m(c *gin.Context) {
	o.ucs.SystemCronTick5m()
}

func (o *St) hSystemCronTick15m(c *gin.Context) {
	o.ucs.SystemCronTick15m()
}

func (o *St) hSystemCronTick30m(c *gin.Context) {
	o.ucs.SystemCronTick30m()
}

func (o *St) hSystemGetPerms(c *gin.Context) {
	c.JSON(http.StatusOK, o.ucs.GetPerms())
}
