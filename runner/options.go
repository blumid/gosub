package runner

import (
	"bufio"
	"os"

	"github.com/projectdiscovery/goflags"
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

	flagSet.StringVarP(&options.Domains, "domains", "l", "", "list wildcard domains resolve (file or stdin)")
	flagSet.StringVar(&options.Output, "o", "scopes", "output folder")
	flagSet.StringVar(&options.Wordlist, "w", "~/BugBounty/wordlist/sort_subs12.txt", "wordlist file")
	flagSet.StringVar(&options.Resolver, "r", "~/BugBounty/wordlist/resolvers.txt", "resolver file")
	flagSet.IntVar(&options.Concurency, "c", 3, "maximum number of concurency process - max:5")
	flagSet.BoolVar(&options.Silent, "s", false, "silent mode - no banner")
	options.validateOptions()
	return options
}

func (options *Options) validateOptions() {
	if options.Concurency > 5 || options.Concurency <= 0 {
		options.Concurency = 3
	}

	var sc *bufio.Scanner

	if fileutil.FileExists(options.Wordlist) {
		f, _ := os.Open(options.Domains)
		sc = bufio.NewScanner(f)
	} else if fileutil.HasStdin() {
		sc = bufio.NewScanner(os.Stdin)
	}

	for sc.Scan() {

	}
}
