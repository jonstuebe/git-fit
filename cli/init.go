package cli

import (
    "os"
    "fmt"
    "github.com/dailymuse/git-fit/util"
)

func getStdinString(prompt string, def string) string {
    var in string
    hasDefault := def != ""

    if hasDefault {
        fmt.Printf("%s (default %s): ", prompt, def)
    } else {
        fmt.Printf("%s: ", prompt)
    }
    
    fmt.Scanf("%s", &in)

    if in == "" && hasDefault {
        in = def
    }

    return in
}

func Init() {
    awsAccessKey := getStdinString("Enter your AWS access key", util.GitConfig("git-fit.aws.access-key"))
    awsSecretKey := getStdinString("Enter your AWS secret key", util.GitConfig("git-fit.aws.secret-key"))
    awsBucket := getStdinString("Enter your AWS S3 bucket", util.GitConfig("git-fit.aws.bucket"))

    util.SetGitConfig("git-fit.aws.access-key", awsAccessKey)
    util.SetGitConfig("git-fit.aws.secret-key", awsSecretKey)
    util.SetGitConfig("git-fit.aws.bucket", awsBucket)

    err := os.MkdirAll(".git/fit", os.ModePerm)

    if err != nil {
        util.Fatal("Could not create the asset staging directory (.git/fit): %s\n", err.Error())
    }
}
