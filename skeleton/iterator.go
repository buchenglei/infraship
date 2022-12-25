package skeleton

import (
	"errors"
	"io"
)

var (
	ErrEmptyIterator = errors.New("empty iterator")
	ErrEOF           = io.EOF
)

// Iterator 迭代器
// 用于遍历一批数据中的每一个元素
type Iterator[T any] interface {
	Next() (T, error)
}

// UseIter 通过匿名函数处理迭代器的返回的每一个元素
// deferCalls 当迭代器遍历完成后，依次执行起到cleanup的作用
func UseIter[T any](iter Iterator[T], handle func(T) error, deferCalls ...func()) error {
	var (
		err  error
		item T
	)

	// register defer calls
	if len(deferCalls) > 0 {
		for _, call := range deferCalls {
			defer call()
		}
	}

	for {
		item, err = iter.Next()
		if err != nil {
			return err
		}
		err = handle(item)
		if err != nil {
			return err
		}
	}
}
