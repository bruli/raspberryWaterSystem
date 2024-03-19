//go:build functional

package functional

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"
)

//func TestFunctional(t *testing.T) {
//	runHomepage(t)
//	runZones(t)
//	runStatus(t)
//	runWeather(t)
//	runPrograms(t)
//	runExecutionLogs(t)
//	runPkg(t)
//}

const serverURL = "http://localhost:8083"

var (
	ctx context.Context
	cl  http.Client
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	cl = http.Client{Timeout: 3 * time.Second}

	code := m.Run()
	os.Exit(code)
}
