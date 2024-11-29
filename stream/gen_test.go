package stream_test

import (
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	cloneStreamForLogPkg("safeStream.go")
	cloneStreamForLogPkg("constraints.go")

	path := "constraints.go"
	g := stream.NewGeneratedFile()
	g.Write(stream.NewBuffer(path).Bytes())
	g.ReplaceAll("package stream", "package pretty")
	g.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/mylog"), "")
	g.ReplaceAll("mylog.", "")
	stream.WriteTruncate("../mylog/pretty/"+stream.BaseName(path)+"_gen.go", g.Buffer)
}

func cloneStreamForLogPkg(path string) {
	g := stream.NewGeneratedFile()
	g.Write(stream.NewBuffer(path).Bytes())
	g.ReplaceAll("package stream", "package mylog")
	g.ReplaceAll(strconv.Quote("github.com/ddkwork/golibrary/mylog"), "")
	g.ReplaceAll("mylog.", "")
	stream.WriteTruncate("../mylog/"+stream.BaseName(path)+"_gen.go", g.Buffer)
}

func TestGeneratedFile_enum(t *testing.T) {
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
	stream.NewGeneratedFile().EnumTypes("app", m)
}
