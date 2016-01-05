package main

var cmdDiff = &Command{
	Run:       runDiff,
	UsageLine: "diff ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdDiff.Flag.BoolVar(&flagA, "a", false, "")
}

// runDiff executes diff command and return exit code.
func runDiff(args []string) int {

	return 0
}
