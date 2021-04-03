package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	before "github.com/bradleyjkemp/memviz"
	after "github.com/seamia/memviz"
)

type inner struct {
	text string
	data int
}
type tree struct {
	id    int
	left  *tree
	right *tree

	list []*tree

	dict map[string]*tree
	wild interface{}

	inner
}

func (t *tree) Custom() {
	fmt.Println("tree.Custom")
}

type stump struct {
	left  *tree
	right *tree
}

func main() {
	root := &tree{
		id: 0,
		left: &tree{
			id: 1,
		},
		right: &tree{
			id: 2,
		},
	}
	/*
		leaf := &tree{
			id: 3,
		}
		leaf2 := &tree{
			id: 4,
		}

		leaf.list = make([]*tree, 0)
		leaf2.list = make([]*tree, 1)

		for i:=0; i<10; i++ {
			root.list = append(root.list, leaf)
			root.list = append(root.list, leaf2)
			root.list = append(root.list, nil)
		}

		root.dict = make(map[string]*tree)

		root.dict[""] = leaf
		root.dict["foo"] = leaf
		root.dict["barfkdfksdjksjdk.sfskdfskdfksdl;lsflsdflsjgjfjddj\nfldkgd\rjkgjjk54tjkjgf"] = leaf

		root.left.right = leaf
		root.right.left = leaf

		s := stump{
			left:  root,
			right: leaf,
		}
		root.left.wild = &s

		var inter interface{} = *root
		if inter, converts := inter.(custoM); converts {
			inter.Custom()
		}

		inter = &inter
		if inter, converts := inter.(custoM); converts {
			inter.Custom()
		}

	*/

	/*
		buf := &bytes.Buffer{}
		memviz.Map(buf, &root)
		err := ioutil.WriteFile("example-tree-data.dot", buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	*/
	dump(&root, ".dot/new.dot", ".dot/old.dot")
}

func dump(what interface{}, one, two string) error {

	if len(two) > 0 && one != two {
		buf := &bytes.Buffer{}
		before.Map(buf, what)
		err := ioutil.WriteFile(two, buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}

	if len(one) > 0 {
		buf := &bytes.Buffer{}
		after.Map(buf, what, "--name here--")
		err := ioutil.WriteFile(one, buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

type custoM interface {
	Custom()
}
