package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

const directory = "pictures"
const bucketName = "pictures-20220710"

func main() {
	infos, err := ioutil.ReadDir(directory)
	if err != nil {
		panic("Unable to read directory")
	}

	fmt.Printf("Getting existing files in bucket...\n")
	findCommandOutput, err := exec.Command("b2", "ls", bucketName).Output()
	if err != nil {
		panic(err)
	}

	existingFileNames := buildExistingFilesMap(findCommandOutput)
	fmt.Println(existingFileNames)

	length := len(infos)
	for i, info := range infos {
		fileName := info.Name()
		fmt.Printf("Handling %s (%d / %d)...\n", fileName, i+1, length)
		_, exists := existingFileNames[fileName]
		if exists {
			fmt.Println("File already uploaded. Skipping...")
			continue
		}

		fmt.Println("Uploading file...")
		filePath := path.Join(directory, fileName)
		_, err := exec.Command("b2", "upload-file", bucketName, filePath, fileName).Output()
		if err != nil {
			panic(err)
		}

		fmt.Println("File uploaded")
	}
}

func buildExistingFilesMap(output []byte) map[string]any {
	existingFileNames := strings.Split(string(output), "\n")
	existingMap := make(map[string]any)
	marker := struct{}{}
	for _, name := range existingFileNames {
		existingMap[name] = marker
	}

	return existingMap
}
