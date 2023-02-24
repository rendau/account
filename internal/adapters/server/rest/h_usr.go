package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rendau/account/internal/domain/entities"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/dop/dopTypes"
)

// @Router  /usr [get]
// @Tags    usr
// @Param   query query entities.UsrListParsSt false "query"
// @Produce json
// @Success 200 {object} dopTypes.PaginatedListRep{results=[]entities.UsrListSt}
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hUsrList(c *gin.Context) {
	pars := &entities.UsrListParsSt{}
	if !dopHttps.BindQuery(c, pars) {
		return
	}

	result, tCount, err := o.ucs.UsrList(o.getRequestContext(c), pars)
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

// @Router  /usr [post]
// @Tags    usr
// @Param   body body entities.UsrCUSt false "body"
// @Produce json
// @Success 200 {object} dopTypes.CreateRep{id=int}
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hUsrCreate(c *gin.Context) {
	reqObj := &entities.UsrCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	newId, err := o.ucs.UsrCreate(o.getRequestContext(c), reqObj)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": newId})
}

// @Router  /usr/:id [get]
// @Tags    usr
// @Param   id path integer true "id"
// @Produce json
// @Success 200 {object} entities.UsrSt
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hUsrGet(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	result, err := o.ucs.UsrGet(o.getRequestContext(c), id)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Router  /usr/:id [put]
// @Tags    usr
// @Param   id   path string           true  "id"
// @Param   body body entities.UsrCUSt false "body"
// @Produce json
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hUsrUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	reqObj := &entities.UsrCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.UsrUpdate(o.getRequestContext(c), id, reqObj))
}

// @Router  /usr/:id/generate_and_save_access_token [put]
// @Tags    usr
// @Param   id   path string           true  "id"
// @Produce json
// @Success 200 {object} entities.AuthByTokenRepSt
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hUsrGenerateAndSaveAccessToken(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	reqObj := &entities.UsrCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	result, err := o.ucs.UsrGenerateAndSaveAccessToken(o.getRequestContext(c), id)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, entities.AuthByTokenRepSt{AccessToken: result})
}

// @Router  /usr/:id [delete]
// @Tags    usr
// @Param   id path string true "id"
// @Produce json
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hUsrDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	dopHttps.Error(c, o.ucs.UsrDelete(o.getRequestContext(c), id))
}
