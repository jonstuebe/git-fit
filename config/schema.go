package config

import (
    "io/ioutil"
    "os"
    "errors"
    "fmt"
    "github.com/dailymuse/git-fit/util"
    "encoding/json"
    "strings"
    "regexp"
)

type Config map[string]ConfigEntry

type ConfigEntry struct {
    Commit string `json:"commit"`
}

const FILE_PATH = ".git-fit.json"

var NEWLINE_CHECKING_PATTERN = regexp.MustCompile("\n$")

func SaveConfig(config Config) error {
    bytes, err := json.MarshalIndent(config, "", "    ")

    if err != nil {
        return err
    }

    if err = ioutil.WriteFile(FILE_PATH, bytes, os.ModePerm); err != nil {
        return err
    }

    return ensureIgnoreEntries(config)
}

func LoadConfig() (Config, error) {
    if util.FileExists(FILE_PATH) {
        file, err := ioutil.ReadFile(FILE_PATH)

        if err != nil {
            return nil, err
        }

        config := make(Config)
        json.Unmarshal(file, &config)

        if err = validateConfig(config); err != nil {
            return nil, err
        } else if err = ensureIgnoreEntries(config); err != nil {
            return nil, err
        }

        return config, nil
    } else {
        config := make(Config)

        if err := SaveConfig(config); err != nil {
            return nil, err
        }

        return config, nil
    }
}

func validateConfig(config Config) error {
    for path, value := range config {
        if len(value.Commit) != 40 {
            return errors.New(fmt.Sprintf("Invalid SHA hash for config entry %s: %s", path, value.Commit))
        }
    }

    return nil
}

func ensureIgnoreEntries(config Config) error {
    candidateEntries := make(map[string]bool)
    endsWithNewLine := true

    for path := range config {
        candidateEntries[fmt.Sprintf("/%s", path)] = true
    }

    if util.FileExists(".gitignore") {
        contents, err := ioutil.ReadFile(".gitignore")

        if err != nil {
            return err
        }

        existingEntries := strings.Split(string(contents), "\n")

        for _, entry := range existingEntries {
            delete(candidateEntries, entry)
        }

        if len(contents) > 0 && !NEWLINE_CHECKING_PATTERN.Match(contents) {
            endsWithNewLine = false
        }
    }

    newEntries := make([]string, 0)

    for entry := range candidateEntries {
        newEntries = append(newEntries, entry)
    }

    file, err := os.OpenFile(".gitignore", os.O_WRONLY | os.O_APPEND | os.O_CREATE, os.ModePerm)

    if err != nil {
        return err
    }

    defer file.Close()

    if len(newEntries) > 0 {
        if !endsWithNewLine {
            if _, err := file.WriteString("\n"); err != nil {
                return err
            }
        }

        _, err = file.WriteString(strings.Join(newEntries, "\n") + "\n"); if err != nil {
            return err
        }
    }

    return nil
}
