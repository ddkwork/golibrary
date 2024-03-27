package downloader_test

import (
	"testing"

	"github.com/ddkwork/golibrary/widget/downloader"
)

func TestDownloader(t *testing.T) {
	return
	downloader.Run("https://github.com/charmbracelet/bubbletea/archive/refs/tags/v0.22.0.zip")
}
