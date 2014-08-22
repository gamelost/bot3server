package util

import (
	"fmt"
	"strings"
)

const COMMAND_PREFIX = "!"

func TrimCommandFromString(input string, command string) string {

	if StringIsCommand(input) {
		lineTxt := strings.TrimPrefix(input, fmt.Sprintf("%s%s", COMMAND_PREFIX, command))
		return strings.TrimSpace(lineTxt)
	} else {
		return input
	}
}

func StringIsCommand(rawstring string) bool {

	return strings.HasPrefix(rawstring, COMMAND_PREFIX) && !strings.HasPrefix(rawstring, COMMAND_PREFIX+" ")
}

func GetCommandFromString(rawstring string) string {

	commandStr := ""

	if StringIsCommand(rawstring) {
		index := strings.Index(rawstring, " ")
		if index > 0 {
			commandStr = rawstring[1:index]
		} else if index == -1 {
			// no space, so return entirety of string except first char
			commandStr = rawstring[1:len(rawstring)]
		}
	}

	return commandStr
}
