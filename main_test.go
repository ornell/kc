package main

import (
	"bufio"
	"os"
	"testing"
)

func TestExists(t *testing.T){
	folder  := createTestFolder()
	exists := Exists(folder)
	if exists != true{
	t.Errorf("Expected result to be true while actuall result was %t", exists)
	}
}

func TestListConfig(t *testing.T){
	testFolder := createTestFolder()
	createTestConfig(testFolder)

	output := ListConfig(testFolder)
	for _, output := range output {
		t.Errorf("Expected result to be testconfig while actuall result was %s", output.Name())
	}

}

func createTestFolder()string{
	testFolder := "/tmp/testfolder"
	os.MkdirAll(testFolder, os.ModePerm)
return testFolder
}

func createTestConfig(testFolder string){
	content := "teststring"
	f, _ := os.Create(testFolder+"testfile")
	defer f.Close()
	f.WriteString(content)
	f.Sync()
	w := bufio.NewWriter(f)
	w.Flush()
}