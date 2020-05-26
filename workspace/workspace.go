package workspace

import (
	"JudgerServer/container"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// WorkSpace 创建Runner工作环境
type WorkSpace struct {
	BaseDir     string
	InputFile   string
	OutputFile  string
	ErrorFile   string
	SourceFile  string
	Language    string
	NeedCompile bool
}

func getPWD() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

// MakeFiles 创建文件与文件夹
func (w *WorkSpace) MakeFiles() {
	if err := os.MkdirAll(w.BaseDir, 0777); err != nil {
		panic(err)
	}
	if _, err := os.Create(w.BaseDir + "/" + w.InputFile); err != nil {
		panic(err)
	}
	if _, err := os.Create(w.BaseDir + "/" + w.OutputFile); err != nil {
		panic(err)
	}
	if _, err := os.Create(w.BaseDir + "/" + w.ErrorFile); err != nil {
		panic(err)
	}
	if _, err := os.Create(w.BaseDir + "/" + w.SourceFile); err != nil {
		panic(err)
	}
	if _, err := os.Create(w.BaseDir + "/output.compile"); err != nil {
		panic(err)
	}
	if _, err := os.Create(w.BaseDir + "/input.compile"); err != nil {
		panic(err)
	}
	if _, err := os.Create(w.BaseDir + "/error.compile"); err != nil {
		panic(err)
	}
}

// WriteInput 封装了对输入文件的写入
func (w *WorkSpace) WriteInput(data []byte) {
	if err := ioutil.WriteFile(w.BaseDir+"/"+w.InputFile, data, 0777); err != nil {
		panic(err)
	}
}

// WriteSourceCode 写入源代码
func (w *WorkSpace) WriteSourceCode(data string) {
	data = strings.Replace(data, " ", "", -1)
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	decodedValue, err := url.QueryUnescape(string(decodeBytes))
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(w.BaseDir+"/"+w.SourceFile, []byte(decodedValue), 0777); err != nil {
		panic(err)
	}
}

// CompileSource 编译源码
func (w *WorkSpace) CompileSource() *container.RunResult {
	c := NewCompiler(w)
	if c == nil {
		panic(errors.New("target compiler not exist"))
	}
	res, err := c.Compile(nil)
	if err != nil {
		panic(err)
	}
	return res
}

// NewWorkSpace 创建一个新的工作环境
func NewWorkSpace(dir string, input string, output string, errorp string, language string, sourceFile string, needCompile bool) *WorkSpace {
	if dir == "" {
		dir = uuid.New().String()
	}
	if input == "" {
		input = "in"
	}
	if output == "" {
		output = "out"
	}
	if errorp == "" {
		errorp = "err"
	}
	if sourceFile == "" {
		sourceFile = "src.compile"
	}
	return &WorkSpace{
		BaseDir:     dir,
		InputFile:   input,
		OutputFile:  output,
		ErrorFile:   errorp,
		SourceFile:  sourceFile,
		Language:    language,
		NeedCompile: needCompile,
	}
}
