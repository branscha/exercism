package secret

import (
    "slices"
)

const (
    wink uint = 1
    dblink = 2
    close = 4
    jump = 8
    reverse = 16
)

func Handshake(code uint) []string {
	hs := []string{}

    if code & wink > 0 {
        hs = append(hs, "wink")
    }

    if code & dblink > 0 {
        hs = append(hs, "double blink")
    }

    if code & close > 0 {
        hs = append(hs, "close your eyes")
    }

    if code & jump > 0 {
        hs = append(hs, "jump")
    }

    if code & reverse > 0 {
        slices.Reverse(hs)
    }

    return hs
}
