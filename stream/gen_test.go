package stream_test

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	g := stream.NewGeneratedFile()
	g.Write(stream.NewBuffer("safeStream.go").Bytes())
	g.ReplaceAll("package stream", "package mylog")
	g.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/mylog"), "")
	g.ReplaceAll("mylog.", "")
	stream.WriteTruncate("../mylog/safeStream.gen.go", g.Buffer)
}

func TestGeneratedFile_Iota(t *testing.T) {
	stream.NewGeneratedFile().Types("app", []string{
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
