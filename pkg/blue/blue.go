package blue

import (
	"iter"
	"math"
	"math/rand/v2"

	"gonum.org/v1/gonum/dsp/fourier"
)

func GenerateBlueNoiseSignal(samples int) []float64 {
	whiteNoise := make([]float64, samples)
	for i := range whiteNoise {
		whiteNoise[i] = rand.NormFloat64()
	}

	fft := fourier.NewFFT(samples)
	spectrum := fft.Coefficients(nil, whiteNoise)

	for i := range spectrum {
		freq := fft.Freq(i)
		spectrum[i] *= complex(math.Sqrt(freq), 0)
	}

	blueNoise := fft.Sequence(nil, spectrum)
	mean, std := meanStd(blueNoise)
	for i := range blueNoise {
		blueNoise[i] = (blueNoise[i] - mean) / std
	}

	return blueNoise
}

func GenerateBluePoints(width, height, r float64, k int) [][2]float64 {
	cellSize := r / math.Sqrt2
	gridW := int(math.Ceil(width / cellSize))
	gridH := int(math.Ceil(height / cellSize))

	grid := make([][]*[2]float64, gridW)
	for i := range grid {
		grid[i] = make([]*[2]float64, gridH)
	}

	var points [][2]float64
	var activeList [][2]float64

	p0 := [2]float64{rand.Float64() * width, rand.Float64() * height}
	points = append(points, p0)
	activeList = append(activeList, p0)
	grid[int(p0[0]/cellSize)][int(p0[1]/cellSize)] = &points[len(points)-1]

	for len(activeList) > 0 {
		idx := rand.IntN(len(activeList))
		source := activeList[idx]
		found := false

		for range k {
			angle := rand.Float64() * 2 * math.Pi
			dist := r + rand.Float64()*r
			cx := source[0] + dist*math.Cos(angle)
			cy := source[1] + dist*math.Sin(angle)

			if cx < 0 || cx >= width || cy < 0 || cy >= height {
				continue
			}

			cellX := int(cx / cellSize)
			cellY := int(cy / cellSize)

			tooClose := false
			for x := max(0, cellX-2); x <= min(gridW-1, cellX+2); x++ {
				for y := max(0, cellY-2); y <= min(gridH-1, cellY+2); y++ {
					if neighbor := grid[x][y]; neighbor != nil {
						dx := cx - neighbor[0]
						dy := cy - neighbor[1]
						if dx*dx+dy*dy < r*r {
							tooClose = true
							break
						}
					}
				}
				if tooClose {
					break
				}
			}

			if !tooClose {
				points = append(points, [2]float64{cx, cy})
				grid[cellX][cellY] = &points[len(points)-1]
				activeList = append(activeList, points[len(points)-1])
				found = true
				break
			}
		}

		if !found {
			activeList = append(activeList[:idx], activeList[idx+1:]...)
		}
	}

	return points
}

func StreamBlueSignal(beta float64) iter.Seq[float64] {
	return func(yield func(float64) bool) {
		lastWhite := rand.NormFloat64()
		lastBlue := 0.0
		for {
			white := rand.NormFloat64()
			blue := (white - lastWhite) + (beta * lastBlue)
			lastWhite = white
			lastBlue = blue
			if !yield(blue) {
				return
			}
		}
	}
}

func StreamBluePoints(width, height, r float64, k int) iter.Seq[[2]float64] {
	return func(yield func([2]float64) bool) {
		cellSize := r / math.Sqrt2
		gridW := int(math.Ceil(width / cellSize))
		gridH := int(math.Ceil(height / cellSize))

		grid := make([][]*[2]float64, gridW)
		for i := range grid {
			grid[i] = make([]*[2]float64, gridH)
		}

		var points [][2]float64
		var activeList [][2]float64

		p0 := [2]float64{rand.Float64() * width, rand.Float64() * height}
		points = append(points, p0)
		activeList = append(activeList, p0)
		grid[int(p0[0]/cellSize)][int(p0[1]/cellSize)] = &points[len(points)-1]

		if !yield(p0) {
			return
		}

		for len(activeList) > 0 {
			idx := rand.IntN(len(activeList))
			source := activeList[idx]
			found := false

			for range k {
				angle := rand.Float64() * 2 * math.Pi
				dist := r + rand.Float64()*r
				cx := source[0] + dist*math.Cos(angle)
				cy := source[1] + dist*math.Sin(angle)

				if cx < 0 || cx >= width || cy < 0 || cy >= height {
					continue
				}

				cellX := int(cx / cellSize)
				cellY := int(cy / cellSize)

				tooClose := false
				for x := max(0, cellX-2); x <= min(gridW-1, cellX+2); x++ {
					for y := max(0, cellY-2); y <= min(gridH-1, cellY+2); y++ {
						if neighbor := grid[x][y]; neighbor != nil {
							dx := cx - neighbor[0]
							dy := cy - neighbor[1]
							if dx*dx+dy*dy < r*r {
								tooClose = true
								break
							}
						}
					}
					if tooClose {
						break
					}
				}

				if !tooClose {
					points = append(points, [2]float64{cx, cy})
					grid[cellX][cellY] = &points[len(points)-1]
					activeList = append(activeList, points[len(points)-1])
					found = true
					if !yield(points[len(points)-1]) {
						return
					}
					break
				}
			}

			if !found {
				activeList = append(activeList[:idx], activeList[idx+1:]...)
			}
		}
	}
}

func meanStd(data []float64) (float64, float64) {
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))

	var sq float64
	for _, v := range data {
		d := v - mean
		sq += d * d
	}
	std := math.Sqrt(sq / float64(len(data)))
	return mean, std
}
