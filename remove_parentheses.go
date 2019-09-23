package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//Node has a value string.
type Node struct {
	value string
	left  *Node
	right *Node
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

// precedenceLevel : return operator precedence level
// In this function, if the operant's precedence is higher, then the returned integer is smaller)
func precedenceLevel(op string) int {
	lowerOp := []string{"-", "+"}
	higherOp := []string{"*", "/"}
	if stringInArray(op, lowerOp) {
		return 1
	}
	if stringInArray(op, higherOp) {
		return 2
	}
	return -1
}

// shuntingYardAlgo: refer to "Shunting Yard Algorithm".
// Make the expression to a notation in which operators follow their operands.
func shuntingYardAlgo(arithmeticExpression string) *Stack {
	var opStack = NewStack()     //Stack for operators
	var outputStack = NewStack() //Stack for operands
	var operator = []string{"+", "-", "*", "/"}
	negSign := false

	for index, value := range arithmeticExpression { //Each character in the arithmetic expression
		newNode := Node{value: string(value)}
		// Negative value
		if negSign { //negSign is true - Becuase current node was stored in stack, so continue to the next char
			negSign = false
			continue
		}
		if string(value) == "-" {
			if _, err := strconv.Atoi(string(arithmeticExpression[index+1])); err == nil {
				if index > 0 {
					if _, err := strconv.Atoi(string(arithmeticExpression[index-1])); err != nil {
						negSign = true
						newNode.value = newNode.value + string(arithmeticExpression[index+1])
						outputStack.Push(&newNode)
						continue
					}
				}
			} else {
				re := regexp.MustCompile(`[a-z A-Z]`)
				if re.MatchString(string(arithmeticExpression[index+1])) {
					negSign = true
					newNode.value = newNode.value + string(arithmeticExpression[index+1])
					outputStack.Push(&newNode)
					continue
				}
			}
		}
		if stringInArray(string(value), operator) {
		PrecedenceCondition:
			for {
				if opStack.Toppest() == nil {
					break PrecedenceCondition
				}
				if precedenceLevel(opStack.Toppest().value) < precedenceLevel(string(value)) {
					break PrecedenceCondition
				}
				// top of the operator stack is of lower or equal precedence,add opStack.pop() to output
				outputStack.Push(opStack.Pop())
			}

			opStack.Push(&newNode)
		} else if string(value) == "(" { // left paranthesis
			opStack.Push(&newNode) //push to the operator stack
		} else if string(value) == ")" { //right paranthesis
		ParentheseCondition:
			for {
				if opStack.Toppest() == nil || opStack.Toppest().value == "(" {
					break ParentheseCondition
				}
				//top of the operator stack is not a left paranthesis
				outputStack.Push(opStack.Pop()) //add operator.pop to the output stack
			}
			opStack.Pop()
		} else {
			if index > 0 && !stringInArray(string(arithmeticExpression[index-1]), operator) &&
				string(arithmeticExpression[index-1]) != "(" &&
				string(arithmeticExpression[index-1]) != ")" { //more than 1 digit
				outputStack.Toppest().value = outputStack.Toppest().value + string(arithmeticExpression[index])
			} else {
				outputStack.Push(&newNode)
			}
		}
		// for _, value := range outputStack.node {
		// 	fmt.Printf("outputStack: %s, ", value.value)
		// }
		// fmt.Println("")
		// for _, value := range opStack.node {
		// 	fmt.Printf("opStack: %s, ", value.value)
		// }
		// fmt.Println("")
	}

	for {
		if opStack.Toppest() == nil {
			break
		}
		outputStack.Push(opStack.Pop())
	}
	return outputStack
}

//expressionTree: to represent the expression to a binary expression tree.
func expressionTree(shingYStack *Stack) *Node {
	var expressionTreeStack = NewStack()
	var operator = []string{"+", "-", "*", "/"}
	for _, node := range shingYStack.node {
		newNode := Node{value: string(node.value), left: nil, right: nil}
		if stringInArray(string(node.value), operator) { //op
			newNode.right = expressionTreeStack.Pop()
			newNode.left = expressionTreeStack.Pop()
		}
		expressionTreeStack.Push(&newNode)
	}
	return expressionTreeStack.Toppest()
}

//inorderTraversal: L -> V -> R
func inorderTraversal(expressionTreeStack *Node) string {
	curValue := ""
	leftValue := ""
	rightValue := ""
	var operator = []string{"+", "-", "*", "/"}
	if expressionTreeStack != nil {
		curValue += expressionTreeStack.value
		leftValue += inorderTraversal(expressionTreeStack.left)
		rightValue += inorderTraversal(expressionTreeStack.right)
		if stringInArray(expressionTreeStack.value, operator) { // Right nodes should using
			if i, err := strconv.Atoi(rightValue); err == nil && i < 0 {
				rightValue = "(" + string(rightValue) + ")"
			} else {
				re := regexp.MustCompile(`^-[a-z A-Z]$`)
				if re.MatchString(rightValue) {
					rightValue = "(" + string(rightValue) + ")"
				}
			}
		}
		if expressionTreeStack.value == "-" && // Right nodes should using () if + or -
			(expressionTreeStack.right.value == "+" || expressionTreeStack.right.value == "-") {
			rightValue = "(" + inorderTraversal(expressionTreeStack.right) + ")"
		}
		if expressionTreeStack.value == "*" { // Nodes of both sides should using () if + or -
			if expressionTreeStack.right.value == "+" || expressionTreeStack.right.value == "-" {
				rightValue = "(" + rightValue + ")"
			}
			if expressionTreeStack.left.value == "+" || expressionTreeStack.left.value == "-" {
				leftValue = "(" + leftValue + ")"
			}
			if i, err := strconv.Atoi(rightValue); err == nil && i < 0 {
				rightValue = "(" + string(rightValue) + ")"
			}
		}
		if expressionTreeStack.value == "/" { // Nodes of both sides should using () if + or -, // Right nodes should using () if *
			if expressionTreeStack.right.value == "+" || expressionTreeStack.right.value == "-" {
				rightValue = "(" + rightValue + ")"
			}
			if expressionTreeStack.left.value == "+" || expressionTreeStack.left.value == "-" {
				leftValue = "(" + leftValue + ")"
			}
			if expressionTreeStack.right.value == "*" || expressionTreeStack.right.value == "/" {
				rightValue = "(" + rightValue + ")"
			}
		}
	}
	return leftValue + curValue + rightValue
}
func f(arithmeticExpression string) string {
	expression := strings.Replace(arithmeticExpression, " ", "", -1) //Trim spaces in string
	shutingYardOutput := shuntingYardAlgo(expression)
	expressionTopNode := expressionTree(shutingYardOutput)
	simplyExpression := inorderTraversal(expressionTopNode)

	return simplyExpression
}
func main() {
	test := []string{
		"a/(b*-c)",
		"a/(b/c)",
		"a*(b*c)",
		"a*(b/c)",
		"(d/a)*(b+c)",
		"(d*a)/(b+c)",
		"(d*a)*(b+c)",
		"2*(1-3)/(1/2)",
		"2*(1-3)*(1/2)",
		"(a*b)*(c/d)",
		"2*(1-3)",
		"(((-1+(2*(-1-(-2))))))",
		"(1+(2))",
		"x+(y+z)+(t+a+(v+w))",
		"2-(2+3)",
		"(2*(3+4)*5)/6",
		"1*(2+(3*(4+5)))",
		"2 + (3 / -5)",
		"x+(y+z)+(t+(v+w))",
		"-6+(3*(x+(y*z)))",
		"2*(2+3-(4*6))+8+7*4", //unpassed start
		"-(2)-(2+3)",
		"-(2+3)",
		"1+(-1)",
		"((2*((2+3)-(4*6))+(8+(7*4))))",
		"((2*((2*3)-(4+6))+(8+(7*4))))",
		"1-(-1)",
		"1*(-1)",
		"1/(-1)",
	}

	for _, testItem := range test {
		fmt.Println("Before: " + testItem + ", Result: " + f(testItem))
	}
}
