package app

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func Ok(c *gin.Context, data ...interface{}) {
	if len(data) == 1 {
		c.JSON(http.StatusOK, data[0])
		return
	}

	c.JSON(http.StatusOK, data)
}

func Created(c *gin.Context, id interface{}) {
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func DbError(c *gin.Context, err error) {
	//TODO better error reporting
	Log.Errorf("database error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func BadRequest(c *gin.Context, err error) {
	//TODO better error reporting
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	c.Abort()
}
