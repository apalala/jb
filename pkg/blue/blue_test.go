package blue

import (
	"math"
	"testing"
)

const tol = 0.15

func TestGenerateBlueNoiseSignalLength(t *testing.T) {
	samples := 256
	signal := GenerateBlueNoiseSignal(samples)
	if len(signal) != samples {
		t.Fatalf("expected length %d, got %d", samples, len(signal))
	}
}

func TestGenerateBlueNoiseSignalStats(t *testing.T) {
	samples := 4096
	signal := GenerateBlueNoiseSignal(samples)
	mean, std := meanStd(signal)
	if math.Abs(mean) > tol {
		t.Errorf("mean ≈ %f, expected near 0", mean)
	}
	if math.Abs(std-1) > tol {
		t.Errorf("std ≈ %f, expected near 1", std)
	}
}

func TestGenerateBluePointsBounds(t *testing.T) {
	w, h, r := 100.0, 100.0, 4.0
	points := GenerateBluePoints(w, h, r, 30)
	if len(points) == 0 {
		t.Fatal("expected at least one point")
	}
	for i, p := range points {
		if p[0] < 0 || p[0] >= w || p[1] < 0 || p[1] >= h {
			t.Fatalf("point %d out of bounds: %v", i, p)
		}
	}
}

func TestGenerateBluePointsMinDistance(t *testing.T) {
	w, h, r := 100.0, 100.0, 4.0
	points := GenerateBluePoints(w, h, r, 30)
	r2 := r * r
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dx := points[i][0] - points[j][0]
			dy := points[i][1] - points[j][1]
			if dx*dx+dy*dy < r2-1e-9 {
				t.Fatalf("points %d and %d too close: dist^2=%f < %f", i, j, dx*dx+dy*dy, r2)
			}
		}
	}
}

func TestStreamBlueSignalFinite(t *testing.T) {
	var count int
	for range StreamBlueSignal(0.65) {
		count++
		if count >= 10000 {
			break
		}
	}
	if count != 10000 {
		t.Fatalf("expected 10000 values, got %d", count)
	}
}

func TestStreamBlueSignalStats(t *testing.T) {
	var values []float64
	for v := range StreamBlueSignal(0.65) {
		values = append(values, v)
		if len(values) >= 10000 {
			break
		}
	}
	mean, std := meanStd(values)
	if math.Abs(mean) > tol {
		t.Errorf("mean ≈ %f, expected near 0", mean)
	}
	if math.Abs(std-1) > tol {
		t.Errorf("std ≈ %f, expected near 1", std)
	}
}

func TestStreamBluePointsBounds(t *testing.T) {
	w, h, r := 50.0, 50.0, 3.0
	var count int
	for p := range StreamBluePoints(w, h, r, 30) {
		if p[0] < 0 || p[0] >= w || p[1] < 0 || p[1] >= h {
			t.Fatalf("point %v out of bounds", p)
		}
		count++
		if count >= 200 {
			break
		}
	}
	if count == 0 {
		t.Fatal("expected at least one point")
	}
}

func TestStreamBluePointsMinDistance(t *testing.T) {
	w, h, r := 50.0, 50.0, 3.0
	var points [][2]float64
	for p := range StreamBluePoints(w, h, r, 30) {
		points = append(points, p)
		if len(points) >= 200 {
			break
		}
	}
	r2 := r * r
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dx := points[i][0] - points[j][0]
			dy := points[i][1] - points[j][1]
			if dx*dx+dy*dy < r2-1e-9 {
				t.Fatalf("points %d and %d too close: dist^2=%f < %f", i, j, dx*dx+dy*dy, r2)
			}
		}
	}
}
