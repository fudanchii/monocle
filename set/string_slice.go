package set

import (
	"github.com/deckarep/golang-set"
)

type Set interface {
	mapset.Set
	ToStringSlice() []string
}

type mset struct {
	mapset.Set
}

func NewSet(s ...interface{}) Set {
	rset := mapset.NewSetFromSlice(s)
	return &mset{rset}
}

func (ms *mset) ToStringSlice() []string {
	var resset []string
	it := ms.Iterator()
	for elt := range it.C {
		if rs, ok := elt.(string); ok {
			resset = append(resset, rs)
		}
	}
	return resset
}
