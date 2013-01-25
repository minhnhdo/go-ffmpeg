package ffmpeg

/*
#include <stdlib.h>
#include <libavformat/avformat.h>
*/
import "C"
import (
    "errors"
    "runtime"
    "unsafe"
)

type MediaFile struct {
    Name string
    fmtctx *C.AVFormatContext
}

func Open(name string) (*MediaFile, error) {
    file := &MediaFile{Name: name}

    cName := C.CString(name)
    defer C.free(unsafe.Pointer(cName))

    if C.avformat_open_input(&file.fmtctx, cName, nil, nil) < 0 {
        return nil, errors.New("cannot open file " + name)
    }

    if C.avformat_find_stream_info(file.fmtctx, nil) < 0 {
        file.Close()
        return nil, errors.New("cannot find stream info for file " + name)
    }

    runtime.SetFinalizer(file, (*MediaFile).Close)

    return file, nil
}

func (file *MediaFile) Close() error {
    if file.fmtctx != nil {
        C.avformat_close_input(&file.fmtctx)
    }

    runtime.SetFinalizer(file, nil)

    return nil
}
