package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "crypto/sha1"
    "encoding/hex"
)

type JotOps struct {
    curDir, dataDir string
}

func (jo JotOps) exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func (jo JotOps) makeSha1(dirpath string)  string {
    h := sha1.New()
    h.Write([]byte(dirpath))
    return hex.EncodeToString(h.Sum(nil))
}

func (jo JotOps) makeDir(dirPath string) (bool, error) {
    dexists, err := jo.exists(dirPath)

    if err != nil { return false, err }

    if !dexists {
        os.Mkdir(dirPath, os.ModePerm)
    }

    return true, err
}

func (jo JotOps) makeDataDir() {
    _, err := jo.makeDir(jo.dataDir)

    if err != nil {
        fmt.Fprintf(os.Stderr, "Coud not create data dir: %s", err)
    }
}

func (jo JotOps) makeProjDir() {
    _, err := jo.makeDir(jo.GetProjDir())

    if err != nil {
        fmt.Fprintf(os.Stderr, "Coud not create project dir: %s", err)
    }
}

func (jo JotOps) GetDataDir() string {
    return jo.dataDir
}

func (jo JotOps) GetProjDir() string {
    pathHash := jo.makeSha1(jo.curDir)
    return filepath.Join(jo.dataDir, pathHash)
}

func (jo JotOps) Init() {

    // create data dir if it doesn't exist
    jo.makeDataDir()

    // make project dir if it doesn't exist!
    jo.makeProjDir()
}

func listDir(dirPath string) error {
    d, err := os.Open(dirPath)

    if err != nil {
        return err
    }

    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        fstats, err := os.Stat(filepath.Join(dirPath, name))

        if err != nil {
            return err
        }

        fmt.Println(fstats.ModTime(), fstats.Name())
    }

    return err
}

func editFile(filePath string) {
    cmd := exec.Command("editor", filePath)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout

    err := cmd.Run()
    if err != nil { // add this to the case stmt???
        fmt.Fprintf(os.Stderr, "Coud not open file for editing: %s", err)
    }
}
