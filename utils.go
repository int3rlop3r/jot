package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type JotOps struct {
	curDir, dataDir string
}

func (jo JotOps) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func (jo JotOps) noSuchJot(path, jot string) error {
	jpath := filepath.Join(path, jot)
	return errors.New("No such jot path:" + jpath)
}

func (jo JotOps) makeSha1(dirpath string) string {
	h := sha1.New()
	h.Write([]byte(dirpath))
	return hex.EncodeToString(h.Sum(nil))
}

func (jo JotOps) makeDir(dirPath string) (bool, error) {
	dexists, err := jo.exists(dirPath)

	if err != nil {
		return false, err
	}

	if !dexists {
		err = os.Mkdir(dirPath, os.ModePerm)
	}

	return true, err
}

func (jo JotOps) makeDataDir() error {
	_, err := jo.makeDir(jo.dataDir)
	return err
}

func (jo JotOps) makeProjDir() error {
	_, err := jo.makeDir(jo.GetProjDir())
	return err
}

func (jo JotOps) JotExists(jotName string) bool {
	jotPath := filepath.Join(jo.GetProjDir(), jotName)
	jexists, err := jo.exists(jotPath)

	if err != nil {
		return false
	}

	return jexists
}

func (jo JotOps) JotDirExists(dirPath string) bool {
	jotDir := jo.GetJotDir(dirPath)
	jexists, err := jo.exists(jotDir)

	if err != nil {
		return false
	}

	return jexists
}

func (jo JotOps) GetJotDir(dirPath string) string {
	pathHash := jo.makeSha1(dirPath)
	return filepath.Join(jo.dataDir, pathHash)
}

func (jo JotOps) CopyJot(srcJot, dstPath string) error {
	if !jo.JotExists(srcJot) {
		return jo.noSuchJot(jo.curDir, srcJot)
	}
	srcj, err := os.Open(filepath.Join(jo.GetProjDir(), srcJot))
	if err != nil {
		return err
	}

	defer srcj.Close()

	f, err := os.Open(dstPath)
	defer f.Close()
	if err != nil {
		f, err = os.Open(filepath.Dir(dstPath))
		if err != nil {
			return err
		}
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	var newJot string
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// dst is a dir path
		jpath := jo.GetJotDir(dstPath)
		jo.makeDir(jpath)
		newJot = filepath.Join(jpath, srcJot)
	case mode.IsRegular():
		// dst is a file path
		jpath := jo.GetJotDir(filepath.Dir(dstPath))
		jo.makeDir(jpath)
		newJot = filepath.Join(jpath, filepath.Base(dstPath))
	}

	df, err := os.Create(newJot)
	defer df.Close()
	if err != nil {
		return err
	}

	if _, err := io.Copy(df, srcj); err != nil {
		return err
	}

	if err = df.Sync(); err != nil {
		return err
	}

	return err
}

func (jo JotOps) MoveJot(srcJot, dstPath string) error {
	err := jo.CopyJot(srcJot, dstPath)
	return err
}

func (jo JotOps) GetDataDir() string {
	return jo.dataDir
}

func (jo JotOps) GetProjDir() string {
	return jo.GetJotDir(jo.curDir)
}

func (jo JotOps) Init() {
	// create data dir if it doesn't exist
	jo.makeDataDir()

}

func (jo JotOps) ListDir(dirPath string, cb func(os.FileInfo)) error {
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

		cb(fstats)
	}

	return err
}

func (jo JotOps) RemoveFile(fileName string) {
	projDir := jo.GetProjDir()
	os.Remove(filepath.Join(projDir, fileName))

	// delete folder if no more files present
	d, err := os.Open(projDir)
	defer d.Close()

	if _, err = d.Readdirnames(2); err != nil {
		os.Remove(projDir)
	}
}

func (jo JotOps) EditFile(filePath string) error {
	// make project dir if it doesn't exist!
	err := jo.makeProjDir()

	if err != nil {
		return err
	}

	cmd := exec.Command("editor", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()

	return err
}
