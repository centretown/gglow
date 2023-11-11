package store

import (
	"fyne.io/fyne/v2"
)

type Stack struct {
	root     fyne.ListableURI
	elements []fyne.ListableURI
}

func NewStack(base fyne.ListableURI) *Stack {
	stack := &Stack{
		root:     base,
		elements: make([]fyne.ListableURI, 0, 8),
	}
	return stack
}

func (stack *Stack) Push(element fyne.ListableURI) {
	stack.elements = append(stack.elements, element)
}

func (stack *Stack) Current() (fyne.ListableURI, bool) {
	length := len(stack.elements)
	isRoot := length < 2
	return stack.elements[length-1], isRoot
}

func (stack *Stack) Pop() (fyne.ListableURI, bool) {
	element, isRoot := stack.Current()
	if !isRoot {
		stack.elements = stack.elements[:len(stack.elements)-1]
	}
	return element, isRoot
}

func (stack *Stack) Dump() []fyne.ListableURI {
	return stack.elements
}
