// Package palettor provides a way to extract the color palette from an image
// using k-means clustering.
package palettor

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

// Extract finds the k most dominant colors in the given image using the
// "standard" k-means clustering algorithm. It returns a Palette, after running
// the algorithm up to maxIterations times.
func Extract(k, maxIterations int, img image.Image) (*Palette, error) {
	return ClusterColors(k, maxIterations, GetColors(img))
}

// ExtractByCentroids ...
func ExtractByCentroids(th int, img image.Image, centroids map[string][]color.Color) (*PaletteCentroid, error) {
	return clusterColorsByCentroids(th, GetColors(img), centroids)
}

// GetColors from an image
func GetColors(img image.Image) []color.Color {
	bounds := img.Bounds()
	pixelCount := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	colors := make([]color.Color, pixelCount)
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colors[i] = img.At(x, y)
			i++
		}
	}
	return colors
}

func inBetween(i, min, max uint32) bool {
	return i >= min && i <= max
}

// ColorEq ...
func ColorEq(th uint32, a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	// return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
	return inBetween(r1, r2-th, r2+th) &&
		inBetween(g1, g2-th, g2+th) &&
		inBetween(b1, b2-th, b2+th) &&
		a1 == a2
}

// ColorsXor ...
func ColorsXor(th uint32, src, space []color.Color) []color.Color {
	result := []color.Color{}
	for _, c := range src {
		found := false
		for _, s := range space {
			if ColorEq(th, c, s) {
				found = true
				break
			}
		}
		for _, r := range result {
			if ColorEq(th, c, r) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, c)
		}
	}

	return result
}

// ColorsToImage ...
func ColorsToImage(src []color.Color) image.Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{1, len(src)}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	bounds := img.Bounds()
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, src[i])
			i++
		}
	}
	return img
}

// ReadImage ...
func ReadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(f)
	return image, err
}

// WriteImage ...
func WriteImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
	return nil
}
