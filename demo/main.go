package main

import (
	"fmt"
	"github.com/AndyGreenwell94/docxt"
)

type TestStruct struct {
	FileName string
	Data     map[string]string
	Items    []TestItemStruct
	Items1   []TestItemStruct
}

type TestItemStruct struct {
	Column1 string
	Column2 string
	Column3 string
	Column4 string
}

func main() {
	template, err := docxt.OpenTemplate("./demo/example.docx")
	if err != nil {
		fmt.Println(err)
		return
	}
	test := new(TestStruct)
	test.FileName = "example.docx"
	test.Items = []TestItemStruct{
		{"1", "2", "5", "6"},
		{"3", "4", "7", "8"},
	}
	test.Items1 = []TestItemStruct{
		{"6", "6", "6", "6"},
		{"6", "6", "6", "6"},
	}
	test.Data = map[string]string{
		"S1": "123",
		"S2": "321",
	}

	if err := template.RenderTemplate(test); err != nil {
		fmt.Println(err)
		return
	}
	if err := template.Save("./demo/result.docx"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Success")
}
