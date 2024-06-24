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
Language: Cpp
BasedOnStyle: webkit
AccessModifierOffset: -4

AlignAfterOpenBracket: Align
AlignConsecutiveAssignments: true 
AlignConsecutiveDeclarations: true

AlignConsecutiveMacros: true

AlignEscapedNewlines: Left
AlignOperands: true

AlignTrailingComments: true

AllowAllArgumentsOnNextLine: false
AllowAllParametersOfDeclarationOnNextLine: false

AllowShortBlocksOnASingleLine: false
AllowShortCaseLabelsOnASingleLine: false
AllowShortFunctionsOnASingleLine: Inline
AllowShortIfStatementsOnASingleLine: false
AllowShortLoopsOnASingleLine: false
AlwaysBreakAfterReturnType: TopLevel
AlwaysBreakBeforeMultilineStrings: false

AlwaysBreakTemplateDeclarations: true #false

BinPackArguments: false
BinPackParameters: false

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
BreakConstructorInitializers: AfterColon
BreakStringLiterals: false

ColumnLimit: 0
CommentPragmas: '^begin_wpp|^end_wpp|^FUNC |^USESUFFIX |^USESUFFIX '

ConstructorInitializerAllOnOneLineOrOnePerLine: true
ConstructorInitializerIndentWidth: 4
ContinuationIndentWidth: 4
Cpp11BracedListStyle: true

DerivePointerAlignment: false
ExperimentalAutoDetectBinPacking: false

IndentCaseLabels: false
IndentPPDirectives: AfterHash
IndentWidth: 4

KeepEmptyLinesAtTheStartOfBlocks: false
Language: Cpp

MacroBlockBegin: '^BEGIN_MODULE$|^BEGIN_TEST_CLASS$|^BEGIN_TEST_METHOD$'
MacroBlockEnd: '^END_MODULE$|^END_TEST_CLASS$|^END_TEST_METHOD$'

MaxEmptyLinesToKeep: 1
NamespaceIndentation: None #All
PointerAlignment: Middle
ReflowComments: true
SortIncludes: false

SpaceAfterCStyleCast: false
SpaceBeforeAssignmentOperators: true
SpaceBeforeCtorInitializerColon: true
SpaceBeforeCtorInitializerColon: true
SpaceBeforeParens: ControlStatements
SpaceBeforeRangeBasedForLoopColon: true
SpaceInEmptyParentheses: false
SpacesInAngles: false
SpacesInCStyleCastParentheses: false
SpacesInParentheses: false
SpacesInSquareBrackets: false

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
    '_Guarded_by_',
    '__drv_preferredFunction',
    '__drv_allocatesMem',
    '__drv_freesMem',
    ]
    
TabWidth: '4'
UseTab: Never
`
