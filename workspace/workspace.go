package workspace

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// WorkSpace 创建Runner工作环境
type WorkSpace struct {
	BaseDir    string
	InputFile  string
	OutputFile string
	ErrorFile  string
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
	if err := os.MkdirAll(getPWD()+"/"+w.BaseDir, 0777); err != nil {
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
}

// WriteInput 封装了对输入文件的写入
func (w *WorkSpace) WriteInput(data []byte) {
	if err := ioutil.WriteFile(w.BaseDir+"/"+w.InputFile, data, 0777); err != nil {
		panic(err)
	}
}

// NewWorkSpace 创建一个新的工作环境
func NewWorkSpace(dir string, input string, output string, errorp string) *WorkSpace {
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
	return &WorkSpace{
		BaseDir:    dir,
		InputFile:  input,
		OutputFile: output,
		ErrorFile:  errorp,
	}
}
