package main

import (
	"fmt"
	"strings"
)

//Node has a value string.
type Node struct {
	value string
}

//Stack -  LIFO
type Stack struct {
	node  []*Node
	count int
}

// Push - adds a node to the stack.
func (s *Stack) Push(n *Node) {
	s.node = append(s.node[:s.count], n)
	s.count++
}

// Pop - removes last node.
func (s *Stack) Pop() *Node {
	if s.count == 0 {
		return nil
	}
	popItem := s.Toppest()

	s.count--
	s.node = s.node[:s.count]
	return popItem
}

// Toppest - toppest node of the stack.
func (s *Stack) Toppest() *Node {
	if s.count == 0 {
		return nil
	}
	return s.node[s.count-1]
}

// NewStack -  careate a new stack.
func NewStack() *Stack {
	return &Stack{}
}

//stringInArray - check target string exists in list or not
func stringInArray(target string, list []string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// f function: remove unneccessary parentheses
func f(arithmeticExpression string) string {
	var operator = []string{"+", "-", "*", "/"}
	var operatorHigh = []string{"*", "/"}
	var operatorLow = []string{"+", "-"}
	var operatorStack = NewStack() //Stack for operators
	var operandStack = NewStack()  //Stack for operands
	negSign := false
	expression := strings.Replace(arithmeticExpression, " ", "", -1) //Trim spaces in string

	for index, value := range expression { //Each character in the arithmetic expression
		newNode := Node{value: string(value)}

		if negSign { //negSign is true - Becuase current node was stored in stack, so continue to the next char
			negSign = false
			continue
		}
		if stringInArray(newNode.value, operator) || newNode.value == "(" {
			if newNode.value == "-" && (index == 0 || // If "-" is negative sign, push this sign and next char
				stringInArray(string(expression[index-1]), operator) ||
				string(expression[index-1]) == "(") {
				negSign = true
				newNode.value = newNode.value + string(expression[index+1])
				operandStack.Push(&newNode)
				continue
			}
			operatorStack.Push(&newNode)
		} else if newNode.value == ")" {
			if stringInArray(operatorStack.Toppest().value, operatorHigh) {
				fisrtNum := operandStack.Pop()
				newExpression := Node{value: operandStack.Pop().value + operatorStack.Pop().value + fisrtNum.value}
				operatorStack.Pop() //Remove "("
				operandStack.Push(&newExpression)

			} else if stringInArray(operatorStack.Toppest().value, operatorLow) {
				fisrtNum := operandStack.Pop()
				newExpression := Node{value: operandStack.Pop().value + operatorStack.Pop().value + fisrtNum.value}
				lastOperator := operatorStack.Pop() //Remove "("
				if lastOperator != nil && stringInArray(operatorStack.Toppest().value, operatorHigh) {
					newExpression.value = "(" + newExpression.value + ")"
				}
				operandStack.Push(&newExpression)
			} else {
				panic("The input string is invalid expression.")
			}
		} else {
			var notNumbers = []string{"+", "-", "*", "/", ")", "("}
			if index > 0 && !stringInArray(string(expression[index-1]), notNumbers) { //more than 1 digit
				operandStack.Toppest().value = operandStack.Toppest().value + string(expression[index])
			} else {
				operandStack.Push(&newNode)
			}
		}
	}

	operatorTurn := false
	result := []string{}
	for {
		if operatorStack.count == 0 && operandStack.count == 0 {
			break
		}
		if operatorTurn {
			result = append(result, operatorStack.Toppest().value)
			operatorStack.Pop()
			operatorTurn = false
		} else {
			result = append(result, operandStack.Toppest().value)
			operandStack.Pop()
			operatorTurn = true
		}
	}

	resExpression := ""
	for i := (len(result) - 1); i >= 0; i-- {
		resExpression += result[i]
	}
	return resExpression
}

func main() {
	test := []string{
		"1*(2+(3*(4+5)))",
		"2 + (3 / -5)",
		"x+(y+z)+(t+(v+w))",
		"-6+(3*(x+(y*z)))",
	}

	for _, testItem := range test {
		fmt.Println("Before: " + testItem + ", Result: " + f(testItem))
	}
}
