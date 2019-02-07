/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package transform

import (
	"math"
	"testing"

	"github.com/unidoc/unidoc/common"
)

func init() {
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
}

// TestAngle tests the Matrix.Angle() function.
func TestAngle(t *testing.T) {
	extraTests := []angleCase{}
	for theta := 0.01; theta <= 360.0; theta *= 1.1 {
		extraTests = append(extraTests, makeAngleCase(2.0, theta))
	}

	const angleTol = 1.0e-10

	for _, test := range append(angleTests, extraTests...) {
		p := test.params
		m := NewMatrix(p.a, p.b, p.c, p.d, p.tx, p.ty)
		theta := m.Angle()
		if math.Abs(theta-test.theta) > angleTol {
			t.Fatalf("Bad angle: m=%s expected=%g° actual=%g°", m, test.theta, theta)
		}
	}
}

type params struct{ a, b, c, d, tx, ty float64 }
type angleCase struct {
	params         // Affine transform.
	theta  float64 // Rotation of affine transform in degrees.
}

var angleTests = []angleCase{
	{params: params{1, 0, 0, 1, 0, 0}, theta: 0},
	{params: params{0, -1, 1, 0, 0, 0}, theta: 90},
	{params: params{-1, 0, 0, -1, 0, 0}, theta: 180},
	{params: params{0, 1, -1, 0, 0, 0}, theta: 270},
	{params: params{1, -1, 1, 1, 0, 0}, theta: 45},
	{params: params{-1, -1, 1, -1, 0, 0}, theta: 135},
	{params: params{-1, 1, -1, -1, 0, 0}, theta: 225},
	{params: params{1, 1, -1, 1, 0, 0}, theta: 315},
}

// makeAngleCase makes an angleCase for a Matrix with scale `r` and angle `theta` degrees.
func makeAngleCase(r, theta float64) angleCase {
	radians := theta / 180.0 * math.Pi
	a := r * math.Cos(radians)
	b := -r * math.Sin(radians)
	c := -b
	d := a
	return angleCase{params{a, b, c, d, 0, 0}, theta}
}

// TestInverse tests the Matrix.Inverse() function.
func TestInverse(t *testing.T) {
	m := NewMatrix(1, 1, 1, 1, 0, 0)
	_, hasInverse := m.Inverse()
	if hasInverse {
		t.Fatalf("%s has inverse", m)
	}

	testInverse(t, NewMatrix(1, 0, 0, 1, 0, 0))
	testInverse(t, NewMatrix(1, 0, 0, -1, 0, 0))
	testInverse(t, NewMatrix(0, -1, -1, 0, 0, 0))
	testInverse(t, NewMatrix(1, 0, 0, 1, 2, 5))
	testInverse(t, NewMatrix(1, 0, 0, 2, 2, 5))
	testInverse(t, NewMatrix(2, 0, 4, 5, 0, 0))
	testInverse(t, NewMatrix(2, 3, 4, 5, 0, 0))
	testInverse(t, NewMatrix(2, 0, 0, 5, 0.1, 0))
	testInverse(t, NewMatrix(2, 6, 6, 5, 0.1, 0))
	testInverse(t, NewMatrix(2, 6, 6, 5, 0.1, 0.3))
	testInverse(t, NewMatrix(1, 1, -1, 1, 0.1, 0))
	testInverse(t, NewMatrix(2, 3, 4, 5, 0.1, 0))
	testInverse(t, NewMatrix(2, 3, 4, 5, 0.1, 0.2))
	testInverse(t, NewMatrix(1e8, 0, 0, 1, 0.1, 0.2))
	testInverse(t, NewMatrix(1e8, 0, 0, 1e-8, 0.1, 0.2))
	testInverse(t, NewMatrix(0, 1e8, 1e-8, 0, 0.1, 0.2))
	testInverse(t, NewMatrix(0, 1e8, 1e-8, 0, 1e-8, 1e8))
	testInverse(t, NewMatrix(1e8, -1e8, 1e-8, -2e-8, 0, 0))
	testInverse(t, NewMatrix(1e8, -1e8, 1e-8, -2e-8, 5, 5))
	testInverse(t, NewMatrix(1e8, -1e8, 1e-8, -2e-8, 5, 1e-8))
	testInverse(t, NewMatrix(1, 1-1e5, 1, 1+1e5, 0, 0))
}

// testInverse tests if `m`.Inverse() is the inverse of `m`.
func testInverse(t *testing.T, m Matrix) {
	inv, hasInverse := m.Inverse()
	if !hasInverse {
		t.Fatalf("No inverse for %s", m)
	}
	pre := m.Mult(inv)
	if !isIdentity(pre) {
		t.Fatalf("Not pre-inverse:\n"+
			"\t   m=%s\n"+
			"\t inv=%s\n"+
			"\t pre=%s\n\t", m, inv, pre)
	}
	post := inv.Mult(m)
	if !isIdentity(post) {
		t.Fatalf("Not post-inverse:\n"+
			"\t   m=%s\n"+
			"\t inv=%s\n"+
			"\tpost=%s\n\t", m, inv, post)
	}
}

// isIdentity returns true if `m` approximates the identity matrix.
func isIdentity(m Matrix) bool {
	return isOne(m[0]) && isZero(m[1]) && isZero(m[2]) &&
		isZero(m[3]) && isOne(m[4]) && isZero(m[5]) &&
		isZero(m[6]) && isZero(m[7]) && isOne(m[8])
}

// isZero returns true if `x` is approximately one.
func isOne(x float64) bool {
	return isZero(x - 1.0)
}

// isZero returns true if `x` is approximately zero.
func isZero(x float64) bool {
	return math.Abs(x) <= tolerance
}

const tolerance = 1.0e-10
