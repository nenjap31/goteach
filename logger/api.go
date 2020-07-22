package logger

import (
	"encoding/json"
	"goteach/utils/httpdump"
	"path"
	"strings"
	"time"
	"runtime"
	"path/filepath"
	"strconv"
	"os"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

// Constant variable in API log. 2nd level logging we put into Event.
// Example : log.Level_1().Msg(Level2)

// APILogHandler : handle something who need to do
func APILogHandler(c echo.Context, req, res []byte) {
	c.Response().Header().Set("X-goteach-ResponseTime", time.Now().Format(time.RFC3339))
	reqTime, err := time.Parse(time.RFC3339, c.Request().Header.Get("X-goteach-RequestTime"))
	var elapstime time.Duration
	if err == nil {
		elapstime = time.Since(reqTime)
	}

	var handler string
	r := c.Echo().Routes()
	cpath := strings.Replace(c.Path(), "/", "", -1)
	for _, v := range r {
		vpath := strings.Replace(v.Path, "/", "", -1)
		if vpath == cpath && v.Method == c.Request().Method {
			handler = v.Name
			// Handler for wrong route.
			if strings.Contains(handler, "func1") {
				handler = "UndefinedRoute"
			}
			break
		}
	}

	// Get Handler Name
	dir, file := path.Split(handler)
	fileStrings := strings.Split(file, ".")
	packHandler := dir + fileStrings[0]
	funcHandler := strings.Replace(handler, packHandler+".", "", -1)

	respHeader, _ := json.Marshal(c.Response().Header())
	reqHeader := httpdump.DumpRequest(c.Request())

	Info().
		Str("Identifier", viper.GetString("log_identifier")+"_http").
		Str("package", packHandler).
		Int64("elapsed_time", elapstime.Nanoseconds()/int64(time.Millisecond)).
		Str("handler", funcHandler).
		Str("ip", c.RealIP()).
		Str("host", c.Request().Host).
		Str("method", c.Request().Method).
		Str("url", c.Request().RequestURI).
		Str("request_time", c.Request().Header.Get("X-goteach-RequestTime")).
		Str("request_header", reqHeader).
		Str("request", string(req)).
		Int("httpcode", c.Response().Status).
		Str("response_time", c.Response().Header().Get("X-goteach-ResponseTime")).
		Str("response_header", string(respHeader)).
		Str("response", string(res)).
		Msg("")
}

// APILogSkipper : rules for APILogHandler
func APILogSkipper(c echo.Context) bool {
	// bool, is this url request include "/api"?
	rules1 := strings.Contains(c.Request().RequestURI, "/api")

	// bool, is this request using method "GET"?
	rules2 := c.Request().Method != "GET"

	// bool, is this url request include "/login"?
	rules3 := strings.Contains(c.Request().RequestURI, "/login")

	if rules1 {
		return false
	}

	if rules2 {
		if !rules3 {
			return false
		}
	}

	return true
}


// LogHandler : handle something who need to do
func InfoLogHandler(msg string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "?"
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	name, _ := os.Hostname()

	Info().
		Str("Identifier", viper.GetString("log_identifier")+"_info").
		Str("file", file).
		Str("handler", fnName).
		Str("line", strconv.Itoa(line)).
		Str("host", name).
		Int("port", viper.GetInt("grpc.port")).
		Msg(msg)
}
