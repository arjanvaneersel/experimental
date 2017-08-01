package gockan

import (
	"errors"
	"fmt"
	"sync"

	"github.com/arjanvaneersel/gockan/model"
)

type MemRepo struct {
	mutex    sync.RWMutex
	packages map[string]*model.Package
}

func NewMemRepo() (repo Repository) {
	mr := &MemRepo{}
	mr.packages = make(map[string]*model.Package)
	return mr
}

func (mr *MemRepo) Count() int {
	return len(mr.packages)
}

func (mr *MemRepo) Packages() (ch chan *model.Package, err error) {
	ch = make(chan *model.Package)
	go func() {
		mr.mutex.RLock()
		for _, pkg := range mr.packages {
			ch <- pkg
		}
		mr.mutex.RUnlock()
		close(ch)
	}()
	return
}

func (mr *MemRepo) GetPackage(id string) (pkg *model.Package, err error) {
	mr.mutex.RLock()
	pkg, ok := mr.packages[id]
	mr.mutex.RUnlock()
	if !ok {
		errs := fmt.Sprintf("no such package: %s", id)
		err = errors.New(errs)
	}
	return
}

func (mr *MemRepo) PutPackage(pkg *model.Package) (err error) {
	mr.mutex.Lock()
	mr.packages[pkg.Id] = pkg
	mr.mutex.Unlock()
	return
}
