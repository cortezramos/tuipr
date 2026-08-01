package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Flags represents command-line arguments.
type Flags struct {
	CreatePR   bool
	MergePR    bool
	MergePRNum int // 0 means not specified.
	Help       bool
}

// resetFlags creates a fresh flag set for testing.
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet("tuipr", flag.ExitOnError)
}

// parseFlags parses command-line arguments.
func parseFlags() Flags {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "tuipr - PR Lifecycle TUI Manager\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  tuipr              Open main dashboard\n")
		fmt.Fprintf(os.Stderr, "  tuipr -c           Open Create PR view\n")
		fmt.Fprintf(os.Stderr, "  tuipr -m           Open Merge PR view (select PR)\n")
		fmt.Fprintf(os.Stderr, "  tuipr -m <num>     Direct merge of PR #<num>\n")
		fmt.Fprintf(os.Stderr, "  tuipr --help       Show this help\n")
		fmt.Fprintf(os.Stderr, "\nNavigation:\n")
		fmt.Fprintf(os.Stderr, "  1-3    Switch between panels\n")
		fmt.Fprintf(os.Stderr, "  j/k    Navigate up/down\n")
		fmt.Fprintf(os.Stderr, "  i      Insert mode\n")
		fmt.Fprintf(os.Stderr, "  Esc    Return to Normal / Back\n")
		fmt.Fprintf(os.Stderr, "\nActions:\n")
		fmt.Fprintf(os.Stderr, "  c      Create PR\n")
		fmt.Fprintf(os.Stderr, "  m      Merge PR\n")
		fmt.Fprintf(os.Stderr, "  e      Open nvim to resolve conflicts\n")
		fmt.Fprintf(os.Stderr, "  r      Refresh list\n")
		fmt.Fprintf(os.Stderr, "  q      Quit\n")
	}

	createPtr := flag.Bool("c", false, "Open directly in Create PR view")
	mergePtr := flag.Bool("m", false, "Open directly in Merge PR view")
	helpPtr := flag.Bool("help", false, "Show help")

	flag.Parse()

	// Check if there's a number after -m.
	mergeNum := 0
	if flag.NArg() > 0 {
		numStr := flag.Arg(0)
		if num, err := strconv.Atoi(numStr); err == nil {
			mergeNum = num
		}
	}

	return Flags{
		CreatePR:   *createPtr,
		MergePR:    *mergePtr || mergeNum > 0,
		MergePRNum: mergeNum,
		Help:       *helpPtr,
	}
}
