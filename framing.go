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
	fmt.Print("]\n")
	return ""
}

func frame(y []float32, frameLength, hopLength int) *Frame {

	froot := Frame{}
	froot.List = []*Frame{}

	f1 := Frame{}
	f1.List = []*Frame{}
	f2 := Frame{}
	f2.List = []*Frame{}

	f1.Data = []float32{1, 2}
	f2.Data = []float32{4, 5}

	wrapper := Frame{}
	wrapper.List = []*Frame{}
	wrapper.List = append(wrapper.List, &f1)
	wrapper.List = append(wrapper.List, &f2)

	froot.List = append(froot.List, &wrapper)
	return &froot
}
