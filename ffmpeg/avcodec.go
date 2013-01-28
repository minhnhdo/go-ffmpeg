package ffmpeg

/*
#include <libavcodec/avcodec.h>
*/
import "C"
import (
    "errors"
)

type decodeFunc func(*C.AVCodecContext, *Frame, *C.AVPacket) (bool, error)

func avcodec_decode_audio(cdcctx *C.AVCodecContext, frame *Frame, packet *C.AVPacket) (bool, error) {
    var gotFrame C.int
    if C.avcodec_decode_audio4(cdcctx, &frame.avframe, &gotFrame, packet) < 0 {
        return false, errors.New("error decoding audio")
    }
    if gotFrame != 0 {
        frame.PTS = int64(C.av_frame_get_best_effort_timestamp(&frame.avframe))
    }
    return gotFrame != 0, nil
}

func avcodec_decode_video(cdcctx *C.AVCodecContext, frame *Frame, packet *C.AVPacket) (bool, error) {
    var gotFrame C.int
    if C.avcodec_decode_video2(cdcctx, &frame.avframe, &gotFrame, packet) < 0 {
        return false, errors.New("error decoding video")
    }
    if gotFrame != 0 {
        frame.PTS = int64(C.av_frame_get_best_effort_timestamp(&frame.avframe))
    }
    return gotFrame != 0, nil
}
