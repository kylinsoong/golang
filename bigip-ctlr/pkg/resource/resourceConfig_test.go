package resource

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/labels"
)

func TestDefaultPartition(t *testing.T) {
	dgPath := strings.Join([]string{DEFAULT_PARTITION, "Shared"}, "/")
	want := "k8s/Shared"
	t.Logf("dgPath: %s", dgPath)
	if dgPath != want {
		t.Errorf("dgPath = %q, want %q", dgPath, want)
	}
}

func TestDefaultConfigMapLabel(t *testing.T) {
	label := DefaultConfigMapLabel
	var l labels.Selector
	var err error
	l, err = labels.Parse(label)
	if err != nil {
		t.Errorf("Failed to create label: %v", err)
	} else {
		t.Logf("labels.Selector: %v", l)
	}

}
