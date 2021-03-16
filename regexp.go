package gogp

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	txtReplaceKeyFmt  = "<%s>"
	txtSectionReverse = "GOGP_REVERSE" //gpg section prefix that for gogp reverse only
	txtSectionIgnore  = "GOGP_IGNORE"  //gpg section prefix that for gogp never process

	keyReservePrefix  = "<GOGP_"            //reserved key, who will not use repalce action
	rawKeyIgnore      = "GOGP_Ignore"       //ignore this section
	rawKeyProductName = "GOGP_CodeFileName" //code file name part
	rawKeySrcPathName = "GOGP_GpFilePath"   //gp file path and name
	rawKeyDontSave    = "GOGP_DontSave"     //do not save
	rawKeyKeyType     = "KEY_TYPE"          //key_type
	rawKeyValueType   = "VALUE_TYPE"        //value_type

	// // #GOGP_COMMENT
	expTxtGogpComment = `(?sm:(?P<COMMENT>/{2,}[ |\t]*#GOGP_COMMENT))`

	//generic-programming flag <XXX>
	expTxtTodoReplace = `(?P<P>.?)(?P<W>\<[[:alpha:]_][[:word:]]*\>)(?P<S>.?)`

	// ignore all text format:
	// //#GOGP_IGNORE_BEGIN <content> //#GOGP_IGNORE_END
	expTxtIgnore = `(?sm:\s*//#GOGP_IGNORE_BEGIN(?P<IGNORE>.*?)(?://)??#GOGP_IGNORE_END.*?$[\r|\n]*)`
	expTxtGPOnly = `(?sm:\s*//#GOGP_GPONLY_BEGIN(?P<GPONLY>.*?)(?://)??#GOGP_GPONLY_END.*?$[\r|\n]*)`

	// select by condition <cd> defines in gpg file:
	// //#GOGP_IFDEF <cd> <true_content> //#GOGP_ELSE <false_content> //#GOGP_ENDIF
	// "<key> || ! <key> || <key> == xxx || <key> != xxx"
	expTxtIf  = `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF[ |\t]+(?P<CONDK>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE(?:.*?$[\r|\n]?)[\r|\n]*(?P<F>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF.*?$[\r|\n]?)`
	expTxtIf2 = `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF2[ |\t]+(?P<CONDK2>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T2>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE2(?:.*?$[\r|\n]?)[\r|\n]*(?P<F2>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF2.*?$[\r|\n]?)`

	// " <key> || !<key> || <key> == xxx || <key> != xxx "
	// [<NOT>] <KEY> [<OP><VALUE>]
	expCondition = `(?sm:^[ |\t]*(?P<NOT>!)?[ |\t]*(?P<KEY>[[:word:]<>]+)[ |\t]*(?:(?P<OP>==|!=)[ |\t]*(?P<VALUE>[[:word:]]+))?[ |\t]*)`

	//#GOGP_SWITCH [<SWITCHKEY>] <CASES> #GOGP_GOGP_ENDSWITCH
	expTxtSwitch = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:#GOGP_SWITCH)(?:[ |\t]+(?P<SWITCHKEY>[[:word:]<>]+))?(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASES>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDSWITCH.*?$[\r|\n]?)`

	//#GOGP_CASE <COND> <CASE> #GOGP_ENDCASE
	//#GOGP_DEFAULT <CASE> #GOGP_ENDCASE
	expTxtCase = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:(?:#GOGP_CASE[ |\t]+(?P<COND>[[:word:]<>\|!]+))|(?:#GOGP_DEFAULT))(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASE>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDCASE.*?$[\r|\n]*)`

	// require another gp file:
	// //#GOGP_REQUIRE(<gpPath> [, <gpgSection>])
	expTxtRequire   = `(?sm:\s*(?P<REQ>^[ |\t]*(?://)?#GOGP_REQUIRE\((?P<REQP>[^\n\r,]*?)(?:[ |\t]*?,[ |\t]*?(?:(?P<REQN>[[:word:]|#|@]*)|#GOGP_GPGCFG\((?P<REQGPG>[[:word:]]+)\)))??(?:[ |\t]*?\))).*?$[\r|\n]*(?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\([^\n\r,]*?\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\([^\n\r,]*?\))?[\r|\n]*)`
	expTxtEmptyLine = `(?sm:(?P<EMPTY_LINE>[\r|\n]{3,}))`

	//must be "-sm", otherwise it with will repeat every line
	expTxtTrimEmptyLine = `(?-sm:^[\r|\n]*(?P<CONTENT>.*?)[\r|\n]*$)`

	// get gpg config string:
	// #GOGP_GPGCFG(<cfgName>)
	expTxtGetGpgCfg = `(?sm:(?://)?#GOGP_GPGCFG\((?P<GPGCFG>[[:word:]]+)\))`

	// #GOGP_REPLACE(<src>,<dst>)
	expTxtReplaceKey = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_REPLACE\((?P<REPSRC>\S+)[ |\t]*,[ |\t]*(?P<REPDST>\S+)\))`
	expTxtMapKey     = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_MAP\((?P<MAPSRC>\S+)[ |\t]*,[ |\t]*(?P<MAPDST>\S+)\))`

	//remove "*" from value type such as "*string -> string"
	// #GOGP_RAWNAME(<strValueType>)
	//gsExpTxtRawName = "(?-sm:(?://)?#GOGP_RAWNAME\((?P<RAWNAME>\S+)\))"

	// only generate <content> once from a gp file:
	// //#GOGP_ONCE <content> //#GOGP_END_ONCE
	expTxtOnce = `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)//#GOGP_ONCE(?:[ |\t]*?//.*?$)?[\r|\n]*(?P<ONCE>.*?)[\r|\n]*[ |\t]*?(?://)??#GOGP_END_ONCE.*?$[\r|\n]*)`

	expTxtFileBegin = `(?sm:\s*(?P<FILEB>//#GOGP_FILE_BEGIN(?:[ |\t]+(?P<OPEN>[[:word:]]+))?).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_begin.*?(?://)?#GOGP_IGNORE_END ///gogp_file_begin.*?$)?[\r|\n]*)`
	expTxtFileEnd   = `(?sm:\s*(?P<FILEE>//#GOGP_FILE_END).*?$[\r|\n]*(?://#GOGP_IGNORE_BEGIN ///gogp_file_end.*?(?://)?#GOGP_IGNORE_END ///gogp_file_end.*?$)?[\r|\n]*)`

	// "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END"
	txtRequireResultFmt   = "//#GOGP_IGNORE_BEGIN ///require begin from(%s)\n%s\n//#GOGP_IGNORE_END ///require end from(%s)"
	txtRequireAtResultFmt = "///require begin from(%s)\n%s\n///require end from(%s)"
	txtGogpIgnoreFmt      = "//#GOGP_IGNORE_BEGIN%s%s//#GOGP_IGNORE_END%s"
)

