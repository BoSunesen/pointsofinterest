package poidata

import (
	"math"
	"testing"
)

func TestConvertDegreesToRadians(t *testing.T) {
	precision := 1 << 30
	compareFloats(t, math.Pi/6, ConvertDegreesToRadians(30), precision)
	compareFloats(t, math.Pi/4, ConvertDegreesToRadians(45), precision)
	compareFloats(t, math.Pi/2, ConvertDegreesToRadians(90), precision)
	compareFloats(t, math.Pi, ConvertDegreesToRadians(180), precision)
	compareFloats(t, math.Pi*2, ConvertDegreesToRadians(360), precision)
}

func TestConvertRadiansToDegrees(t *testing.T) {
	precision := 1 << 30
	compareFloats(t, 30, ConvertRadiansToDegrees(math.Pi/6), precision)
	compareFloats(t, 45, ConvertRadiansToDegrees(math.Pi/4), precision)
	compareFloats(t, 90, ConvertRadiansToDegrees(math.Pi/2), precision)
	compareFloats(t, 180, ConvertRadiansToDegrees(math.Pi), precision)
	compareFloats(t, 360, ConvertRadiansToDegrees(math.Pi*2), precision)
}

func TestCalculateBoundingBox(t *testing.T) {
	precision := 1 << 20
	//TODO Test more
	box := CalculateBoundingBox(37.7760487, -122.423939, 100)
	compareFloats(t, 37.7769482, box.MaxLatitude, precision)
	compareFloats(t, 37.7751492, box.MinLatitude, precision)
	compareFloats(t, -122.422801, box.MaxLongitude, precision)
	compareFloats(t, -122.425077, box.MinLongitude, precision)
}

func compareFloats(t *testing.T, expected, actual float64, precision int) {
	if !areEqual(expected, actual, precision) {
		t.Errorf("Expected: %v - actual: %v - diff: %v ", expected, actual, math.Abs(expected-actual))
	}
}

func areEqual(float1, float2 float64, precision int) bool {
	var limit float64 = 1 / float64(precision)
	if math.Abs(float1-float2) < limit {
		return true
	}
	return false
}
