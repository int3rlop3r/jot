# Jot

Jot is a tool that helps you take notes in your terminal without leaving files lying around in your project's directory. It does this by grouping all your notes (a.k.a. jots) under a "tracked" workspace. Your "jots" are actually stored in an sqlite db and can be viewed inside directories that they were created in using `jot`'s commands. All jots are opened in the default editor (you can change it using `update-alternatives`)

### Usage

    Jot - jot stuff down without messing up your workspace!
    
    usage: jot [file]             edit jot file in working directory
       or: jot <command> [<args>] perform misc operations
    
    commands:
        -t                 Track current directory
        -u                 Untrack current directory. Note: this will delete all jots in the dir
        -l                 List jot files in the working directory
        -o                 Print file contents on the standard output
        -d                 Delete a jot from the working directory
        -m                 Import a file. The resulting file's name will be the jot's name
        -clean-all         Remove all jot files in the system (dangerous)
        -list-tracked      List all tracked dirs
        -help              Print Help (this message) and exit

### Install

You can download a binary (for linux) from [here](https://github.com/int3rlop3r/jot/releases) (don't forget to mark it as 'executable').
