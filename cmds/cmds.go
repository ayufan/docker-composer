package cmds

import (
	"sort"

	"github.com/urfave/cli"
)

type CommandSlice []cli.Command

func (p CommandSlice) Len() int           { return len(p) }
func (p CommandSlice) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p CommandSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var Commands CommandSlice

func registerCommand(cmd cli.Command) {
	Commands = append(Commands, cmd)
	sort.Sort(Commands)
}
