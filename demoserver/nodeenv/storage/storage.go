package storage

import (
	"context"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/velavokr/dsplayground/ifaces"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"path/filepath"
)

func NewStorage(rt *runner.Runtime) ifaces.Storage {
	return &storage{rt: rt, tables:make(map[string]ifaces.DiskTable)}
}

func (s *storage) OpenTable(path string) ifaces.DiskTable {
	ot, ok := s.tables[path]
	if ok {
		return ot
	}

	f := filepath.Join(s.rt.Cfg.DbDir, fmt.Sprintf("%s.%d", path, s.rt.Cfg.Self))
	t := &table{
		fname: f,
		rt:    s.rt,
	}
	s.rt.Run(func() {
		db, err := leveldb.OpenFile(f, &opt.Options{Strict:opt.StrictAll})
		if err != nil && errors.IsCorrupted(err) {
			db, err = leveldb.RecoverFile(f, nil)
			if err != nil {
				panic(err)
			}
		}
		s.rt.RunAsync(func(ctx context.Context) {
			select {
			case <-ctx.Done():
				s.rt.RunGuarded(func() {
					t.db = nil
					if err := db.Close(); err != nil {
						panic(err)
					}
				}, runner.ExitOnPanic|runner.VerboseLog,"closing table", f)
			}
		}, runner.ExitOnPanic|runner.VerboseLog, "table closer", f)
		t.db = db
	}, runner.ExitOnPanic|runner.VerboseLog, "open table", f)
	s.tables[path] = t
	return t
}

func (t *table) StoreValue(key []byte, val []byte) {
	t.rt.Run(func() {
		err := t.db.Put(key, val, &opt.WriteOptions{Sync:true})
		if err != nil {
			panic(err)
		}
	}, runner.VerboseLog|runner.ExitOnPanic, "store to", t.fname, key, "=", val)
}

func (t *table) LoadValue(key []byte) []byte {
	var res []byte
	t.rt.Run(func() {
		val, err := t.db.Get(key, nil)
		if err != nil && err != leveldb.ErrNotFound {
			panic(err)
		}
		res = val
	}, runner.VerboseLog|runner.ExitOnPanic, "load from", t.fname, key)
	return res
}

func (t *table) DeleteKey(key []byte) {
	t.rt.Run(func() {
		if err := t.db.Delete(key, nil); err != nil {
			panic(err)
		}
	}, runner.VerboseLog|runner.ExitOnPanic, "delete from", t.fname, key)
}

func (t *table) LoadKeys() [][]byte {
	keys := make([][]byte, 0, 1)
	t.rt.Run(func() {
		iter := t.db.NewIterator(nil, nil)
		for iter.Next() {
			buf := make([]byte, len(iter.Key()))
			copy(buf, iter.Key())
			keys = append(keys, buf)
		}
	}, runner.VerboseLog|runner.ExitOnPanic, "load keys from", t.fname)
	return keys
}

type table struct {
	fname string
	rt    *runner.Runtime
	db    *leveldb.DB
}

type storage struct {
	rt *runner.Runtime
	tables map[string]ifaces.DiskTable
}
