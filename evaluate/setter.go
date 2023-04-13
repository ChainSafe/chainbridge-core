package evaluate

import "time"

func SetT0(t time.Time) {
	T0 = t
	WriteBreakToFile()
	
	WriteToFile("T0 (Trigger deposit): %s", T0)
	// WriteToFileEvaluate("Debug T0: %s", T0)
}

func SetT0a(t time.Time) {
	T0a = t

	WriteToFile("T0a (Finish deposit): %s", t)

	// WriteToFileEvaluate("Debug T0a: %s", T0a)

	// if !T0.IsZero() {
	// 	WriteBreakToFileEvaluate()
	// 	WriteToFileEvaluate("Step 0 (finish deposit): %s", T0a.Sub(T0))
	// }

}

func SetT1(t time.Time) {
	T1 = t
	WriteToFile("T1 (Relayer caught deposit event): %s", t)

	WriteBreakToFileEvaluate()
	WriteToFileEvaluate("Debug T0: %s", T0)
	WriteToFileEvaluate("Debug T0a: %s", T0a)
	WriteToFileEvaluate("Debug T1: %s", T1)

	if !T0a.IsZero() {
		WriteBreakToFileEvaluate()
		WriteToFileEvaluate("Step 0 (finish deposit): %s", T1.Sub(T0a))
	}
}

func SetT2(t time.Time) {
	T2 = t
	WriteToFile("T2 (Trigger/Start vote): %s", t)

	WriteToFileEvaluate("Debug T2: %s", T2)

}

func SetT2a(t time.Time) {
	T2 = t
	WriteToFile("T2 (Trigger/Start vote): %s", t)

	WriteToFileEvaluate("Debug T2a: %s", T2a)

}

// func SetT2a(t time.Time) {
// 	T2a = t
// 	WriteToFile("T2a (Finish vote): %s", t)
// }

// func SetTimeT2a(t time.Time) {
// 	T2a = t
// 	WriteToFile("This is the last vote -> executing ... (T2a): %s", T2)
// }

func SetT3(t time.Time) {
	T3 = t
	WriteToFileWithCustomPath("./../example/log_new.txt", "T3 (Finish execute - Executed): %s", T3)

	WriteToFileEvaluateWithCustomPath("Debug T3: %s", T3)

	// if !T1.IsZero() && IsMet {
	// 	WriteToFileEvaluateWithCustomPath("Step 3 (first vote until executed successfully - threshold met): %s", T3.Sub(T1))
	// }
}

func SetT4(t time.Time) {
	T4 = t
	WriteToFile("Time finish deposit (T4): %s", T4)
}

func SetIsMet(value bool) {
	IsMet = value
	// WriteToFile("Time finish deposit (T4): %s", T4)
}
