// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"log"
	"os"

	"github.com/ChainSafe/chainbridge-core/example"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "chainbridge",
		Usage:  "refactoring research",
		Action: example.Run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
