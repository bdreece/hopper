package iter_test

import (
	"fmt"
	"testing"

	"github.com/bdreece/hopper/pkg/services/utils/iter"
)

func TestMap1(t *testing.T) {
	out := iter.Collect(iter.NewMap(
		iter.FromSlice(&[]int{1, 2, 3}),
		func(i *int) string {
			return fmt.Sprint(i)
		},
	))

	for i, str := range out {
		expected := fmt.Sprint(i + 1)
		if str != expected {
			t.Logf("%s != %s", str, expected)
			t.Fail()
		}
	}
}
