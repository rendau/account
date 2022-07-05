package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rendau/account/internal/domain/entities"
	dopHttps "github.com/rendau/dop/adapters/server/https"
)

// @Router   /profile/send_validating_code [post]
// @Tags     profile
// @Param    body  body  entities.SendPhoneValidatingCodeReqSt  false  "body"
// @Produce  json
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileSendPhoneValidatingCode(c *gin.Context) {
	reqObj := &entities.SendPhoneValidatingCodeReqSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.ProfileSendPhoneValidatingCode(o.getRequestContext(c), reqObj.Phone, reqObj.ErrNE))
}

// @Router   /profile/auth [post]
// @Tags     profile
// @Param    body  body  entities.PhoneAndSmsCodeSt  false  "body"
// @Produce  json
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileAuth(c *gin.Context) {
	reqObj := &entities.PhoneAndSmsCodeSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	accessToken, refreshToken, err := o.ucs.ProfileAuth(o.getRequestContext(c), reqObj)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, entities.AuthRepSt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Router   /profile/auth/token [post]
// @Tags     profile
// @Param    body  body  entities.AuthByTokenReqSt  false  "body"
// @Produce  json
// @Success  200  {object}  entities.AuthByTokenRepSt
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileAuthToken(c *gin.Context) {
	reqObj := &entities.AuthByTokenReqSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	accessToken, err := o.ucs.ProfileAuthByRefreshToken(o.getRequestContext(c), reqObj.RefreshToken)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, entities.AuthByTokenRepSt{AccessToken: accessToken})
}

// @Router   /profile/reg [post]
// @Tags     profile
// @Param    body  body  entities.UsrRegReqSt  false  "body"
// @Produce  json
// @Success  200  {object}  entities.AuthRepSt
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileReg(c *gin.Context) {
	reqObj := &entities.UsrRegReqSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	accessToken, refreshToken, err := o.ucs.ProfileReg(o.getRequestContext(c), reqObj)
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, entities.AuthRepSt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Router   /profile/logout [post]
// @Tags     profile
// @Produce  json
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileLogout(c *gin.Context) {
	dopHttps.Error(c, o.ucs.ProfileLogout(o.getRequestContext(c)))
}

// @Router   /profile [get]
// @Tags     profile
// @Produce  json
// @Success  200  {object}  entities.UsrProfileSt
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileGet(c *gin.Context) {
	profile, err := o.ucs.ProfileGet(o.getRequestContext(c))
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, profile)
}

// @Router   /profile [put]
// @Tags     profile
// @Param    body  body  entities.UsrCUSt  false  "body"
// @Produce  json
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileUpdate(c *gin.Context) {
	reqObj := &entities.UsrCUSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.ProfileUpdate(o.getRequestContext(c), reqObj))
}

// @Router   /profile/change_phone [put]
// @Tags     profile
// @Param    body  body  entities.PhoneAndSmsCodeSt  false  "body"
// @Produce  json
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hProfileChangePhone(c *gin.Context) {
	reqObj := &entities.PhoneAndSmsCodeSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.ProfileChangePhone(o.getRequestContext(c), reqObj))
}

func (o *St) hProfileDelete(c *gin.Context) {
	dopHttps.Error(c, o.ucs.ProfileDelete(o.getRequestContext(c)))
}
