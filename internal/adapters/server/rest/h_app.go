package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rendau/account/internal/domain/entities"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/dop/dopTypes"
)

// @Router  /app [get]
// @Tags    app
// @Param   query query entities.AppListParsSt false "query"
// @Produce json
// @Success 200 {object} dopTypes.PaginatedListRep{results=[]entities.AppSt}
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hAppList(c *gin.Context) {
	pars := &entities.AppListParsSt{}
	if !dopHttps.BindQuery(c, pars) {
		return
	}

	result, tCount, err := o.ucs.AppList(o.getRequestContext(c), pars)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, dopTypes.PaginatedListRep{
		Page:       pars.Page,
		PageSize:   pars.PageSize,
		TotalCount: tCount,
		Results:    result,
	})
}

// @Router  /app [post]
// @Tags    app
// @Param   body body     entities.AppCUSt false "body"
// @Success 200  {object} dopTypes.CreateRep{id=integer}
// @Failure 400  {object} dopTypes.ErrRep
func (o *St) hAppCreate(c *gin.Context) {
	reqObj := &entities.AppCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	result, err := o.ucs.AppCreate(o.getRequestContext(c), reqObj)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, dopTypes.CreateRep{Id: result})
}

// @Router  /app/:id [get]
// @Tags    app
// @Param   id path integer true "id"
// @Produce json
// @Success 200 {object} entities.AppSt
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hAppGet(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	result, err := o.ucs.AppGet(o.getRequestContext(c), id)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Router  /app/:id [put]
// @Tags    app
// @Param   id   path integer          true  "id"
// @Param   body body entities.AppCUSt false "body"
// @Produce json
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hAppUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	reqObj := &entities.AppCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.AppUpdate(o.getRequestContext(c), id, reqObj))
}

// @Router  /app/:id [delete]
// @Tags    app
// @Param   id path integer true "id"
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hAppDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	dopHttps.Error(c, o.ucs.AppDelete(o.getRequestContext(c), id))
}

// @Router  /app/fetch_perms [post]
// @Tags    app
// @Param   body body entities.AppFetchPermsReqSt false "body"
// @Produce json
// @Success 200 {object} entities.SystemGetPermsRepSt
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hAppFetchPerms(c *gin.Context) {
	reqObj := &entities.AppFetchPermsReqSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	result, err := o.ucs.AppFetchPerms(o.getRequestContext(c), reqObj.Uri)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Router  /app/:id/sync_perms [put]
// @Tags    app
// @Param   id path integer true "id"
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hAppSyncPerms(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	dopHttps.Error(c, o.ucs.AppSyncPerms(o.getRequestContext(c), id))
}
