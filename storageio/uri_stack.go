package storageio

import (
	"fyne.io/fyne/v2"
)

type URIStack struct {
	root     fyne.ListableURI
	elements []fyne.ListableURI
}

func NewStack(base fyne.ListableURI) *URIStack {
	stack := &URIStack{
		root:     base,
		elements: make([]fyne.ListableURI, 0, 8),
	}
	return stack
}

func (stack *URIStack) Push(element fyne.ListableURI) {
	stack.elements = append(stack.elements, element)
}

func (stack *URIStack) Current() (fyne.ListableURI, bool) {
	length := len(stack.elements)
	isRoot := length < 2
	return stack.elements[length-1], isRoot
}

func (stack *URIStack) Pop() (fyne.ListableURI, bool) {
	element, isRoot := stack.Current()
	if !isRoot {
		stack.elements = stack.elements[:len(stack.elements)-1]
	}
	return element, isRoot
}

func (stack *URIStack) Dump() []fyne.ListableURI {
	return stack.elements
}

func (stack *URIStack) Route() (route []string) {
	route = make([]string, len(stack.elements))
	for i, uri := range stack.Dump() {
		route[i] = MakeTitle(uri)
	}
	return
}
