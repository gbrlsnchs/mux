package radix_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/gbrlsnchs/mux/internal/mocks"
	. "github.com/gbrlsnchs/mux/internal/radix"
	"github.com/golang/mock/gomock"
)

func TestTree(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	handler := NewMockHandler(mockCtrl)
	testTable := []struct {
		paths  []string
		reqs   []string
		params map[string][]byte
		expect bool
	}{
		{
			paths: []string{
				"/test",
			},
			reqs: []string{
				"/test",
			},
			expect: true,
		},
		{
			paths: []string{
				"/test",
				"/testing",
			},
			reqs: []string{
				"/test",
				"/testing",
			},
			expect: true,
		},
		{
			paths: []string{
				"/test",
				"/team",
				"/testing",
			},
			reqs: []string{
				"/test",
				"/team",
				"/testing",
			},
			expect: true,
		},
		{
			paths: []string{
				"/test/:name",
			},
			reqs: []string{
				"/test/hello",
			},
			params: map[string][]byte{
				"name": []byte("hello"),
			},
			expect: true,
		},
		{
			paths: []string{
				"/test/:name1/:name2/:name3",
			},
			reqs: []string{
				"/test/hello1/hello2/hello3",
			},
			params: map[string][]byte{
				"name1": []byte("hello1"),
				"name2": []byte("hello2"),
				"name3": []byte("hello3"),
			},
			expect: true,
		},
		{
			paths: []string{
				"/test/:name/testing",
			},
			reqs: []string{
				"/test/hello/testing",
			},
			params: map[string][]byte{
				"name": []byte("hello"),
			},
			expect: true,
		},
		{
			paths: []string{
				"",
			},
			reqs: []string{
				"",
			},
			expect: false,
		},
	}
	for _, tt := range testTable {
		t.Run(fmt.Sprintf("%#v", tt.reqs), func(t *testing.T) {
			tree := New()
			// First, add all labels to the tree,
			// forcing every kind of insertion case.
			for i := range tt.paths {
				tree.Add([]byte(tt.paths[i]), handler)
			}

			pmap := make(map[string][]byte)
			for i := range tt.reqs {
				n, params := tree.Get([]byte(tt.reqs[i]))
				if want, got := tt.expect, n != nil; want != got {
					t.Errorf("want %t, got %t", want, got)
				}
				for k := range params {
					pmap[k] = params[k]
				}
			}
			if len(pmap) == 0 {
				pmap = nil
			}
			if want, got := tt.params, pmap; !reflect.DeepEqual(want, got) {
				t.Errorf("want %#v, got %#v", want, got)
			}
		})
	}
}
