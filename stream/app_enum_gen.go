package stream

import (
	"strings"

	"golang.org/x/exp/constraints"
)

// Code generated by GeneratedFile enum - DO NOT EDIT.

type AppKind byte

const (
	SuperRecovery4Kind AppKind = iota
	AneData6Kind
	ChaoQiangZhaoPianKind
	DataExploreKind
	DiskGetorKind
	DocxBuildKind
	ExcelHuiFuZhuanJiaKind
	ExcelRebuild1Kind
	ExcelRebuild3Kind
	ExcelScanKind
	HykHfKind
	Jphf3Kind
	MailScanKind
	OfficeBuildKind
	OfficeZipBuildKind
	PptxBuildKind
	SuperRecovery2Kind
	WordRepairerKind
	WordScanKind
	XlsxBuildKind
	ZipBuildKind
	InvalidAppKind
)

func ConvertInteger2AppKind[T constraints.Integer](v T) AppKind {
	return AppKind(v)
}

func (k AppKind) AssertKind(kinds string) AppKind {
	for _, kind := range k.Kinds() {
		if strings.ToLower(kinds) == strings.ToLower(kind.String()) {
			return kind
		}
	}
	return InvalidAppKind
}

func (k AppKind) String() string {
	switch k {
	case SuperRecovery4Kind:
		return "SuperRecovery4"
	case AneData6Kind:
		return "AneData6"
	case ChaoQiangZhaoPianKind:
		return "ChaoQiangZhaoPian"
	case DataExploreKind:
		return "DataExplore"
	case DiskGetorKind:
		return "DiskGetor"
	case DocxBuildKind:
		return "DocxBuild"
	case ExcelHuiFuZhuanJiaKind:
		return "ExcelHuiFuZhuanJia"
	case ExcelRebuild1Kind:
		return "ExcelRebuild1"
	case ExcelRebuild3Kind:
		return "ExcelRebuild3"
	case ExcelScanKind:
		return "ExcelScan"
	case HykHfKind:
		return "HykHf"
	case Jphf3Kind:
		return "Jphf3"
	case MailScanKind:
		return "MailScan"
	case OfficeBuildKind:
		return "OfficeBuild"
	case OfficeZipBuildKind:
		return "OfficeZipBuild"
	case PptxBuildKind:
		return "PptxBuild"
	case SuperRecovery2Kind:
		return "SuperRecovery2"
	case WordRepairerKind:
		return "WordRepairer"
	case WordScanKind:
		return "WordScan"
	case XlsxBuildKind:
		return "XlsxBuild"
	case ZipBuildKind:
		return "ZipBuild"
	default:
		return "InvalidAppKind"
	}
}

func (k AppKind) Keys() []string {
	return []string{
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
	}
}

func (k AppKind) Kinds() []AppKind {
	return []AppKind{
		SuperRecovery4Kind,
		AneData6Kind,
		ChaoQiangZhaoPianKind,
		DataExploreKind,
		DiskGetorKind,
		DocxBuildKind,
		ExcelHuiFuZhuanJiaKind,
		ExcelRebuild1Kind,
		ExcelRebuild3Kind,
		ExcelScanKind,
		HykHfKind,
		Jphf3Kind,
		MailScanKind,
		OfficeBuildKind,
		OfficeZipBuildKind,
		PptxBuildKind,
		SuperRecovery2Kind,
		WordRepairerKind,
		WordScanKind,
		XlsxBuildKind,
		ZipBuildKind,
	}
}

func (k AppKind) PngFileName() string {
	switch k {
	case SuperRecovery4Kind:
		return "SuperRecovery4"
	case AneData6Kind:
		return "AneData6"
	case ChaoQiangZhaoPianKind:
		return "ChaoQiangZhaoPian"
	case DataExploreKind:
		return "DataExplore"
	case DiskGetorKind:
		return "DiskGetor"
	case DocxBuildKind:
		return "DocxBuild"
	case ExcelHuiFuZhuanJiaKind:
		return "ExcelHuiFuZhuanJia"
	case ExcelRebuild1Kind:
		return "ExcelRebuild1"
	case ExcelRebuild3Kind:
		return "ExcelRebuild3"
	case ExcelScanKind:
		return "ExcelScan"
	case HykHfKind:
		return "HykHf"
	case Jphf3Kind:
		return "Jphf3"
	case MailScanKind:
		return "MailScan"
	case OfficeBuildKind:
		return "OfficeBuild"
	case OfficeZipBuildKind:
		return "OfficeZipBuild"
	case PptxBuildKind:
		return "PptxBuild"
	case SuperRecovery2Kind:
		return "SuperRecovery2"
	case WordRepairerKind:
		return "WordRepairer"
	case WordScanKind:
		return "WordScan"
	case XlsxBuildKind:
		return "XlsxBuild"
	case ZipBuildKind:
		return "ZipBuild"
	default:
		return "InvalidAppKind"
	}
}
