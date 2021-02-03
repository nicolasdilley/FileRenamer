package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sqweek/dialog"
)

var (
	reader = bufio.NewReader(os.Stdin)
)

// FileRenamer takes a folder and two strings (old, new) as input and replaces the occurences of 'old' in the name of all the files contained in the folder by 'new'
// This is not done recursively. Meaning that the files contained in a folder inside the parent folder will not be renamed.
func main() {

	parentFolder := getFolder()
	parent_info, _ := parentFolder.Stat()
	old, n := getOldAndNew()
	filepath.Walk(parentFolder.Name(), func(path string, info os.FileInfo, err error) error {

		if err == nil && !os.SameFile(parent_info, info) {
			return replaceFileName(path, info, err, old, n)
		} else {
			fmt.Println("Could not read file: ", path)
		}

		return nil
	})

}

func getFolder() (parent_folder *os.File) {
	var err error
	if len(os.Args) > 1 {
		// the name of the folder has been given

		parent_folder, err = os.Open(os.Args[1])

		if err != nil {
			fmt.Println("Could not find folder : ", os.Args[1], "\n", err)
			parent_folder = nil
		}
	}

	if parent_folder == nil {

		var dir *dialog.DirectoryBuilder

		dir.StartDir = "."
		var dirName string

		// ask the user until a valid folder is given
		for parent_folder == nil {

			dirName, err = dir.Browse()

			if err != nil {
				fmt.Println("Could not find folder : ", os.Args[1], "\n", err)
				parent_folder = nil
			}

			parent_folder, err = os.Open(dirName)

			if err != nil {
				fmt.Println("Could not find folder : ", os.Args[1], "\n", err)
				parent_folder = nil
			}
		}
	}

	return parent_folder
}

func getOldAndNew() (old string, n string) {

	if len(os.Args) > 2 {
		old = os.Args[2]
	} else {
		fmt.Print("Enter the old pattern: ")
		old, _ = reader.ReadString('\n')
	}
	if len(os.Args) > 3 {
		n = os.Args[3]
	} else {
		fmt.Print("Enter the new pattern: ")
		n, _ = reader.ReadString('\n')
	}

	return old, n
}

func replaceFileName(path string, info os.FileInfo, err error, old string, n string) error {

	fmt.Println("renaming ", path, " with ", old, " to ", n)
	os.Rename(path, strings.Replace(path, old, n, -1))

	if info.IsDir() {
		return filepath.SkipDir
	}
	return err

}
