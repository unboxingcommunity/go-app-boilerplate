package rating

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// log.Println("Do stuff BEFORE the tests!")
	// setup() if any initializations are required , make them here
	t := m.Run()
	// log.Println("Do stuff AFTER the tests!")
	os.Exit(t)
}
