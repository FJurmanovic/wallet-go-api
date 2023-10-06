package controller

import (
	"net/http"
	"wallet-api/pkg/filter"
	"wallet-api/pkg/middleware"
	"wallet-api/pkg/model"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

/*
NewUserController

Initializes UserController.

	Args:
		*services.UserService: User service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*UserController: Controller for "auth" interactions
*/
func NewUserController(rs *service.UserService, routeGroups *common.RouteGroups) *UserController {
	rc := &UserController{
		service: rs,
	}

	routeGroups.Auth.POST("login", rc.PostLogin)
	routeGroups.Auth.POST("register", rc.PostRegister)
	routeGroups.Auth.DELETE("deactivate", middleware.Auth, rc.Delete)
	routeGroups.Auth.GET("check-token", rc.CheckToken)

	return rc
}

/*
PostLogin
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /auth/login).
func (rc *UserController) PostLogin(c *gin.Context) {
	body := new(model.Login)
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Email == "" {
		c.JSON(400, "Email cannot be empty!")
		return
	}
	if body.Password == "" {
		c.JSON(400, "Password cannot be empty!")
		return
	}
	returnedUser, exceptionReturn := rc.service.Login(c, body)

	if exceptionReturn.Message != "" {
		c.JSON(exceptionReturn.StatusCode, exceptionReturn)
		return
	} else {
		c.JSON(200, returnedUser)
	}
}

/*
PostRegister
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (POST /auth/register).
func (rc *UserController) PostRegister(c *gin.Context) {
	body := new(model.User)
	body.Init()
	body.IsActive = true
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	returnedUser, exceptionReturn := rc.service.Create(c, body)

	if exceptionReturn.Message != "" {
		c.JSON(exceptionReturn.StatusCode, exceptionReturn)
		return
	} else {
		c.JSON(200, returnedUser.Payload())
	}
}

/*
Delete
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (DELETE /auth/deactivate).
func (rc *UserController) Delete(c *gin.Context) {
	authGet := c.MustGet("auth")
	userId := authGet.(*model.Auth).Id

	flt := filter.NewUserFilter(model.Params{})
	flt.UserId = userId

	mr, er := rc.service.Deactivate(c, flt)

	if er.Message != "" {
		c.JSON(er.StatusCode, er)
		return
	} else {
		c.JSON(200, mr)
	}
}

/*
CheckToken
	Args:
		*gin.Context: Gin Application Context
*/
// ROUTE (GET /auth/check-token).
func (rc *UserController) CheckToken(c *gin.Context) {
	token, _ := c.GetQuery("token")
	re := new(model.CheckToken)

	_, err := middleware.CheckToken(token)

	if err != nil {
		re.Valid = false
		c.AbortWithStatusJSON(400, re)
		return
	}

	re.Valid = true

	c.JSON(200, re)
}
