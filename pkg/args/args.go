package args

import (
	"flag"
	"os"
)

// Arguments is a struct that holds the parsed arguments.
type Arguments struct {
	SrcRootPath             string
	DstRootPath             string
	ReplaceNotMatchingFiles bool
	RemoveDstLeftover       bool
}

// Parse parses the arguments.
func Parse(osArgs []string) Arguments {
	args := Arguments{}
	flagSet := flag.NewFlagSet(osArgs[0], flag.ExitOnError)

	flagSet.StringVar(&args.SrcRootPath, "src", "", "The source root path (required)")
	flagSet.StringVar(&args.DstRootPath, "dst", "", "The source root path (required)")
	flagSet.BoolVar(&args.ReplaceNotMatchingFiles, "replace", false, "Replace file on dst when different")
	flagSet.BoolVar(&args.RemoveDstLeftover, "remove", false, "Remove files and directories in dst not included in src")

	if err := flagSet.Parse(osArgs[1:]); err != nil ||
		len(args.SrcRootPath) == 0 || len(args.DstRootPath) == 0 {
		flagSet.Usage()
		os.Exit(1)
	}

	return args
}
