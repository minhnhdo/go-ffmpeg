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
    "fmt"
    "io"
    "runtime"
    "unsafe"
)

type MediaFile struct {
    Name string
    fmtctx *C.AVFormatContext
    Streams []Stream
    StreamIndex int
    CurrentFrame *Frame
}

func Open(name string) (*MediaFile, error) {
    file := &MediaFile{
        Name: name,
        StreamIndex: -1,
        CurrentFrame: &Frame{},
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
            stream: C.AVFormatContext_GetStream(file.fmtctx, C.int(i)),
        }
        file.Streams[i].init()
    }

    return file, nil
}

func (file *MediaFile) IndexBestStream(mediaType MediaType) int {
    return int(C.av_find_best_stream(file.fmtctx, C.enum_AVMediaType(mediaType), -1, -1, nil, 0))
}

func (file *MediaFile) IndexFirstStream(mediaType MediaType) int {
    for i, s := range file.Streams {
        if s.stream.codec.codec_type == C.enum_AVMediaType(mediaType) {
            return i
        }
    }
    return -1
}

func (file *MediaFile) NextFrame() error {
    var (
        gotFrame C.int = 0
        packet C.AVPacket
    )

    for {
        if C.av_read_frame(file.fmtctx, &packet) < 0 {
            return io.EOF
        }
        defer C.av_free_packet(&packet)

        if index := packet.stream_index; file.StreamIndex < 0 || index == C.int(file.StreamIndex) {
            file.CurrentFrame.Defaults()
            switch file.Streams[index].stream.codec.codec_type {
            case C.AVMEDIA_TYPE_AUDIO:
                if C.avcodec_decode_audio4(file.Streams[index].cdcctx, &file.CurrentFrame.frame, &gotFrame, &packet) < 0 {
                    return errors.New(fmt.Sprintf("cannot decode audio packet from stream %d", index))
                }
            case C.AVMEDIA_TYPE_VIDEO:
                if C.avcodec_decode_video2(file.Streams[index].cdcctx, &file.CurrentFrame.frame, &gotFrame, &packet) < 0 {
                    return errors.New(fmt.Sprintf("cannot decode video packet from stream %d", index))
                }
            default:
                // unsupported media type
                continue
            }
            if gotFrame != 0 {
                file.CurrentFrame.PTS = int64(C.av_frame_get_best_effort_timestamp(&file.CurrentFrame.frame))
                return nil
            }
        }
    }
    return errors.New("should not reach end of NextFrame method")
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
