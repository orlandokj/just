package ringlist

import (
	"testing"

    "github.com/stretchr/testify/assert"
)

func TestRingList(t *testing.T) {
    rl := NewRingList(10)
    for i := 0; i < 10; i++ {
        rl.Add(i)
    }

    result := []int{}
    rl.Do(func(x interface{}) {
        result = append(result, x.(int))
    })

    expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    assert.ElementsMatch(t, expected, result)
    assert.Equal(t, 10, rl.Len())
}

func TestRingList_Overflowing(t *testing.T) {
    rl := NewRingList(10)
    for i := 0; i < 15; i++ {
        rl.Add(i)
    }

    result := []int{}
    rl.Do(func(x interface{}) {
        result = append(result, x.(int))
    })

    expected := []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
    assert.ElementsMatch(t, expected, result)
    assert.Equal(t, 10, rl.Len())
}

func TestRingList_Clear(t *testing.T) {
    rl := NewRingList(10)
    for i := 0; i < 15; i++ {
        rl.Add(i)
    }

    assert.Equal(t, 10, rl.Len())

    rl.Clear()

    result := []int{}
    rl.Do(func(x interface{}) {
        result = append(result, x.(int))
    })

    assert.Equal(t, 0, rl.Len())
    assert.Empty(t, result)
}

func TestRingList_SpareSpace(t *testing.T) {
    rl := NewRingList(10)
    for i := 0; i < 5; i++ {
        rl.Add(i)
    }

    result := []int{}
    rl.Do(func(x interface{}) {
        result = append(result, x.(int))
    })

    expected := []int{0, 1, 2, 3, 4}
    assert.ElementsMatch(t, expected, result)
    assert.Equal(t, 5, rl.Len())
}
