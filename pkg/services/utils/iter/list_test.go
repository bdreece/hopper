package iter_test

import (
	"testing"

	"github.com/bdreece/hopper/pkg/services/utils/iter"
)

func TestListIterator1(t *testing.T) {
	slice := []int{1, 2, 3}
	list := iter.FromSlice(&slice)

	var err error = nil
	for i := 1; err == nil; i++ {
		item, err := list.Next()
		if err != nil {
			break
		}

		if *item != i {
			t.Logf("%d is not equal to %d\n", *item, i)
			t.Fail()
		}
	}
}

func TestListIterator2(t *testing.T) {
	out := iter.Collect(
		iter.FromSlice(&[]int{1, 2, 3}))

	for i, val := range out {
		if val != i+1 {
			t.Logf("%d != %d", i+1, val)
			t.Fail()
		}
	}
}
