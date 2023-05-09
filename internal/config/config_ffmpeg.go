package config

import (
	"fmt"

	"github.com/photoprism/photoprism/internal/ffmpeg"
)

// FFmpegBin returns the ffmpeg executable file name.
func (c *Config) FFmpegBin() string {
	return findBin(c.options.FFmpegBin, "ffmpeg")
}

// FFmpegEnabled checks if FFmpeg is enabled for video transcoding.
func (c *Config) FFmpegEnabled() bool {
	return !c.DisableFFmpeg()
}

// FFmpegEncoder returns the FFmpeg AVC encoder name.
func (c *Config) FFmpegEncoder() ffmpeg.AvcEncoder {
	if c.options.FFmpegEncoder == "" || c.options.FFmpegEncoder == ffmpeg.SoftwareEncoder.String() {
		return ffmpeg.SoftwareEncoder
	} else if c.NoSponsor() {
		log.Infof("ffmpeg: hardware transcoding is available to members only")
		return ffmpeg.SoftwareEncoder
	}

	return ffmpeg.FindEncoder(c.options.FFmpegEncoder)
}

// FFmpegBitrate returns the ffmpeg bitrate limit in MBit/s.
func (c *Config) FFmpegBitrate() int {
	switch {
	case c.options.FFmpegBitrate <= 0:
		return 50
	case c.options.FFmpegBitrate >= 960:
		return 960
	default:
		return c.options.FFmpegBitrate
	}
}

// FFmpegBitrateExceeded tests if the ffmpeg bitrate limit is exceeded.
func (c *Config) FFmpegBitrateExceeded(mbit float64) bool {
	if mbit <= 0 {
		return false
	} else if max := c.FFmpegBitrate(); max <= 0 {
		return false
	} else {
		return mbit > float64(max)
	}
}

// FFmpegMapVideo returns the video streams to be transcoded as string.
func (c *Config) FFmpegMapVideo() string {
	if c.options.FFmpegMapVideo == "" {
		return ffmpeg.MapVideoDefault
	}

	return c.options.FFmpegMapVideo
}

// FFmpegMapAudio returns the audio streams to be transcoded as string.
func (c *Config) FFmpegMapAudio() string {
	if c.options.FFmpegMapAudio == "" {
		return ffmpeg.MapAudioDefault
	}

	return c.options.FFmpegMapAudio
}

// FFmpegOptions returns the FFmpeg transcoding options.
func (c *Config) FFmpegOptions(encoder ffmpeg.AvcEncoder, bitrate string) (ffmpeg.Options, error) {
	// Transcode all other formats with FFmpeg.
	opt := ffmpeg.Options{
		Bin:      c.FFmpegBin(),
		Encoder:  encoder,
		Bitrate:  bitrate,
		MapVideo: c.FFmpegMapVideo(),
		MapAudio: c.FFmpegMapAudio(),
	}

	// Check
	if opt.Bin == "" {
		return opt, fmt.Errorf("ffmpeg is not installed")
	} else if c.DisableFFmpeg() {
		return opt, fmt.Errorf("ffmpeg is disabled")
	} else if bitrate == "" {
		return opt, fmt.Errorf("bitrate must not be empty")
	} else if encoder.String() == "" {
		return opt, fmt.Errorf("encoder must not be empty")
	}

	return opt, nil
}
