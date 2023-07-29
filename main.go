package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/blumid/tools/gosub/runner"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
)

func worker(domain string, commands map[int]string, pw progress.Writer, queue chan string) {
	var item result

	tracker := progress.Tracker{Message: domain, Total: int64(len(commands)), Units: progress.UnitsDefault}
	tracker.Reset()
	pw.AppendTracker(&tracker)

	//mkdir folder for each domain
	outdir := output + "/" + domain
	os.MkdirAll(outdir, os.ModePerm)

	item.domain = domain

	for i := 0; i < len(commands); i++ {

		// cmd := fmt.Sprintf(commands[i], domain)
		// item.runCommand(cmd)
		time.Sleep(time.Millisecond * 150)
		tracker.Increment(1)
	}
	time.Sleep(time.Second * 2)
	// wg.Done()
	if pw.LengthActive() == 0 {
		pw.Stop()
	}
	<-queue
}

func (i *result) runCommand(command string) {
	com := exec.Command("bash", "-c", command)
	if err := com.Run(); err != nil {
		fmt.Println("fuck! you have an error. maybe you didn't install requirement tools.")
		fmt.Println("error is:", err)
		os.Exit(1)
	}
}

type result struct {
	domain string
}

func initialCommands(outdir string, wordlist string) map[int]string {
	commands := map[int]string{
		// round1
		0: "assetfinder -subs-only  %[1]s | anew > " + outdir + "/%[1]s" + "/assetfinder",
		1: "subfinder -d %[1]s -o " + outdir + "/%[1]s" + "/subfinder",
		2: "amass enum -passive -d %[1]s > " + outdir + "/%[1]s" + "/amass",
		3: "cat " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass | deduplicate --sort > " + outdir + "/%[1]s" + "/round1",

		// delete assetfinder subfinder amass
		4: "rm -f " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass 2>/dev/null",

		// step1
		5: "cp " + wordlist + " " + outdir + "/%[1]s" + "/dnsx",
		6: "sed -e \"s/$/.%[1]s/\"  -i " + outdir + "/%[1]s" + "/dnsx",
		7: "dnsx -list " + outdir + "/%[1]s" + "/dnsx -silent -o " + outdir + "/%[1]s" + "/step1",

		// add new things to round1 && delete dnsx
		8: "cat " + outdir + "/%[1]s" + "/step1 | anew -q " + outdir + "/%[1]s" + "/round1 && rm -f " + outdir + "/%[1]s" + "/dnsx",

		// gotator - depth 2
		9: "gotator -silent -sub " + outdir + "/%[1]s" + "/round1 -depth 2 -mindup > " + outdir + "/%[1]s" + "/gotator",

		// step2 && delete gotator
		10: "dnsx -list " + outdir + "/%[1]s" + "/gotator -r " + resolver + " -silent -o " + outdir + "/%[1]s" + "/step2 && rm -f " + outdir + "/%[1]s" + "/gotator",

		// round2
		11: "cat " + outdir + "/%[1]s" + "/step1 " + outdir + "/%[1]s" + "/step2 | deduplicate --sort > " + outdir + "/%[1]s" + "/round2",

		// gotator2 - depth 3
		12: "gotator -silent -sub " + outdir + "/%[1]s" + "/round2 -depth 3 -mindup > " + outdir + "/%[1]s" + "/gotator2",

		// dnsx - tempararily stop deleting gotator2
		13: "dnsx -list " + outdir + "/%[1]s" + "/gotator2 -r " + resolver + " -silent -o " + outdir + "/%[1]s" + "/step3 && rm -f " + outdir + "/%[1]s" + "/gotator",

		// httpx - json file
		14: "cat " + outdir + "/%[1]s" + "/step3 | httpx -silent -sc -location -td -json -o " + outdir + "/%[1]s" + "/final",
	}

	return commands
}

func Style() progress.Style {
	stylecol := progress.StyleColors{
		Message: text.Colors{text.FgHiWhite},
		Stats:   text.Colors{text.FgRed},
		Time:    text.Colors{text.FgRed},
		Percent: text.Colors{text.FgYellow},
		Value:   text.Colors{text.FgYellow},
		Tracker: text.Colors{text.FgRed},
	}

	style := progress.Style{
		Name:       "colorful shit",
		Colors:     stylecol,
		Options:    progress.StyleOptionsDefault,
		Chars:      progress.StyleCharsDefault,
		Visibility: progress.StyleVisibilityDefault,
	}
	return style
}

var wordlist string
var resolver string
var output string
var max int
var silent bool

func main() {

	flag.StringVar(&output, "o", "output", "output directory")
	flag.StringVar(&wordlist, "w", "sort_subs12.txt", "wordlist path")
	flag.StringVar(&resolver, "r", "resolvers.txt", "resolver path")
	flag.IntVar(&max, "m", 5, "maximum number of Synchronized process, <=5")
	flag.BoolVar(&silent, "s", false, "silent mode")

	flag.Parse()

	if !silent {
		runner.ShowBanner()
	}

	pw := progress.NewWriter()
	pw.SetStyle(Style())
	pw.SetOutputWriter(os.Stdout)

	commands := initialCommands(output, wordlist)

	sc := bufio.NewScanner(os.Stdin)
	// initialize verified domains
	var domains []string
	for sc.Scan() {
		r, _ := regexp.Compile(`^\*\.`)

		if r.MatchString(sc.Text()) {
			domain := r.ReplaceAllString(sc.Text(), "")
			domains = append(domains, domain)
		}
	}

	// using buffered channel:
	if max > 5 || max <= 0 {
		max = 5
	}
	queue := make(chan string, max)

	go func() {
		for _, v := range domains {
			queue <- v
			// start go worker on v
			go worker(v, commands, pw, queue)

		}
		close(queue)

	}()
	pw.Render()

}
