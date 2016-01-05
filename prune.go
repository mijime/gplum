package main

var cmdPrune = &Command{
	Run:       runPrune,
	UsageLine: "prune ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdPrune.Flag.BoolVar(&flagA, "a", false, "")
}

// runPrune executes prune command and return exit code.
func runPrune(args []string) int {

	return 0
}
