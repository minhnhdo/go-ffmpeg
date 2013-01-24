package ffmpeg

/*
#include <stdlib.h>
#include <libavformat/avformat.h>
*/
import "C"
import (
    "errors"
    "unsafe"
)

type MediaFile struct {
    Filename string
    fmtctx *C.AVFormatContext
}

func OpenFile(filename string) (f MediaFile, err error) {
    f.Filename = filename

    cFilename := C.CString(filename)
    defer C.free(unsafe.Pointer(cFilename))

    if C.avformat_open_input(&f.fmtctx, cFilename, nil, nil) < 0 {
        return MediaFile{}, errors.New("Cannot open file " + filename)
    }

    return
}

func (f MediaFile) Close() {
    if f.fmtctx != nil {
        C.avformat_close_input(&f.fmtctx)
    }
}
