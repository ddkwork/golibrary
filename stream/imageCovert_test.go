package stream

import "testing"

func TestIco2PngAll(t *testing.T) {
	Ico2PngAll("ICON")
	Ico2PngAll("D:\\workspace\\workspace\\gui\\pageIco")
}

func TestEmbedImage(t *testing.T) {
	Png2SvgAll("../ICON")
	Png2SvgAll("D:\\workspace\\workspace\\gui\\pageIco")
}
