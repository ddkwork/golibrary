package clang

import (
	"bytes"
	_ "embed"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

type  object struct{
	once sync.Once
}
func New() *object { return &object{} }

func (o *object) writeClangFormatBody(rootPath string) {
	join := filepath.Join(rootPath, ".clang-format")
	stream.WriteTruncate(join, clangFormatBody)
}

func Walk(root string)  {
	New().Walk(root)
}

func (o *object) Walk(root string) {
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		switch filepath.Ext(path) {
		case ".h", ".c", ".cpp":
			o.Format(path)
		}
		return err
	})
}
func (o *object) Format(absPath string) {
	o.once.Do(func() {
		o.writeClangFormatBody(filepath.Dir(absPath))
	})

	g := removeCppComments(mylog.Check2(os.ReadFile(absPath)))
	stream.WriteTruncate(absPath, g.String())

	if strings.Contains(absPath, `\`) {
		absPath = strings.ReplaceAll(absPath, `\`, `\\`)
	}
	command := "clang-format -i --style=file " + absPath
	stream.RunCommand(command)
}

////////////////////////
// 核心修复：确保函数返回类型和签名在同一行
func fixFunctionSignatures(content string) string {
	// 正则模式匹配返回类型单独一行的情况
	pattern := regexp.MustCompile(`(\w[\w\s:]*?\s*[*&]*)\s*\r?\n\s*(\w[\w:]*\s*\([^)]*\)\s*(\{|;))`)

	// 替换为同行的形式
	result := pattern.ReplaceAllString(content, "$1 $2")

	// 处理没有参数的函数
	pattern2 := regexp.MustCompile(`(\w[\w\s:]*?\s*[*&]*)\s*\r?\n\s*(\w[\w:]*\s*\(\s*\)\s*(\{|;))`)
	result = pattern2.ReplaceAllString(result, "$1 $2")

	return result
}

func removeComments2(code string) {
	g := stream.NewGeneratedFile()
	skip := false
	for s := range strings.Lines(code) {
		if strings.HasPrefix(strings.TrimSpace(s), "/*") {
			skip = true
		}
		if strings.TrimSpace(s) == "*/" {
			skip = false
			continue
		}
		if skip {
			continue
		}
		g.P(s)
	}
}

func removeCppComments(source []byte) bytes.Buffer {
	removeComments2("") //working
	var out bytes.Buffer
	const (
		stateCode = iota
		stateLineComment
		stateBlockComment
		stateString
		stateChar
	)

	state := stateCode
	prev := byte(0)

	for _, ch := range source {
		switch state {
		case stateCode:
			switch {
			case ch == '/' && prev == '/':
				out.Truncate(out.Len() - 1)
				state = stateLineComment
			case ch == '*' && prev == '/':
				out.Truncate(out.Len() - 1)
				state = stateBlockComment
			case ch == '"':
				state = stateString
				out.WriteByte(ch)
			case ch == '\'':
				state = stateChar
				out.WriteByte(ch)
			default:
				out.WriteByte(ch)
			}

		case stateLineComment:
			if ch == '\n' {
				state = stateCode
				out.WriteByte(ch)
			}

		case stateBlockComment:
			if ch == '/' && prev == '*' {
				state = stateCode
				prev = 0
				continue
			}

		case stateString:
			out.WriteByte(ch)
			if ch == '"' && prev != '\\' {
				state = stateCode
			}

		case stateChar:
			out.WriteByte(ch)
			if ch == '\'' && prev != '\\' {
				state = stateCode
			}
		}

		// Update prev for all states
		if state == stateString || state == stateChar {
			if ch == '\\' {
				prev = 0
			} else {
				prev = ch
			}
		} else {
			prev = ch
		}
	}
	s := out.String()
	s = fixFunctionSignatures(s)
	out.Reset()
	out.WriteString(s)
	return out
}


const clangFormatBody = `

# Generated from CLion C/C++ Code Style settings
#Language: Cpp
BasedOnStyle: LLVM
AccessModifierOffset: -4
AlignAfterOpenBracket: Align
AlignConsecutiveAssignments: None
AlignOperands: Align
AllowAllConstructorInitializersOnNextLine: false
AllowShortCaseLabelsOnASingleLine: false
AllowShortIfStatementsOnASingleLine: Always
AllowShortLambdasOnASingleLine: All
AllowShortLoopsOnASingleLine: true
#AlwaysBreakTemplateDeclarations: Yes
BreakBeforeBraces: Custom
BraceWrapping:
  AfterCaseLabel: false
  AfterClass: false
  AfterControlStatement: Never
  AfterEnum: false
  AfterFunction: false
  AfterNamespace: false
  AfterUnion: false
  BeforeCatch: false
  BeforeElse: false
  IndentBraces: false
  SplitEmptyFunction: false
  SplitEmptyRecord: true
BreakBeforeBinaryOperators: None
BreakBeforeTernaryOperators: true
BreakConstructorInitializers: BeforeColon
BreakInheritanceList: BeforeColon
CompactNamespaces: false
ContinuationIndentWidth: 8
IndentCaseLabels: true
IndentPPDirectives: None
IndentWidth: 4
NamespaceIndentation: All
ObjCSpaceAfterProperty: false
ObjCSpaceBeforeProtocolList: true
PointerAlignment: Right
ReflowComments: false
SpaceAfterCStyleCast: true
SpaceAfterLogicalNot: false
SpaceAfterTemplateKeyword: false
SpaceBeforeAssignmentOperators: true
SpaceBeforeCpp11BracedList: false
SpaceBeforeCtorInitializerColon: true
SpaceBeforeInheritanceColon: true
SpaceBeforeParens: ControlStatements
SpaceBeforeRangeBasedForLoopColon: false
SpaceInEmptyParentheses: false
SpacesBeforeTrailingComments: 0
SpacesInAngles: false
SpacesInCStyleCastParentheses: false
SpacesInContainerLiterals: false
SpacesInParentheses: false
SpacesInSquareBrackets: false
AlignConsecutiveDeclarations: true
AlignConsecutiveMacros: true
AlignEscapedNewlines: Left
AlignTrailingComments: true
AlwaysBreakBeforeMultilineStrings: false
BinPackArguments: false
BreakStringLiterals: false
CommentPragmas: '^begin_wpp|^end_wpp|^FUNC |^USESUFFIX |^USESUFFIX '
ConstructorInitializerAllOnOneLineOrOnePerLine: true
ConstructorInitializerIndentWidth: 4
Cpp11BracedListStyle: true
DerivePointerAlignment: false
ExperimentalAutoDetectBinPacking: false
SortIncludes: false
MacroBlockBegin: '^BEGIN_MODULE$|^BEGIN_TEST_CLASS$|^BEGIN_TEST_METHOD$'
MacroBlockEnd: '^END_MODULE$|^END_TEST_CLASS$|^END_TEST_METHOD$'

Standard: Cpp11
StatementMacros: [
  'EXTERN_C',
  'PAGED',
  'PAGEDX',
  'NONPAGED',
  'PNPCODE',
  'INITCODE',
  '_At_',
  '_When_',
  '_Success_',
  '_Check_return_',
  '_Must_inspect_result_',
  '_IRQL_requires_same_',
  '_IRQL_requires_',
  '_IRQL_requires_max_',
  '_IRQL_requires_min_',
  '_IRQL_saves_',
  '_IRQL_restores_',
  '_IRQL_saves_global_',
  '_IRQL_restores_global_',
  '_IRQL_raises_',
  '_IRQL_lowers_',
  '_Acquires_lock_',
  '_Releases_lock_',
  '_Acquires_exclusive_lock_',
  '_Releases_exclusive_lock_',
  '_Acquires_shared_lock_',
  '_Releases_shared_lock_',
  '_Requires_lock_held_',
  '_Use_decl_annotations_',
  'DEBUGLOGGER_INIT',
  'DEBUGLOGGER_PAGED',
  '_Guarded_by_',
  '__drv_preferredFunction',
  '__drv_allocatesMem',
  '__drv_freesMem',
]

# 控制空行
AllowAllParametersOfDeclarationOnNextLine: false  # 参数较多时换行方式
AlwaysBreakAfterDefinitionReturnType: None        # 禁止返回类型后换行

# 核心空行控制配置
MaxEmptyLinesToKeep: 1                  # 全局最多保留连续1个空行
KeepEmptyLinesAtTheStartOfBlocks: false   # 删除函数体开头的空行
SeparateDefinitionBlocks: Always         # 强制函数定义间加空行
ColumnLimit: 0                  # 禁用行宽限制，防止自动换行
BreakAfterReturnType: None    # Clang-Format 10+ 替代选项
TypeNames: [Param]             # 类型名单独一行
BreakAfterReturnType: None
AllowShortFunctionsOnASingleLine: Empty  # 允许空函数单行显示
AllowShortBlocksOnASingleLine: Empty     # 允许空代码块单行显示
BinPackParameters: true
AllowAllArgumentsOnNextLine: false
PenaltyReturnTypeOnItsOwnLine: 10000
TabWidth: 4
UseTab: Never


`
