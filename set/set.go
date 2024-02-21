package set

// Copyright (c) 2024, Lee.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

type empty struct{}

type Set[T comparable] struct {
	setMap map[T]empty
}

var _EMPTY = empty{}

// NewSetOf create a new set,
// the type of the set will be the same as the type of t
func NewSetOf[T comparable](t T) Set[T] {
	return newSet[T]()
}

// NewSet create a new set,
// the type of the set is custom
func NewSet[T comparable]() Set[T] {
	return newSet[T]()
}

// NewSetWithSlice create a new set,
// the type of the set will be the same as the type of l,
// and new set will be init with l
func NewSetWithSlice[T comparable](l ...T) Set[T] {
	s := newSet[T]()

	for _, v := range l {
		s.setMap[v] = _EMPTY
	}

	return s
}

// Len get the size of set
func (s *Set[T]) Len() int {
	return s.len()
}

// Add elements to the set
func (s *Set[T]) Add(l ...T) {

	for _, e := range l {
		s.add(e)
	}
	return
}

// Remove elements from the set
func (s *Set[T]) Remove(l ...T) {
	for _, e := range l {
		s.remove(e)
	}
	return
}

// Clear get an empty set with origin type
func (s *Set[T]) Clear() {
	s.setMap = make(map[T]empty, 0)
	return
}

// Clone create a new set which has the same elements as s
func (s *Set[T]) Clone() Set[T] {
	return s.clone()
}

// List get all elements as a slice
func (s *Set[T]) List() []T {
	var l []T
	for k := range s.setMap {
		l = append(l, k)
	}
	return l
}

// Contains  whether e belongs to s
func (s *Set[T]) Contains(e T) bool {
	return s.contains(e)
}

// IsSubSetOf  whether s is subset of d
func (s *Set[T]) IsSubSetOf(d Set[T]) bool {
	return s.isSubSetOf(d)
}

// IsSuperSetOf  whether s is super set of d
func (s *Set[T]) IsSuperSetOf(d Set[T]) bool {
	return s.isSuperSetOf(d)
}

// IsEqual  whether s is the same as d
func (s *Set[T]) IsEqual(d Set[T]) bool {
	if s.isSubSetOf(d) && s.isSuperSetOf(d) {
		return true
	}
	return false
}

// Union 交集：return set = s * d
func (s *Set[T]) Union(d Set[T]) Set[T] {
	//创建交集
	union := newSet[T]()
	//遍历d中元素
	for k := range d.setMap {
		//若s中存在，则加入集合
		if s.contains(k) {
			union.add(k)
		}
	}

	return union

}

// Intersect 并集：return set = s + d
func (s *Set[T]) Intersect(d Set[T]) Set[T] {
	intersect := s.clone()
	for k := range d.setMap {
		intersect.add(k)
	}
	return intersect
}

// Difference 差集：return set = s - d
func (s *Set[T]) Difference(d Set[T]) Set[T] {
	diff := newSet[T]()
	for k := range s.setMap {
		//在s中而不在d中
		if !d.contains(k) {
			diff.add(k)
		}
	}
	//输出差集结果
	return diff
}

// Complement 补集：return set = d - s
func (s *Set[T]) Complement(d Set[T]) Set[T] {
	comp := newSet[T]()
	for k := range d.setMap {
		//在d中而不在s中
		if !s.contains(k) {
			comp.add(k)
		}
	}
	//输出差集结果
	return comp
}

func newSet[T comparable]() Set[T] {
	s := Set[T]{}
	s.setMap = make(map[T]empty)
	return s
}

func (s *Set[T]) add(e T) {
	s.setMap[e] = _EMPTY
	return
}

func (s *Set[T]) remove(e T) {

	if _, ok := s.setMap[e]; ok {
		delete(s.setMap, e)
	}

	return
}

func (s *Set[T]) clone() Set[T] {
	c := newSet[T]()
	for k := range s.setMap {
		c.setMap[k] = _EMPTY
	}
	return c
}

func (s *Set[T]) len() int {
	return len(s.setMap)
}

func (s *Set[T]) isSubSetOf(d Set[T]) bool {
	for v := range s.setMap {
		if _, ok := d.setMap[v]; !ok {
			return false
		}
	}
	return true
}

func (s *Set[T]) isSuperSetOf(d Set[T]) bool {
	for v := range d.setMap {
		if _, ok := s.setMap[v]; !ok {
			return false
		}
	}
	return true
}

func (s *Set[T]) contains(e T) bool {
	if _, ok := s.setMap[e]; ok {
		return true
	}
	return false
}
