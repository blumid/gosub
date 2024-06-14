package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/progress"
	"github.com/jedib0t/go-pretty/text"
	fileutil "github.com/projectdiscovery/utils/file"
)

var MenuShown bool = true
var pw progress.Writer

const (
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Reset   = "\033[0m"
)

type result struct {
	domain string
}

type ContextWithID struct {
	item    string
	context context.Context
	cancel  context.CancelFunc
}

func init() {
	pw = progress.NewWriter()
}

func runCommand(command string) {
	com := exec.Command("bash", "-c", command)
	if err := com.Run(); err != nil {
		fmt.Println("runCommand() - error:", err)
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

func (options *Options) worker(domain string, pw progress.Writer, queue chan string) {
	var item result
	commands := initialCommands(options.Output, options.Wordlist, options.Resolver)

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
	// pw := progress.NewWriter()
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

	Contexts := make([]*ContextWithID, 0, options.Concurency)

	go func() {
		for _, v := range domains {
			queue <- v

			// make a new ctx
			ctx, cancel := context.WithCancel(context.Background())
			Contexts = append(Contexts, &ContextWithID{
				item:    v,
				context: ctx,
				cancel:  cancel,
			})

			// start go worker on v
			go options.worker(v, pw, queue)

		}
		close(queue)

	}()

	pw.Render()

}

// phase 3:
func DisplayMenu() {
	tempFile, err := os.CreateTemp("", "temp_stdout")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	pw.SetOutputWriter(tempFile)

	clearTerminal()

	options := []string{
		Magenta + "Nothing! I just have ass's worm :)" + Reset,
		Magenta + "delete - list" + Reset,
		Magenta + "quit" + Reset,
	}
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	fmt.Print("Select an option: ")
	// Read user input
	// reader := bufio.NewReader(os.Stdin)
	// input, _ := reader.ReadString('\n')
	var input string = ""
	fmt.Scanln(&input)
	fmt.Printf("input is: %s,%s\n", input, "A")

	// Clean up the input
	input = strings.TrimSpace(input)

	// Convert input to integer
	// index, err := strconv.Atoi(input)
	// if err != nil || index < 1 || index > len(options) {
	// 	fmt.Println("Invalid option selected.")
	// 	return
	// }

	switch input {
	case "1":
		time.Sleep(time.Millisecond * 500)
		return
	case "2":
		displayTargets()
	case "3":
		os.Exit(1)
	default:
		os.Exit(2)
	}

	time.Sleep(time.Second * 3)
	pw.SetOutputWriter(os.Stdout)
	MenuShown = true

}

func displayTargets() {
	fmt.Println("displaying targets...")
}
func clearTerminal() {
	cmd := exec.Command("clear") // Use "cls" on Windows instead of "clear"
	cmd.Stdout = os.Stdout
	cmd.Run()
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
