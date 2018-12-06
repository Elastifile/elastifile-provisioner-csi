package helputils

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"encoding/json"
	"github.com/fatih/structs"
	"github.com/go-errors/errors"
	"github.com/kardianos/osext"
)

func init() {
	StartupFolder()
}

// RangeInt returns slice of integers initialized with a range of numbers from n to m(excluded) -> [n,m)
// Example: RangeInt(n, m) -> []int{n, n+1, n+2..m-1}
func RangeInt(n, m int) (rng []int) {
	for i := n; i < m; i++ {
		rng = append(rng, i)
	}
	return
}

// ShowElapsed Will print context and elapsed time every given interval
// until signaled to stop.
//
// e.g. "(Add Node) Elapsed 1.15m"
func ShowElapsed(context string, interval time.Duration) chan<- bool {
	done := make(chan bool)

	start := time.Now().Round(time.Second)
	go func() {
		for {
			select {
			case <-done:
				break
			default:
				fmt.Printf("(%v) Elapsed %v\n", context, time.Now().Round(time.Second).Sub(start))
			}
			time.Sleep(interval)
		}
	}()

	return done
}

func MustStructToKeyValueInterfaces(strct interface{}) (kvIfs []interface{}) {
	if !IsReflectingStruct(strct) {
		panic(fmt.Sprintf("must have kind struct or *struct, received: %v", strct))
	}

	StructExportedFieldIter(strct, func(field *structs.Field) {
		if !field.IsZero() {
			if IsReflectingStruct(field.Value()) {
				kvIfs = append(kvIfs, MustStructToKeyValueInterfaces(field.Value())...)
			} else {
				kvIfs = append(kvIfs, field.Name(), field.Value())
			}
		}
	})

	return kvIfs
}

func IsReflectingStruct(f interface{}) bool {
	v := reflect.ValueOf(f)
	return v.Kind() == reflect.Struct ||
		v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct
}

func StructExportedFieldIter(strctIf interface{}, fn func(field *structs.Field)) {
	strct := structs.New(strctIf)
	for _, field := range strct.Fields() {
		if field.IsExported() {
			fn(field)
		}
	}
}

type StringSet []string

func (ss *StringSet) Add(values ...string) {
	for _, val := range values {
		if !ContainsStr(*ss, val) {
			*ss = append(*ss, val)
		}
	}
}

func (ss *StringSet) Contains(value string) bool {
	return ContainsStr(*ss, value)
}

func (ss *StringSet) Sub(other *StringSet) (result StringSet) {
	for _, val := range *ss {
		if !ContainsStr(*other, val) {
			result = append(result, val)
		}
	}
	return
}

func (ss *StringSet) Union(other *StringSet) (result StringSet) {
	result = make(StringSet, len(*ss))
	copy(result, *ss)
	result.Add(*other...)
	return
}

func (ss *StringSet) Intersect(other *StringSet) (result StringSet) {
	for _, val := range *ss {
		if ContainsStr(*other, val) {
			result = append(result, val)
		}
	}
	return
}

func RemoveMatches(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			return err
		}
	}
	return nil
}

var startupFolderPath string

func StartupFolder() string {
	if startupFolderPath == "" {
		var err error
		startupFolderPath, err = osext.ExecutableFolder()
		if err != nil {
			panic("Unexpected error while trying to get executable folder, err=" + err.Error())
		}
	}
	return startupFolderPath
}

// MinMaxInt returns the lowest and the highest values in an int slice
func MinMaxInt(values ...int) (min int, max int) {
	sliceLen := len(values)
	if sliceLen == 0 {
		return 0, 0
	}

	min = values[0]
	max = values[0]
	for i := 0; i < sliceLen; i++ {
		if values[i] < min {
			min = values[i]
		}
		if values[i] > max {
			max = values[i]
		}
	}

	return
}

// MinInt returns the lowest value in an int slice
func MinInt(values ...int) (min int) {
	min, _ = MinMaxInt(values...)
	return
}

// MaxInt returns the highest value in an int slice
func MaxInt(values ...int) (max int) {
	_, max = MinMaxInt(values...)
	return
}

// DeepCopy creates a copy of a struct - as opposed to shallow copy that will copy pointers instead of actual values
func DeepCopy(dst interface{}, src interface{}) (err error) {
	var bytes []byte

	// Serialize
	if src == nil {
		return fmt.Errorf("source struct is nil")
	}
	bytes, err = json.Marshal(src)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to marshal src: %v", err), 0)
	}

	// Deserialize
	if dst == nil {
		return fmt.Errorf("destination struct is nil")
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to marshal dst: %v", err), 0)
	}
	return
}

func RoundSeconds(dur time.Duration) time.Duration {
	return time.Duration(time.Second * time.Duration(int(time.Duration(dur).Seconds())))
}

func RoundMinutes(dur time.Duration) time.Duration {
	return time.Duration(time.Minute * time.Duration(int(time.Duration(dur).Minutes())))
}

func RoundHours(dur time.Duration) time.Duration {
	return time.Duration(time.Hour * time.Duration(int(time.Duration(dur).Hours())))
}

func HandleSignal(sig os.Signal, handle func(), report func(msg string)) {
	if report == nil {
		report = func(msg string) {
			println(msg)
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, sig)

	go func() {
		for range signalChan {
			report(fmt.Sprintf("Received interrupt %s, stack traces:", sig))
			ReportFullStack(report)

			if handle != nil {
				report("Handling signal...")
				handle()
				report("Done.")
			}
		}
	}()
}

func ReportFullStack(report func(msg string)) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	for _, line := range strings.Split(string(buf), "\n") {
		report(line)
	}
}
