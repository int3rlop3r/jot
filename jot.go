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

    switch jotArgs[1] {
    case "clean-all":
        fmt.Println("Delete all files this system?")
        fmt.Scanln(&confrm)

        if p := strings.ToLower(confrm); p != "yes" && p != "y" {
            fmt.Println("Aborted")
            return
        }

        fmt.Println("Deleted all files on this system!")
    case "clean":
        fmt.Println("Delete all files in this project?")
        fmt.Scanln(&confrm)

        if p := strings.ToLower(confrm); p != "yes" && p != "y" {
            fmt.Println("Aborted")
            return
        }

        fmt.Println("Deleted all files in this project!")
    default:
        fmt.Println("Creating file:", jotArgs[1])
        curdir, err := os.Getwd()
        homedir := os.Getenv("HOME")
        datadir := filepath.Join(homedir, ".jot")

        if err != nil {
            log.Fatal(err)
        }

        jo := JotOps{jotArgs[1], curdir, datadir, ""}
        jo.Start()
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: jot [options]")
    } else {
        fmt.Println("Opening editor")
        procCmd(os.Args)
        fmt.Println("Done")
    }
}
