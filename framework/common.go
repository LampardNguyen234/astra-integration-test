package framework

type NodeType uint

const (
	ContainerNodeType NodeType = iota
	SubjectNodeType
)

type ProcessFunc func()
