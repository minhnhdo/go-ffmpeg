package ffmpeg

/*
#include <libavformat/avformat.h>
*/
import "C"
import (
    "errors"
    "fmt"
)

type Stream struct {
    avstream *C.AVStream
    cdcctx *C.AVCodecContext
    decodeF decodeFunc
    frame *Frame
    Frames chan *Frame
    packets chan *C.AVPacket
    noMorePackets chan bool
}

func (stream *Stream) init() error {
    // need to allocate this first so that it can be closed on error
    stream.Frames = make(chan *Frame)

    if stream.avstream == nil {
        close(stream.Frames)
        return errors.New("nil avstream")
    }

    if stream.cdcctx != nil {
        stream.FreeCodecContext()
    }
    stream.cdcctx = stream.avstream.codec

    if decoder := C.avcodec_find_decoder(stream.cdcctx.codec_id); decoder == nil || C.avcodec_open2(stream.cdcctx, decoder, nil) < 0 {
        stream.cdcctx = nil
        close(stream.Frames)
        return errors.New(fmt.Sprintf("Cannot find decoder for %s", C.GoString(C.avcodec_get_name(stream.cdcctx.codec_id))))
    }

    stream.packets = make(chan *C.AVPacket)
    stream.frame = &Frame{}
    switch stream.avstream.codec.codec_type {
    case C.AVMEDIA_TYPE_AUDIO:
        stream.decodeF = avcodec_decode_audio
    case C.AVMEDIA_TYPE_VIDEO:
        stream.decodeF = avcodec_decode_video
    default:
        stream.FreeCodecContext()
        close(stream.Frames)
        return errors.New("unsupported codec")
    }

    return nil
}

func (stream *Stream) decode() {
    for packet := range stream.packets {
        if packet == nil {
            close(stream.Frames)
            return
        }
        stream.frame.Defaults()
        gotFrame, err := stream.decodeF(stream.cdcctx, stream.frame, packet)
        C.av_free_packet(packet)
        if err != nil {
            // ignore frame
            continue
        }
        if gotFrame {
            stream.Frames <- stream.frame
        }
    }
}

func (stream *Stream) FreeCodecContext() {
    if stream.cdcctx != nil {
        C.avcodec_close(stream.cdcctx)
        stream.cdcctx = nil
    }
}
