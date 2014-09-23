package main

import (
    "os"
    "fmt"
    "path/filepath"
    "github.com/mitchellh/goamz/s3"
    "github.com/mitchellh/goamz/aws"
    "github.com/dailymuse/git-fit/util"
    "github.com/dailymuse/git-fit/cli"
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
)

func help(code int) {
    fmt.Printf(`usage: git-fit push|pull|rm

    init
        Initializes git-fit for a repo by adding the necessary configs.

    push [file 1] ... [file n]
        Pushes the specified files. If no arguments are given, all of the
        files in .git-fit.json are pushed.

    pull [file 1] ... [file n]
        Pulls the specified files. If no arguments are given, all of the files
        in .git-fit.json are pulled.

    rm [file 1] ... [file n]
        Removes the specified files from .git-fit.json and consequently from
        git-fit source control.

`)

    os.Exit(code)
}

func affixToGitRepo() {
    gitDirectory, err := filepath.Abs(util.GitDir())

    if err != nil {
        panic(err)
    } else if !util.IsDirectory(gitDirectory) {
        panic("Not in a git repository")
    }

    if err = os.Chdir(filepath.Join(gitDirectory, "..")); err != nil {
        panic(err)
    }
}

func getTransport() transport.Transport {
    auth, err := aws.GetAuth(util.GitConfig("git-fit.aws.access-key"), util.GitConfig("git-fit.aws.secret-key"))

    if err != nil {
        panic(err)
    }

    bucket := s3.New(auth, aws.USEast).Bucket(util.GitConfig("git-fit.aws.bucket"))
    return transport.NewS3Transport(bucket)
}

func main() {
    if len(os.Args) < 2 {
        help(0)
    }

    affixToGitRepo()

    schema, err := config.LoadConfig()

    if err != nil {
        panic(err)
    }

    if os.Args[1] == "init" {
        cli.Init()
    } else {
        trans := getTransport()

        switch os.Args[1] {
        case "help":
            help(0)
        case "push":
            cli.Push(schema, trans, os.Args[2:])
        case "rm":
            cli.Remove(schema, trans, os.Args[2:])
        case "pull":
            cli.Pull(schema, trans, os.Args[2:])
        default:
            fmt.Printf("Unknown command")
            help(1)
        }
    }

    if err = config.SaveConfig(schema); err != nil {
        panic(err)
    }
}
