package ffmpeg

/*
#include <libavformat/avformat.h>
*/
import "C"

type Frame struct {
    frame C.AVFrame
    PTS int64
}

func (frame *Frame) Defaults() {
    C.avcodec_get_frame_defaults(&frame.frame)
}
