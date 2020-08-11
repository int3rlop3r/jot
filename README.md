# Jot
Jot stuff down without messing up your workspace!

[![asciicast](https://asciinema.org/a/bqlsbmokx5zdc0ti4y901krde.png)](https://asciinema.org/a/bqlsbmokx5zdc0ti4y901krde)

### Usage

    Jot - jot stuff down without messing up your workspace!
    
    usage: jot [file]             edit jot file in working directory
       or: jot <command> [<args>] perform misc operations
    
    commands:
        -t              Track current directory
        -u              Untrack current directory. Note: this will delete all jots in the dir
        -l              List jot files in the working directory
        -o              Print file contents on the standard output
        -d              Delete a jot from the working directory
        -clean-all      Remove all jot files in the system (dangerous)
        -list-jots      List all jot dirs in the system
        -help           Print Help (this message) and exit

### Install

    $ go get github.com/int3rlop3r/jot

