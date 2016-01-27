package atlas

import (
	"errors"
	"github.com/gonutz/binpacker"
	"image"
	"image/draw"
)

// Atlas is itself a draw.Image and contains a number of sub-images.
type Atlas struct {
	draw.Image
	SubImages []SubImage
	packer    *binpacker.Packer
}

// SubImage is a draw.Image with a string ID. The image is always a reference
// back to the Atlas image, just the bounds are different.
type SubImage struct {
	draw.Image
	ID string
}

// New creates a new square image atlas of the given size. The bounds will have
// (0,0) as the minimum and the image is of type NRGBA.
func New(size int) *Atlas {
	packer := binpacker.New(size, size)
	return &Atlas{
		Image:  image.NewNRGBA(image.Rect(0, 0, size, size)),
		packer: packer,
	}
}

// NewFromImage uses the given image as the destination for all calls to Add.
// It is assumed to be empty at the beginning so all the available space will be
// used for sub-images.
func NewFromImage(atlas draw.Image) *Atlas {
	packer := binpacker.New(atlas.Bounds().Dx(), atlas.Bounds().Dy())
	return &Atlas{
		Image:  atlas,
		packer: packer,
	}
}

// GetSubImageByID returns the first SubImage with the given ID. If none is
// found, an error is returned.
func (a *Atlas) GetSubImageByID(id string) (SubImage, error) {
	for i := range a.SubImages {
		if a.SubImages[i].ID == id {
			return a.SubImages[i], nil
		}
	}
	return SubImage{}, errors.New("no SubImage with ID '" + id + "' found")
}

// Add finds a position for the given image in the atlas and copies the image
// there. It returns the new sub-image. If there is no more space in the atlas,
// it returns an error.
func (a *Atlas) Add(id string, img image.Image) (SubImage, error) {
	// determine the position of the new image
	rect, err := a.packer.Insert(img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		return SubImage{}, errors.New("unable to add image to atlas: " + err.Error())
	}
	var bounds image.Rectangle
	bounds.Min = a.Image.Bounds().Min.Add(image.Pt(rect.X, rect.Y))
	bounds.Max = bounds.Min.Add(image.Pt(rect.Width, rect.Height))

	// copy the image data into the atlas
	draw.Draw(a.Image, bounds, img, img.Bounds().Min, draw.Src)

	sub := SubImage{
		Image: subImage{
			a.Image,
			bounds,
		},
		ID: id,
	}
	a.SubImages = append(a.SubImages, sub)

	return sub, nil
}

type subImage struct {
	draw.Image
	bounds image.Rectangle
}

func (i subImage) Bounds() image.Rectangle {
	return i.bounds
}
