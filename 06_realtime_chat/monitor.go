package main

import (
	"os"
)

func runMonitor(brok broker) error {
	return handleSub(brok, os.Stdout)
}
