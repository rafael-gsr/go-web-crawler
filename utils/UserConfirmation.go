// Package utils contains utilities used across the module
package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func DoUserWantsToContinue(confirmationMessage string) (string, error) {
	YES := "yes"
	NO := "no"
	QUIT := "quit"

	valid := false
	answer := ""

	for !valid {
		fmt.Printf("%s [Y/n/q]\n", confirmationMessage)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		switch strings.ToLower(strings.TrimSpace(response)) {
		case "y":
			answer = YES
			valid = true
			continue

		case "":
			answer = YES
			valid = true
			continue

		case "n":
			answer = NO
			valid = true
			continue

		case "q":
			answer = QUIT
			valid = true
			continue
		}

	}

	return answer, nil
}
