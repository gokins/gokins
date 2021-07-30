package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gokins/core/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

var pth string

func main() {
	flag.StringVar(&pth, "d", "", "ui dir")
	flag.Parse()
	if pth == "" {
		pth = "../web/dist"
	}
	err := gengo()
	if err != nil {
		println("bdzip err:" + err.Error())
	}
}

func gengo() error {
	zipfl := filepath.Join(utils.HomePath(), "dist.zip")
	os.RemoveAll(zipfl)
	defer os.RemoveAll(zipfl)
	err := utils.Zip(pth, zipfl, true)
	if err != nil {
		return err
	}
	bts, err := ioutil.ReadFile(zipfl)
	if err != nil {
		return err
	}
	cont := base64.StdEncoding.EncodeToString(bts)
	err = ioutil.WriteFile("comm/uis.go",
		[]byte(fmt.Sprintf("package comm\n\nconst StaticPkg = \"%s\"", cont)),
		0644)
	if err != nil {
		return err
	}
	println("ui insert go ok!!!")
	return nil
}
