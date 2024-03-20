package skeleton

// Anti-Corruption Layder Declaration
// 以Entity为核心的各种类型转换操作的定义

// FormEntity 从entity对象中派生出一个新的对象
type FromEntity[E any] interface {
	From(E) error
}

// IntoEntity 将自己(指定对象)转换成entity对象
type IntoEntity[E any] interface {
	Into() (E, error)
}

// IntoEntityWrapper IntoEntity Wrapper
// Example: RawData -> Wrap1 -> Wrap2 -> Wrap3 -> Entity
type IntoEntityWrapper[E any] struct {
	F func(E) (E, error)

	_placeholder E
	raw          IntoEntity[E]
}

func (w IntoEntityWrapper[E]) Wrap(r IntoEntity[E]) IntoEntity[E] {
	w.raw = r
	return w
}

func (w IntoEntityWrapper[E]) Into() (E, error) {
	e, err := w.raw.Into()
	if err != nil {
		return w._placeholder, err
	}

	return w.F(e)
}
