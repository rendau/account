package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rendau/account/internal/domain/usecases"
	"github.com/rendau/dop/adapters/logger"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	swagFiles "github.com/swaggo/files"
	ginSwag "github.com/swaggo/gin-swagger"
)

type St struct {
	lg  logger.Lite
	ucs *usecases.St
}

func GetHandler(lg logger.Lite, ucs *usecases.St, withCors bool) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// middlewares

	r.Use(dopHttps.MwRecovery(lg, nil))

	if withCors {
		r.Use(dopHttps.MwCors())
	}

	// handlers

	// doc
	r.GET("/doc/*any", ginSwag.WrapHandler(swagFiles.Handler, func(c *ginSwag.Config) {
		c.DefaultModelsExpandDepth = 0
		c.DocExpansion = "none"
	}))

	s := &St{lg: lg, ucs: ucs}

	// healthcheck
	r.GET("/healthcheck", func(c *gin.Context) { c.Status(http.StatusOK) })

	// system
	r.PUT("/mss/fs/filter_unused_files", s.hSystemFilterUnusedFiles)
	r.GET("/mss/cron/tick5m", s.hSystemCronTick5m)
	r.GET("/mss/cron/tick15m", s.hSystemCronTick15m)
	r.GET("/mss/cron/tick30m", s.hSystemCronTick30m)

	// dic
	r.GET("/dic", s.hDicGet)

	// profile
	r.POST("/profile/send_validating_code", s.hProfileSendPhoneValidatingCode)
	r.POST("/profile/auth", s.hProfileAuth)
	r.POST("/profile/auth/token", s.hProfileAuthToken)
	r.POST("/profile/reg", s.hProfileReg)
	r.POST("/profile/logout", s.hProfileLogout)
	r.GET("/profile", s.hProfileGet)
	r.PUT("/profile", s.hProfileUpdate)
	r.PUT("/profile/change_phone", s.hProfileChangePhone)

	// role
	r.GET("/role", s.hRoleList)
	r.POST("/role", s.hRoleCreate)
	r.GET("/role/:id", s.hRoleGet)
	r.PUT("/role/:id", s.hRoleUpdate)
	r.DELETE("/role/:id", s.hRoleDelete)

	// usr
	r.GET("/usr", s.hUsrList)
	r.POST("/usr", s.hUsrCreate)
	r.GET("/usr/:id", s.hUsrGet)
	r.PUT("/usr/:id", s.hUsrUpdate)
	r.DELETE("/usr/:id", s.hUsrDelete)

	// perm
	r.GET("/perm", s.hPermList)
	r.POST("/perm", s.hPermCreate)
	r.GET("/perm/:id", s.hPermGet)
	r.PUT("/perm/:id", s.hPermUpdate)
	r.DELETE("/perm/:id", s.hPermDelete)

	return r
}

func (o *St) getRequestContext(c *gin.Context) context.Context {
	token := dopHttps.GetAuthToken(c)

	return o.ucs.SessionSetToContextByToken(nil, token)
}
