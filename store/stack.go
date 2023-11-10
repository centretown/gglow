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

func (stack *Stack) Pop() (fyne.ListableURI, bool) {
	length := len(stack.elements)
	if length == 0 {
		return stack.base, true
	}

	length--
	element := stack.elements[length]
	stack.elements = stack.elements[:length]
	return element, false
}
