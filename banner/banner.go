package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.1"

func PrintVersion() {
	fmt.Printf("Current querygen version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
  ____ _ __  __ ___   _____ __  __ ____ _ ___   ____ 
 / __  // / / // _ \ / ___// / / // __  // _ \ / __ \
/ /_/ // /_/ //  __// /   / /_/ // /_/ //  __// / / /
\__, / \__,_/ \___//_/    \__, / \__, / \___//_/ /_/ 
  /_/                    /____/ /____/
`
    fmt.Printf("%s\n%55s\n\n", banner, "Current querygen version "+version)
}
