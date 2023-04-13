package evaluate

import "time"

var T0 time.Time  // Trigger deposit
var T0a time.Time // Finish deposit
var T1 time.Time  // Relayer catched deposit event
var T2 time.Time  // Trigger vote
var T2a time.Time // Finish vote - threshold met
var T3 time.Time  // Finish Execute
var T4 time.Time  // Received func

var IsMet = false

var currDepositHash string
