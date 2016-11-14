package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func procCmd(jotArgs []string) {
	var confrm string

	curdir, err := os.Getwd()
	homedir := os.Getenv("HOME")
	datadir := filepath.Join(homedir, ".jot")

	if err != nil {
		log.Fatal(err)
	}

	jo := &JotOps{curdir, datadir}
	jo.Init()

	switch jotArgs[1] {
	case "ls":
		fmt.Print("Jot files in this folder")
		err := jo.ListDir(jo.GetProjDir())

		if err != nil {
			fmt.Fprintf(os.Stderr, "No jot files present for: %s", curdir)
		}
	case "clean-all":
		fmt.Print("Delete all 'jot' files from this system? [Y/n] ")
		fmt.Scanln(&confrm)

		if p := strings.ToLower(confrm); p != "yes" && p != "y" {
			fmt.Println("Aborted")
			return
		}

		os.RemoveAll(datadir)
		os.Mkdir(datadir, os.ModePerm)

		fmt.Println("Deleted all 'jot' files on this system!")
	case "clean":
		fmt.Print("Delete all 'jot' files in this project? [Y/n] ")
		fmt.Scanln(&confrm)

		if p := strings.ToLower(confrm); p != "yes" && p != "y" {
			fmt.Println("Aborted")
			return
		}
		os.RemoveAll(jo.GetProjDir())

		fmt.Println("Deleted all 'jot' files in this project!")
	default:
		absFilePath := filepath.Join(jo.GetProjDir(), jotArgs[1])
		jo.EditFile(absFilePath)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jot [options]")
	} else {
		procCmd(os.Args)
	}
}
