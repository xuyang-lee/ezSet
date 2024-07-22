package syncSet

import (
	"github.com/xuyang-lee/ezSet/internal/consts"
	"sync"
	"sync/atomic"
)

var (
	uuid uint64 = 0
)

type SyncSet[T comparable] struct {
	setMap map[T]consts.Empty
	m      sync.RWMutex
	u      uint64 //集合全局uuid,加锁时使用，避免重复加锁
}

// NewSetOf create a new set,
// the type of the set will be the same as the type of t
func NewSyncSetOf[T comparable](t T) *SyncSet[T] {
	return newSet[T]()
}

// NewSet create a new set,
// the type of the set is custom
func NewSyncSet[T comparable]() *SyncSet[T] {
	return newSet[T]()
}

// NewSetWithSlice create a new set,
// the type of the set will be the same as the type of l,
// and new set will be init with l
func NewSyncSetWithSlice[T comparable](l []T) *SyncSet[T] {
	s := newSet[T]()
	s.m.Lock()
	defer s.m.Unlock()

	for _, v := range l {
		s.setMap[v] = consts.EMPTY
	}

	return s
}

// Len get the size of set
func (s *SyncSet[T]) Len() int {
	return s.len()
}

// Add elements to the set
func (s *SyncSet[T]) Add(l ...T) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, e := range l {
		s.add(e)
	}
	return
}

// Remove elements from the set
func (s *SyncSet[T]) Remove(l ...T) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, e := range l {
		s.remove(e)
	}
	return
}

// Clear get an empty set with origin type
func (s *SyncSet[T]) Clear() {
	s.m.Lock()
	defer s.m.Unlock()
	s.setMap = make(map[T]consts.Empty)
	return
}

// Clone create a new set which has the same elements as s
func (s *SyncSet[T]) Clone() *SyncSet[T] {
	return s.clone()
}

// List get all elements as a slice
func (s *SyncSet[T]) List() []T {
	return s.list()
}

// Contains  whether e belongs to s
func (s *SyncSet[T]) Contains(e T) bool {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.contains(e)
}

// IsSubSetOf  whether s is subset of d
func (s *SyncSet[T]) IsSubSetOf(d *SyncSet[T]) bool {
	return s.isSubSetOf(d)
}

// IsSuperSetOf  whether s is super set of d
func (s *SyncSet[T]) IsSuperSetOf(d *SyncSet[T]) bool {
	return s.isSuperSetOf(d)
}

// IsEqual  whether s is the same as d
func (s *SyncSet[T]) IsEqual(d *SyncSet[T]) bool {
	return s.isSubSetOf(d) && s.isSuperSetOf(d)
}

// Union 交集：return set = s * d
func (s *SyncSet[T]) Union(d *SyncSet[T]) *SyncSet[T] {
	//创建交集
	union := newSet[T]()

	less, great := order(s, d)
	less.m.RLock()
	defer less.m.RUnlock()
	great.m.RLock()
	defer great.m.RUnlock()

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
func (s *SyncSet[T]) Intersect(d *SyncSet[T]) *SyncSet[T] {
	intersect := s.clone()
	d.m.RLock()
	defer d.m.RUnlock()

	for k := range d.setMap {
		intersect.add(k)
	}
	return intersect
}

// Difference 差集：return set = s - d
func (s *SyncSet[T]) Difference(d *SyncSet[T]) *SyncSet[T] {
	diff := newSet[T]()

	less, great := order(s, d)
	less.m.RLock()
	defer less.m.RUnlock()
	great.m.RLock()
	defer great.m.RUnlock()

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
func (s *SyncSet[T]) Complement(d *SyncSet[T]) *SyncSet[T] {
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

func newSet[T comparable]() *SyncSet[T] {
	return &SyncSet[T]{
		setMap: make(map[T]consts.Empty),
		u:      atomic.AddUint64(&uuid, 1),
	}
}

func (s *SyncSet[T]) add(e T) {
	s.setMap[e] = consts.EMPTY
	return
}

func (s *SyncSet[T]) remove(e T) {

	if _, ok := s.setMap[e]; ok {
		delete(s.setMap, e)
	}

	return
}

func (s *SyncSet[T]) clone() *SyncSet[T] {
	c := newSet[T]()
	s.m.RLock()
	defer s.m.RUnlock()
	for k := range s.setMap {
		c.setMap[k] = consts.EMPTY
	}
	return c
}

func (s *SyncSet[T]) len() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return len(s.setMap)
}

func (s *SyncSet[T]) list() []T {
	s.m.RLock()
	defer s.m.RUnlock()
	var l = make([]T, 0)
	for k := range s.setMap {
		l = append(l, k)
	}
	return l
}

func (s *SyncSet[T]) isSubSetOf(d *SyncSet[T]) bool {
	if s.len() > d.len() {
		return false
	}

	less, great := order(s, d)
	less.m.RLock()
	defer less.m.RUnlock()
	great.m.RLock()
	defer great.m.RUnlock()

	for v := range s.setMap {
		if _, ok := d.setMap[v]; !ok {
			return false
		}
	}
	return true
}

func (s *SyncSet[T]) isSuperSetOf(d *SyncSet[T]) bool {
	if s.len() < d.len() {
		return false
	}

	less, great := order(s, d)
	less.m.RLock()
	defer less.m.RUnlock()
	great.m.RLock()
	defer great.m.RUnlock()

	for v := range d.setMap {
		if _, ok := s.setMap[v]; !ok {
			return false
		}
	}
	return true
}

func (s *SyncSet[T]) contains(e T) bool {
	if _, ok := s.setMap[e]; ok {
		return true
	}
	return false
}

func order[T comparable](a, b *SyncSet[T]) (less, great *SyncSet[T]) {
	//排序拥有有序加锁，避免死锁
	less, great = a, b
	if a.u > b.u {
		less, great = b, a
	}
	return
}
