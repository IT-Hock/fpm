package utils

import (
	"fmt"
	"golang.org/x/term"
	"os"
)

// AskYesNo BUG: Pressing up causes weird issues. -> https://i.imgur.com/xwUJ6C1.png
func AskYesNo(question string) bool {
	if FlagYes() {
		Println("%s (<green>y</green>/<red>n</red>): <green>y</green>", question)
		return true
	}

	if FlagNo() {
		Println("%s (<green>y</green>/<red>n</red>): <red>n</red>", question)
		return false
	}

	var answer byte

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	Print(question + " (<green>y</green>/<red>n</red>): ")
	for answer != 'y' && answer != 'n' {
		if answer == 3 {
			os.Exit(0)
		}

		b := make([]byte, 1)
		_, err = os.Stdin.Read(b)
		if err != nil {
			fmt.Println(err)
			return false
		}
		answer = b[0]
	}

	if answer == 'y' {
		Print("<green>y</green>\n")
	} else {
		Print("<red>n</red>\n")
	}
	defer func(fd int, oldState *term.State) {
		err := term.Restore(fd, oldState)
		if err != nil {
			panic(err)
		}
	}(int(os.Stdin.Fd()), oldState)

	// Move cursor to the left
	print("\033[500D")

	return answer == 'y' || answer == 'Y'
}
