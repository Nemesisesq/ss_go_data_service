package streamsavvy

import (
	"testing"
	"os"
)

func TestSearch(t *testing.T) {
	os.Setenv("SS_DJANGO_DATA_SERVICE", "http://localhost:8001")
	res := Search("Orange")
	//
	if res == nil {
		t.Error("Expected Results but 0 were found")
	}
}
