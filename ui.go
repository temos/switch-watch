package main

import "fmt"

const (
	EscMoveHome        = "\x1b[H"
	EscEraseRestOfLine = "\x1b[0K"
	EscEraseScreen     = "\x1b[2J"
)

func toReadable(speed uint) string {
	const (
		K = 1000
		M = K * 1000
		G = M * 1000
	)

	if speed >= G {
		return fmt.Sprintf("%.1f GBits/s", float32(speed)/G)
	}

	if speed >= M {
		return fmt.Sprintf("%.1f Mbits/s", float32(speed)/M)
	}

	if speed >= K {
		return fmt.Sprintf("%.1f Kbits/s", float32(speed)/K)
	}

	return fmt.Sprintf("%d bits/s", speed)
}
