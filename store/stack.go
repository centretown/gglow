package store

import (
	"fyne.io/fyne/v2"
)

type Stack struct {
	base     fyne.ListableURI
	elements []fyne.ListableURI
}

func NewStack(base fyne.ListableURI) *Stack {
	stack := &Stack{
		base:     base,
		elements: make([]fyne.ListableURI, 0, 8),
	}
	return stack
}

func (stack *Stack) Push(element fyne.ListableURI) {
	stack.elements = append(stack.elements, element)
}

func (stack *Stack) Current() (fyne.ListableURI, bool) {
	length := len(stack.elements)
	isBase := length < 2
	return stack.elements[length-1], isBase
}

func (stack *Stack) Pop() (fyne.ListableURI, bool) {
	element, isBase := stack.Current()
	if !isBase {
		stack.elements = stack.elements[:len(stack.elements)-1]
	}
	return element, isBase
}
