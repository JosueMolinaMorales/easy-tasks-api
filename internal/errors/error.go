package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestError struct {
	code    int
	message string
}

func NewRequestError(code int, message string) *RequestError {
	return &RequestError{
		code, message,
	}
}

func (re RequestError) Error() string {
	return re.message
}

func HandleError(ctx *gin.Context, err error) {
	reqErr, ok := err.(*RequestError)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}

	switch reqErr.code {
	case http.StatusInternalServerError:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
	default:
		ctx.JSON(reqErr.code, gin.H{"error": reqErr.message})
	}
}
