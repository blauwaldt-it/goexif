package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blauwaldt-it/goexif/exif"
	"github.com/blauwaldt-it/goexif/mknote"
	"github.com/blauwaldt-it/goexif/tiff"
)

func main() {

	fn := "..\\pixdb_go\\tmp\\DSC_8508.NEF"

	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	exif.RegisterParsers(mknote.NikonV3)

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	var p Printer
	x.Walk(p)

	// fmt.Println("++++++++++++++++++++++++++++++++")
	// fmt.Println(len(x.Tiff.Dirs))
	// fmt.Println(x.Tiff.Dirs[0])
	// for _, v := range x.Tiff.Dirs[0].Tags {
	// 	if v.Id == 0x14a {
	// 		fmt.Println(v)
	// 	}
	// }
	// fmt.Println()

	// tn, err2 := x.JpegThumbnail()
	// if err2 == nil {
	// 	fmt.Println(tn)
	// } else {
	// 	mknote, err2 := x.Get(exif.MakerNote)
	// 	if err2 == nil {
	// 		fmt.Printf("~~~~~~~~~~~~~~~~~~ %s %T", mknote, mknote)
	// 		for i, v := range mknote.Val {
	// 			fmt.Printf("%02x ", v)
	// 			if i > 400 {
	// 				break
	// 			}
	// 		}
	// 		fmt.Printf("\n")
	// 		for i, v := range mknote.Val[10+0x9d42:] {
	// 			fmt.Printf("%02x ", v)
	// 			if i > 400 {
	// 				break
	// 			}
	// 		}
	// 		fmt.Printf("\n")
	// 		for _, v := range mknote.Val[10+0x9d9e : 8+10+0x9d9e] {
	// 			fmt.Printf("%02x ", v)
	// 		}
	// 		fmt.Printf("\n")
	// 		for _, v := range mknote.Val[10+0x9da6 : 8+10+0x9da6] {
	// 			fmt.Printf("%02x ", v)
	// 		}
	// 		fmt.Printf("\n")

	// 		// fmt.Println(mknote.Val)
	// 	}
	// }

	return
}

type Printer struct{}

func (p Printer) Walk(name exif.FieldName, tag *tiff.Tag) error {
	fmt.Println(name, tag)
	return nil
}
