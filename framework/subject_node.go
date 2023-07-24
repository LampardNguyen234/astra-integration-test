package framework

type ISubjectNode interface {
	ITestNode
}

type SubjectNode struct {
	BaseNode
	run ProcessFunc
}

func (n *SubjectNode) Type() NodeType {
	return SubjectNodeType
}

func (n *SubjectNode) Run() {
	if n.before != nil {
		n.before()
	}
	defer n.recover()
	if n.run != nil {
		n.run()
	}
	if n.after != nil {
		n.after()
	}
}

func It(text string, args ...interface{}) ITestNode {
	ret := &SubjectNode{
		BaseNode: BaseNode{
			text: text,
		},
	}
	for _, arg := range args {
		if arg == nil {
			continue
		}
		switch arg.(type) {
		case ProcessFunc:
			ret.run = arg.(ProcessFunc)
		case AfterFunc:
			ret.after = arg.(AfterFunc)
		case BeforeFunc:
			ret.before = arg.(BeforeFunc)
		case func():
			ret.run = Process(arg.(func()))
		default:
			continue
		}
	}

	return ret
}
