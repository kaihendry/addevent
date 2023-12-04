package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"log/slog"

	"github.com/apex/gateway/v2"
	ics "github.com/arran4/golang-ical"
)

var GoVersion = runtime.Version()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", addevent)
	var err error
	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
		err = gateway.ListenAndServe("", mux)
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
		slog.Info("local development", "port", os.Getenv("PORT"))
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), mux)
	}
	slog.Error("error listening", err)
}

func addevent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Version", fmt.Sprintf("%s %s", os.Getenv("version"), GoVersion))
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetUrl("https://addevent.dabase.com/")
	event := cal.AddEvent(fmt.Sprintf("hendry+%s@iki.fi", "a test"))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(time.Now())
	event.SetEndAt(time.Now().Add(1 * time.Hour))
	event.SetSummary(fmt.Sprintf("Just testing %s", time.Now().Format("Monday")))

	// get the UA string from the request headers
	ua := r.Header.Get("User-Agent")

	event.SetLocation("A wood " + ua)

	w.Write([]byte(cal.Serialize()))
}
