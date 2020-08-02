package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	i         = flag.Bool("i", false, "Initialize a new project directory")
	l         = flag.Bool("l", false, "List jot files in the working directory")
	o         = flag.String("o", "", "Print file contents on the standard output")
	d         = flag.Bool("d", false, "Delete a jot from the working directory")
	D         = flag.Bool("D", false, "Delete all jots in the working directory")
	clean_all = flag.Bool("clean-all", false, "Remove all jot files in the system (dangerous)")
	list_all  = flag.Bool("list-all", false, "List all jot files in the system")
	help      = flag.Bool("help", false, "Print Help (this message) and exit")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Jot - jot stuff down without messing up your workspace!

usage: jot [file]             edit jot file in working directory
   or: jot <command> [<args>] perform misc operations

commands:
	-i         %v
	-l         %v
	-o         %v
	-d         %v
	-D         %v

	-clean-all %v
	-help      %v
`, flag.Lookup("i").Usage,
			flag.Lookup("l").Usage,
			flag.Lookup("o").Usage,
			flag.Lookup("d").Usage,
			flag.Lookup("D").Usage,
			flag.Lookup("clean-all").Usage,
			flag.Lookup("help").Usage,
		)
	}

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
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB error: %s", err)
		return
	}
	switch {
	case *i:
		if err := db.initialize(curDir); err != nil {
			fmt.Fprint(os.Stderr, "Directory already tracked\n")
		} else {
			fmt.Println("Directory initialized")
		}
	case *l:
		res, err := db.listByPath(curDir)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB err:", err)
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
		jotName := strings.Join(os.Args[2:], " ")
		var content string
		err := db.get(curDir, jotName).Scan(&content)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB err:", err)
			return
		}
		fmt.Println(content)
	case *d:
		fmt.Println("deleting")
	case *D:
		fmt.Println("deleting all current")
	case *list_all:
		fmt.Println("deleting all jot on computer")
	case *clean_all:
		fmt.Println("deleting all jot on computer")
	default:
		fmt.Println("creating new jot:", os.Args[1])
	}
}
