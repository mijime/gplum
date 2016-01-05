package main

var cmdSync = &Command{
	Run:       runSync,
	UsageLine: "sync ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdSync.Flag.BoolVar(&flagA, "a", false, "")
}

// runSync executes sync command and return exit code.
func runSync(args []string) int {

	return 0
}
