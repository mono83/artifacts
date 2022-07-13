package main

import (
	"github.com/mono83/artifacts/cmd"
	"github.com/mono83/xray/std/xcobra"
)

func main() {
	xcobra.Start(cmd.MainCmd)
}
