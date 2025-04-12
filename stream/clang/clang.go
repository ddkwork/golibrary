package clang

import (
	_ "embed"
	"path/filepath"
	"strings"

	"github.com/ddkwork/golibrary/stream"
)

type (
	Interface interface {
		WriteClangFormatBody(rootPath string)
		Format(absPath string)
	}
	object struct{}
)

func (o *object) WriteClangFormatBody(rootPath string) {
	join := filepath.Join(rootPath, ".clang-format")
	stream.WriteTruncate(join, clangFormatBody)
}

func (o *object) Format(absPath string) {
	if strings.Contains(absPath, `\`) {
		absPath = strings.ReplaceAll(absPath, `\`, `\\`)
	}
	command := "clang-format -i --style=file " + absPath
	stream.RunCommand(command)
}

func New() Interface { return &object{} }

var clangFormatBody = `
# Generated from CLion C/C++ Code Style settings
#Language: Cpp
BasedOnStyle: LLVM
AccessModifierOffset: -4
AlignAfterOpenBracket: Align
AlignConsecutiveAssignments: None
AlignOperands: Align
AllowAllConstructorInitializersOnNextLine: false
AllowShortBlocksOnASingleLine: Always
AllowShortCaseLabelsOnASingleLine: false
AllowShortFunctionsOnASingleLine: All
AllowShortIfStatementsOnASingleLine: Always
AllowShortLambdasOnASingleLine: All
AllowShortLoopsOnASingleLine: true
AlwaysBreakTemplateDeclarations: Yes
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
KeepEmptyLinesAtTheStartOfBlocks: true
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
MaxEmptyLinesToKeep: 0         # 函数体内最多保留 0 空行，函数间保留 1 空行
SeparateDefinitionBlocks: Always  # 强制函数间空行（需 clang-format 14+）

# 参数在同一行
ColumnLimit: 180
AllowAllParametersOfDeclarationOnNextLine: false
BinPackParameters: true
AllowAllArgumentsOnNextLine: false

# 强制返回类型与函数名同行
AlwaysBreakAfterDefinitionReturnType: None
AlwaysBreakAfterReturnType: None
PenaltyReturnTypeOnItsOwnLine: 1000000

TabWidth: 4
UseTab: Never

`
