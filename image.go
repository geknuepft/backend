package main

import (
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

func ImageGet(format int, path string) (oupImg *image.NRGBA, err error) {

	inpDir := os.Getenv("IMAGE_INPUT_DIR")
	if inpDir == "" {
		log.Fatal("env variable IMAGE_INPUT_DIR must be set")
		return
	}
	filePath := inpDir + path
	log.Printf("reading %v", filePath)

	inpImgF, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer inpImgF.Close()

	// read and decode exif data to figure out the orientation
	x, err := exif.Decode(inpImgF)
	if err != nil {
		return
	}

	rotate := false
	if tag, err := x.Get(exif.Orientation); err != nil && tag != nil {
		var orientation = tag.Val
		rotate = orientation[0] == 8
	}

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
