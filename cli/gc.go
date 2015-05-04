package cli

import (
	"fmt"
	"github.com/dailymuse/git-fit/config"
	"github.com/dailymuse/git-fit/transport"
	"github.com/dailymuse/git-fit/util"
	"io/ioutil"
	"os"
)

func Gc(schema *config.Config, trans transport.Transport, args []string) {
	savedFiles := make(map[string]bool, len(schema.Files)*2)

	for _, hash := range schema.Files {
		savedFiles[hash] = true
	}

	allFiles, err := ioutil.ReadDir(".git/fit")

	if err != nil {
		util.Fatal("Could not read .git/fit: %s\n", err.Error())
	}

	for _, file := range allFiles {
		_, ok := savedFiles[file.Name()]

		if !ok {
			path := fmt.Sprintf(".git/fit/%s", file.Name())
			err = os.Remove(path)

			if err != nil {
				util.Error("Could not delete cached file %s: %s\n", path, err.Error())
			}
		}
	}
}
