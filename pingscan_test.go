package main

import (
	"testing"
)

func TestPing(t *testing.T) {
    timeout := 5

	r := *ping(&timeout, &[]string{"127.0.0.1"})
    if r[0].RTT() <= 0 || r[0].Err != nil {
        t.Error("Ping localhost failed")
    }

    r = *ping(&timeout, &[]string{"google.com"})
    if r[0].RTT() <= 0 || r[0].Err != nil {
        t.Error("Ping google.com failed")
    }

    r = *ping(&timeout, &[]string{"yahoo.com"})
    if r[0].RTT() <= 0 || r[0].Err != nil {
        t.Error("Ping yahoo.com failed")
    }

    r = *ping(&timeout, &[]string{"owjdfiojsfdjfsoijeifojweoifjdsiojciosdjc.czs"})
    if r[0].Err == nil {
        t.Error("Ping nonsense not failed")
    }

    r = *ping(&timeout, &[]string{"266.266.266.266"})
    if r[0].Err == nil {
        t.Error("Ping nonsense not failed")
    }
}
