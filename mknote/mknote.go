// Package mknote provides makernote parsers that can be used with goexif/exif.
package mknote

import (
	"bytes"
	"errors"

	"github.com/blauwaldt-it/goexif/exif"
	"github.com/blauwaldt-it/goexif/tiff"
)

var (
	// Canon is an exif.Parser for canon makernote data.
	Canon = &canon{}
	// NikonV3 is an exif.Parser for nikon makernote data.
	NikonV3 = &nikonV3{}
	// All is a list of all available makernote parsers
	All = []exif.Parser{Canon, NikonV3}
)

type canon struct{}

// Parse decodes all Canon makernote data found in x and adds it to x.
func (_ *canon) Parse(x *exif.Exif) error {
	m, err := x.MakerNote()
	if err != nil {
		return nil
	}

	mk, err := x.Get(exif.Make)
	if err != nil {
		return nil
	}

	if val, err := mk.StringVal(); err != nil || val != "Canon" {
		return nil
	}

	// Canon notes are a single IFD directory with no header.
	// Reader offsets need to be w.r.t. the original tiff structure.
	buf := bytes.NewReader(append(make([]byte, m.ValOffset), m.Val...))
	buf.Seek(int64(m.ValOffset), 0)

	mkNotesDir, _, err := tiff.DecodeDir(buf, x.Tiff.Order)
	if err != nil {
		return err
	}
	x.LoadTagsPref(mkNotesDir, makerNoteCanonFields, false, "mkcanon.")
	return nil
}

type nikonV3 struct{}

// Parse decodes all Nikon makernote data found in x and adds it to x.
func (_ *nikonV3) Parse(x *exif.Exif) error {
	m, err := x.MakerNote()
	if err != nil {
		return nil
	}

	if bytes.Compare(m.Val[:6], []byte("Nikon\000")) != 0 {
		return nil
	}

	// Nikon v3 maker note is a self-contained IFD (offsets are relative
	// to the start of the maker note)
	mkNotes, err := tiff.Decode(bytes.NewReader(m.Val[10:]))
	if err != nil {
		return err
	}
	x.LoadTagsPref(mkNotes.Dirs[0], makerNoteNikon3Fields, false, "mknikon.")

	// Read Preview SubIFD, if present
	tag, err := x.Get("mknikon." + Preview)
	if err == nil {
		offset, err := tag.Int64(0)
		if err == nil {
			r := bytes.NewReader(m.Val[10:])
			_, err := r.Seek(offset, 0)
			if err == nil {
				subDir, _, err := tiff.DecodeDir(r, x.Tiff.Order)
				if err == nil {
					x.LoadTagsPref(subDir, exif.ExifFields, false, "mknprev.")
				}
			}
		}
	}

	return nil
}

func (_ *nikonV3) MakerNotePreview(x *exif.Exif) ([]byte, error) {

	mn, err := x.MakerNote()
	if err != nil {
		return nil, err
	}
	if bytes.Compare(mn.Val[:6], []byte("Nikon\000")) != 0 {
		return nil, errors.New("No Nikon makernote found")
	}

	pstt, err := x.Get("mknprev." + exif.PreviewImageStart)
	if err != nil {
		return nil, err
	}
	pst, err := pstt.Int(0)
	if err != nil {
		return nil, err
	}

	plent, err := x.Get("mknprev." + exif.PreviewImageLength)
	if err != nil {
		return nil, err
	}
	plen, err := plent.Int(0)
	if err != nil {
		return nil, err
	}

	return mn.Val[10+pst : 10+pst+plen], nil
}
