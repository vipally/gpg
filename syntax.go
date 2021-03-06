// MIT License
//
// Copyright (c) 2021 @gxlb
// Url:
//     https://github.com/gxlb
//     https://gitee.com/gxlb
// AUTHORS:
//     Ally Dale <vipally@gamil.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gogp

import (
	"bytes"
	"fmt"

	"regexp"
	//regexp "github.com/dlclark/regexp2" //back-reference works too slow
)

var allSyntax = []*syntax{
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#comment",
		usage: "make an in line comment in fake .go file.",
		expr:  `(?sm:(?P<COMMENT>/{2,}[ \t]*#GOGP_COMMENT))`,
		syntax: `
// #GOGP_COMMENT {expected code}
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#if",
		usage: "double-way branch selector by condition",
		expr:  `(?sm:^(?:[ \t]*/{2,}[ \t]*)#GOGP_IFDEF[ \t]+(?P<IFCOND>[[:word:]<>\|!= \t]+)(?:.*?$[\r\n]?)(?P<IFT>.*?)(?:(?:[ \t]*/{2,}[ \t]*)#GOGP_ELSE(?:(?:[ \t].*?)?$[\r\n]?)(?P<IFF>.*?))?(?:[ \t]*/{2,}[ \t]*)#GOGP_ENDIF(?:[ \t].*?)?$[\r\n]?)`,
		syntax: `
// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
[// #GOGP_ELSE
	{else content}]
// #GOGP_ENDIF

// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
// #GOGP_ENDIF
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		ignoreInList: false,
		name:         "#if2",
		usage:        "double-way branch selector by condition, to nested with #if",
		expr:         `(?sm:^(?:[ \t]*/{2,}[ \t]*)#GOGP_IFDEF2[ \t]+(?P<IFCOND2>[[:word:]<>\|!= \t]+)(?:.*?$[\r\n]?)(?P<IFT2>.*?)(?:(?:[ \t]*/{2,}[ \t]*)#GOGP_ELSE2(?:(?:[ \t].*?)?$[\r\n]?)(?P<IFF2>.*?))?(?:[ \t]*/{2,}[ \t]*)#GOGP_ENDIF2(?:[ \t].*?)?$[\r\n]?)`,
		syntax: `
// #GOGP_IFDEF2 <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
[// #GOGP_ELSE2
	{else content}]
// #GOGP_ENDIF2

// #GOGP_IFDEF2 <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
// #GOGP_ENDIF2
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#switch",
		usage: "multi-way branch selector by condition. It is one-switch logic(only one case brantch can trigger out)",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)(?:#GOGP_SWITCH)(?:[ \t]+(?P<SWITCHKEY>[[:word:]<>]+))?(?:.*?$)[\r\n]?(?P<SWITCHCONTENT>.*?)(?:^[ \t]*/{2,}[ \t]*)#GOGP_ENDSWITCH(?:[ \t].*?)?$[\r\n]?)`,
		syntax: `
// #GOGP_SWITCH [<SwitchKey>] 
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
// #GOGP_ENDSWITCH
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#multi-switch",
		usage: "multi-way branch selector by condition. It is multi-switch logic(more than one case brantch can trigger out)",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)(?:#GOGP_MULTISWITCH)(?:[ \t]+(?P<MULTISWITCHKEY>[[:word:]<>]+))?(?:.*?$)[\r\n]?(?P<MULTISWITCHCONTENT>.*?)(?:^[ \t]*/{2,}[ \t]*)#GOGP_ENDMULTISWITCH(?:[ \t].*?)?$[\r\n]?)`,
		syntax: `
// #GOGP_MULTISWITCH [<MultiSwitchKey>] 
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
// #GOGP_ENDMULTISWITCH
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		ignoreInList: false,
		name:         "#case",
		usage:        "branches of #switch/#multi-switch syntax",
		expr:         `(?sm:(?:^[ \t]*/{2,}[ \t]*)(?:(?:#GOGP_CASE[ \t]+(?P<CASEKEY>[[:word:]<>\!]+))|(?:#GOGP_DEFAULT))(?:[ \t]*?.*?$)[\r\n]*(?P<CASECONTENT>.*?)(?:^[ \t]*/{2,}[ \t]*)#GOGP_ENDCASE.*?$[\r\n]*)`,
		syntax: `
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#require",
		usage: "require another .gp file",
		expr:  `(?sm:(?P<REQ>(?:^[ \t]*/{2,}[ \t]*)#GOGP_REQUIRE\((?P<REQP>[^\n\r,]*?)(?:[ \t]*?,[ \t]*?(?:(?P<REQN>[[:word:]#@]*)|#GOGP_GPGCFG\((?P<REQGPG>[[:word:]]+)\)))??(?:[ \t]*?\))).*?$[\r\n](?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\([^\n\r,]*?\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\([^\n\r,]*?\))?[\r\n]*)`,
		syntax: `
// #GOGP_REQUIRE(<gp-path> [, <gpgSection>])
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#replace",
		usage: "<src> -> <dst>, declare build-in key-value replace command for generating .gp file",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)#GOGP_REPLACE\((?P<REPSRC>\S+)[ \t]*,[ \t]*(?P<REPDST>\S+)\))`,
		syntax: `
// #GOGP_REPLACE(<src>, <dst>)
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#map",
		usage: "build-in key-value define for generating .gp file. Which can affect brantch of #if and #switch after this code.",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)#GOGP_MAP\((?P<MAPSRC>\S+)[ \t]*,[ \t]*(?P<MAPDST>\S+)\))`,
		syntax: `
****<src> -> <dst>, which can affect brantch of #GOGP_IFDEF and #GOGP_SWITCH after this code****
// #GOGP_MAP(<src>, <dst>)
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#ignore",
		usage: "txt that will ignore by gogp tool.",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)#GOGP_IGNORE_BEGIN(?P<IGNORE>.*?)(?://)??#GOGP_IGNORE_END.*?$[\r\n]*)`,
		syntax: `
// #GOGP_IGNORE_BEGIN 
     {ignore-content} 
// #GOGP_IGNORE_END
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#gp-only",
		usage: "txt that will stay at .gp file only. Which will ignored at final .go file.",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)#GOGP_GPONLY_BEGIN(?P<GPONLY>.*?)(?://)??#GOGP_GPONLY_END.*?$[\r\n]*)`,
		syntax: `
// #GOGP_GPONLY_BEGIN 
     {gp-only content} 
// #GOGP_GPONLY_END
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#empty-line",
		usage: "empty line.",
		expr:  `(?sm:(?P<EMPTY_LINE>[\r\n]{3,}))`,
		syntax: `
{empty-lines} 
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		ignoreInList: true,
		name:         "#trim-empty-line",
		usage:        "trim empty line",
		// must be "s-m", otherwise it with will repeat every line, or not match
		expr: `(?s-m:^[\r\n]*(?P<CONTENT>.*?)[\r\n]*$)`,
		syntax: `
{empty-lines} 
{contents}
{empty-lines} 
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#gpg-config",
		usage: "refer .gpg config",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)?#GOGP_GPGCFG\((?P<GPGCFG>[[:word:]<\->]+)\))`,
		syntax: `
[//] #GOGP_GPGCFG(<GPGCFG>)
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#once",
		usage: "code that will generate once during one .gp file processing.",
		expr:  `(?sm:(?:^[ \t]*/{2,}[ \t]*)#GOGP_ONCE(?:[ \t]*?//.*?$)?[\r\n]*(?P<ONCE>.*?)[\r\n]?(?:^[ \t]*/{2,}[ \t]*)#GOGP_END_ONCE.*?$[\r\n]?)`,
		syntax: `
// #GOGP_ONCE 
    {only generate once from a gp file} 
// #GOGP_END_ONCE 
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#file-begin",
		usage: "file head of a fake .go file.",
		expr:  `(?sm:(?P<FILEB>(?:^[ \t]*/{2,}[ \t]*)#GOGP_FILE_BEGIN(?:[ \t]+(?P<OPEN>[[:word:]]+))?).*?$[\r\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_begin.*?(?://)?#GOGP_IGNORE_END ///gogp_file_begin.*?$)?[\r\n]*)`,
		syntax: `
// #GOGP_FILE_BEGIN
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		name:  "#file-end",
		usage: "file tail of a fake .go file.",
		expr:  `(?sm:(?P<FILEE>(?:^[ \t]*/{2,}[ \t]*)#GOGP_FILE_END).*?$[\r\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_end.*?(?://)?#GOGP_IGNORE_END ///gogp_file_end.*?$)?[\r\n]*)`,
		syntax: `
// #GOGP_FILE_END
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		ignoreInList: true,
		name:         "#to-replace",
		usage:        "literal that waiting to replacing.",
		expr:         `(?P<REPLACEKEY>\<[[:alpha:]_][[:word:]]*\>)`,
		syntax: `
<{to-replace}>
`,
	},
	//--------------------------------------------------------------------------
	&syntax{
		ignoreInList: true,
		name:         "#condition",
		usage:        "txt that for #if or #case condition field parser.",
		expr:         `(?sm:^[ \t]*(?P<NOT>!)?[ \t]*(?P<KEY>[[:word:]<>]+)[ \t]*(?:(?P<OP>==|!=)[ \t]*(?P<VALUE>[[:word:]]+))?[ \t]*)`,
		syntax: `
<key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
`,
	},
}

// syntax regexp descriptor
type syntax struct {
	name         string
	usage        string
	expr         string
	syntax       string
	ignoreInList bool
}

func compileMultiRegexps(res ...*syntax) *regexp.Regexp {
	var b bytes.Buffer
	var exp = `\Q#GOGP_DO_NOT_HAVE_ANY_REGEXP_SYNTAX#\E`
	if len(res) > 0 {
		for _, v := range res {
			if !v.ignoreInList {
				b.WriteString(v.expr)
				b.WriteByte('|')
			}
		}
		if b.Len() > 0 {
			b.Truncate(b.Len() - 1) //remove last '|'
			exp = b.String()
		}
	}
	return regexp.MustCompile(exp)
}

func (st *syntax) MustCompile() *regexp.Regexp {
	return regexp.MustCompile(st.expr)
}

func findSyntax(name string) *syntax {
	for _, v := range allSyntax {
		if v.name == name {
			return v
		}
	}
	panic(fmt.Errorf("findSyntax(%s) not found", name))
	return nil
}

var (
	gogpExpTodoReplace   = findSyntax("#to-replace").MustCompile()
	gogpExpIgnore        = findSyntax("#ignore").MustCompile()
	gogpExpCases         = findSyntax("#case").MustCompile()
	gogpExpEmptyLine     = findSyntax("#empty-line").MustCompile()
	gogpExpTrimEmptyLine = findSyntax("#trim-empty-line").MustCompile()
	gogpExpRequire       = findSyntax("#require").MustCompile()
	gogpExpCondition     = findSyntax("#condition").MustCompile()
	gogpExpComment       = findSyntax("#comment").MustCompile()

	gogpExpCodeSelector = compileMultiRegexps(
		findSyntax("#ignore"),
		findSyntax("#gp-only"),
		findSyntax("#map"),
		findSyntax("#switch"),
		findSyntax("#multi-switch"),
		findSyntax("#if"),
		findSyntax("#if2"),
	)

	gogpExpPretreatAll = compileMultiRegexps(
		findSyntax("#ignore"),
		findSyntax("#require"),
		findSyntax("#gpg-config"),
		findSyntax("#once"),
		findSyntax("#replace"),
		findSyntax("#comment"),
	)

	gogpExpRequireAll = compileMultiRegexps(
		findSyntax("#require"),
		findSyntax("#file-begin"),
		findSyntax("#file-end"),
	)

	gogpExpReverseIgnoreAll = compileMultiRegexps(
		findSyntax("#file-begin"),
		findSyntax("#file-end"),
		findSyntax("#ignore"),
	)
)
