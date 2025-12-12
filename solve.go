package solver

import (
	"errors"
	"fmt"
	"iter"
	"math"
	"slices"

	"golang.org/x/exp/constraints"
)

/*
type Segment[T constraints.Integer] struct {
	closedStart, openEnd T
	value T
}

type Segments[T constraints.Integer] struct {
	parts []Segment
}

func (s *Segments[T]) Split(at T, value T) {
	for _
}

func NewSegments() *Segments {
	return &Segments{parts: make([]Segment, 0, 10)}
}
*/

type Vector[T constraints.Integer] []T

func (v Vector[T]) IsZero() bool {
	for _, x := range v {
		if x != 0 {
			return false
		}
	}
	return true
}

// DivideBy(x) returns the number N such than
// N * v.DivideBy(x) has the minimum possible Size
func (v Vector[T]) DivideBy(x Vector[T]) T {
	var result T
	var count T
	for i, n := range v {
		if x[i] != 0 {
			result += n / x[i]
			count++
		}
	}
	if count > 0 {
		result /= count
	}
	return result
}

func (v *Vector[T]) AddTo(other Vector[T]) {
	me := *v
	for i, n := range other {
		me[i] = me[i] + n
	}
}

func (v *Vector[T]) Scale(x T) Vector[T] {
	me := *v
	result := make(Vector[T], len(me))
	for i, n := range me {
		result[i] = n * x
	}
	return result
}

func (v Vector[T]) Add(other Vector[T]) Vector[T] {
	result := make(Vector[T], len(v))
	for i, n := range other {
		result[i] = v[i] + n
	}
	return result
}

func (v *Vector[T]) SubtractFrom(other Vector[T]) {
	me := *v
	for i, n := range other {
		me[i] -= n
	}
}

func (v Vector[T]) Subtract(other Vector[T]) Vector[T] {
	result := make(Vector[T], len(v))
	for i, n := range other {
		result[i] = v[i] - n
	}
	return result
}

// Magnitude returns the cartesian magnitude of this vector
func (v Vector[T]) Magnitude() float64 {
	var result float64
	for _, n := range v {
		result += float64(n) * float64(n)
	}
	return math.Sqrt(result)
}

// Size returns the sum of the absolute value of the dimensions
func (v Vector[T]) Size() int64 {
	var result int64
	for _, n := range v {
		if n > 0 {
			result += int64(n)
		} else {
			result -= int64(n)
		}
	}
	return result
}

func VectorCmp[T constraints.Integer](from, to Vector[T]) int {
	f := from.Size()
	t := to.Size()
	if f < t {
		return -1
	}
	if t > f {
		return 1
	}
	return 0
}

type Input[T constraints.Integer] struct {
	Target        Vector[T]
	Ivs           []Vector[T]
	AllowNegative bool
}

func (i Input[T]) Validate() error {
	if i.Target == nil {
		return errors.New("target is not defined")
	}
	if i.Ivs == nil {
		return errors.New("independent vectors (Ivs) are not defined")
	}
	n := len(i.Target)
	if n == 0 {
		return errors.New("empty target")
	}
	var atLeastOneNonzeroVector bool
	for _, iv := range i.Ivs {
		if len(iv) != n {
			return errors.New("inconsistent vector size")
		}
		if iv.Magnitude() != 0 {
			atLeastOneNonzeroVector = true
		}
	}
	if !atLeastOneNonzeroVector {
		return errors.New("all vectors are zero")
	}
	return nil
}

var ErrNope = errors.New("no solution found")

func Solve[T constraints.Integer](i Input[T]) (Vector[T], error) {
	if err := i.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	sequence := make([]int, len(i.Ivs))
	for i := range len(i.Ivs) {
		sequence[i] = i
	}
	slices.SortFunc(sequence, func(a, b int) int {
		// use negative as largest size should come first
		return -VectorCmp(i.Ivs[a], i.Ivs[b])
	})
	solution := Vector[T](make([]T, len(i.Ivs)))
	err := solve(i.Ivs, i.Target, i.AllowNegative, sequence, solution)
	return solution, err
}

func solve[T constraints.Integer](ivs []Vector[T], target Vector[T], allowNegative bool, sequence []int, solution Vector[T]) error {
	if len(sequence) == 0 {
		fmt.Println("last", ivs, "target:", target, target.IsZero())
		if target.IsZero() {
			return nil
		}
		return ErrNope
	}
	this := ivs[sequence[0]]
	d := target.DivideBy(this)
	fmt.Println("sequence:", sequence, "target:", target, "this:", this, "d:", d)
	for step := range StepOut[T](allowNegative) {
		delta := target.Subtract(this.Scale(d - step))
		fmt.Println("d-step", d-step, "delta", delta)
		if err := solve(ivs, delta, allowNegative, sequence[1:], solution); err == nil {
			solution[sequence[0]] = d - step
			return nil
		}
	}
	return ErrNope
}

func StepOut[T constraints.Integer](allowNegative bool) iter.Seq[T] {
	var i T
	return func(yield func(T) bool) {
		for i < 2 {
			if !yield(i) {
				return
			}
			i++
			if allowNegative {
				if !yield(-i) {
					return
				}
			}
		}
	}
}
