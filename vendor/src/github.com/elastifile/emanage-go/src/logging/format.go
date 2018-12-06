package logging

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	log "gopkg.in/inconshreveable/log15.v2"
)

const (
	logMsgLen     = 40
	logTimeFormat = "2006-01-02 15:04:05.000"
	htmlBody      = `<body style="font-family:Monospace">
`
	htmlEntryFmt = `%s</span><br/>
`
)

var lvlColorCodeMap = map[log.Lvl]int{
	log.LvlCrit:   tm.MAGENTA,
	log.LvlError:  tm.RED,
	log.LvlWarn:   tm.YELLOW,
	log.LvlInfo:   tm.GREEN,
	log.LvlDebug:  tm.CYAN,
}

var htmlColorMap = map[int]string{
	30 + tm.MAGENTA: "Magenta",
	30 + tm.RED:     "Red",
	30 + tm.YELLOW:  "Gold",
	30 + tm.BLUE:    "Blue",
	30 + tm.GREEN:   "LimeGreen",
	30 + tm.CYAN:    "Turquoise",
	37:              "Gray",
}

func recordToBytes(rec *log.Record, color int, paddingFunc func(string) int) []byte {
	buf := getLogEntryBuffer(rec, color, paddingFunc)
	logKeyValues(buf, rec.Ctx, color)

	return buf.Bytes()
}

func coloredEntryFormat(paddingFunc func(string) int) log.Format {
	return log.FormatFunc(func(rec *log.Record) []byte {
		color := lvlColorCodeMap[rec.Lvl]
		return recordToBytes(rec, color, paddingFunc)
	})
}

func monochromeEntryFormat(paddingFunc func(string) int) log.Format {
	return log.FormatFunc(func(rec *log.Record) []byte {
		color := 0
		return recordToBytes(rec, color, paddingFunc)
	})
}

func htmlEntryFormat(paddingFunc func(string) int) log.Format {
	return log.FormatFunc(func(rec *log.Record) []byte {
		color := lvlColorCodeMap[rec.Lvl]
		buf := recordToBytes(rec, color, paddingFunc)
		entry := fmt.Sprintf(htmlEntryFmt, htmlColor(string(buf)))
		return []byte(entry)
	})
}

// formatters utils
func noPadding(msg string) int {
	return 0
}

func msgPadding(msg string) int {
	lenMsg := lenNoColor(msg)
	if lenMsg < logMsgLen {
		return logMsgLen - lenMsg
	}
	return 0
}

func getLogEntryBuffer(rec *log.Record, color int, paddingFunc func(string) int) *bytes.Buffer {
	buf := &bytes.Buffer{}

	lvl := strings.ToUpper(rec.Lvl.String())
	fmt.Fprintf(buf, "%s %s: %s ", rec.Time.Format(logTimeFormat), ifColor(lvl, color), noColor(rec.Msg, color))
	//fmt.Fprintf(buf, "%s %.10s@%s %s: %s ", rec.Time.Format(logTimeFormat), os.Args[0], helputils.MyExternalIP().String(), ifColor(lvl, color), noColor(rec.Msg, color))

	padSize := paddingFunc(rec.Msg)
	buf.Write(bytes.Repeat([]byte{' '}, padSize))

	return buf
}

func escapeString(s string) string {
	needQuotes := false
	e := bytes.Buffer{}
	e.WriteByte('"')
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			needQuotes = true
		}

		switch r {
		case '\\', '"':
			e.WriteByte('\\')
			e.WriteByte(byte(r))
		case '\n':
			e.WriteByte('\\')
			e.WriteByte('n')
		case '\r':
			e.WriteByte('\\')
			e.WriteByte('r')
		case '\t':
			e.WriteByte('\\')
			e.WriteByte('t')
		default:
			e.WriteRune(r)
		}
	}
	e.WriteByte('"')
	start, stop := 0, e.Len()
	if !needQuotes {
		start, stop = 1, stop-1
	}
	return string(e.Bytes()[start:stop])
}

const (
	floatFormat = 'f'
	timeFormat  = "2006-01-02T15:04:05-0700"
)

func formatShared(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

// formatValue formats a value for serialization
func formatLogfmtValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	value = formatShared(value)
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, 3, 64)
	case float64:
		return strconv.FormatFloat(v, floatFormat, 3, 64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case string:
		return escapeString(v)
	default:
		return escapeString(fmt.Sprintf("%+v", value))
	}
}

func logKeyValues(buf *bytes.Buffer, ctx []interface{}, color int) {
	pkg := ""
	var caller *string

	for i := 0; i < len(ctx); i += 2 {
		val := formatLogfmtValue(ctx[i+1])
		key, ok := ctx[i].(string)
		if !ok {
			key, val = "EROR", formatLogfmtValue(key)
		}

		switch key {
		case "package":
			if pkg == "" {
				pkg = val
			}
			continue
		case "caller":
			if !(caller != nil && *caller == "") {
				caller = &val
			}
			continue
		}

		if val != "" {
			fmt.Fprintf(buf, "%s=%s ", ifColor(key, color), val)
		}
	}

	if caller != nil && *caller != "" {
		fmt.Fprint(buf, "  ", ifColor("<", color), pkg, ifColor("|", color), *caller, ifColor(">", color))
	}

	fmt.Fprintln(buf)
}

func lenNoColor(str string) int {
	return len(noColor(str, 0))
}

func noColor(str string, color int) string {
	if color == 0 {
		re := regexp.MustCompile("\x1b\\[(\\d|\\.)+m")
		return re.ReplaceAllString(str, "")
	}
	return str
}

func htmlColor(str string) string {
	for code, name := range htmlColorMap {
		reStr := fmt.Sprintf("\x1b\\[%dm", code)
		re := regexp.MustCompile(reStr)
		str = re.ReplaceAllString(str, fmt.Sprintf(`<span style="color:%s">`, name))
	}

	re := regexp.MustCompile("\x1b\\[1m")
	str = re.ReplaceAllString(str, `<span style="font-weight:bold">`)
	re = regexp.MustCompile("\x1b\\[0m")
	str = re.ReplaceAllString(str, "</span>")

	str = strings.Replace(str, "  ", "&nbsp;&nbsp;", -1)
	str = strings.Replace(str[:len(str)-1], "\n", "<br>", -1)

	return str
}

func ifColor(s string, color int) string {
	if color > 0 {
		return tm.Color(s, color)
	}
	return s
}

func Framed(msg string, char byte) string {
	if msg != "" {
		msg = " " + msg + " "
	}
	frameLen := (logMsgLen - len(msg)) / 2
	if frameLen < 3 {
		frameLen = 3
	}
	frame := bytes.Repeat([]byte{char}, frameLen)
	return fmt.Sprintf("%s%s%s", frame, msg, frame)
}
