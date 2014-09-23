package cli

import (
    "github.com/dailymuse/git-fit/transport"
    "github.com/dailymuse/git-fit/config"
)

func Remove(schema config.Config, trans transport.Transport, args []string) {
    for _, path := range args {
        delete(schema, path)
    }
}
