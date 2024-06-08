package main

import (
	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	graph := gorgonnx.NewGraph()
	model := onnx.NewModel(graph)
	b, _ := ioutil.ReadFile("example/onnx_model/onnx_file/resnet-preproc-v1-18.onnx")
	err := model.UnmarshalBinary(b)
	if err != nil {
		logrus.Errorf("init onnx model fail, err = %v", err)
		return
	}
	imgBytes, err := ioutil.ReadFile("onnx_file/dog.jpeg")
	if err != nil {
		logrus.Errorf("init img file fail, err = %v", err)
		return
	}
	tensorInput, err := onnx.NewTensor(imgBytes)
	if err != nil {
		logrus.Errorf("init tensor fail, err = %v", err)
		return
	}
	model.SetInput(0, tensorInput)
	err = graph.Run()
	if err != nil {
		logrus.Errorf("graph run fail, err = %v", err)
		return
	}
}

func GetCurrentDirectory() string {
	//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Errorf("err = %v", err)
	}

	//将\替换成/
	return strings.Replace(dir, "\\", "/", -1)
}
