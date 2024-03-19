//go:build functional

package functional

import (
	"testing"
)

func TestFunctional(t *testing.T) {
	runHomepage(t)
	runZones(t)
	runStatus(t)
	runWeather(t)
	runPrograms(t)
	runExecutionLogs(t)
	runPkg(t)
}
