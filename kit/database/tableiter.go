package database

import (
	"io"

	"github.com/buchenglei/infraship/skeleton"
)

var _ skeleton.Iterator[any] = &TableIterator[any]{}

type DataLoader[T any] func(limit, offset int) ([]T, error)

type TableIterator[T any] struct {
	_placeholder T

	limit      int
	offset     int
	totalCount int64
	elements   []T
	eleChan    chan []T
	loadFunc   DataLoader[T]
}

func NewTableIterator[T any](total int64, pageSize int, loadFunc DataLoader[T]) (*TableIterator[T], error) {
	iter := &TableIterator[T]{
		limit:      pageSize,
		offset:     0,
		totalCount: total,
		elements:   make([]T, 0, pageSize),
		eleChan:    make(chan []T, 1),
		loadFunc:   loadFunc,
	}

	// 异步加载数据
	go iter.loadElements()
	return iter, nil
}

func (t *TableIterator[T]) Next() (T, error) {
	// 判断当前的数据是否已经处理完
	var ok bool
	if len(t.elements) == 0 {
		t.elements, ok = <-t.eleChan
		if !ok || len(t.elements) == 0 {
			return t._placeholder, io.EOF
		}
	}

	// 永远只取models中的第一个， 然后pop出去
	// 上面的逻辑会判断models是否为空， 为空则会从数据库中加载一批
	e := t.elements[0]
	t.elements = t.elements[1:] // pop first

	return e, nil
}

func (t *TableIterator[T]) loadElements() {
	iterdCount := 0
	for {
		elements, err := t.loadFunc(t.limit, t.offset)
		if err != nil {
			break
		}
		t.offset += len(elements)
		t.eleChan <- elements
		iterdCount += len(elements)
		if iterdCount >= int(t.totalCount) {
			break
		}
	}

	close(t.eleChan)
}
