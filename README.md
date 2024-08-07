# gosub

Accept line-delimited wild domains on stdin and execute some tools on them and store output in a directory which you determined using `-o` flag.

at the end you have `final` and `hidden` for each wild domain

- `hidden` : contains subs which have explored with random words and methods
- `final` : contains all results

<!-- <h1 align="center"> -->
  <img src="static/gosub_run_v2.jpg" alt="gosub"></a>
<!-- </h1> -->


# Usage

```bash
cat list | gosub -w wordlist.txt -r resolver.txt -o scopes
```
press `esc` key on the keyboard to show you prompt, there you can stop some domain progress.

### prompt:
<img src="static/prompt.png" alt="gosub"></a>


in order to get out the of prompt and then go to process screen just press `ctrl+D`.

### result:
<img src="static/canceled_abbas3.png" alt="gosub"></a>

# Install

```bash
go install github.com/blumid/gosub@latest
```

<br>

```yaml
Flags:
   -l string  list wildcard domains resolve (file or stdin)
   -o string  output folder (default "scopes")
   -w string  wordlist file (default "~/BugBounty/wordlist/sort_subs12.txt")
   -r string  resolver file (default "~/BugBounty/wordlist/resolvers.txt")
   -c int     maximum number of concurency processes - max:5 (default 3)
   -s         silent mode - no banner

```


**Notice:** you have to install below tools before run this tool
Requirement tools:
* amass
* subfinder
* assetfinder
* dnsx
* gotator
* anew
* shuffledns


# Map
<h1 align="center">
  <img src="static/Map.png" alt="gosub"></a>
</h1>