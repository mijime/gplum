package main

var cmdStatus = &Command{
	Run:       runStatus,
	UsageLine: "status ",
	Short:     "",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdStatus.Flag.BoolVar(&flagA, "a", false, "")
}

// runStatus executes status command and return exit code.
func runStatus(args []string) int {

	return 0
}
