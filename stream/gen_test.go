package stream_test

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/mylog"

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
	mylog.FormatAllFiles()
	m := stream.NewOrderedMap("", "")
	m.Set("SuperRecovery4", "SuperRecovery4")
	m.Set("AneData6", "AneData6")
	m.Set("ChaoQiangZhaoPian", "ChaoQiangZhaoPian")
	m.Set("DataExplore", "DataExplore")
	m.Set("DiskGetor", "DiskGetor")
	m.Set("DocxBuild", "DocxBuild")
	m.Set("ExcelHuiFuZhuanJia", "ExcelHuiFuZhuanJia")
	m.Set("ExcelRebuild1", "ExcelRebuild1")
	m.Set("ExcelRebuild3", "ExcelRebuild3")
	m.Set("ExcelScan", "ExcelScan")
	m.Set("HykHf", "HykHf")
	m.Set("Jphf3", "Jphf3")
	m.Set("MailScan", "MailScan")
	m.Set("OfficeBuild", "OfficeBuild")
	m.Set("OfficeZipBuild", "OfficeZipBuild")
	m.Set("PptxBuild", "PptxBuild")
	m.Set("SuperRecovery2", "SuperRecovery2")
	m.Set("WordRepairer", "WordRepairer")
	m.Set("WordScan", "WordScan")
	m.Set("XlsxBuild", "XlsxBuild")
	m.Set("ZipBuild", "ZipBuild")
	stream.NewGeneratedFile().Types("app", m)
}
