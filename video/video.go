package video

import (
	"fmt"
	"io"

	"github.com/kkdai/youtube/v2"
)

type Video struct {
	Filename string
	Stream *io.ReadCloser
	
}

func (v *Video) NameFile(name string) {
	v.Filename = name + ".mp4"
}

type Client struct {
	youtube.Client
}

func (c *Client) Stream(Id string) (*Video, error) {
	// get stream of bytes by video id

	video, err := c.GetVideo(Id)
	if err != nil {
		return nil, fmt.Errorf("cannot get id %s: %w", Id, err)
	}

	format, err := BestQuality(video.Formats) // only get videos with audio
	if err != nil {
		return nil, fmt.Errorf("cannot find best quality id %s: %w", Id, err)
	}

	stream, _, err := c.GetStream(video, format)
	if err != nil {
		return nil, fmt.Errorf("cannot get stream of id %s: %w", Id, err)
	}

	v := new(Video)
	v.Stream = &stream
	v.NameFile(Id)

	return v, nil
}

func BestQuality(list youtube.FormatList) (*youtube.Format, error) {
	better := new(youtube.Format)
	for i := range list {
		if list[i].AudioChannels > 0 && better.Height < list[i].Height {
			better = &list[i]
		}
	}
	if better.Height == 0 {
		return nil, ErrNoBestQuality
	}
	return better, nil
}