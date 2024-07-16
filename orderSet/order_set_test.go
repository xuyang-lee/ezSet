package orderSet

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestOrderSet_Remove(t *testing.T) {
	s := NewOrderSetWithSlice([]int{1, 2, 3})
	s.Remove(3)
	fmt.Printf("%v", s.List())
}

func TestRemoveCost(t *testing.T) {
	setSize := 10000
	delSize := 5000
	delList := make([]int, 0, delSize)

	s := NewOrderSetOf(0)

	for i := 0; i < setSize; i++ {
		s.Add(i)
	}

	for i := 0; i < delSize; i++ {
		delList = append(delList, rand.Int()%setSize)
	}

	c := s.Clone()

	t1 := time.Now()
	s.Remove(delList...)
	tc1 := time.Since(t1).Seconds()

	t2 := time.Now()
	c.RemoveMany(delList...)
	tc2 := time.Since(t2).Seconds()

	fmt.Printf("删除百分比%v  ；一个个删除：%v     ；  一起删除：%v  \n", float64(delSize)/float64(setSize), tc1, tc2)
}
