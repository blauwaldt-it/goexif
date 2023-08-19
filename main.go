package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/blauwaldt-it/goexif/exif"
	"github.com/blauwaldt-it/goexif/mknote"
	"github.com/blauwaldt-it/goexif/tiff"
)

func main() {

	// OPen any NEF file
	fn := "..\\pixdb_go\\tmp\\DSC_8508.NEF"

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
			_, _ = f.Write(im)
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
