package stream

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func testReadFileSubTimeWithCallBack() {
	var (
		m11 = `C:\Users\Admin\go\pkg\mod\fyne.io\fyne\v2@v2.4.4\theme\bundled-emoji.go`
		m62 = `D:\Desktop\packets.json`
	)
	ReadFileSubTimeWithCallBack("m11 os.file.read 1024 byte every time(block)", m11, func(s string) { readBlock(m11) })
	ReadFileSubTimeWithCallBack("m11 os.file.read byte file size once", m11, func(s string) { readAllBuff(m11) })
	ReadFileSubTimeWithCallBack("m11 os.ReadFile", m11, func(s string) { readAll(m11) })
	ReadFileSubTimeWithCallBack("m11 bufio ReadLine", m11, func(s string) { readEachLineReader(m11) })
	ReadFileSubTimeWithCallBack("m11 bufio Scanner ------------> seems fast", m11, func(s string) { readEachLineScanner(m11) })

	ReadFileSubTimeWithCallBack("m62 os.file.read 1024 byte every time(block)", m62, func(s string) { readBlock(m62) })
	ReadFileSubTimeWithCallBack("m62 os.file.read byte file size once", m62, func(s string) { readAllBuff(m62) })
	ReadFileSubTimeWithCallBack("m62 os.ReadFile", m62, func(s string) { readAll(m62) })
	ReadFileSubTimeWithCallBack("m62 bufio ReadLine", m62, func(s string) { readEachLineReader(m62) })
	ReadFileSubTimeWithCallBack("m62 bufio Scanner ------------> seems fast", m62, func(s string) { readEachLineScanner(m62) })
}

func readBlock(filePath string) {
	FileHandle, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer FileHandle.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := FileHandle.Read(buffer)
		if err != nil && err != io.EOF {
			log.Println(err)
		}
		if n == 0 {
			break
		}
	}
}

func readEachLineReader(filePath string) {
	FileHandle, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer FileHandle.Close()
	lines := make([]byte, 0)
	lineReader := bufio.NewReader(FileHandle)
	for {
		// func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
		// func (b *Reader) ReadBytes(delim byte) (line []byte, err error)
		// func (b *Reader) ReadString(delim byte) (line string, err error)
		line, _, err := lineReader.ReadLine()
		if err == io.EOF {
			break
		}
		lines = append(lines, line...)
	}
}

func readEachLineScanner(filePath string) {
	FileHandle, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer FileHandle.Close()
	lineScanner := bufio.NewScanner(FileHandle)
	lines := make([]byte, 0)
	for lineScanner.Scan() {
		// func (s *Scanner) Bytes() []byte
		// func (s *Scanner) Text() string
		lines = append(lines, lineScanner.Bytes()...)
	}
}

func readAll(filePath string) {
	os.ReadFile(filePath)
}

func readAllBuff(filePath string) {
	FileHandle, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer FileHandle.Close()
	fileInfo, err := FileHandle.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = FileHandle.Read(buffer)
	if err != nil {
		log.Println(err)
	}
}

func ReadFileSubTimeWithCallBack(title, path string, callBack func(string)) {
	now := time.Now()
	callBack(path)
	sub := time.Since(now)
	fmt.Printf(" %-48s spend | %s \n", title, sub)
}
