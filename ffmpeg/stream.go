package ffmpeg

/*
#include <libavformat/avformat.h>
*/
import "C"

type Stream struct {
    stream *C.AVStream
    cdcctx *C.AVCodecContext
}

func (stream *Stream) init() {
    if stream.stream == nil {
        return
    }

    if stream.cdcctx != nil {
        stream.FreeCodecContext()
    }
    stream.cdcctx = stream.stream.codec

    if decoder := C.avcodec_find_decoder(stream.cdcctx.codec_id); decoder == nil || C.avcodec_open2(stream.cdcctx, decoder, nil) < 0 {
        stream.cdcctx = nil
        return
    }
}

func (stream *Stream) FreeCodecContext() {
    if stream.cdcctx != nil {
        C.avcodec_close(stream.cdcctx)
        stream.cdcctx = nil
    }
}
