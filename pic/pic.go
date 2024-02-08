package pic

import (
	"image"
	"os"

	"github.com/disintegration/imaging"
)

type ResampleItem uint

const (
	NearestNeighbor ResampleItem = iota
	Box
	Linear
	Hermite
	MitchellNetravali
	CatmullRom
	BSpline
	Gaussian
	Bartlett
	Lanczos
	Hann
	Hamming
	Blackman
	Welch
	Cosine
	RESAMPLE_ITEM_COUNT
)

var ResampleList = []string{
	"Nearest Neighbor",
	"Box",
	"Linear",
	"Hermite",
	"Mitchell / Netravali",
	"Catmull / Rom",
	"BSpline",
	"Gaussian",
	"Bartlett",
	"Lanczos",
	"Hann",
	"Hamming",
	"Blackman",
	"Welch",
	"Cosine",
}

func (i ResampleItem) String() string {
	if i >= RESAMPLE_ITEM_COUNT {
		return ""
	}
	return ResampleList[i]
}

func (i ResampleItem) Filter() imaging.ResampleFilter {
	switch i {
	case NearestNeighbor:
		return imaging.NearestNeighbor
	case Box:
		return imaging.Box
	case Linear:
		return imaging.Linear
	case Hermite:
		return imaging.Hermite
	case MitchellNetravali:
		return imaging.MitchellNetravali
	case CatmullRom:
		return imaging.CatmullRom
	case BSpline:
		return imaging.BSpline
	case Gaussian:
		return imaging.Gaussian
	case Bartlett:
		return imaging.Bartlett
	case Lanczos:
		return imaging.Lanczos
	case Hann:
		return imaging.Hann
	case Hamming:
		return imaging.Hamming
	case Blackman:
		return imaging.Blackman
	case Welch:
		return imaging.Welch
	case Cosine:
		return imaging.Cosine
	}
	return imaging.NearestNeighbor
}

func ResamplePath(path string, rows, cols int, filter imaging.ResampleFilter) (pic *image.NRGBA, err error) {
	var r *os.File
	r, err = os.Open(path)
	if err != nil {
		return
	}
	defer r.Close()

	var tmp image.Image
	tmp, err = imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return
	}

	pic = imaging.Resize(tmp, cols, rows, filter)
	return
}

func LoadPicPath(path string, rows, cols int) (pic *image.NRGBA, err error) {
	return ResamplePath(path, rows, cols, imaging.Box)
}
