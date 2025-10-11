package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vky5/mailcat/internal/db"
	"github.com/vky5/mailcat/internal/db/models"
)

func RegisterAccountRoutes(r *gin.Engine) {
	acc := r.Group("/account")
	{ // no special syntax just limiting the acc in scope could be removed no affect
		acc.POST("/", CreateAccount)
		acc.GET("/", GetAccounts)
	}
}

func CreateAccount(c *gin.Context) {
	var account models.Account
	/*
	ShouldBindBodyWithJSON
		Gin reads http body from c.Request.Body 
		Decodes JSON into struct field using encoding/json
		automatically handles errors like malformed JSON or missing fields

	usually reading the request once consume it (cant read it twice)
	but ShouldBindWithJSON stores a copy internally to bind it again
	*/

	if err := c.ShouldBindBodyWithJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func GetAccounts(c *gin.Context) {
	var accounts []models.Account
	if err := db.DB.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}


/*
c is a request lifecycle container a struct that wraps
- c.Request : Incoming http request (type *http.Request)
- c.Writer : http response writer (type *gin.ResponseWriter)

but 

c.Writer.WriteHeader(200)
c.Writer.Write([]byte(`{"msg":"ok"}`))
this is how response will look like: ou’d have to handle JSON encoding, headers, and error handling yourself.

instead use this
c.JSON(http.StatusOK, gin.H{"msg": "ok"}) // sets Content-Type: application/json header
*/

