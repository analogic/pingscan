package main

import (
	"testing"
)

func TestPing(t *testing.T) {
	_, err := ping("127.0.0.1", 1)
	if err != nil {
		t.Error("Ping localhost failed")
	}

	_, err = ping("google.com", 5)
	if err != nil {
		t.Error("Ping google.com failed")
	}

	_, err = ping("266.266.266.266", 1)
	if err == nil {
		t.Error("Ping to 266.266.266.266 not failed")
	}

	_, err = ping("owjdfiojsfdjfsoijeifojweoifjdsiojciosdjc.cz", 1)
	if err == nil {
		t.Error("Ping to owjdfiojsfdjfsoijeifojweoifjdsiojciosdjc.cz not failed")
	}
}
