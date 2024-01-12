package appmanager

import (
    "testing"
)

func TestHello(t *testing.T) {
    result := 1
    if result != 1 {
        t.Fatalf("result not as expected")
    }
}


