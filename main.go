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
	// redirect to addevent
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/addevent", http.StatusTemporaryRedirect)
	})

	mux.HandleFunc("/addevent", addevent)
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

	now := time.Now()
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetUrl("https://addevent.dabase.com/")
	event := cal.AddEvent(fmt.Sprintf("hendry+%s@iki.fi", "addevent"))
	event.SetCreatedTime(now)
	event.SetDtStampTime(now)
	event.SetModifiedAt(now)
	event.SetStartAt(now)
	event.SetEndAt(now.Add(1 * time.Hour))
	event.SetSummary(fmt.Sprintf("Last fetched %s, hour %s", now.Format("Monday"), now.Format("15:04")))

	event.SetDescription(r.Header.Get("User-Agent"))
	event.SetURL("https://github.com/kaihendry/addevent")
	event.SetLocation("A wood")

	w.Write([]byte(cal.Serialize()))
}
