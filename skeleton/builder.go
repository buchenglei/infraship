package skeleton

type Builder[T any] interface {
	Build() (T, error)
}
