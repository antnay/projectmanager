package internal

import (
	"fmt"
	"os"
	"sync"
)

type Project struct {
	name     string
	desc     string
	path     string
	language string
	active   bool
}

func NewProj(args []string) {
	// TODO: args as name, desc, path, ...
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
	newProjBytes := newProject.newEntrytoYamlBytes()

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
	// for i, project := range projects.ProjectList {
	// }
	// NOTE: count 5 or however many newlines there are and remove them?

}

func EditProj() {
	// for i, project := range projects.ProjectList {
	// }

}

// TODO: Add option to display todos and to only show a certain project specified with args
// TODO: Make output pretier and not a yaml dump
func DisplayProj() {
	// NOTE: when do i crawl
	file := retrieveData()
	l := len(file)
	buf := make([]byte, l)
	copy(buf, file)
	fmt.Println(string(buf))

}
