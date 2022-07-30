//go:build functional
// +build functional

package functional_test

import "testing"

func TestFunctional(t *testing.T) {
	runHomepage(t)
	runZones(t)
}
