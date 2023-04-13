// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"github.com/ChainSafe/chainbridge-core/example/cmd"
)

func main() {
	// file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// // Set log output to the file
	// log.SetOutput(file)

	// // Log a message with a formatted string
	// log.Printf("This is a log message with a formatted string: %s", "Hello, World!")

	cmd.Execute()
}
