package utils

import (
	"fmt"
	"fpm/src/build"
	"regexp"
	"strings"
)

const (
	Black  = "\033[30m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	White  = "\033[37m"
)

func PrintDebug(format string, a ...any) {
	if build.Debug {
		fmt.Println(Colorize(Blue, format, a...))
	}
}

func Print(format string, a ...any) {
	fmt.Print(ColorizeHtml(format, a...))
}

func Println(format string, a ...any) {
	fmt.Println(ColorizeHtml(format, a...))
}

func ColorizeHtml(format string, a ...any) string {
	// This has an issue. It messes with the following:
	// <red>Test</blue></red>

	colorTagRegexp, err := regexp.Compile("<(red|green|blue|yellow|white|black)>(.*?)</(red|green|blue|yellow|white|black)>")
	if err != nil {
		panic(err)
	}

	matches := colorTagRegexp.FindAllStringSubmatch(format, -1)
	for _, match := range matches {
		switch match[1] {
		case "red":
			format = strings.Replace(format, match[0], PlainColorize(Red, match[2]), 1)
		case "green":
			format = strings.Replace(format, match[0], PlainColorize(Green, match[2]), 1)
		case "blue":
			format = strings.Replace(format, match[0], PlainColorize(Blue, match[2]), 1)
		case "yellow":
			format = strings.Replace(format, match[0], PlainColorize(Yellow, match[2]), 1)
		case "white":
			format = strings.Replace(format, match[0], PlainColorize(White, match[2]), 1)
		case "black":
			format = strings.Replace(format, match[0], PlainColorize(Black, match[2]), 1)
		default:
			panic("Unknown color: " + match[1])
		}
	}

	return fmt.Sprintf(format, a...)
}

func PlainColorize(color string, str string) string {
	return color + str + "\033[0m"
}

func Colorize(color string, format string, a ...interface{}) string {
	var newA []interface{}
	for _, v := range a {
		if str, ok := v.(string); ok {
			if strings.Contains(str, "\033[") {
				newA = append(newA, str+color)
			} else {
				newA = append(newA, str)
			}
		}
	}

	return color + fmt.Sprintf(format, newA...) + "\033[0m"
}

func ColorizeRGB(red int, green int, blue int, text string) string {
	red = Clamp(red, 0, 255)
	green = Clamp(green, 0, 255)
	blue = Clamp(blue, 0, 255)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", red, green, blue, text)
}

func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func Italic(text string) string {
	return "\033[3m" + text + "\033[0m"
}

func Underline(text string) string {
	return "\033[4m" + text + "\033[0m"
}

func Strike(text string) string {
	return "\033[9m" + text + "\033[0m"
}
