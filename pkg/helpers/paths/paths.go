package paths

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

// Application standard paths
var (
	Home    string
	Session string
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Mkdir(home+"/.owl", 0644)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	Home = path.Join(home, ".owl")
	Session = path.Join(Home, "session.yaml")
}
