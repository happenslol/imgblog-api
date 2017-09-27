package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, data ...interface{}) {
	if len(data) == 1 {
		c.JSON(http.StatusOK, data[0])
		return
	}

	c.JSON(http.StatusOK, data)
}

func Created(c *gin.Context, entity interface{}) {
	c.JSON(http.StatusCreated, entity)
}

func Deleted(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func DbError(c *gin.Context, err error) {
	//TODO better error reporting
	Log.Errorf("database error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func ServerError(c *gin.Context, err error) {
	Log.Errorf("server error: %v", err)
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
