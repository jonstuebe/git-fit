package cli

import (
	"github.com/dailymuse/git-fit/config"
	"github.com/dailymuse/git-fit/transport"
)

func Remove(schema *config.Config, trans transport.Transport, args []string) {
	for _, path := range args {
		delete(schema.Files, path)
	}
}
