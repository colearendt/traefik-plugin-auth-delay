package traefik_plugin_auth_delay

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// AuthDelay holds one auth-delay configuration section
type AuthDelay struct {
	MinCode  int    `json:"min-code,omitempty"`
	MaxCode  int    `json:"max-code,omitempty"`
	MinDelay string `json:"min-delay,omitempty"`
	MaxDelay string `json:"max-delay,omitempty"`
}

type authDelayInternal struct {
	minCode  int
	maxCode  int
	minDelay time.Duration
	maxDelay time.Duration
}

// Config holds the plugin configuration
type Config struct {
	AuthDelays []AuthDelay `json:"auth-delay,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type addAuthDelay struct {
	name       string
	next       http.Handler
	authDelays []authDelayInternal
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	authDelays := make([]authDelayInternal, len(config.AuthDelays))

	for i, authDelayConfig := range config.AuthDelays {
		// convert timestamp string into time.Duration
		minDelay, err := time.ParseDuration(authDelayConfig.MinDelay)
		if err != nil {
			return nil, fmt.Errorf("error parsing minDelay duration %q: %w\n", authDelayConfig.MinDelay, err)
		}
		maxDelay, err := time.ParseDuration(authDelayConfig.MaxDelay)
		if err != nil {
			return nil, fmt.Errorf("error parsing maxDelay duration %q: %w\n", authDelayConfig.MinDelay, err)
		}

		if minDelay > maxDelay {
			return nil, fmt.Errorf("error parsing configuration. min-delay greater than max-delay")
		}
		if minDelay < 0 {
			return nil, fmt.Errorf("error parsing configuration. min-delay is a negative duration")
		}

		authDelays[i] = authDelayInternal{
			minCode:  authDelayConfig.MinCode,
			maxCode:  authDelayConfig.MaxCode,
			minDelay: minDelay,
			maxDelay: maxDelay,
		}
	}

	return &addAuthDelay{
		name:       name,
		next:       next,
		authDelays: authDelays,
	}, nil
}

func (r *addAuthDelay) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedWriter := &responseWriter{
		writer:     rw,
		authDelays: r.authDelays,
	}
	r.next.ServeHTTP(wrappedWriter, req)
}

type responseWriter struct {
	writer     http.ResponseWriter
	authDelays []authDelayInternal
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	for _, authDelay := range r.authDelays {
		if statusCode >= authDelay.minCode && statusCode <= authDelay.maxCode {
			// TODO: add random delay in range
			minNSec := authDelay.minDelay.Nanoseconds()
			maxNSec := authDelay.maxDelay.Nanoseconds()

			rand.Seed(time.Now().UnixNano())
			randDelayNSec := rand.Int63n(maxNSec-minNSec+1) + minNSec
			// log.Printf("Max nanosec: %v\n", maxNSec)
			// log.Printf("Min nanosec: %v\n", minNSec)
			// log.Printf("Rand nanosec: %v\n", randDelayNSec)
			randDelay := time.Duration(randDelayNSec)
			log.Printf("Adding %v delay for status: %v\n", randDelay.String(), statusCode)
			time.Sleep(randDelay)
		}
	}
	r.writer.WriteHeader(statusCode)
}
