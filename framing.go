package main

import "fmt"

type Frame struct {
	List []*Frame
	Data []float32
}

func findKids(f *Frame, level int) {
	fmt.Print("[")
	if len(f.List) == 0 {
		fmt.Print(f.Data)
	} else {
		for _, kid := range f.List {
			findKids(kid, level+1)
		}
	}
	if len(f.List) == 0 {
		fmt.Print("]\n")
	}
}

func (f *Frame) String() string {
	findKids(f, 0)
	fmt.Print("]]\n")
	return ""
}

func frame(y []float32, frameLength, hopLength int) *Frame {

	froot := Frame{}
	froot.List = []*Frame{}

	w := Frame{}
	w.List = []*Frame{}
	framePart(&w, 2, len(y)-1)
	froot.List = append(froot.List, &w)

	return &froot
}
func framePart(w *Frame, size, many int) {

	items := []*Frame{}
	for {
		f := Frame{}
		f.List = []*Frame{}
		f.Data = []float32{}
		for {
			f.Data = append(f.Data, 1)
			if len(f.Data) == many {
				break
			}
		}
		items = append(items, &f)
		if len(items) == size {
			break
		}
	}

	for _, f := range items {
		w.List = append(w.List, f)
	}
}
