package paths

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

// LocalDir where local files will be stored
var LocalDir string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	LocalDir = path.Join(home, ".owl")

	err = os.Mkdir(LocalDir, 0644)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}
}
