package middleware

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "net/http/httputil"
    "os"
    "runtime"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/silenceper/log"
)

var (
    dunno     = []byte("???")
    centerDot = []byte("·")
    dot       = []byte(".")
    slash     = []byte("/")
)

func stack(skip int) []byte {
    buf := new(bytes.Buffer) // the returned data
    // As we loop, we open files and read them. These variables record the currently
    // loaded file.
    var lines [][]byte
    var lastFile string
    for i := skip; ; i++ { // Skip the expected number of frames
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        // Print this much at least.  If we can't find the source, it won't show.
        fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
        if file != lastFile {
            data, err := ioutil.ReadFile(file)
            if err != nil {
                continue
            }
            lines = bytes.Split(data, []byte{'\n'})
            lastFile = file
        }
        fmt.Fprintf(buf, "    %s: %s\n", function(pc), source(lines, line))
    }
    return buf.Bytes()
}

func function(pc uintptr) []byte {
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return dunno
    }
    name := []byte(fn.Name())
    if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
        name = name[lastSlash+1:]
    }
    if period := bytes.Index(name, dot); period >= 0 {
        name = name[period+1:]
    }
    name = bytes.Replace(name, centerDot, dot, -1)
    return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
    n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
    if n < 0 || n >= len(lines) {
        return dunno
    }
    return bytes.TrimSpace(lines[n])
}

// ErrorTrace 异常错误捕获
func ErrorTrace() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                var brokenPipe bool
                if ne, ok := err.(*net.OpError); ok {
                    if se, ok := ne.Err.(*os.SyscallError); ok {
                        if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
                            brokenPipe = true
                        }
                    }
                }
                stack := stack(3)
                httpRequest, _ := httputil.DumpRequest(c.Request, false)
                headers := strings.Split(string(httpRequest), "\r\n")
                for idx, header := range headers {
                    current := strings.Split(header, ":")
                    if current[0] == "Authorization" {
                        headers[idx] = current[0] + ": *"
                    }
                }

                log.Infof("err:%s, stack:%s, brokenPipe:%s", err, string(stack), brokenPipe)
                switch err.(type) {
                case error:
                    c.JSON(http.StatusInternalServerError, gin.H{
                        "msg":   "Sorry，服务器累瘫了",
                        "debug": err.(error).Error(),
                    })
                default:
                    c.JSON(http.StatusInternalServerError, gin.H{
                        "msg":   "Sorry，服务器累瘫了",
                        "debug": err,
                    })
                }
            }
        }()
        c.Next()
    }
}
