package main

import (
	"flag"
	. "github.com/rande/goapp"
	"github.com/rande/gonode/explorer/helper"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
)

func init() {
	bind.WithFlag()
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
	//  graceful.DoubleKickWindow(2 * time.Second)
}

func main() {
	app := NewApp()

	helper.BuildApp(app, "./config.toml")

	mux := app.Get("goji.mux").(*web.Mux)

	if !flag.Parsed() {
		flag.Parse()
	}

	mux.Compile()

	// Install our handler at the root of the standard net/http default mux.
	// This allows packages like expvar to continue working as expected.
	http.Handle("/", mux)

	listener := bind.Default()
	log.Println("Starting Goji on", listener.Addr())

	graceful.HandleSignals()
	bind.Ready()

	graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
	graceful.PostHook(func() { log.Printf("Goji stopped") })

	err := graceful.Serve(listener, http.DefaultServeMux)

	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}
