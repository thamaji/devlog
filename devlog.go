package devlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var Enable = true

var Writer io.Writer = os.Stderr

var TimeFormat = "2006-01-02T15:04:05.000"

func Warn(a ...interface{}) {
	if !Enable {
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		relpath, err := filepath.Rel(gopath, file)
		if err == nil {
			file = relpath
			break
		}
	}
	funcname := runtime.FuncForPC(pc).Name()

	color.New(color.FgHiYellow).Fprintln(Writer, append([]interface{}{"[DEVLOG]", time.Now().Format(TimeFormat), file + ":" + strconv.Itoa(line), path.Base(funcname)}, a...)...)
}

func Warnf(f string, a ...interface{}) {
	if !Enable {
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		relpath, err := filepath.Rel(gopath, file)
		if err == nil {
			file = relpath
			break
		}
	}
	funcname := runtime.FuncForPC(pc).Name()

	color.New(color.FgHiYellow).Fprintln(Writer, "[DEVLOG]", time.Now().Format(TimeFormat), file+":"+strconv.Itoa(line), path.Base(funcname), fmt.Sprintf(f, a...))
}

func Error(a ...interface{}) {
	if !Enable {
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		relpath, err := filepath.Rel(gopath, file)
		if err == nil {
			file = relpath
			break
		}
	}
	funcname := runtime.FuncForPC(pc).Name()

	color.New(color.FgHiRed).Fprintln(Writer, append([]interface{}{"[DEVLOG]", time.Now().Format(TimeFormat), file + ":" + strconv.Itoa(line), path.Base(funcname)}, a...)...)
}

func Errorf(f string, a ...interface{}) {
	if !Enable {
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		relpath, err := filepath.Rel(gopath, file)
		if err == nil {
			file = relpath
			break
		}
	}
	funcname := runtime.FuncForPC(pc).Name()

	color.New(color.FgHiRed).Fprintln(Writer, "[DEVLOG]", time.Now().Format(TimeFormat), file+":"+strconv.Itoa(line), path.Base(funcname), fmt.Sprintf(f, a...))
}

func Log(a ...interface{}) {
	if !Enable {
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		relpath, err := filepath.Rel(gopath, file)
		if err == nil {
			file = relpath
			break
		}
	}
	funcname := runtime.FuncForPC(pc).Name()

	color.New(color.FgHiBlue).Fprintln(Writer, append([]interface{}{"[DEVLOG]", time.Now().Format(TimeFormat), file + ":" + strconv.Itoa(line), path.Base(funcname)}, a...)...)
}

func Logf(f string, a ...interface{}) {
	if !Enable {
		return
	}

	pc, file, line, _ := runtime.Caller(1)
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		relpath, err := filepath.Rel(gopath, file)
		if err == nil {
			file = relpath
			break
		}
	}
	funcname := runtime.FuncForPC(pc).Name()

	color.New(color.FgHiBlue).Fprintln(Writer, "[DEVLOG]", time.Now().Format(TimeFormat), file+":"+strconv.Itoa(line), path.Base(funcname), fmt.Sprintf(f, a...))
}

func Password(v string) string {
	n := len(v) / 2
	return v[:n] + strings.Repeat("*", len(v)-n)
}

func JSON(v interface{}) string {
	if !Enable {
		return ""
	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return "\n" + string(bytes) + "\n"
}

var TableSeparator = " | "

func Table(v interface{}) string {
	if !Enable {
		return ""
	}

	tbl := table(reflect.ValueOf(v))

	maxs := []int{}
	for _, row := range tbl {
		if len(maxs) < len(row) {
			maxs = append(maxs, make([]int, len(row)-len(maxs))...)
		}

		for i, v := range row {
			if len(v) > maxs[i] {
				maxs[i] = len(v)
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})
	fmt.Fprintln(buf)

	for _, row := range tbl {
		if len(row) < len(maxs) {
			row = append(row, make([]string, len(maxs)-len(row))...)
		}

		for i, v := range row {
			fmt.Fprint(buf, v+strings.Repeat(" ", maxs[i]-len(v)))
			if i < (len(row) - 1) {
				fmt.Fprint(buf, TableSeparator)
			}
		}

		fmt.Fprintln(buf)
	}

	return buf.String()
}

func table(v reflect.Value) [][]string {
	if !v.IsValid() {
		return [][]string{[]string{"nil"}}
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return [][]string{[]string{"nil"}}
		}
		return table(v.Elem())

	case reflect.Interface:
		return table(v.Elem())

	case reflect.Array, reflect.Slice:
		n := v.Len()

		tbl := [][]string{}
		for i := 0; i < n; i++ {
			key := strconv.Itoa(i)
			for _, row := range table(v.Index(i)) {
				tbl = append(tbl, append([]string{key}, row...))
				key = ""
			}
		}

		return tbl

	case reflect.Map:
		tbl := [][]string{}
		mapkeys := v.MapKeys()
		sort.Slice(mapkeys, func(i, j int) bool {
			return fmt.Sprint(mapkeys[i]) < fmt.Sprint(mapkeys[j])
		})
		for _, mapkey := range mapkeys {
			key := fmt.Sprint(mapkey)
			for _, row := range table(v.MapIndex(mapkey)) {
				tbl = append(tbl, append([]string{key}, row...))
				key = ""
			}
		}

		return tbl

	case reflect.Struct:
		typ := v.Type()
		n := typ.NumField()

		tbl := [][]string{}
		for i := 0; i < n; i++ {
			key := typ.Field(i).Name
			for _, row := range table(v.Field(i)) {
				tbl = append(tbl, append([]string{key}, row...))
				key = ""
			}
		}

		return tbl

	default:
		return [][]string{[]string{fmt.Sprint(v)}}
	}
}
