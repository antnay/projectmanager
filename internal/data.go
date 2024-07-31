package internal

import (
	"fmt"
	"os"
	"sync"
)

// TODO: no point in storing lastModified in yaml
type Project struct {
	name     string
	desc     string
	path     string
	language string
	active   bool
}

type Projects struct {
	ProjectList []Project `yaml:"projects"`
}

func NewProj(args []string) {
	// TODO: Add args as struct things
	var waitGroup sync.WaitGroup
	fileBytesChan := make(chan []byte, 1)
	projectChan := make(chan Project, 1)
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		fileBytesChan <- retrieveData()
	}()
	go func() {
		defer waitGroup.Done()
		projectChan <- makeProject()
	}()
	waitGroup.Wait()
	fileBytes := <-fileBytesChan
	newProject := <-projectChan

	newProjBytes := newProject.toYamlBytes()

	fileBytes = append(fileBytes, newProjBytes...)

	err := os.WriteFile(dataPathTemp, fileBytes, 0644)
	if err != nil {
		fmt.Println("Could not create new yaml file")
		os.Exit(1)
	}
	_ = os.Remove(dataPath)
	os.Rename(dataPathTemp, dataPath)

}

func RemoveProj() {

}

func EditProj() {

}

// TODO: Add option to display todos and to only show a certain project specified with args
// TODO: Make output pretier and not a yaml dump
func DisplayProj() {
	file := retrieveData()
	l := len(file)
	buf := make([]byte, l)
	// offset := 0
	copy(buf, file)
	// offset += len(itemType)
	// copy(buf[offset:], ":")
	fmt.Println(string(buf))

}

//  // edit
// for i, project := range projects.ProjectList {
// }
