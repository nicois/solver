package solver_test

import (
	"testing"

	"github.com/nicois/solver"
	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	input := solver.Input[int]{
		Target:        solver.Vector[int]{3, 2, 5},
		Ivs:           []solver.Vector[int]{{1, 0, 0}, {0, 2, 0}, {0, 0, 1}},
		AllowNegative: false,
	}

	result, err := solver.Solve(input)
	require.NoError(t, err)
	require.Equal(t, solver.Vector[int]{3, 1, 5}, result)
}

func TestVectorAddSubtract(t *testing.T) {
	a := solver.Vector[int]{1, 2, 3}
	b := solver.Vector[int]{1, 0, -5}

	a.AddTo(b)
	require.Equal(t, solver.Vector[int]{2, 2, -2}, a)
	require.Equal(t, solver.Vector[int]{1, 0, -5}, b)
	a.SubtractFrom(b)
	require.Equal(t, solver.Vector[int]{1, 2, 3}, a)
	require.Equal(t, solver.Vector[int]{1, 0, -5}, b)
}

func TestVectorMagnitude(t *testing.T) {
	require.Equal(t, float64(3.7416573867739413), solver.Vector[int]{1, 2, 3}.Magnitude())
	require.Equal(t, float64(2.23606797749979), solver.Vector[int]{1, 2, 0}.Magnitude())
	require.Equal(t, float64(0), solver.Vector[int]{0, 0, 0}.Magnitude())
	require.Equal(t, float64(3.7416573867739413), solver.Vector[int]{1, 2, -3}.Magnitude())
}

func TestVectorSize(t *testing.T) {
	require.Equal(t, int64(6), solver.Vector[int]{1, 2, 3}.Size())
	require.Equal(t, int64(3), solver.Vector[int]{1, 2, 0}.Size())
	require.Equal(t, int64(0), solver.Vector[int]{0, 0, 0}.Size())
	require.Equal(t, int64(6), solver.Vector[int]{1, 2, -3}.Size())
}

func TestVectorDiv(t *testing.T) {
	require.Equal(t, 5, solver.Vector[int]{5, 0}.DivideBy(solver.Vector[int]{1, 0}))
	require.Equal(t, 5, solver.Vector[int]{5, -6}.DivideBy(solver.Vector[int]{1, 0}))
	require.Equal(t, 0, solver.Vector[int]{5, -6}.DivideBy(solver.Vector[int]{1, 1}))
	require.Equal(t, 5, solver.Vector[int]{5, -6}.DivideBy(solver.Vector[int]{1, -1}))
	require.Equal(t, -5, solver.Vector[int]{5, -6}.DivideBy(solver.Vector[int]{-1, 1}))

	a := solver.Vector[int]{1, 2, 3}
	b := solver.Vector[int]{10, 25, 50}
	require.Equal(t, 12, b.DivideBy(a))
}
