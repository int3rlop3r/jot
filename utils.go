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
    fileName, curDir, dataDir, projDir string
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

func (jo JotOps) makeProjDir() string {
    pathHash := jo.makeSha1(jo.curDir)
    projDir := filepath.Join(jo.dataDir, pathHash)

    _, err := jo.makeDir(jo.projDir)

    if err != nil {
        fmt.Fprintf(os.Stderr, "Coud not create project dir: %s", err)
    }

    return projDir
}

func (jo JotOps) Start() {

    // create data dir if it doesn't exist
    jo.makeDataDir()

    // make project dir if it doesn't exist!
    projDir := jo.makeProjDir()
    filePath := filepath.Join(projDir, jo.fileName)

    // fire up the default editor and start editing the file !!
    jo.editFile()
}

func (jo JotOps) editFile(filePath string) {
    cmd := exec.Command("editor", filePath)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout

    err := cmd.Run()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Coud not open file for editing: %s", err)
    }
}
