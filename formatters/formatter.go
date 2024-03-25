package formatters

type Formatter[T any] interface {
	Format([]T) (string, error)
}
