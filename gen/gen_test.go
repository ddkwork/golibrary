package gen

import (
	"testing"
)

func TestGeneratedFile_Iota(t *testing.T) {
	New().Enum("ki", []string{
		"SliceViewBase",
		"WidgetBase",
	}, nil)
	New().Enum("app", []string{
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
	New().FileAction()
}

func TestGeneratedFile_ReadTemplates(t *testing.T) {
	t.Skip()
	New().ReadTemplates("fileAction.go")
}
