package formatters

import "github.com/gin-gonic/gin"

type Formatter[T any] interface {
	Format([]T, *gin.Context) (string, error)
}
