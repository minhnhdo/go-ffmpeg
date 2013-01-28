package ffmpeg

/*
#include <stdlib.h>
#include <libavformat/avformat.h>

AVStream* AVFormatContext_GetStream(AVFormatContext* fmtctx, int streamid) {
    return fmtctx->streams[streamid];
}
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
    Streams []Stream
    DecodedStreams []int
    packets chan *C.AVPacket
}

func Open(name string) (*MediaFile, error) {
    file := &MediaFile{
        Name: name,
        packets: make(chan *C.AVPacket, 8),
    }

    cName := C.CString(name)
    defer C.free(unsafe.Pointer(cName))

    if C.avformat_open_input(&file.fmtctx, cName, nil, nil) < 0 {
        return nil, errors.New("cannot open file " + name)
    }
    runtime.SetFinalizer(file, (*MediaFile).Close)

    if C.avformat_find_stream_info(file.fmtctx, nil) < 0 {
        return nil, errors.New("cannot find stream info for file " + name)
    }

    file.Streams = make([]Stream, file.fmtctx.nb_streams)
    for i := range file.Streams {
        file.Streams[i] = Stream{
            avstream: C.AVFormatContext_GetStream(file.fmtctx, C.int(i)),
        }
    }

    return file, nil
}

func (file *MediaFile) DecodeStream(index int) {
    if index < 0 || index >= len(file.Streams) {
        return
    }

    file.DecodedStreams = append(file.DecodedStreams, index)

    file.Streams[index].init()
}

func (file *MediaFile) StartDecoding() {
    for _, i := range file.DecodedStreams {
        go file.Streams[i].decode()
    }

    go func() {
        outer:
        for packet := range file.packets {
            if packet == nil {
                break outer
            }
            for _, i := range file.DecodedStreams {
                if packet.stream_index == C.int(i) {
                    file.Streams[i].packets <- packet
                    continue outer
                }
            }
            C.av_free_packet(packet)
        }
        for _, i := range file.DecodedStreams {
            file.Streams[i].packets <- nil
        }
    }()

    go func() {
        for {
            packet := new(C.AVPacket)
            C.av_init_packet(packet)
            if C.av_read_frame(file.fmtctx, packet) < 0 {
                // assume EOF
                file.packets <- nil
                return
            }
            file.packets <- packet
        }
    }()
}

func (file *MediaFile) IndexBestStream(mediaType MediaType) int {
    return int(C.av_find_best_stream(file.fmtctx, C.enum_AVMediaType(mediaType), -1, -1, nil, 0))
}

func (file *MediaFile) IndexFirstStream(mediaType MediaType) int {
    for i := range file.Streams {
        if file.Streams[i].avstream.codec.codec_type == C.enum_AVMediaType(mediaType) {
            return i
        }
    }
    return -1
}

func (file *MediaFile) Close() error {
    for i := range file.Streams {
        file.Streams[i].FreeCodecContext()
    }

    if file.fmtctx != nil {
        C.avformat_close_input(&file.fmtctx)
    }

    runtime.SetFinalizer(file, nil)

    return nil
}
