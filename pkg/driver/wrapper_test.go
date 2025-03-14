package driver

import (
	"fmt"
	"testing"

	"github.com/0xbytejay/mediadevices/pkg/io/audio"
	"github.com/0xbytejay/mediadevices/pkg/io/video"
	"github.com/0xbytejay/mediadevices/pkg/prop"
)

var (
	recordErr = fmt.Errorf("failed to start recording")
)

type adapterMock struct{}

func (a *adapterMock) Open() error              { return nil }
func (a *adapterMock) Close() error             { return nil }
func (a *adapterMock) Properties() []prop.Media { return []prop.Media{prop.Media{}} }

type videoAdapterMock struct{ adapterMock }

func (a *videoAdapterMock) VideoRecord(p prop.Media) (r video.Reader, err error) { return nil, nil }

type videoAdapterBrokenMock struct{ adapterMock }

func (a *videoAdapterBrokenMock) VideoRecord(p prop.Media) (r video.Reader, err error) {
	return nil, recordErr
}

type audioAdapterMock struct{ adapterMock }

func (a *audioAdapterMock) AudioRecord(p prop.Media) (r audio.Reader, err error) { return nil, nil }

type audioAdapterBrokenMock struct{ adapterMock }

func (a *audioAdapterBrokenMock) AudioRecord(p prop.Media) (r audio.Reader, err error) {
	return nil, recordErr
}

type availabilityAdapterMock struct{ videoAdapterMock }

func (a *availabilityAdapterMock) IsAvailable() (bool, error) { return true, nil }

func TestVideoWrapperState(t *testing.T) {
	var a videoAdapterMock
	d := wrapAdapter(&a, Info{})

	if d.Properties() != nil {
		t.Errorf("expected nil, but got %v", d.Properties())
	}

	vr := d.(VideoRecorder)
	_, err := vr.VideoRecord(prop.Media{})
	if err == nil {
		t.Errorf("expected to get an invalid state")
	}

	err = d.Open()
	if err != nil {
		t.Errorf("expected to successfully open, but got %v", err)
	}

	_, err = vr.VideoRecord(prop.Media{})
	if err != nil {
		t.Errorf("expected to successfully start recording, but got %v", err)
	}
}

func TestVideoWrapperWithBrokenRecorderState(t *testing.T) {
	var a videoAdapterBrokenMock
	d := wrapAdapter(&a, Info{})

	err := d.Open()
	if err != nil {
		t.Errorf("expected to open successfully")
	}

	vr := d.(VideoRecorder)
	_, err = vr.VideoRecord(prop.Media{})
	if err == nil {
		t.Errorf("expected to get an error")
	}

	if err != recordErr {
		t.Errorf("expected to get %v, but got %v", recordErr, err)
	}

	if d.Status() != StateClosed {
		t.Errorf("expected the status to be %v, but got %v", StateClosed, d.Status())
	}
}

func TestAudioWrapperState(t *testing.T) {
	var a audioAdapterMock
	d := wrapAdapter(&a, Info{})

	if d.Properties() != nil {
		t.Errorf("expected nil, but got %v", d.Properties())
	}

	ar := d.(AudioRecorder)
	_, err := ar.AudioRecord(prop.Media{})
	if err == nil {
		t.Errorf("expected to get an invalid state")
	}

	err = d.Open()
	if err != nil {
		t.Errorf("expected to successfully open, but got %v", err)
	}

	_, err = ar.AudioRecord(prop.Media{})
	if err != nil {
		t.Errorf("expected to successfully start recording, but got %v", err)
	}
}

func TestAudioWrapperWithBrokenRecorderState(t *testing.T) {
	var a audioAdapterBrokenMock
	d := wrapAdapter(&a, Info{})

	err := d.Open()
	if err != nil {
		t.Errorf("expected to open successfully")
	}

	ar := d.(AudioRecorder)
	_, err = ar.AudioRecord(prop.Media{})
	if err == nil {
		t.Errorf("expected to get an error")
	}

	if err != recordErr {
		t.Errorf("expected to get %v, but got %v", recordErr, err)
	}

	if d.Status() != StateClosed {
		t.Errorf("expected the status to be %v, but got %v", StateClosed, d.Status())
	}
}

func TestWrapperAvailabilityAdapter(t *testing.T) {
	var aa availabilityAdapterMock
	d := wrapAdapter(&aa, Info{})

	ok, err := IsAvailable(d)
	if err != nil {
		t.Errorf("expected nil, but got %v", err)
	}
	if !ok {
		t.Errorf("expected true, but got %v", ok)
	}

	var v videoAdapterMock
	d = wrapAdapter(&v, Info{})

	ok, err = IsAvailable(d)
	if err == nil {
		t.Errorf("expected err, but got %v", err)
	}
	if ok {
		t.Errorf("expected false, but got %v", ok)
	}

	var a audioAdapterMock
	d = wrapAdapter(&a, Info{})

	ok, err = IsAvailable(d)
	if err == nil {
		t.Errorf("expected err, but got %v", err)
	}
	if ok {
		t.Errorf("expected false, but got %v", ok)
	}
}
