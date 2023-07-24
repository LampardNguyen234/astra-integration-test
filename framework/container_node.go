package framework

import (
	"fmt"
	"strings"
)

type IContainerNode interface {
	ITestNode
}

type ContainerNode struct {
	BaseNode
	children    []ITestNode
	failedCount int
	beforeEach  BeforeEachFunc
	afterEach   AfterEachFunc
}

func (n *ContainerNode) Type() NodeType {
	return ContainerNodeType
}

func (n *ContainerNode) Run() {
	defer n.recover()
	if n.before != nil {
		n.before()
	}

	for _, child := range n.children {
		if n.beforeEach != nil {
			n.beforeEach()
		}
		child.Run()
		if n.afterEach != nil {
			n.afterEach()
		}
		if !child.Succeeded() {
			n.failedCount++
		}
	}

	if n.after != nil {
		n.after()
	}
}

func (n *ContainerNode) Succeeded() bool {
	return n.failedCount == 0
}

func (n *ContainerNode) Logs() []string {
	ret := make([]string, 0)
	if n.Succeeded() {
		ret = append(ret, fmt.Sprintf("%v PASSED", n.text))
	} else {
		ret = append(ret, fmt.Sprintf("%v FAILED %v/%v test(s)", n.text, n.failedCount, len(n.children)))
	}
	for _, child := range n.children {
		tmp := child.Logs()
		for _, d := range tmp {
			ret = append(ret, fmt.Sprintf("\t%v", d))
		}
	}
	return ret
}

func (n *ContainerNode) Report() {
	fmt.Println(strings.Join(n.Logs(), "\n"))
}

func Describe(text string, args ...interface{}) IContainerNode {
	ret := &ContainerNode{
		BaseNode: BaseNode{
			text: text,
		},
	}

	nodes := make([]ITestNode, 0)
	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch arg.(type) {
		case ITestNode:
			nodes = append(nodes, arg.(ITestNode))
		case BeforeFunc:
			ret.before = arg.(BeforeFunc)
		case BeforeEachFunc:
			ret.beforeEach = arg.(BeforeEachFunc)
		case AfterFunc:
			ret.after = arg.(AfterFunc)
		case AfterEachFunc:
			ret.afterEach = arg.(AfterEachFunc)
		default:
			continue
		}
	}
	ret.children = nodes

	return ret
}

func When(text string, args ...interface{}) IContainerNode {
	return Describe(fmt.Sprintf("When %v", text), args...)
}

func Context(text string, args ...interface{}) IContainerNode {
	return Describe(fmt.Sprintf("Context: %v", text), args...)
}
