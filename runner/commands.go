package runner

func initialCommands(outdir string, wordlist string, resolver string) map[int]string {
	commands := map[int]string{
		// round1
		0: "assetfinder -subs-only  %[1]s | anew > " + outdir + "/%[1]s" + "/assetfinder",
		1: "subfinder -d %[1]s -o " + outdir + "/%[1]s" + "/subfinder",
		2: "amass enum -passive -d %[1]s > " + outdir + "/%[1]s" + "/amass",
		3: "cat " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass | deduplicate --sort > " + outdir + "/%[1]s" + "/round1",

		// delete assetfinder subfinder amass
		4: "rm -f " + outdir + "/%[1]s" + "/assetfinder " + outdir + "/%[1]s" + "/subfinder " + outdir + "/%[1]s" + "/amass 2>/dev/null",

		// step1
		5: "cp " + wordlist + " " + outdir + "/%[1]s" + "/shuffledns",
		6: "sed -e \"s/$/.%[1]s/\"  -i " + outdir + "/%[1]s" + "/shuffledns",
		7: "shuffledns -sw -list " + outdir + "/%[1]s" + "/shuffledns -silent -o " + outdir + "/%[1]s" + "/step1",

		// add new things to round1 && delete shuffledns
		8: "cat " + outdir + "/%[1]s" + "/step1 | anew -q " + outdir + "/%[1]s" + "/round1 && rm -f " + outdir + "/%[1]s" + "/shuffledns",

		// gotator - depth 2
		9: "gotator -silent -sub " + outdir + "/%[1]s" + "/round1 -depth 2 -mindup > " + outdir + "/%[1]s" + "/gotator",

		// step2 && delete gotator
		10: "shuffledns -sw -list " + outdir + "/%[1]s" + "/gotator -r " + resolver + " -silent -o " + outdir + "/%[1]s" + "/step2 && rm -f " + outdir + "/%[1]s" + "/gotator",

		// round2
		11: "cat " + outdir + "/%[1]s" + "/step1 " + outdir + "/%[1]s" + "/step2 | deduplicate --sort > " + outdir + "/%[1]s" + "/round2",

		// adding * at the begining of the resolved sub and make another gotator depth 2:

		// gotator2 - depth 3
		12: "gotator -silent -sub " + outdir + "/%[1]s" + "/round2 -depth 2 -mindup > " + outdir + "/%[1]s" + "/gotator2",

		// dnsx - tempararily stop deleting gotator2
		13: "dnsx -list " + outdir + "/%[1]s" + "/gotator2 -r " + resolver + " -silent -o " + outdir + "/%[1]s" + "/step3 && rm -f " + outdir + "/%[1]s" + "/gotator",

		// httpx - json file
		14: "cat " + outdir + "/%[1]s" + "/step3 | httpx -silent -sc -location -td -json -o " + outdir + "/%[1]s" + "/final",
	}

	return commands
}
