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
    Name string
    fmtctx *C.AVFormatContext
}

func Open(name string) (*MediaFile, error) {
    f := &MediaFile{Name: name}

    cName := C.CString(name)
    defer C.free(unsafe.Pointer(cName))

    if C.avformat_open_input(&f.fmtctx, cName, nil, nil) < 0 {
        return nil, errors.New("cannot open file " + name)
    }

    return f, nil
}

func (f *MediaFile) Close() error {
    if f.fmtctx != nil {
        C.avformat_close_input(&f.fmtctx)
    }

    return nil
}
