package orderSet

import (
	"github.com/xuyang-lee/ezSet/internal/consts"
)

// Copyright (c) 2024, Lee.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Attention please!
// T is a type that can be compared.
// Including int, float, string and structures which cannot contain incomparable field types(such as slice, map, chan).
// If the structure contains pointers, note that if the pointer fields in the two structures are different,
// then these two structures are different elements, even the values of the two pointers point are the same
// e.g.
// 	type A struct {
// 		age  *int
//	}
//  a,b := 1,1
//  set1 := utils.NewSetWithSlice(A{age: &a})
//	set2 := utils.NewSetWithSlice(A{age: &b})
//  fmt.Println(set1.IsEqual(set2))
// ============
//  false
// ============

// 请注意！ T是可比较类型，包括int，float，string 和 不包括不可比较字段的结构体（结构体中没有slice、 map、 chan）。
// 如果结构体包含指针，要注意只要结构体中的指针不同，不论这两个指针所指向的值是不是一样的，这两个结构体就认为是不同的元素。
// 例子：
// 	type A struct {
// 		age  *int
//	}
//  a,b := 1,1
//  set1 := utils.NewSetWithSlice(A{age: &a})
//	set2 := utils.NewSetWithSlice(A{age: &b})
//  fmt.Println(set1.IsEqual(set2))
// ============
//  false
// ============

type OrderSet[T comparable] struct {
	setMap map[T]consts.Empty
	keys   []T
}

// NewOrderSetOf create a new set,
// the type of the set will be the same as the type of t
func NewOrderSetOf[T comparable](t T) *OrderSet[T] {
	return newSet[T]()
}

// NewOrderSet create a new set,
// the type of the set is custom
func NewOrderSet[T comparable]() *OrderSet[T] {
	return newSet[T]()
}

// NewOrderSetWithSlice create a new set,
// the type of the set will be the same as the type of l,
// and new set will be init with l
func NewOrderSetWithSlice[T comparable](l []T) *OrderSet[T] {
	s := newSet[T]()
	s.keys = make([]T, 0, len(l))

	for _, v := range l {
		if _, ok := s.setMap[v]; ok {
			continue
		}
		s.setMap[v] = consts.EMPTY
		s.keys = append(s.keys, v)
	}
	return s
}

// Len get the size of set
func (s *OrderSet[T]) Len() int {
	return s.len()
}

// Add elements to the set
func (s *OrderSet[T]) Add(l ...T) {

	for _, e := range l {
		s.add(e)
	}
	return
}

// Remove elements from the set
func (s *OrderSet[T]) Remove(l ...T) {
	//粗略估计如果元素数量大于10000或者删除元素数量大于50，直接使用removeMany更高效
	if s.len() > 10000 || len(l) > 50 {
		s.removeMany(l)
		return
	}

	for _, e := range l {
		s.remove(e)
	}
	return
}

func (s *OrderSet[T]) RemoveMany(l ...T) {
	s.removeMany(l)
	return
}

// Clear get an empty set with origin type
func (s *OrderSet[T]) Clear() {
	s.setMap = make(map[T]consts.Empty)
	s.keys = make([]T, 0)
	return
}

// Clone create a new set which has the same elements as s
func (s *OrderSet[T]) Clone() *OrderSet[T] {
	return s.clone()
}

// List get all elements as a slice
func (s *OrderSet[T]) List() []T {
	return s.list()
}

// Contains  whether e belongs to s
func (s *OrderSet[T]) Contains(e T) bool {
	return s.contains(e)
}

// IsSubSetOf  whether s is subset of d
func (s *OrderSet[T]) IsSubSetOf(d *OrderSet[T]) bool {
	return s.isSubSetOf(d)
}

// IsSuperSetOf  whether s is super set of d
func (s *OrderSet[T]) IsSuperSetOf(d *OrderSet[T]) bool {
	return s.isSuperSetOf(d)
}

// IsEqual  whether s is the same as d
func (s *OrderSet[T]) IsEqual(d *OrderSet[T]) bool {
	return s.isSubSetOf(d) && s.isSuperSetOf(d)
}

func newSet[T comparable]() *OrderSet[T] {
	s := new(OrderSet[T])
	s.setMap = make(map[T]consts.Empty)
	return s
}

func (s *OrderSet[T]) add(e T) {
	if s.contains(e) {
		return
	}
	s.keys = append(s.keys, e)
	s.setMap[e] = consts.EMPTY

	return
}

func (s *OrderSet[T]) remove(e T) {
	if s.contains(e) {
		delete(s.setMap, e)
		s.find(e)
		copy(s.keys[s.find(e):], s.keys[s.find(e)+1:])
		s.keys = s.keys[:len(s.keys)-1]
	}
	return
}

func (s *OrderSet[T]) removeMany(l []T) {
	indicesToDelete := make(map[T]consts.Empty, len(l)) // 使用map为了快速查找
	for _, e := range l {
		indicesToDelete[e] = consts.EMPTY
		//删除集合中的元素
		delete(s.setMap, e)
	}

	//删除顺序切片中的元素
	tmp := s.keys[:0] // 使用相同的底层数组创建新的切片
	for _, key := range s.keys {
		if _, found := indicesToDelete[key]; !found {
			tmp = append(tmp, key)
		}
	}
	s.keys = tmp
}

func (s *OrderSet[T]) clone() *OrderSet[T] {
	c := newSet[T]()
	for _, key := range s.keys {
		c.setMap[key] = consts.EMPTY
	}
	c.keys = make([]T, s.len())
	copy(c.keys, s.keys)
	return c
}

func (s *OrderSet[T]) len() int {
	return len(s.setMap)
}

func (s *OrderSet[T]) list() []T {
	return s.keys
}

func (s *OrderSet[T]) isSubSetOf(d *OrderSet[T]) bool {
	if s.len() > d.len() {
		return false
	}

	for v := range s.setMap {
		if _, ok := d.setMap[v]; !ok {
			return false
		}
	}
	return true
}

func (s *OrderSet[T]) isSuperSetOf(d *OrderSet[T]) bool {
	if s.len() < d.len() {
		return false
	}

	for v := range d.setMap {
		if _, ok := s.setMap[v]; !ok {
			return false
		}
	}
	return true
}

func (s *OrderSet[T]) contains(e T) (ok bool) {
	_, ok = s.setMap[e]
	return
}

func (s *OrderSet[T]) find(e T) int {
	for i, v := range s.keys {
		if v == e {
			return i
		}
	}
	return -1
}
