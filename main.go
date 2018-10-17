package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

// SubImager xxx
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// ImageSize represents the dimensions of an image
type ImageSize struct {
	AspectWidth  int
	AspectHeight int
	FinalWidth   uint
	Name         string
}

var sizes = [3]ImageSize{
	ImageSize{AspectWidth: 640, AspectHeight: 470, FinalWidth: 640, Name: "md"},
	ImageSize{AspectWidth: 200, AspectHeight: 200, FinalWidth: 200, Name: "sm"},
	ImageSize{AspectWidth: 80, AspectHeight: 80, FinalWidth: 80, Name: "xs"},
}

func main() {
	err := cropImage("./images", "test-1.jpg")
	if err != nil {
		fmt.Println(err)
	}
}

func cropImage(pathName, fileName string) error {
	fpath := path.Join(pathName, fileName)
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, size := range sizes {
		analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
		topCrop, err := analyzer.FindBestCrop(img, size.AspectWidth, size.AspectHeight)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// The crop will have the requested aspect ratio, but you need to copy/scale it yourself
		croppedimg := img.(SubImager).SubImage(topCrop)
		m := resize.Resize(size.FinalWidth, 0, croppedimg, resize.Lanczos3)
		f, err = os.Create(path.Join("./cropped", fmt.Sprintf("%s-%s", size.Name, fileName)))
		if err != nil {
			return errors.Wrap(err, "Cannot create file: "+fpath)
		}
		err = jpeg.Encode(f, m, &jpeg.Options{Quality: 85})
		if err != nil {
			return errors.Wrap(err, "Failed to encode the image as JPEG")
		}
	}

	return nil
}
