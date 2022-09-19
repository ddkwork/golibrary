package cpp2go_test

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/cpp2go"

	"github.com/ddkwork/golibrary/src/mydoc"
	"testing"
)

func TestCpp2go(t *testing.T) {
	return
	d := cpp2go.New()
	assert := mylog.Assert(t)
	assert.True(d.Translate("cfile"))
	//assert.True(d.Translate("tt"))
}

func TestDocCpp2go(t *testing.T) {
	doc := mydoc.New()
	doc.Append(mydoc.Row{
		Api:      "TranslateCFile(path, pkg string)",
		Function: "convert c to go", // https://github.com/goplus/c2go.git
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "Translate(root string) (ok bool)",
		Function: "Translate cpp or c to go, not full part convert",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "RemoveComment(root string) (ok bool)",
		Function: "remove comment",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	body := doc.Gen()
	println(body)
}

func TestDocScanner(t *testing.T) {
	doc := mydoc.New()
	doc.Append(mydoc.Row{
		Api:      "Translate(root string) (ok bool)",
		Function: "call Scan(root string) (ok bool)",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "RemoveComment(root string) (ok bool)",
		Function: "Scan(root string) (ok bool)",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "scan() (ok bool)",
		Function: "walk root dir and filter file ext for scan lexer",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "generateNoCommentFile(body string) (ok bool)",
		Function: "generate No Comment File for checking",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "translate() (ok bool)",
		Function: "FindAllBlock,makeBlock,translateBlock,bindBlockType,handlePkgOrApiName,generateGoCodes,reset scanner",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "reset()",
		Function: "reset scanner ctx when finished every file convert work",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "FindAllBlock() (ok bool)",
		Function: "MergeLineElem,FindTypedefs,FindEnums,FindDefines,FindExterns,FindMethods",
		Note:     "first merge every line text elem lexer to line text, and reset allLines when cpp type was found evey time",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "MergeLineElem() (ok bool)",
		Function: "Merge Line Elem (every lexer word) in to line text",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "FindTypedefs()",
		Function: "Find Typedefs",
		Note:     "it will be include struct or point reType,enum,",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "FindEnums()",
		Function: "Find Enums",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "FindDefines()",
		Function: "Find Defines",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "FindExterns()",
		Function: "Find Externs",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "FindMethods()",
		Function: "Find Methods",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "ReSetAllLines(block []LineInfo)",
		Function: "when founded struct,method etc reset lines to null",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "makeBlock() (ok bool)",
		Function: "all line merge in to block,for example all struct,all define",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "translateBlock() (ok bool)",
		Function: "range all bock and convert to go for get go object name,type etc",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "bindBlockType() (ok bool)",
		Function: "bind all cpp type to go type",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "handlePkgOrApiName() (ok bool)",
		Function: "check go syntax",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	doc.Append(mydoc.Row{
		Api:      "generateGoCodes() (ok bool)",
		Function: "last work",
		Note:     "",
		Chinese:  "",
		Todo:     "",
	})
	body := doc.Gen()
	println(body)
}
