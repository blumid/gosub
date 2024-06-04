package runner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/progress"
	"github.com/jedib0t/go-pretty/text"
	fileutil "github.com/projectdiscovery/utils/file"
)

var MenuShown bool = true

type result struct {
	domain string
}

func runCommand(command string) {
	com := exec.Command("bash", "-c", command)
	if err := com.Run(); err != nil {
		fmt.Println("fuck! you have an error. maybe you didn't install requirement tools.")
		fmt.Println("error is:", err)
		os.Exit(1)
	}
}

func style() progress.Style {
	stylecol := progress.StyleColors{
		Message: text.Colors{text.FgHiWhite},
		Stats:   text.Colors{text.FgRed},
		Time:    text.Colors{text.FgRed},
		Percent: text.Colors{text.FgYellow},
		Value:   text.Colors{text.FgYellow},
		Tracker: text.Colors{text.FgRed},
	}

	style := progress.Style{
		Name:    "colorful shit",
		Colors:  stylecol,
		Options: progress.StyleOptionsDefault,
		Chars:   progress.StyleCharsDefault,
	}
	return style
}

func (options *Options) worker(domain string, commands map[int]string, pw progress.Writer, queue chan string) {
	var item result

	tracker := progress.Tracker{Message: domain, Total: int64(len(commands)), Units: progress.UnitsDefault}
	tracker.Reset()
	pw.AppendTracker(&tracker)

	//mkdir folder for each domain
	outdir := options.Output + "/" + domain
	os.MkdirAll(outdir, os.ModePerm)

	item.domain = domain

	for i := 0; i < len(commands); i++ {

		cmd := fmt.Sprintf(commands[i], domain)
		runCommand(cmd)
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

func Run(options *Options) {

	commands := initialCommands(options.Output, options.Wordlist, options.Resolver)

	pw := progress.NewWriter()
	pw.SetStyle(style())
	pw.SetOutputWriter(os.Stdout)

	// initialize verified domains
	var domains []string

	var sc *bufio.Scanner
	if fileutil.FileExists(options.Domains) {
		f, _ := os.Open(options.Domains)
		sc = bufio.NewScanner(f)
	} else if fileutil.HasStdin() {
		sc = bufio.NewScanner(os.Stdin)
	}

	for sc.Scan() {
		r, _ := regexp.Compile(`^\*\.`)

		if r.MatchString(sc.Text()) {
			domain := r.ReplaceAllString(sc.Text(), "")
			domains = append(domains, domain)
		}
	}

	// using buffered channel:
	queue := make(chan string, options.Concurency)

	go func() {
		for _, v := range domains {
			queue <- v
			// start go worker on v
			go options.worker(v, commands, pw, queue)

		}
		close(queue)

	}()
	pw.Render()
}

// phase 3:
func DisplayMenu() {
	MenuShown = false
	fmt.Println("Options:")
	fmt.Println("1. Nothing")
	fmt.Println("2. Targets")
	fmt.Println("3. Quit")
	fmt.Println("Enter your choice:")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')

	if choice == "2" {
		displayTargets()
		return
	}

	if choice == "\n" || choice == "1" {
		fmt.Println("Continuing...")
		MenuShown = false
		return
	}
	index, err := strconv.Atoi(choice[:len(choice)-1])
	if err != nil || index < 1 || index > 10 {
		fmt.Println("Invalid choice. Continuing...")
		MenuShown = false
		return
	}

	MenuShown = false
}

func displayTargets() {
	fmt.Println("displaying targets...")
}

// I should to check it out this part later:
// for _, cmdArgs := range commandList {
//     ctx, cancel := context.WithCancel(context.Background())
//     cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
//     cmd.Stdout = os.Stdout
//     cmd.Stderr = os.Stderr
//     commands = append(commands, cmd)
//     cancelFuncs = append(cancelFuncs, cancel)
//     go runCommand(cmd)
// }
