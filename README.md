An extension of https://github.com/rwcarlsen/goexif.

The extension includes in detail:
* All tags in SubIfds (e.g. tag 0x14a) or IFDs other than IFD0 get a prefix as TagName (e.g. "subN.", "gps.", "exif." etc). This allows storing tags with the same name that are referenced in different (sub-)IFDs.
* Definition of functions that extract images stored in sub-IFDs or as preview. This allows e.g. to extract the low resolution JPG images as well as the JPG preview image stored in Nikon NEF files (see extension_example/main.go)

References for Exif/Tags:
* exif file format <https://www.media.mit.edu/pia/Research/deepview/exif.html>
* exif tags (Phil Haryvey) <https://exiftool.org/TagNames/EXIF.html>
* Nikon NEF format <http://lclevy.free.fr/nef/>
* Nikon Tags <https://exiftool.org/TagNames/Nikon.html>
