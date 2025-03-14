//go:build !linux
// +build !linux

package screen

import (
	"fmt"
	"image"
	"io"

	"github.com/0xbytejay/mediadevices/pkg/driver"
	"github.com/0xbytejay/mediadevices/pkg/frame"
	"github.com/0xbytejay/mediadevices/pkg/io/video"
	"github.com/0xbytejay/mediadevices/pkg/prop"
	"github.com/kbinani/screenshot"
)

type screen struct {
	displayIndex int
	doneCh       chan struct{}
}

func init() {
	Initialize()
}

// Initialize finds and registers active displays. This is part of an experimental API.
func Initialize() {
	activeDisplays := screenshot.NumActiveDisplays()
	for i := 0; i < activeDisplays; i++ {
		priority := driver.PriorityNormal
		if i == 0 {
			priority = driver.PriorityHigh
		}

		s := newScreen(i)
		driver.GetManager().Register(s, driver.Info{
			Label:      fmt.Sprint(i),
			DeviceType: driver.Screen,
			Priority:   priority,
		})
	}
}

func newScreen(displayIndex int) *screen {
	s := screen{
		displayIndex: displayIndex,
	}
	return &s
}

func (s *screen) Open() error {
	s.doneCh = make(chan struct{})
	return nil
}

func (s *screen) Close() error {
	close(s.doneCh)
	return nil
}

func (s *screen) VideoRecord(selectedProp prop.Media) (video.Reader, error) {
	r := video.ReaderFunc(func() (img image.Image, release func(), err error) {
		select {
		case <-s.doneCh:
			return nil, nil, io.EOF
		default:
		}

		img, err = screenshot.CaptureDisplay(s.displayIndex)
		release = func() {}
		return
	})
	return r, nil
}

func (s *screen) Properties() []prop.Media {
	resolution := screenshot.GetDisplayBounds(s.displayIndex)
	supportedProp := prop.Media{
		Video: prop.Video{
			Width:       resolution.Dx(),
			Height:      resolution.Dy(),
			FrameFormat: frame.FormatRGBA,
		},
	}
	return []prop.Media{supportedProp}
}
