package counting

import (
  "testing"
)

var fibTests = []struct {
  n int
  expected int
} {
  {1,1},
  {2,1},
  {3,2},
  {4,3},
  {5,5},
  {6,8},
  {7,13},
}

func TestCount(t *testing.T) {
  if Count() != 1 {
    t.Error("Expected 1")
  }
  if Count() != 2 {
    t.Error("Expected 2")
  }
  if Count() != 3 {
    t.Error("Expected 3")
  }
}

func TestFib(t *testing.T) {
  for _, tt := range fibTests {
    actual := Fib(tt.n)
    if actual != tt.expected {
      t.Errorf("Fib(%d): expected %d, actual %d", tt.n, tt.expected, actual)
    }
  }
}
