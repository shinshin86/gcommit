package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("git", "diff")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running git diff: %v\n", err)
		return
	}

	diff := strings.TrimSpace(string(output))
	fmt.Println(diff)
}
