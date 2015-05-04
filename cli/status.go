package cli

import (
	"github.com/dailymuse/git-fit/config"
	"github.com/dailymuse/git-fit/transport"
	"github.com/dailymuse/git-fit/util"
)

func Status(schema *config.Config, trans transport.Transport, args []string) {
	paths := args

	if len(paths) == 0 {
		paths = make([]string, 0)

		for path := range schema.Files {
			paths = append(paths, path)
		}
	}

	for _, path := range paths {
		expectedHash, ok := schema.Files[path]

		if !ok {
			util.Message("%s: Does not exist in git-fit.json\n", path)
		} else if !util.FileExists(path) {
			util.Message("%s: Does not exist\n", path)
		} else {
			actualHash, err := util.FileHash(path)

			if err != nil {
				util.Error("%s: Could not get file hash: %s\n", path, err.Error())
			} else if actualHash != expectedHash {
				util.Message("%s: Out of sync; push the file, or delete it and pull\n", path)
			} else {
				util.Message("%s: In sync\n", path)
			}
		}
	}
}
