package pattern

import (
	"image"
	"image/png"
	"os"
)

var WorleyNoiseOutputFilename = "worley.png"
var WorleyNoiseZoomLevels = []int{}

const WorleyNoiseBaseLabel = "WorleyNoise"

// WorleyNoise Pattern
// Generates Worley (cellular) noise.
func ExampleNewWorleyNoise() {
	// Standard F1 Euclidean Worley Noise
	i := NewWorleyNoise(
		SetFrequency(0.05),
		SetSeed(1),
	)
	f, err := os.Create(WorleyNoiseOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateWorleyNoise(b image.Rectangle) image.Image {
	return NewWorleyNoise(SetBounds(b), SetFrequency(0.05), SetSeed(1))
}

func GenerateWorleyNoiseReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	refs := make(map[string]func(image.Rectangle) image.Image)
	var keys []string

	add := func(key string, ops ...func(any)) {
		keys = append(keys, key)
		refs[key] = func(b image.Rectangle) image.Image {
			baseOps := []func(any){SetBounds(b), SetFrequency(0.05), SetSeed(1)}
			return NewWorleyNoise(append(baseOps, ops...)...)
		}
	}

	add("F1_Euclidean", SetWorleyMetric(MetricEuclidean), SetWorleyOutput(OutputF1))
	add("F2_Euclidean", SetWorleyMetric(MetricEuclidean), SetWorleyOutput(OutputF2))
	add("F2MinusF1_Euclidean", SetWorleyMetric(MetricEuclidean), SetWorleyOutput(OutputF2MinusF1))
	add("CellID", SetWorleyOutput(OutputCellID))

	add("F1_Manhattan", SetWorleyMetric(MetricManhattan), SetWorleyOutput(OutputF1))
	add("F2_Manhattan", SetWorleyMetric(MetricManhattan), SetWorleyOutput(OutputF2))
	add("F2MinusF1_Manhattan", SetWorleyMetric(MetricManhattan), SetWorleyOutput(OutputF2MinusF1))

	add("F1_Chebyshev", SetWorleyMetric(MetricChebyshev), SetWorleyOutput(OutputF1))
	add("F2_Chebyshev", SetWorleyMetric(MetricChebyshev), SetWorleyOutput(OutputF2))
	add("F2MinusF1_Chebyshev", SetWorleyMetric(MetricChebyshev), SetWorleyOutput(OutputF2MinusF1))

	return refs, keys
}

func init() {
	RegisterGenerator(WorleyNoiseBaseLabel, GenerateWorleyNoise)
	RegisterReferences(WorleyNoiseBaseLabel, GenerateWorleyNoiseReferences)
}
