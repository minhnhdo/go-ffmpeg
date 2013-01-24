package ffmpeg

/*
#cgo pkg-config: libavformat libavcodec
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
*/
import "C"

func init() {
    C.av_register_all()
    C.avcodec_register_all()
}
