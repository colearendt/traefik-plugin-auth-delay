package traefik_plugin_auth_delay

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeHTTP(t *testing.T) {
	t.Run("check that delay added", func(t *testing.T) {
		config := &Config{
			AuthDelays: []AuthDelay{
				genConfig(403, "1ms"),
				genConfig(401, "1ms"),
			},
		}

		reqDuration := testHttpRequest(t, config, http.StatusForbidden)
		targetDuration := time.Duration(1 * 1000 * 1000)
		if reqDuration < targetDuration {
			t.Errorf("Request was faster than expected (%v): %v\n", targetDuration.String(), reqDuration.String())
		}
		if reqDuration > 2*targetDuration {
			t.Errorf("Duration was additive when it should not be: %v\n", reqDuration.String())
		}
	})

	t.Run("check that delay not added", func(t *testing.T) {
		config := &Config{
			AuthDelays: []AuthDelay{
				genConfig(403, "1s"),
			},
		}

		reqDuration := testHttpRequest(t, config, http.StatusOK)
		targetDuration := time.Duration(1 * 1000 * 1000 * 1000)
		if reqDuration > targetDuration {
			t.Errorf("Request was delayed when it should not have been: %v\n", reqDuration.String())
		}
	})

	t.Run("range is used properly", func(t *testing.T) {
		config := &Config{
			AuthDelays: []AuthDelay{
				{
					MinCode:  400,
					MaxCode:  404,
					MinDelay: "5ms",
					MaxDelay: "10ms",
				},
			},
		}

		reqDuration := testHttpRequest(t, config, http.StatusForbidden)
		minDuration := time.Duration(5 * 1000 * 1000)
		if reqDuration <= minDuration {
			t.Errorf("Request duration did not exceed minimum duratino (5ms): %v\n", reqDuration.String())
		}
	})
}

func genConfig(code int, delay string) AuthDelay {
	return AuthDelay{
		MinCode:  code,
		MaxCode:  code,
		MinDelay: delay,
		MaxDelay: delay,
	}
}

func testHttpRequest(t *testing.T, config *Config, responseStatus int) time.Duration {
	next := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(responseStatus)
	}

	authDelay, err := New(context.Background(), http.HandlerFunc(next), config, "authDelay")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	start := time.Now()
	authDelay.ServeHTTP(recorder, req)
	end := time.Now()
	return end.Sub(start)
}
