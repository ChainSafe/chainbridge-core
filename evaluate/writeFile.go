package evaluate

import (
	"log"
	"os"
)

func WriteToFile(format string, v ...interface{}) {
	file, err := os.OpenFile("log_new.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)

	// // Remove the default date and time from the go log output
	// log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Log a message with a formatted string
	log.Printf(format, v...)

}

func WriteToFileEvaluate(format string, v ...interface{}) {
	file, err := os.OpenFile("log_evaluate.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)

	// Remove the default date and time from the go log output
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Log a message with a formatted string

	log.Printf(format, v...)

}

func WriteToFileEvaluateWithCustomPath(format string, v ...interface{}) {
	file, err := os.OpenFile("./../example/log_evaluate.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)

	// Remove the default date and time from the go log output
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Log a message with a formatted string

	log.Printf(format, v...)

}

func WriteToFileWithCustomPath(path string, format string, v ...interface{}) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)

	// // Remove the default date and time from the go log output
	// log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Log a message with a formatted string
	log.Printf(format, v...)

}

func WriteBreakToFile() {
	file, err := os.OpenFile("log_new.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)

	// // Remove the default date and time from the go log output
	// log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Log a message with a formatted string
	log.Println()
	log.Println("============================================================================")

}

func WriteBreakToFileEvaluate() {
	file, err := os.OpenFile("log_evaluate.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)

	// Remove the default date and time from the go log output
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Log a message with a formatted string

	log.Println()
	log.Println()
	// log.Println("============================================================================")

}

// func WriteTimeDurationToFile(format string, elapsedTime time.Duration) {
// 	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	// Set log output to the file
// 	log.SetOutput(file)

// 	// Log a message with a formatted string
// 	log.Printf("This is a log message with a formatted string: %s", "Hello, World!")

// 	log.Printf(format, elapsedTime)

// }

// func WriteTimeToFile(format string, elapsedTime time.Time) {
// 	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	// Set log output to the file
// 	log.SetOutput(file)

// 	// Log a message with a formatted string
// 	log.Printf("This is a log message with a formatted string: %s", "Hello, World!")

// 	log.Printf(format, elapsedTime)

// }

// Update

// type logRequest struct {
// 	format string
// 	v      []interface{}
// }

// var logChannel = make(chan logRequest)

// func init() {
// 	go logWorker()
// }

// func logWorker() {
// 	// Open the log file
// 	file, err := os.OpenFile("log_new_updated.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	// Set log output to the file
// 	log.SetOutput(file)

// 	// Process log requests as they arrive
// 	for request := range logChannel {
// 		log.Printf(request.format, request.v...)
// 	}
// }

// func WriteToFileNew(format string, v ...interface{}) {
// 	// Add the log request to the queue
// 	logChannel <- logRequest{format, v}
// }
