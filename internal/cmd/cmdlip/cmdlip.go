package cmdlip

import (
	"flag"
	"fmt"
	"os"

	"github.com/lippkg/lip/internal/cmd/cmdlipcache"
	"github.com/lippkg/lip/internal/cmd/cmdlipinstall"
	"github.com/lippkg/lip/internal/cmd/cmdliplist"
	"github.com/lippkg/lip/internal/cmd/cmdlipshow"
	"github.com/lippkg/lip/internal/cmd/cmdliptooth"
	"github.com/lippkg/lip/internal/cmd/cmdlipuninstall"
	"github.com/lippkg/lip/internal/context"

	log "github.com/sirupsen/logrus"
)

type FlagDict struct {
	helpFlag    bool
	versionFlag bool
	verboseFlag bool
	quietFlag   bool
}

const helpMessage = `
Usage:
  lip [options] [<command> [subcommand options]] ...

Commands:
  cache                       Inspect and manage lip's cache.
  install                     Install a tooth.
  list                        List installed teeth.
  show                        Show information about installed teeth.
  tooth                       Maintain a tooth.
  uninstall                   Uninstall a tooth.

Options:
  -h, --help                  Show help.
  -V, --version               Show version and exit.
  -v, --verbose               Show verbose output.
  -q, --quiet                 Show only errors.
`

func Run(ctx context.Context, args []string) error {
	var err error

	flagSet := flag.NewFlagSet("lip", flag.ContinueOnError)

	// Rewrite the default messages.
	flagSet.Usage = func() {
		// Do nothing.
	}

	// Parse flags.
	var flagDict FlagDict
	flagSet.BoolVar(&flagDict.helpFlag, "help", false, "")
	flagSet.BoolVar(&flagDict.helpFlag, "h", false, "")
	flagSet.BoolVar(&flagDict.versionFlag, "version", false, "")
	flagSet.BoolVar(&flagDict.versionFlag, "V", false, "")
	flagSet.BoolVar(&flagDict.verboseFlag, "verbose", false, "")
	flagSet.BoolVar(&flagDict.verboseFlag, "v", false, "")
	flagSet.BoolVar(&flagDict.quietFlag, "quiet", false, "")
	flagSet.BoolVar(&flagDict.quietFlag, "q", false, "")
	err = flagSet.Parse(args)
	if err != nil {
		return fmt.Errorf("cannot parse flags: %w", err)
	}

	// Set logging level.
	if flagDict.verboseFlag {
		log.SetLevel(log.DebugLevel)
	} else if flagDict.quietFlag {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Help flag has the highest priority.
	if flagDict.helpFlag {
		fmt.Print(helpMessage)
		return nil
	}

	// Version flag has the second highest priority.
	if flagDict.versionFlag {
		fmt.Printf("lip %v from %v", ctx.LipVersion().String(), os.Args[0])
		return nil
	}

	// Verbose and quiet flags are mutually exclusive.
	if flagDict.verboseFlag && flagDict.quietFlag {
		return fmt.Errorf("verbose and quiet flags are mutually exclusive")
	}

	// If there is a subcommand, run it and exit.
	if flagSet.NArg() >= 1 {
		switch flagSet.Arg(0) {
		case "cache":
			err = cmdlipcache.Run(ctx, flagSet.Args()[1:])
			if err != nil {
				return err
			}
			return nil

		case "install":
			err = cmdlipinstall.Run(ctx, flagSet.Args()[1:])
			if err != nil {
				return err
			}
			return nil

		case "list":
			err = cmdliplist.Run(ctx, flagSet.Args()[1:])
			if err != nil {
				return err
			}
			return nil

		case "show":
			err = cmdlipshow.Run(ctx, flagSet.Args()[1:])
			if err != nil {
				return err
			}
			return nil

		case "tooth":
			err = cmdliptooth.Run(ctx, flagSet.Args()[1:])
			if err != nil {
				return err
			}
			return nil

		case "uninstall":
			err = cmdlipuninstall.Run(ctx, flagSet.Args()[1:])
			if err != nil {
				return err
			}
			return nil

		default:
			return fmt.Errorf("unknown command: lip %v", flagSet.Arg(0))
		}
	}

	return fmt.Errorf("no command specified. See 'lip --help' for more information")
}
