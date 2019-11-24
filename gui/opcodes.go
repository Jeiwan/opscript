package gui

import (
	"regexp"
)

func formatOpcode(opcode string) string {
	opData := regexp.MustCompile(`OP_DATA_\d+ `)
	opcode = opData.ReplaceAllString(opcode, "")

	return opcode
}
