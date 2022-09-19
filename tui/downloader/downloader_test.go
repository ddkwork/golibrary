package downloader_test

import (
	"github.com/ddkwork/golibrary/tui/downloader"
	"testing"
)

func TestDownloader(t *testing.T) {
	return
	downloader.Run("https://github.com/charmbracelet/bubbletea/archive/refs/tags/v0.22.0.zip")
}
