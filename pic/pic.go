package pic

import (
	"image"
	"os"

	"github.com/disintegration/imaging"
)

const MaxImageLength = 100

// var imageCache map[string]*image.NRGBA = make(map[string]*image.NRGBA)
func LoadPicPath(path string, rows, cols int) (pic *image.NRGBA, err error) {
	var r *os.File
	r, err = os.Open(path)
	if err != nil {
		return
	}
	defer r.Close()
	var x, y int
	// cached, found := imageCache[path]
	// if !found {

	var tmp image.Image
	tmp, err = imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return
	}

	// if tmp.Bounds().Dx() > tmp.Bounds().Dy() {
	// 	x = MaxImageLength
	// } else {
	// 	y = MaxImageLength
	// }

	// leaving x or y at 0 maintains aspect ratio
	// cached = imaging.Resize(tmp, x, y, imaging.Lanczos)
	// imageCache[path] = cached
	// }

	if cols > rows {
		x = cols
	} else {
		y = rows
	}

	pic = imaging.Resize(tmp, x, y, imaging.Lanczos)
	return
}