var res = []*re{
	&re{
		name:   "comment",
		exp:    `(?sm:(?P<COMMENT>/{2,}[ |\t]*#GOGP_COMMENT))`,
		syntax: `// #GOGP_COMMENT`,
	},
	&re{
		name: "if",
		exp:  `(?sm:^(?:[ |\t]*/{2,}[ |\t]*)#GOGP_IFDEF[ |\t]+(?P<CONDK>[[:word:]<>\|!= \t]+)(?:.*?$[\r|\n]?)(?P<T>.*?)(?:(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ELSE(?:.*?$[\r|\n]?)[\r|\n]*(?P<F>.*?))?(?:[ |\t]*/{2,}[ |\t]*)#GOGP_ENDIF.*?$[\r|\n]?)`,
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
	&re{
		name: "switch",
		exp:  `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:#GOGP_SWITCH)(?:[ |\t]+(?P<SWITCHKEY>[[:word:]<>]+))?(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASES>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDSWITCH.*?$[\r|\n]?)`,
		syntax: `
**** it is multi-switch logic(more than one case brantch can trigger out) ****
// #GOGP_SWITCH [<SwitchKey>] 
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
// #GOGP_GOGP_ENDSWITCH
`,
	},
	&re{
		name: "case",
		exp:  `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)(?:(?:#GOGP_CASE[ |\t]+(?P<COND>[[:word:]<>\|!]+))|(?:#GOGP_DEFAULT))(?:[ |\t]*?.*?$)[\r|\n]*(?P<CASE>.*?)(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_ENDCASE.*?$[\r|\n]*)`,
		syntax: `
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
`,
	},
	&re{
		name:   "require",
		exp:    `(?sm:\s*(?P<REQ>^[ |\t]*(?://)?#GOGP_REQUIRE\((?P<REQP>[^\n\r,]*?)(?:[ |\t]*?,[ |\t]*?(?:(?P<REQN>[[:word:]|#|@]*)|#GOGP_GPGCFG\((?P<REQGPG>[[:word:]]+)\)))??(?:[ |\t]*?\))).*?$[\r|\n]*(?:(?://#GOGP_IGNORE_BEGIN )?///require begin from\([^\n\r,]*?\)(?P<REQCONTENT>.*?)(?://)?(?:#GOGP_IGNORE_END )?///require end from\([^\n\r,]*?\))?[\r|\n]*)`,
		syntax: `// #GOGP_REQUIRE(<gp-path> [, <gpgSection>])`,
	},
	&re{
		name: "replace",
		exp:  `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_REPLACE\((?P<REPSRC>\S+)[ |\t]*,[ |\t]*(?P<REPDST>\S+)\))`,
		syntax: `
****<src> -> <dst>, literal replacement****
// #GOGP_REPLACE(<src>, <dst>)
`,
	},
	&re{
		name: "map",
		exp:  `(?sm:(?:^[ |\t]*/{2,}[ |\t]*)#GOGP_MAP\((?P<MAPSRC>\S+)[ |\t]*,[ |\t]*(?P<MAPDST>\S+)\))`,
		syntax: `
****<src> -> <dst>, which can affect brantch of #GOGP_IFDEF and #GOGP_SWITCH after this code****
// #GOGP_MAP(<src>, <dst>)
`,
	},
}

var (
	gogpExpTodoReplace      = regexp.MustCompile(expTxtTodoReplace)
	gogpExpPretreatAll      = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtRequire, expTxtGetGpgCfg, expTxtOnce, expTxtReplaceKey, expTxtGogpComment))
	gogpExpIgnore           = regexp.MustCompile(expTxtIgnore)
	gogpExpCodeSelector     = regexp.MustCompile(fmt.Sprintf("%s|%s|%s|%s|%s|%s", expTxtIgnore, expTxtGPOnly, expTxtIf, expTxtIf2, expTxtMapKey, expTxtSwitch))
	gogpExpCases            = regexp.MustCompile(expTxtCase)
	gogpExpEmptyLine        = regexp.MustCompile(expTxtEmptyLine)
	gogpExpTrimEmptyLine    = regexp.MustCompile(expTxtTrimEmptyLine)
	gogpExpRequire          = regexp.MustCompile(expTxtRequire)
	gogpExpRequireAll       = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtRequire, expTxtFileBegin, expTxtFileEnd))
	gogpExpReverseIgnoreAll = regexp.MustCompile(fmt.Sprintf("%s|%s|%s", expTxtFileBegin, expTxtFileEnd, expTxtIgnore))
	gogpExpCondition        = regexp.MustCompile(expTxtRequire)
	gogpExpComment          = regexp.MustCompile(expTxtGogpComment)

	txtFileBeginContent = `//
/*   //This line can be uncommented to disable all this file, and it doesn't effect to the .gp file
//	 //If test or change .gp file required, comment it to modify and compile as normal go file
//
// This is a fake go code file
// It is used to generate .gp file by gogp tool
// Real go code file will be generated from .gp file
//
`
	txtFileBeginContentOpen = strings.Replace(txtFileBeginContent, "/*", "///*", 1)
	txtFileEndContent       = "//*/\n"
)

// regexp object
type re struct {
	name   string
	exp    string
	syntax string
}

func regexpCompile(res ...*re) *regexp.Regexp {
	var b bytes.Buffer
	var exp = `\Q#GOGP_DO_NOT_HAVE_ANY_KEY#\E`
	if len(res) > 0 {
		for _, v := range res {
			b.WriteString(v.exp)
			b.WriteByte('|')
		}
		b.Truncate(b.Len() - 1) //remove last '|'
		exp = b.String()

	}
	return regexp.MustCompile(exp)
}

func (r *re) Regexp() *regexp.Regexp {
	return regexp.MustCompile(r.exp)
}

func (r *re) Syntax() string {
	return r.syntax
}

func (r *re) Name() string {
	return r.name
}

func findRE(name string) *re {
	for _, v := range res {
		if v.name == name {
			return v
		}
	}
	panic(fmt.Errorf("findRE(%s) not found", name))
	return nil
}
