package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "path/filepath"
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
        fmt.Println("Jot files in this folder")
        err := listDir(jo.GetProjDir())

        if err != nil {
            fmt.Fprintf(os.Stderr, "Error listing dir: %s", err)
        }
    case "clean-all":
        fmt.Println("Delete all 'jot' files from this system?")
        fmt.Scanln(&confrm)

        if p := strings.ToLower(confrm); p != "yes" && p != "y" {
            fmt.Println("Aborted")
            return
        }

        os.RemoveAll(datadir)
        os.Mkdir(datadir, os.ModePerm)

        fmt.Println("Deleted all 'jot' files on this system!")
    case "clean":
        fmt.Println("Delete all 'jot' files in this project?")
        fmt.Scanln(&confrm)

        if p := strings.ToLower(confrm); p != "yes" && p != "y" {
            fmt.Println("Aborted")
            return
        }
        os.RemoveAll(jo.GetProjDir())

        fmt.Println("Deleted all 'jot' files in this project!")
    default:
        absFilePath := filepath.Join(jo.GetProjDir(), jotArgs[1])
        editFile(absFilePath)
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: jot [options]")
    } else {
        procCmd(os.Args)
    }
}
