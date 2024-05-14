package assert

import (
	"log"
	"math"
)

func Eq[T comparable](lhs, rhs T) {
	if lhs != rhs {
		log.Fatalf("Assertion failed:\nLeft != Right\nlhs = %v\nrhs = %v", lhs, rhs)
	}
}

func Ne[T comparable](lhs, rhs T) {
	if lhs == rhs {
		log.Fatalf("Assertion failed:\nLeft == Right\nlhs = %v\nrhs = %v", lhs, rhs)
	}
}

func Ok(err error) {
	if err != nil {
		log.Fatalf("Assertion failed:\nError = %v", err)
	}
}

func ApproxEq(lhs, rhs float64) {
	if math.Abs(lhs-rhs) > 0.000_001 {
		log.Fatalf("Assertion failed:\nLeft !~= Right\nLeft = %f\nRight = %f", lhs, rhs)
	}
}
