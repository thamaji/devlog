package devlog

import (
	"fmt"
	"go/build"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/fatih/color"
)

var Enable = true

var Writer io.Writer = os.Stderr

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
	filename := runtime.FuncForPC(pc).Name()

	color.New(color.FgHiBlack).Fprintln(Writer, append([]interface{}{"[DEVLOG]", time.Now().Format(time.RFC3339), file + ":" + strconv.Itoa(line), path.Base(filename)}, a...)...)
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
	filename := runtime.FuncForPC(pc).Name()
	color.New(color.FgHiBlack).Fprintln(Writer, "[DEVLOG]", time.Now().Format(time.RFC3339), file+":"+strconv.Itoa(line), path.Base(filename), fmt.Sprintf(f, a...))
}
