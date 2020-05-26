package workspace

import (
	"JudgerServer/container"
	"strings"
)

// AbstractCompiler 编译抽象 不要直接初始化，使用NewCompiler初始化.
type AbstractCompiler struct {
	Src     string
	Dest    string
	WS      *WorkSpace
	Compile func(a *interface{}) (*container.RunResult, error)
}

// Compile 编译
// func (c *AbstractCompiler) Compile(a *interface{}) (*container.RunResult, error) { return nil, nil }

// NewCompiler 返回一个对应语言的编译器对象. 如果该语言不存在则返回nil
func NewCompiler(workspace *WorkSpace) *AbstractCompiler {
	workspace.Language = strings.ToLower(string(workspace.Language))
	if workspace.Language == "cpp" {
		c := &CppCompiler{}
		c.Src = workspace.SourceFile
		c.Dest = "binary.compile"
		c.WS = workspace
		c.AbstractCompiler.Compile = c.Compile
		return &c.AbstractCompiler
	}
	return nil
}
