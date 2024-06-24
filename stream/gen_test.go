package stream_test

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	stream.NewGeneratedFile().ReadTemplates("safeStream.go", "stream")
	generateIR("../mylog/safeStream.gen.go", func(b *stream.Buffer) {
		b.ReplaceAll("package stream", "package mylog")
		b.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/mylog"), "")
		b.ReplaceAll("mylog.", "")
	})
}

func TestGeneratedFile_Iota(t *testing.T) {
	stream.NewGeneratedFile().Enum("app", []string{
		"SuperRecovery4",
		"AneData6",
		"ChaoQiangZhaoPian",
		"DataExplore",
		"DiskGetor",
		"DocxBuild",
		"ExcelHuiFuZhuanJia",
		"ExcelRebuild1",
		"ExcelRebuild3",
		"ExcelScan",
		"HykHf",
		"Jphf3",
		"MailScan",
		"OfficeBuild",
		"OfficeZipBuild",
		"PptxBuild",
		"SuperRecovery2",
		"WordRepairer",
		"WordScan",
		"XlsxBuild",
		"ZipBuild",
	}, nil)
}
