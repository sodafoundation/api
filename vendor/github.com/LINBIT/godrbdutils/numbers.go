package godrbdutils

import (
	"errors"
	"sort"
	"sync"
)

// GetNumber is used to return a free number within [min,max], where both are >=0
// It can be used to allocate new Port/Minor numbers
func GetNumber(min, max int, used []int) (int, error) {
	if max-min <= 0 || min < 0 || max < 0 {
		return -1, errors.New("min and/or max not valid")
	}

	// all free, use first one
	if len(used) == 0 {
		return min, nil
	}

	if !sort.IntsAreSorted(used) {
		sort.Ints(used)
	}

	// use the next after current max if possible
	curMax := used[len(used)-1]
	if curMax < max {
		return curMax + 1, nil
	}

	// find a hole
	for i := 0; i < len(used)-1; i++ {
		cur := used[i]
		if used[i+1]-cur > 1 { // found hole
			candidate := cur + 1
			if candidate >= min && candidate <= max {
				return candidate, nil
			}
		}
	}

	return -1, errors.New("Could not find a free number")
}

// NumberPool is used as a stateful type to keep track of used numbers
type NumberPool struct {
	min, max int
	used     []int
	sync.Mutex
}

// NewNumberPool is used to allacte a new number pool
func NewNumberPool(min, max int, used []int) *NumberPool {
	return &NumberPool{min: min, max: max, used: used}
}

// Get is used to get a free number
func (n *NumberPool) Get() (int, error) {
	n.Lock()
	defer n.Unlock()

	num, err := GetNumber(n.min, n.max, n.used)
	if err != nil {
		return -1, err
	}
	n.used = append(n.used, num)
	return num, nil
}
