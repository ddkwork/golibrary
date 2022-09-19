package gzip

import (
	"bytes"
	"compress/gzip"
	"github.com/ddkwork/golibrary/mylog"

	"github.com/ddkwork/golibrary/src/stream"
	"io/ioutil"
)

type (
	Interface interface {
		Decode(in []byte) *stream.Stream
	}
	object struct{ s *stream.Stream }
)

func New() Interface { return &object{s: stream.NewNil()} }

func (o *object) Decode(in []byte) *stream.Stream {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if !mylog.Error(err) {
		return nil
	}
	defer func() {
		if reader == nil {
			mylog.Error("gzipReader == nil")
			return
		}
		mylog.Error(reader.Close())
	}()
	all, err2 := ioutil.ReadAll(reader)
	if !mylog.Error(err2) {
		return stream.NewErrorInfo(err2.Error())
	}
	return stream.NewBytes(all)
}
