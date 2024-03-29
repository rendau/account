package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rendau/account/internal/domain/entities"
	dopHttps "github.com/rendau/dop/adapters/server/https"
)

// @Router  /role [get]
// @Tags    role
// @Param   query query entities.RoleListParsSt false "query"
// @Produce json
// @Success 200 {array}  entities.RoleListSt
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hRoleList(c *gin.Context) {
	pars := &entities.RoleListParsSt{}
	if !dopHttps.BindQuery(c, pars) {
		return
	}

	result, err := o.ucs.RoleList(o.getRequestContext(c), pars)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Router  /role [post]
// @Tags    role
// @Param   body body entities.RoleCUSt false "body"
// @Produce json
// @Success 200 {object} dopTypes.CreateRep{id=int}
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hRoleCreate(c *gin.Context) {
	reqObj := &entities.RoleCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	result, err := o.ucs.RoleCreate(o.getRequestContext(c), reqObj)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": result})
}

// @Router  /role/:id [get]
// @Tags    role
// @Param   id path string true "id"
// @Produce json
// @Success 200 {object} entities.RoleSt
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hRoleGet(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	result, err := o.ucs.RoleGet(o.getRequestContext(c), id)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Router  /role/:id [put]
// @Tags    role
// @Param   id   path string            true  "id"
// @Param   body body entities.RoleCUSt false "body"
// @Produce json
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hRoleUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	reqObj := &entities.RoleCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.RoleUpdate(o.getRequestContext(c), id, reqObj))
}

// @Router  /role/:id [delete]
// @Tags    role
// @Param   id path string true "id"
// @Produce json
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (o *St) hRoleDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	dopHttps.Error(c, o.ucs.RoleDelete(o.getRequestContext(c), id))
}
