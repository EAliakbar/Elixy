package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	_ "github.com/motemen/go-loghttp/global"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	upstream := flag.String("upstream", "", "upstream to proxy requests")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// Using Humanize log format
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *upstream == "" {
		log.Fatal().Msg("an Upstream should be defined using \"--upstream\"")
	}

	remote_upstream, err := url.Parse(*upstream)
	if err != nil {
		log.Fatal().Msgf("Error in Parsing upstream: %s", err.Error())
	}
	if remote_upstream.Scheme != "http" && remote_upstream.Scheme != "https" {
		log.Fatal().Msgf("Elixy only Supports \"http\" and \"https\" upstreams")
	}
	log.Info().Msgf("Set Upstream to %s", remote_upstream.Hostname())

	handler := httputil.NewSingleHostReverseProxy(remote_upstream)

	srv := &http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Msgf("Starting to Listen on: %s", "0.0.0.0:8000")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

}
