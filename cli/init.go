package cli

import (
	"fmt"
	"github.com/dailymuse/git-fit/util"
	"os"
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

func Init(args []string) {
	awsAccessKey := ""
	awsSecretKey := ""
	awsBucket := ""

	if len(args) > 0 && args[0] == "env" {
		// Configure via environment variables
		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		awsBucket = os.Getenv("AWS_S3_BUCKET")
		awsRegion = os.Getenv("AWS_S3_REGION")
	} else {
		// Configure via stdin
		awsAccessKey = getStdinString("Enter your AWS access key", util.GitConfig("git-fit.aws.access-key"))
		awsSecretKey = getStdinString("Enter your AWS secret key", util.GitConfig("git-fit.aws.secret-key"))
		awsBucket = getStdinString("Enter your AWS S3 bucket", util.GitConfig("git-fit.aws.bucket"))
		awsRegion := getStdinString("Enter your AWS S3 bucket region", util.GitConfig("git-fit.aws.region"))
	}

	util.SetGitConfig("git-fit.aws.access-key", awsAccessKey)
	util.SetGitConfig("git-fit.aws.secret-key", awsSecretKey)
	util.SetGitConfig("git-fit.aws.bucket", awsBucket)
	util.SetGitConfig("git-fit.aws.region", awsRegion)

	err := os.MkdirAll(".git/fit", os.ModePerm)

	if err != nil {
		util.Fatal("Could not create the asset staging directory (.git/fit): %s\n", err.Error())
	}
}
