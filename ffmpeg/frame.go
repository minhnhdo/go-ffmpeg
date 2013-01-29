package ffmpeg

/*
#include <libavformat/avformat.h>
*/
import "C"

type Frame struct {
    PTS int64
    avframe C.AVFrame
}

func (frame *Frame) Defaults() {
    C.avcodec_get_frame_defaults(&frame.avframe)
}
