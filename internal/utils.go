package internal

import (
	"fmt"
	"os"
	// "path"
	"path/filepath"
	// "sync"
)

// TODO: allow use of $HOME or ~, might already be implemnted idk...

const dataPath = "/Users/anthony/.local/share/ProjectManager/projects.yaml"
const dataPathTemp = "/Users/anthony/.local/share/ProjectManager/projects.yaml.temp"

// TODO: make this faster with fixed buff size
func (proj Project) toYamlBytes() []byte {
	outByteArr := []byte{}
	outByteArr = append(outByteArr, []byte{0x20, 0x20, 0x2d, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x20}...)
	outByteArr = append(outByteArr, []byte(proj.name)...)
	outByteArr = append(outByteArr, []byte{0x0a, 0x20, 0x20, 0x20, 0x20, 0x70, 0x61, 0x74, 0x68, 0x3a, 0x20}...)
	outByteArr = append(outByteArr, []byte(proj.path)...)
	outByteArr = append(outByteArr, []byte{0x0a, 0x20, 0x20, 0x20, 0x20, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x3a, 0x20}...)
	outByteArr = append(outByteArr, []byte(proj.language)...)
	outByteArr = append(outByteArr, []byte{0x0a, 0x20, 0x20, 0x20, 0x20, 0x6c, 0x61, 0x73, 0x74, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x3a, 0x20}...)
	if proj.active {
		outByteArr = append(outByteArr, []byte("true")...)
	} else {
		outByteArr = append(outByteArr, []byte("false")...)
	}
	outByteArr = append(outByteArr, []byte("\n")...)
	return outByteArr
}

func retrieveData() []byte {
	stat, err := os.Stat(dataPath)
	if os.IsNotExist(err) {
		p := filepath.Join(dataPath)
		if _, err := os.Stat(filepath.Dir(p)); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(p), 0700)
		}
		_, _ = os.Create(dataPath)
		err = os.WriteFile(dataPath, []byte("projects:\n"), 0644)
		if err != nil {
			fmt.Println("could not write base to yaml data file")
			os.Exit(1)
		}
	}
	if stat.Size() == 0 {
		err = os.WriteFile(dataPath, []byte("projects:\n"), 0644)
	}
	file, err := os.ReadFile(dataPath)
	if err != nil {
		fmt.Println("error reading data file: %w")
		os.Exit(1)
	}
	return file
}

func makeProject() Project {
	var proj Project
	var in string
	// Name
	fmt.Print("Name: ")
	fmt.Scanln(&in)
	proj.name = in
	in = ""
	// Path
	pathGood := false
	for {
		fmt.Print("Path (Leave blank for current dir or a path starting with . for a relative path): ")
		fmt.Scanln(&in)
		if in == "" {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Println("Could not get current working direcrtory, path will probably be inaccurate")
			} else {
				proj.path = wd
				pathGood = true
			}
		} else if in[0] == '.' {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Println("Could not get current working direcrtory, path will probably be inaccurate")
			} else {
				proj.path = fmt.Sprintf("%s/%s", wd, in[1:])
				pathGood = true
			}
		} else {
			ad, err := filepath.Abs(in)
			if err != nil {
				fmt.Println("Could not get path from your input")
			} else {
				proj.path = ad
				pathGood = true
			}
		}
		_, sErr := os.Stat(proj.path)
		if os.IsNotExist(sErr) {
			fmt.Println("provided path does not exist")
			pathGood = false
		}
		if pathGood {
			break
		}
	}
	in = ""
	// Language
	fmt.Print("language: ")
	fmt.Scanln(&in)
	proj.language = in
	in = ""
	// Last modified
	proj.FindLastModified()

	// Active
	for {
		fmt.Print("active (y or n): ")
		fmt.Scanln(&in)
		in = ""
		if in == "Y" || in == "y" {
			proj.active = true
			break
		} else if in == "N" || in == "n" {
			proj.active = false
			break
		}
	}
	return proj
}

func (proj Project) crawlDir() {
	// TODO: find last modified
	// TODO: find todo, fixme, etc
	// goroutines!!!!!!!!!!!

}

func findLastModified() {
	var lastModified string
	lastModifiedInt := 0
	files, err := os.ReadDir(proj.path)
	if err != nil {
		fmt.Println("error with reading directory")
		os.Exit(1)
	}
	for _, f := range files {
		fi, err := os.Stat(proj.path + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = f.Name()
		}
	}

	proj.lastModified = lastModified
}
