package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func printUsage() {
	fmt.Print(`Jot - jot stuff down without messing up your workspace!

usage: jot [file]             edit jot file in working directory
   or: jot <command> [<args>] perform misc operations

commands:
    ls          List jot files in the working directory
    rm          Remove jot files from the working directory
    clean       Remove all jot files from the working directory
    clean-all   Remove all jot files in the system
    help        Print Help (this message) and exit
    `)
}

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
		fmt.Println("Jot files in this folder")
		err := jo.ListDir(jo.GetProjDir())

		if err != nil {
			fmt.Fprintf(os.Stderr, "No jot files present for: %s", curdir)
		}
	case "rm":
		if 3 > len(jotArgs) {
			fmt.Fprintf(os.Stderr, "Insufficient argumens passed to 'rm'")
			return
		}

		jexists := jo.JotExists(jotArgs[2])
		if !jexists {
			fmt.Fprintf(os.Stderr, "No such jot: %s", jotArgs[2])
			return
		}

		fmt.Printf("Delete %s? [y/N] ", jotArgs[2])
		fmt.Scanln(&confrm)

		if p := strings.ToLower(confrm); p != "yes" && p != "y" {
			fmt.Println("Aborted")
			return
		}
		jo.RemoveFile(jotArgs[2])
	case "clean-all":
		fmt.Print("Delete all 'jot' files from this system? [y/N] ")
		fmt.Scanln(&confrm)

		if p := strings.ToLower(confrm); p != "yes" && p != "y" {
			fmt.Println("Aborted")
			return
		}

		os.RemoveAll(datadir)
		os.Mkdir(datadir, os.ModePerm)

		fmt.Println("Deleted all 'jot' files on this system!")
	case "clean":
		fmt.Print("Delete all 'jot' files in this project? [y/N] ")
		fmt.Scanln(&confrm)

		if p := strings.ToLower(confrm); p != "yes" && p != "y" {
			fmt.Println("Aborted")
			return
		}
		os.RemoveAll(jo.GetProjDir())

		fmt.Println("Deleted all 'jot' files in this project!")
	case "help":
		printUsage()
	default:
		absFilePath := filepath.Join(jo.GetProjDir(), jotArgs[1])
		jo.EditFile(absFilePath)
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	} else {
		procCmd(os.Args)
	}
}
