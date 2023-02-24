package bimg

import (
	"testing"
)

func TestValidAnimatedGifDelay(t *testing.T) {
	delay := 655350
	t.Logf("delay=%d", delay)
	if ValidAnimatedGifDelay(delay) {
		t.Logf("valid")
	} else {
		t.Fatalf("invalid")
	}
}
