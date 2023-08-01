package framework

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"runtime/debug"
	"strings"
)

var (
	logger log.Logger
)

var nodeLog struct {
	level int
}

type ITestNode interface {
	Text() string
	Run()
	Err() error
	Succeeded() bool
	Type() NodeType
	Logs() []string
	Report()
}

type BaseNode struct {
	text      string
	err       error
	processed bool
	before    BeforeFunc
	after     AfterFunc
}

func (n *BaseNode) Text() string {
	return n.text
}

func (n *BaseNode) Err() error {
	return n.err
}

func (n *BaseNode) Processed() bool {
	return n.processed
}

func (n *BaseNode) Succeeded() bool {
	return n.processed && n.err == nil
}

func (n *BaseNode) Logs() []string {
	if !n.processed {
		return []string{fmt.Sprintf("%v NOT PROCESSED", n.text)}
	}
	if n.Succeeded() {
		return []string{fmt.Sprintf("%v PASSED", n.text)}
	}

	return []string{fmt.Sprintf("%v FAILED", n.text)}
}

func (n *BaseNode) Report() {
	fmt.Println(strings.Join(n.Logs(), "\n"))
}

func (n *BaseNode) recover() {
	n.processed = true
	if r := recover(); r != nil {
		n.err = fmt.Errorf("%v\n\t%v", r, string(debug.Stack()))
		color.Red("%v FAILED\n\t\t%v\n", n.text, n.err)
	} else {
		color.Green("%v PASSED\n", n.text)
	}

}
