package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

var (
	// change i to t. -t to track directory and -u to untrack a directory
	t        = flag.Bool("t", false, "Track current directory")
	u        = flag.Bool("u", false, "Untrack current directory. Note: this will delete all jots in the dir")
	l        = flag.Bool("l", false, "List jot files in the working directory")
	o        = flag.String("o", "", "Print file contents on the standard output")
	d        = flag.String("d", "", "Delete a jot from the working directory")
	cleanAll = flag.Bool("clean-all", false, "Remove all jot files in the system (dangerous)")
	listAll  = flag.Bool("list-all", false, "List all jot files in the system")
	help     = flag.Bool("help", false, "Print Help (this message) and exit")
)

func showUsage() {
	var args = []string{"t", "u", "l", "o", "d", "clean-all", "help"}
	fmt.Fprintf(os.Stderr, `Jot - jot stuff down without messing up your workspace!

usage: jot [file]             edit jot file in working directory
   or: jot <command> [<args>] perform misc operations

commands:
`)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()
	var i int
	for i = 0; i < len(args); i++ {
		opt := args[i]
		fmt.Fprintf(w, "    -%s\t\t%v\n", opt, flag.Lookup(opt).Usage)
	}
}

func main() {
	flag.Usage = showUsage
	flag.Parse()
	if len(os.Args) < 2 || *help {
		flag.Usage()
		return
	}

	processArgs()
}

func processArgs() {
	curDir, _ := os.Getwd()
	jotDir := getDBPath()
	db, err := setupDB(jotDir)
	defer db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB error path: %s, error: %s", jotDir, err)
		return
	}
	switch {
	case *t:
		if err := db.initialize(curDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("Directory initialized")
	case *u:
		if !confirm(fmt.Sprintf("Untrack: %s?", curDir)) {
			fmt.Fprintf(os.Stderr, "didn't delete: %s\n", curDir)
			return
		}
		if err := db.uninitialize(curDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println("Removed current dir and all jots")
	case *l:
		res, err := db.listByPath(curDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DB err: %s\n", err)
			return
		}
		var i int
		for i = 0; res.Next(); i++ {
			var t time.Time
			var n string
			res.Scan(&n, &t)
			fmt.Printf("%s\t%s\n", t.Format("01-02-2006 15:04:05"), n)
		}
		if i == 0 {
			fmt.Println("No jots in this dir")
		}
	case *o != "":
		jotName := *o
		contents, err := getJot(db, curDir, jotName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Fprint(os.Stdin, *contents)
	case *d != "":
		if !confirm(fmt.Sprintf("Delete: %s", *d)) {
			fmt.Fprintf(os.Stderr, "didn't delete: %s\n", *d)
			return
		}
		jotName := *d
		id, err := db.getJotDir(curDir)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		err = db.delete(id, jotName)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB err:", err)
			return
		}
		fmt.Println("deleted:", *d)
	case *listAll:
		fmt.Println("list all jots on the system")
	case *cleanAll:
		if !confirm("DELETE ALL JOT FROM THE SYSTEM?") {
			fmt.Fprintf(os.Stderr, "not deleted\n")
			return
		}
		os.Remove(filepath.Join(jotDir, "jot.db"))
		fmt.Println("deleted all jots")
	default:
		jotFile := os.Args[1]
		id, err := db.getJotDir(curDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		tmpFile, err := ioutil.TempFile("", "jot")
		if err != nil {
			fmt.Fprint(os.Stderr, "couldn't make tmp file:", tmpFile)
			return
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		// open the jot in an editor
		cmd := exec.Command("editor", tmpFile.Name())
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("couldn't create jot:", err)
			return
		}

		jotContents, err := ioutil.ReadFile(tmpFile.Name())
		if err != nil {
			fmt.Fprint(os.Stderr, "error reading tmp file:", err)
			return
		}

		_, err = db.createJot(id, jotFile, string(jotContents))
		if err != nil {
			fmt.Fprint(os.Stderr, "error creating jot:", err)
			return
		}
		fmt.Printf("new jot '%s' created\n", jotFile)
	}
}

func confirm(prompt string) bool {
	fmt.Fprintf(os.Stderr, "%s [N/y] ", prompt)
	var userInput string
	fmt.Scanln(&userInput)
	usrInput := strings.TrimSpace(strings.ToLower(userInput))
	return usrInput == "y" || usrInput == "yes"
}

func getJot(db *DB, curDir, jotName string) (*string, error) {
	id, err := db.getJotDir(curDir)
	if err != nil {
		return nil, err
	}
	contents, err := db.get(id, jotName)
	if err != nil {
		return nil, fmt.Errorf("DB err:", err)
	}
	return &contents, nil
}
