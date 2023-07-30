package runner

import (
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
)

type Options struct {
	Domains    string
	Output     string
	Wordlist   string
	Resolver   string
	Concurency int
	Silent     bool
}

func ParseOptions() *Options {
	options := &Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`gosub is a multi-thread tool that run some tools to find subdomains`)

	flagSet.StringVar(&options.Domains, "l", "", "list wildcard domains resolve (file or stdin)")
	flagSet.StringVar(&options.Output, "o", "scopes", "output folder")
	flagSet.StringVar(&options.Wordlist, "w", "~/BugBounty/wordlist/sort_subs12.txt", "wordlist file")
	flagSet.StringVar(&options.Resolver, "r", "~/BugBounty/wordlist/resolvers.txt", "resolver file")
	flagSet.IntVar(&options.Concurency, "c", 3, "maximum number of concurency processes - max:5")
	flagSet.BoolVar(&options.Silent, "s", false, "silent mode - no banner")

	_ = flagSet.Parse()
	options.validateOptions()
	return options
}

func (options *Options) validateOptions() {
	if options.Concurency > 5 || options.Concurency <= 0 {
		options.Concurency = 3
	}

	// validate wordlist
	if !fileutil.FileExists(options.Wordlist) {
		gologger.Fatal().Msg("can't find wordlist file [-w flag]")
	}

	// validate domains
	if fileutil.FileExists(options.Domains) || fileutil.HasStdin() {

	} else {
		gologger.Fatal().Msg("can't find input [-l flag]/StdIn")
	}

	// validate resolver
	if !fileutil.FileExists(options.Resolver) {
		gologger.Fatal().Msg("can't find resolver file [-r flag]")
	}

	// showing banner
	if !options.Silent {
		ShowBanner()
	}

}
