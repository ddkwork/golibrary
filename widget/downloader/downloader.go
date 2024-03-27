package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ddkwork/golibrary/mylog"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64)
}

func (pw *progressWriter) Start() {
	_, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw))
	if err != nil {
		if p != nil {
			p.Send(progressErrMsg{err})
		}
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

func getResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp == nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}

func Run(url string) (ok bool) {
	resp, err := getResponse(url)
	if !mylog.Error(err) {
		return
	}
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	filename := filepath.Base(url)
	file, err := os.Create(filename)
	if !mylog.Error(err) {
		return
	}
	defer file.Close()
	pw := &progressWriter{
		total:  int(resp.ContentLength),
		file:   file,
		reader: resp.Body,
		onProgress: func(ratio float64) {
			if p != nil {
				p.Send(progressMsg(ratio))
			}
		},
	}
	m := model{
		pw:       pw,
		progress: progress.New(progress.WithDefaultGradient()),
	}
	go pw.Start()
	if resp.ContentLength > 0 {
		p = tea.NewProgram(m)
		return mylog.Error(p.Start())
	}
	return
}
