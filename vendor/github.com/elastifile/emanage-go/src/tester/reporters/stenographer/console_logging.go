package stenographer

import (
	"fmt"
	"strings"

	"tester/internal/logging"
)

func (s *consoleStenographer) colorize(colorCode string, format string, args ...interface{}) string {
	var out string

	if len(args) > 0 {
		out = fmt.Sprintf(format, args...)
	} else {
		out = format
	}

	normalStyle := defaultStyle
	if strings.HasSuffix(colorCode, boldStyle) {
		normalStyle += normalStyle
	}
	if s.color {
		return fmt.Sprintf("%s%s%s", colorCode, out, normalStyle)
	} else {
		return out
	}
}

func (s *consoleStenographer) printBanner(text string, bannerCharacter string) {
	logging.Log().Info(text)
	logging.Log().Info(strings.Repeat(bannerCharacter, len(text)))
}

var printBuffer string

func (s *consoleStenographer) getPrintBuffer() string {
	return printBuffer
}

func (s *consoleStenographer) flushLine() {
	s.printNewLine()
}

func (s *consoleStenographer) printNewLine() {
	logging.Log().Info(printBuffer)
	printBuffer = ""
}

func (s *consoleStenographer) printDelimiter() {
	s.info(0, s.colorize(grayColor, "%s", strings.Repeat("-", 30)))
}

func (s *consoleStenographer) print(indentation int, format string, args ...interface{}) {
	printBuffer += s.indent(indentation, format, args...)
}

func (s *consoleStenographer) println(indentation int, format string, args ...interface{}) {
	s.info(indentation, format, args...)
}

func (s *consoleStenographer) info(indentation int, format string, args ...interface{}) {
	logging.Log().Info(printBuffer + s.indent(indentation, format, args...))
	printBuffer = ""
}

func (s *consoleStenographer) error(indentation int, format string, args ...interface{}) {
	logging.Log().Error(printBuffer + s.indent(indentation, format, args...))
	printBuffer = ""
}

func (s *consoleStenographer) warn(indentation int, format string, args ...interface{}) {
	logging.Log().Warn(printBuffer + s.indent(indentation, format, args...))
	printBuffer = ""
}

func (s *consoleStenographer) action(indentation int, format string, args ...interface{}) {
	logging.Log().Info(printBuffer + s.indent(indentation, format, args...))
	printBuffer = ""
}

func (s *consoleStenographer) indent(indentation int, format string, args ...interface{}) string {
	var text string

	if len(args) > 0 {
		text = fmt.Sprintf(format, args...)
	} else {
		text = format
	}

	stringArray := strings.Split(text, "\n")
	padding := ""
	if indentation >= 0 {
		padding = strings.Repeat("  ", indentation)
	}
	for i, s := range stringArray {
		stringArray[i] = fmt.Sprintf("%s%s", padding, s)
	}

	return strings.Join(stringArray, "\n")
}
