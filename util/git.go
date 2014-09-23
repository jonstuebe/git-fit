package util

import (
    "os/exec"
    "io/ioutil"
    "strings"
)

func Stdout(name string, args ...string) string {
    cmd := exec.Command(name, args...)
    stdout, err := cmd.StdoutPipe()
    cmd.Start()

    bytes, err := ioutil.ReadAll(stdout)

    if err != nil {
        panic(err)
    }

    return strings.TrimSpace(string(bytes))
}

func GitDir() string {
    return Stdout("git", "rev-parse", "--git-dir")
}

func LatestCommit() string {
    return Stdout("git", "rev-parse", "HEAD")
}

func GitConfig(name string) string {
    return Stdout("git", "config", "--get", name)
}

func SetGitConfig(name string, value string) {
    Stdout("git", "config", "--add", name, value)
}
