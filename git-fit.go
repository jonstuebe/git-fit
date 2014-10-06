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
    fmt.Printf(`usage: git-fit push|pull|rm|gc

    init
        Initializes git-fit for a repo by adding the necessary configs.

    push [file 1] ... [file n]
        Pushes the specified files. If no arguments are given, all of the
        files in git-fit.json are pushed.

    pull [file 1] ... [file n]
        Pulls the specified files. If no arguments are given, all of the files
        in git-fit.json are pulled.

    rm [file 1] ... [file n]
        Removes the specified files from git-fit.json and consequently from
        git-fit source control.

    gc
        Removes cached assets that aren't being used by the current commit.
        This will save space, but may slow down future pulls while the cache
        is warming back up.

`)

    os.Exit(code)
}

func affixToGitRepo() {
    gitDirectory, err := filepath.Abs(util.GitDir())

    if err != nil {
        util.Fatal("Could not determine the repo root: %s\n", err.Error())
    } else if !util.IsDirectory(gitDirectory) {
        util.Fatal("Not in a git repository\n")
    }

    if err = os.Chdir(filepath.Join(gitDirectory, "..")); err != nil {
        util.Fatal("Could not change the working directory to the repo root (%s): %s\n", gitDirectory, err.Error())
    }
}

func getTransport() transport.Transport {
    auth, err := aws.GetAuth(util.GitConfig("git-fit.aws.access-key"), util.GitConfig("git-fit.aws.secret-key"))

    if err != nil {
        util.Fatal("AWS authentication failed: %s\n", err.Error())
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
        util.Fatal("Could not load config file %s: %s\n", config.FILE_PATH, err.Error())
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
        case "gc":
            cli.Gc(schema, trans, os.Args[2:])
        default:
            util.Error("Unknown command")
            help(-1)
        }
    }

    if err = config.SaveConfig(schema); err != nil {
        util.Fatal("Could not save config file %s: %s\n", config.FILE_PATH, err.Error())
    }
}
