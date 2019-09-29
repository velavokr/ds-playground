package storage

import (
	"bytes"
	"fmt"
	"github.com/velavokr/dsplayground/demoserver/nodeenv"
	"github.com/velavokr/dsplayground/demoserver/runner"
	"github.com/velavokr/dsplayground/demoserver/utils"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"testing"
)

func TestStorage(t *testing.T) {
	dirname, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		if err := os.RemoveAll(dirname); err != nil {
			t.Fatal(err.Error())
		}
	}()

	func() {
		b := bytes.Buffer{}
		rt := runner.NewRuntime(runner.UserCfg{DbDir: dirname}, &b)
		st := nodeenv.NewNodeEnv(rt, NewStorage).Storage()
		t00 := st.OpenTable("test-0")
		t01 := st.OpenTable("test-0")
		t1 := st.OpenTable("test-1")
		t00.StoreValue([]byte("k-00"), []byte("v-00"))
		t01.StoreValue([]byte("k-01"), []byte("v-01"))
		t1.StoreValue([]byte("k-10"), []byte("v-10"))
		t1.StoreValue([]byte("k-11"), []byte("v-11"))
		rt.Cancel()
		rt.WaitAll()
	}()

	runtime.GC()

	func() {
		b := bytes.Buffer{}
		rt := runner.NewRuntime(runner.UserCfg{DbDir: dirname}, &b)
		st := nodeenv.NewNodeEnv(rt, NewStorage).Storage()
		for i := 0; i < 2; i++ {
			tn := fmt.Sprintf("test-%d", i)
			tbl := st.OpenTable(tn)
			k := tbl.LoadKeys()
			sort.Slice(k, utils.Less(k))
			if len(k) != 2 {
				t.Error(k)
			}
			for j := 0; j < 2; j++ {
				suff := fmt.Sprintf("-%d%d", i, j)
				if string(k[j]) != "k"+suff {
					t.Error(tn, j, k)
				}
				if v := tbl.LoadValue([]byte("k" + suff)); string(v) != "v"+suff {
					t.Error(tn, "k"+suff, v)
				}
			}
		}
		rt.Cancel()
		rt.WaitAll()
	}()

	runtime.GC()

	func() {
		b := bytes.Buffer{}
		rt := runner.NewRuntime(runner.UserCfg{DbDir: dirname}, &b)
		st := nodeenv.NewNodeEnv(rt, NewStorage).Storage()
		for i := 0; i < 2; i++ {
			tn := fmt.Sprintf("test-%d", i)
			tbl := st.OpenTable(tn)
			tbl.DeleteKey([]byte(fmt.Sprintf("k-%d0", i)))
		}
		rt.Cancel()
		rt.WaitAll()
	}()

	runtime.GC()

	func() {
		b := bytes.Buffer{}
		rt := runner.NewRuntime(runner.UserCfg{DbDir: dirname}, &b)
		st := nodeenv.NewNodeEnv(rt, NewStorage).Storage()
		for i := 0; i < 2; i++ {
			tn := fmt.Sprintf("test-%d", i)
			tbl := st.OpenTable(tn)
			k := tbl.LoadKeys()
			sort.Slice(k, utils.Less(k))
			if len(k) != 1 {
				t.Error(k)
			}
			suff := fmt.Sprintf("-%d1", i)
			if string(k[0]) != "k"+suff {
				t.Error(tn, k)
			}
			if v := tbl.LoadValue([]byte("k" + suff)); string(v) != "v"+suff {
				t.Error(tn, "k"+suff, v)
			}
		}
		rt.Cancel()
		rt.WaitAll()
	}()

	runtime.GC()

	func() {
		b := bytes.Buffer{}
		rt := runner.NewRuntime(runner.UserCfg{DbDir: dirname}, &b)
		st := nodeenv.NewNodeEnv(rt, NewStorage).Storage()
		for i := 0; i < 2; i++ {
			tn := fmt.Sprintf("test-%d", i)
			tbl := st.OpenTable(tn)
			tbl.DeleteKey([]byte(fmt.Sprintf("k-%d1", i)))
		}
		rt.Cancel()
		rt.WaitAll()
	}()

	runtime.GC()

	func() {
		b := bytes.Buffer{}
		rt := runner.NewRuntime(runner.UserCfg{DbDir: dirname}, &b)
		st := nodeenv.NewNodeEnv(rt, NewStorage).Storage()
		for i := 0; i < 2; i++ {
			tn := fmt.Sprintf("test-%d", i)
			tbl := st.OpenTable(tn)
			k := tbl.LoadKeys()
			sort.Slice(k, utils.Less(k))
			if len(k) != 0 {
				t.Error(k)
			}
		}
		rt.Cancel()
		rt.WaitAll()
	}()
}
