package status

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/config"
	"io/ioutil"
	"os"
)

type State uint

const (
	ENABLED State = iota
	DISABLED
)

type Status struct {
	State     State
	CoAuthors []string
}

func Save(state State, coauthors ...string) error {
	cfg, _ := config.Load()

	status := Status{State: state, CoAuthors: coauthors}
	buf := new(bytes.Buffer)

	err := toml.NewEncoder(buf).Encode(status)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), []byte(buf.String()), 0644)
}

func Print() {
	cfg, _ := config.Load()

	var status Status
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), &status); err != nil {
		color.Red("Team mode disabled.")
		os.Exit(0)
	}

	switch status.State {
	case ENABLED:
		color.Green("Team mode enabled.")
		coauthors := status.CoAuthors
		if len(coauthors) > 0 {
			blackBold := color.New(color.FgBlack).Add(color.Bold)
			fmt.Println()
			blackBold.Println("Co-authors:")
			blackBold.Println("-----------")
			for _, coauthor := range coauthors {
				color.Magenta(coauthor)
			}
		}
	default:
		color.Red("Team mode disabled.")
	}
}
