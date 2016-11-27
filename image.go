package main

import (
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

var inpDir = "/srv/db_tracked/"

func ImageGet(format int, path string) (oupImg *image.NRGBA, err error) {

	inpImgF, err := os.Open(inpDir + path)
	if err != nil {
		return
	}
	defer inpImgF.Close()

	// read and decode exif data to figure out the orientation
	x, err := exif.Decode(inpImgF)
	if err != nil {
		return
	}

	tag, _ := x.Get(exif.Orientation)
	var orientation = tag.Val
	var rotate = orientation[0] == 8

	// read and decode jpg image
	inpImgF.Seek(0, 0)
	inpImg, err := jpeg.Decode(inpImgF)
	if err != nil {
		return
	}

	// compute output image
	var width, height int
	if rotate {
		width = int(format)
		height = 0
	} else {
		width = 0
		height = int(format)
	}

	// resize image
	oupImg = imaging.Resize(inpImg, width, height, imaging.Box)

	// rotate image if necessary
	if rotate {
		oupImg = imaging.Rotate90(oupImg)
	}

	return
}

func ImageWrite(oupWriter io.Writer, oupImg *image.NRGBA) (err error) {
	return jpeg.Encode(oupWriter, oupImg, nil)
}
