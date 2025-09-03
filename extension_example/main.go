package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/blauwaldt-it/goexif/exif"
	"github.com/blauwaldt-it/goexif/mknote"
	"github.com/blauwaldt-it/goexif/tiff"
)

func insert_orientation_exif(image []byte, orientation byte) []byte {

	if len(image) < 12 || image[0] != 0xFF || image[1] != 0xD8 {
		// unknown format
		return image
	}

	exif := []byte{
		0xFF, 0xD8, // SOI
		0xFF, 0xE1, 0x00, 0x22, // APP1 & Length (big-endian)
		0x45, 0x78, 0x69, 0x66, 0x00, 0x00, // Exif Header
		0x49, 0x49, 0x2A, 0x00, // TIFF Header Intel
		0x08, 0x00, 0x00, 0x00, // Offset to IFD0
		0x01, 0x00, // Number of Directory Entries
		0x12, 0x01, // Orientation
		0x03, 0x00, // Datatype Unsigned Short
		0x01, 0x00, 0x00, 0x00, // 1 component
		orientation, 0x00, 0x00, 0x00, // Value
		0x00, 0x00, 0x00, 0x00, // Next IFD Offset (0=End)
	}

	// Check if Exif is already present
	present := 1
	for v := range []int{2, 3, 6, 7, 8, 9, 10, 11} {
		if image[v] != exif[v] {
			present = 0
			break
		}
	}
	if present == 1 {
		// Exif is present
		return image
	}

	return append(exif, image[2:]...)
}

func main() {

	// Open any NEF file
	fn := "DSC_0404.NEF"

	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	// Parse exif tags
	exif.RegisterParsers(mknote.NikonV3)

	x, err := exif.Decode(f)
	if err != nil {
		panic(err)
	}

	// output all tags found
	var p Printer
	x.Walk(p)

	// Save images referenced in SubIfds
	for i := 0; i < 3; i++ {
		im, err := x.SubIfdImage(i)
		if err == nil {
			f, err := os.Create("subifb" + strconv.Itoa(i) + "img.jpg")
			if err != nil {
				continue
			}
			// _, _ = f.Write(im)
			_, _ = f.Write(insert_orientation_exif(im, 6))
			f.Close()
		}
	}

	// Save Preview image referenced in Nikon Makernote
	im, err := mknote.NikonV3.MakerNotePreview(x)
	if err == nil {
		f, err := os.Create("mknoteprev.jpg")
		if err == nil {
			_, _ = f.Write(im)
		}
		f.Close()
	}
}

type Printer struct{}

func (p Printer) Walk(name exif.FieldName, tag *tiff.Tag) error {
	fmt.Println(name, tag)
	return nil
}
