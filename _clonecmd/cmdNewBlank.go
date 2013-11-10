package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type newBlank struct {
	folderName string
}

func getSourcePath() string {
	return SanitizePath(os.ExpandEnv("$GOPATH/src/github.com/francoishill/beegowebapp"))
}

func checkAllPathsExist() (bool, *[]string) {
	sourcePath := getSourcePath()
	pathsToCheckExists := []string{
		sourcePath,
		path.Join(sourcePath, ".git"),
		path.Join(sourcePath, "conf"),
		path.Join(sourcePath, "db"),
		path.Join(sourcePath, "mailer"),
		path.Join(sourcePath, "models"),
		path.Join(sourcePath, "routers"),
		path.Join(sourcePath, "static"),
		path.Join(sourcePath, "static_source"),
		path.Join(sourcePath, "utils"),
		path.Join(sourcePath, "views"),
		path.Join(sourcePath, "bee.json"),
		path.Join(sourcePath, "main.go"),
		path.Join(sourcePath, "README.md"),
	}

	missingPaths := []string{}
	for _, path := range pathsToCheckExists {
		pathExist, _ := PathExists(path)

		if !pathExist {
			missingPaths = append(missingPaths, path)
		}
	}

	if len(missingPaths) == 0 {
		return true, nil
	} else {
		return false, &missingPaths
	}
}

func pullGit() (bool, error) {
	sourcePath := getSourcePath()
	err := os.Chdir(sourcePath)
	if err != nil {
		return false, err
	}

	out, err := exec.Command("git", "pull", "-v", "--progress", "origin", "master").Output()
	if err != nil {
		fmt.Println(string(out))
		return false, err
	}
	if !strings.Contains(strings.ToLower(string(out)), strings.ToLower("Already up-to-date")) &&
		!strings.Contains(strings.ToLower(string(out)), strings.ToLower("Fast-forward")) {
		return false, errors.New("Error, unknown git pull result: " + string(out))
	}
	return true, nil
}

func cloneAllFiles(destinationDir string, appFoldername string) (bool, []error) {
	sourcePath := getSourcePath()

	exclusionList := &[]string{
		path.Join(sourcePath, ".git"),
		path.Join(sourcePath, ".gitignore"),
		path.Join(sourcePath, "_clonecmd"),
		path.Join(sourcePath, "beegowebapp.exe"),
		path.Join(sourcePath, "How to make it standalone.txt"),
		path.Join(sourcePath, "install_cloner.sh"),
		path.Join(sourcePath, "main.exe"),
		path.Join(sourcePath, "README.md"),
	}

	destinationPath := path.Join(destinationDir, appFoldername)

	os.Mkdir(destinationPath, 0755)
	errs := CopyDir(sourcePath, destinationPath, exclusionList)
	if errs != nil {
		return false, errs
	}
	return true, nil
}

func (d *newBlank) Parse(args []string) {
	//var name string
	const cUndefinedFoldername = "noname_newbeego"

	flagSet := flag.NewFlagSet("newBlank command: blank", flag.ExitOnError)
	flagSet.StringVar(&(*d).folderName, "foldername", cUndefinedFoldername, "Requires a foldername for the blank app.")

	err := flagSet.Parse(args)
	if err != nil {
		panic(err)
	} else if (*d).folderName == cUndefinedFoldername {
		fmt.Println(fmt.Errorf("Please specify the foldername:"))
		flagSet.PrintDefaults()
		os.Exit(3)
	}
}

func (d *newBlank) Run(s *Settings) error {
	curpath_destinationDir, _ := os.Getwd()
	curpath_destinationDir = SanitizePath(curpath_destinationDir)

	allPathsExist, missingPaths := checkAllPathsExist()
	if !allPathsExist {
		for _, mis := range *missingPaths {
			fmt.Printf("Missing path: %s", mis)
		}

		fmt.Println("Missing paths, see log above")
		os.Exit(1)
		return errors.New("Missing paths, see log above") //Just for safety
	}

	/* DO NOT do git pull, maybe user does not desire to do this?
	pullGitSuccess, err := pullGit()
	if !pullGitSuccess {
		fmt.Println("Error occurred with git pull: ", err.Error())
		os.Exit(2)
		return errors.New("Error occurred with git pull: " + err.Error()) //Just for safety
	}*/

	copySuccess, errsList := cloneAllFiles(curpath_destinationDir, d.folderName)
	if !copySuccess {
		for _, er := range errsList {
			fmt.Println("Error occurred with copying files: ", er.Error())
		}
		os.Exit(2)
		return errors.New("Error occurred with copying files: " + errsList[0].Error()) //Just for safety return last error
	}

	return nil
}
