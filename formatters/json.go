package formatters

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type JSONFormatter[T any] struct {
}

func (f *JSONFormatter[T]) Format(in []T, ctx *gin.Context) (string, error) {
	ctx.JSON(200, in)
	ret, err := json.Marshal(in)
	return string(ret), err
}
