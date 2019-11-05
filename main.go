package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	clip "github.com/atotto/clipboard"
)

func main() {
	var jobType string
	var configName string
	switch len(os.Args[1:]){
	case 1:
		arg1 := os.Args[1]
		if arg1 == "list" || arg1 == "ls" || arg1 == "-ls" || arg1 == "l" || arg1 == "-l"{
			jobType = os.Args[1]
			job(jobType, configName)
		}else{
			help()
		}
	case 2:
		jobType = strings.ToLower(os.Args[1])
		configName = strings.ToLower(os.Args[2])
		job(jobType, configName)
	case 3:
		jobType = strings.ToLower(os.Args[1])
		configName = strings.ToLower(os.Args[2])
		job(jobType, configName)
	default:
		help()
	}
	}



func job(jobType, configName string){
	homedir, _ := os.UserHomeDir()
	kubeDir :=  homedir+"/.kube/configs/"
	configDir := homedir+"/.kube/"
	if Exists(kubeDir) != true{
		os.MkdirAll(kubeDir, os.ModePerm)
	}
	switch jobType {
	case "add", "a", "-a":
		addConfig(kubeDir, configName)
	case "addfile":
		fmt.Println("add file")
		filepath := os.Args[3]
		addConfigFile(kubeDir, configName, filepath)
	case "remove", "rm", "-rm":
		removeConfig(kubeDir, configName)
	case "switch", "s", "-s":
		removeConfig(configDir, "config")
		switchConfig(kubeDir, configName, configDir)
		fmt.Println("Your current config is",configName)
	case "list", "ls", "-ls", "l", "-l":
		ListConfig(kubeDir)
	default:
		help()
	}
}

func addConfig(kubeDir, configName string){
	if Exists(kubeDir+configName) == true {
		removeConfig(kubeDir, configName)
		writeClipboard(kubeDir, configName)
	}else{
		writeClipboard(kubeDir, configName)
	}
}

func addConfigFile(kubeDir, configName, filepath string){
	if Exists(kubeDir+configName){
		removeConfig(kubeDir, configName)
	}
	if Exists(filepath) == true{
		from, err := os.Open(filepath)
		noConfig(err, configName)
		defer from.Close()

		to, err := os.OpenFile(kubeDir+configName, os.O_RDWR|os.O_CREATE, 0666)
		noConfig(err, configName)
		defer to.Close()

		_, err = io.Copy(to, from)
		noConfig(err, configName)
	}else{
		fmt.Println("Config file does not exist")
	}

}

func writeClipboard(kubeDir, configName string){
	content, _ := clip.ReadAll()
	if strings.Contains(content, "apiVersion: v1"){
		f, _ := os.Create(kubeDir+configName)
		defer f.Close()
		f.WriteString(content)
		f.Sync()
		w := bufio.NewWriter(f)
		w.Flush()
	}else{
		fmt.Println("You don't have a valid kubeconfig file in your clipboard")
	}
}

func removeConfig(kubeDir, configName string){
	var err = os.Remove(kubeDir+configName)
	noConfig(err, configName)
}


func switchConfig(kubeDir, configName, configDir string){
	from, err := os.Open(kubeDir+configName)
	noConfig(err, configName)
	defer from.Close()

	to, err := os.OpenFile(configDir+"config", os.O_RDWR|os.O_CREATE, 0666)
	noConfig(err, configName)
	defer to.Close()

	_, err = io.Copy(to, from)
	noConfig(err, configName)
}


func ListConfig(kubeDir string)[]os.FileInfo{
	var files []os.FileInfo

	files, _ = ReadDir(kubeDir)
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return files
}

func help(){
	fmt.Println("")
	fmt.Println("Usage: kc [option] [configname]")
	fmt.Println("Valid Options:")
	fmt.Println("--------------")
	fmt.Println("")
	fmt.Println("Add,	a,  -a 		Add a new config")
	fmt.Println("Addfile")
	fmt.Println("Update,	u,  -u 		Update existing config")
	fmt.Println("Remove,	rm, -rm 	Remove a config")
	fmt.Println("List,	l,  -l 		Add a new config")
	fmt.Println("Switch, s,  -s 		Switch to different config")
	fmt.Println("")
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	list, err := f.Readdir(-1)
	f.Close()
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, err
}

func noConfig(err error, configName string){
	if err != nil {
		fmt.Println("No config named", configName)
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}