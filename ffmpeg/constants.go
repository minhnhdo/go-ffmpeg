package ffmpeg

/*
#include <libavformat/avformat.h>
*/
import "C"

type MediaType int

const (
    VideoType MediaType = C.AVMEDIA_TYPE_VIDEO
    AudioType MediaType = C.AVMEDIA_TYPE_AUDIO
)
