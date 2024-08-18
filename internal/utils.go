package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"path/filepath"
)

// TODO: allow use of $HOME or ~, might already be implemnted idk...
// TODO: changed with config
const dataPath = "/Users/anthony/.local/share/projectmanager/projects.yaml"
const dataPathTemp = "/Users/anthony/.local/share/projectmanager/projects.yaml.temp"

func (proj Project) newEntrytoYamlBytes() []byte {
	length := 60
	if proj.active {
		length += len("true\n")
	} else {
		length += len("false\n")
	}
	length += len(proj.name) + len(proj.desc) + len(proj.path) + len(proj.language)
	buf := make([]byte, length)
	offset := 0
	copy(buf[offset:], "  - name: ")
	offset += len("  - name: ")
	copy(buf[offset:], []byte(proj.name))
	offset += len([]byte(proj.name))
	copy(buf[offset:], "\n    desc: ")
	offset += len("\n    desc: ")
	copy(buf[offset:], []byte(proj.desc))
	offset += len([]byte(proj.desc))
	copy(buf[offset:], "\n    path: ")
	offset += len("\n    path: ")
	copy(buf[offset:], []byte(proj.path))
	offset += len([]byte(proj.path))
	copy(buf[offset:], "\n    language: ")
	offset += len("\n    language: ")
	copy(buf[offset:], []byte(proj.language))
	offset += len([]byte(proj.language))
	copy(buf[offset:], "\n    active: ")
	offset += len("\n    active: ")
	if proj.active {
		copy(buf[offset:], "true\n")
		offset += len("true\n")
	} else {
		copy(buf[offset:], "false\n")
		offset += len("false\n")
	}
	return buf
}

// TODO: make me faster!
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
	in := bufio.NewReader(os.Stdin)
	// Name
	fmt.Print("Name: ")
	line, _ := in.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	proj.name = line
	// desc
	fmt.Print("Description: ")
	line, _ = in.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	proj.desc = line
	// Path
	pathGood := false
	for !pathGood {
		fmt.Print("Path (Leave blank for current dir or a path starting with . for a relative path): ")
		line, _ = in.ReadString('\n')
		if line == "\n" {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Println("Could not get current working direcrtory, path will probably be inaccurate")
			} else {
				proj.path = wd
				pathGood = true
			}
		} else if line[0] == '.' {
			// TODO: Test me
			wd, err := os.Getwd()
			if err != nil {
				fmt.Println("Could not get current working direcrtory, path will probably be inaccurate")
			} else {
				proj.path = fmt.Sprintf("%s/%s", wd, line[1:])
				pathGood = true
			}
		} else {
			ad, err := filepath.Abs(line)
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
		// fmt.Println(in)
	}
	// Language
	fmt.Print("language: ")
	line, _ = in.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")
	proj.language = line
	// Active
	for {
		fmt.Print("active (y or n): ")
		line, _ = in.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if line == "Y" || line == "y" {
			proj.active = true
			break
		} else if line == "N" || line == "n" {
			proj.active = false
			break
		}
	}
	return proj
}

func crawlDirs(args []int) {
	// TODO: find last modified
	// TODO: find todo, fixme, etc
	// goroutines!!!!!!!!!!!
	// one pulls in paths the other crawls?
	// go routines split into dirs and files, sending data to centeral channel?
	var wg sync.WaitGroup
	wg.Add(1)

}

// NOTE: this should probably be in crawlDirs
// func findLastModified(path string) {
// 	var lastModified string
// 	var lastModifiedInt int64 = 0
// 	files, err := os.ReadDir(path)
// 	if err != nil {
// 		fmt.Println("error with reading directory")
// 		os.Exit(1)
// 	}
// 	for _, f := range files {
// 		fi, err := os.Stat(path + f.Name())
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		currTime := fi.ModTime().Unix()
// 		if currTime > lastModifiedInt {
// 			lastModifiedInt = currTime
// 			lastModified = f.Name()
// 		}
// 	}
// }
