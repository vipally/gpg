# gogp [![GoDoc](https://godoc.org/github.com/vipally/gogp?status.svg)](https://godoc.org/github.com/vipally/gogp) ![Version](https://img.shields.io/badge/version-4.0.0-green.svg)
----
	
package gogp is a generic-programming solution for golang or any other languages.

----

CopyRight 2016 @Ally Dale. All rights reserved.
	
Author  : [Ally Dale(vipally@gmail.com)](mailto://vipally@gmail.com)

Blog    : [http://blog.csdn.net/vipally](http://blog.csdn.net/vipally)

Site    : [https://github.com/vipally](https://github.com/vipally)


----

## todo
1. [x] replace regexp package to [regexp2](https://github.com/dlclark/regexp2) (give up for back-reference works too slow)
1. [x] rebuid and test all regexp syntax
1. [ ] use [cli](https://github.com/urfave/cli) rewrite the command
1. [ ] add command to generate a initialized .go file
1. [x] add syntax spec

## usage of gogp tool:
    1. (Recommend)use cmdline(cmd/gogp):
  
        Tool gogp is a generic-programming solution for golang or any other languages.
        Usage:
          gogp [-e|ext=<Ext>] [-f|force=<force>] [-m|more=<more>] [-remove=<remove>] [<filePath>]
        -e|ext=<Ext>  string
          Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.
        -f|force=<force>
          Force update all products.
        -m|more=<more>
          More information in working process.
        -remove=<remove>
          Only remove all products.
        <filePath>  string
          Path that gogp will work. GoPath and WorkPath is allowed.
  
        usage samples:
           gogp
           gogp gopath
  
    2. package usage:
  
        2.1 (Recommend)import gogp/auto package in test file
          import (
              //"testing"
              _ "github.com/vipally/gogp/auto" //auto runs gogp tool on GoPath when init()
          )
    
        2.2 (Seldom use)import gogp package in test file
          import (
              //"testing"
              "github.com/vipally/gogp"
          )
          func init() {
              gogp.WorkOnGoPath()
          }
----

## Detail desctription:

        Tool Site: https://github.com/vipally/gogp
        Work flow: DummyGoFile  --(GPGFile[1])-->  gp_file  --(GPGFile[2])-->  real_go_files
    
          In this flow, DummyGoFile and GPGFile are hand making, and the rests are products
        of gogp tool.
    
        1. DummyGoFiles
          Sample: https://github.com/vipally/gogp/blob/master/examples/stack.go
          This is a "normal" go file with WELL-DESIGNED structure.
          Texts that matches
                "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END ...\n"
        case will be ingored by gogp tool when loading.
          Any identifier who wants to be replaced with is defines as unique dummy
        word(eg: GOGPStackElem), which is similar to template parameter T in C++.
          GPG file "GOGP_REVERSE_xxx" style sections defines the cases to replacing
        them to <key> format "identifiers" in GP file.
    
        2. GPG files(.gpg)
          GPG file is an ini-format file, that defines key-value replacing cases from
        source to the product.
          "GOGP_IGNORE_xxx" style sections will be ignored by gogp tool.
          "GOGP_REVERSE_xxx" style sections are defined for reverse-mode to generate
        GP file from DummyGoFiles.
          So normal work mode will not generate go code file for these sections.
          "GOGP_xxx" style keys are reserved by gogp tool which will never be replacing with.
          Corresponding GP file may with the same path and name.
          But we can redirect it by key "GOGP_GpFilePath".
          Key "GOGP_Name" is used to specify gp file name in reverse flow.
          And specify go-file-name-suffix in normal flow.
    
        3. GP files(.gp)
          A go-like file, but exists some <xxx> style keys,
          that need to be replaced with which defined in GPG file.
    
        4. GO files(.go)
          gogp tool auto-generated GO files are exactly normal go code files.
          But never modify it manually, you can see this warning infomation at each file head.
          Auto work on GoPath is recmmended.
          gogp tool will deep travel the path to find all gpg files for processing.
          If the generated go code file's body has no changes, this file will not be updated.
          So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
          So any manually modification will be restored by this tool.
          Take care of that.

	    5. Predefined gpg file
```go
		"GOGP_REVERSE"			//gpg section prefix that for gogp reverse only
		"GOGP_IGNORE"			//gpg section prefix that for gogp never process
		"GOGP_xxx" 				//format keys are reserved by gogp tool, who will not be replaced
		"GOGP_Ignore"      		//ignore this section
		"GOGP_DontSave"       	//do not save
		"GOGP_CodeFileName"		//code file name part
		"GOGP_GpFilePathName" 	//gp file path and name
		"KEY_TYPE"         		//key_type
		"VALUE_TYPE"       		//value_type
		
		some predefined gogp grammar
		ignore all text format:
		//#GOGP_IGNORE_BEGIN <content> //#GOGP_IGNORE_END
		
		select by condition <cd> defines in gpg file:
		//#GOGP_IFDEF <cd> <true_content> //#GOGP_ELSE <false_content> //#GOGP_ENDIF
		
		require another gp file:
		//#GOGP_REQUIRE(<gpPath> [, <gpgSection>])
		
		get gpg config string:
		#GOGP_GPGCFG(<cfgName>)
		
		only generate <content> once from a gp file:
		//#GOGP_ONCE <content> //#GOGP_END_ONCE
```

----	
## syntax spec

- **01/19 #comment**<br>
  {make an in line comment in fake .go file.}
```go
// #GOGP_COMMENT {expected code}
```
- **02/19 #if**<br>
  {double-way branch selector by condition}
```go
// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
[// #GOGP_ELSE
	{else content}]
// #GOGP_ENDIF

// #GOGP_IFDEF <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
// #GOGP_ENDIF
```
- **03/19 #if2**<br>
  {double-way branch selector by condition, to nested with #if}
```go
// #GOGP_IFDEF2 <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
[// #GOGP_ELSE2
	{else content}]
// #GOGP_ENDIF2

// #GOGP_IFDEF2 <key> || ! <key> || <key> == xxx || <key> != xxx
	{true content}
// #GOGP_ENDIF2
```
- **04/19 #switch**<br>
  {multi-way branch selector by condition. It is one-switch logic(only one case brantch can trigger out)}
```go
// #GOGP_SWITCH [<SwitchKey>] 
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
// #GOGP_ENDSWITCH
```
- **05/19 #multi-switch**<br>
  {multi-way branch selector by condition. It is multi-switch logic(more than one case brantch can trigger out)}
```go
// #GOGP_MULTISWITCH [<MultiSwitchKey>] 
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
// #GOGP_ENDMULTISWITCH
```
- **06/19 #case**<br>
  {branches of #switch/#multi-switch syntax}
```go
//    #GOGP_CASE <key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
        {case content}
//    #GOGP_ENDCASE
//    #GOGP_DEFAULT
        {default content}
//    #GOGP_ENDCASE
```
- **07/19 #require**<br>
  {require another .gp file}
```go
// #GOGP_REQUIRE(<gp-path> [, <gpgSection>])
```
- **08/19 #replace**<br>
  {<src> -> <dst>, declare build-in key-value replace command for generating .gp file}
```go
// #GOGP_REPLACE(<src>, <dst>)
```
- **09/19 #map**<br>
  {build-in key-value define for generating .gp file. Which can affect brantch of #if and #switch after this code.}
```go
****<src> -> <dst>, which can affect brantch of #GOGP_IFDEF and #GOGP_SWITCH after this code****
// #GOGP_MAP(<src>, <dst>)
```
- **10/19 #ignore**<br>
  {txt that will ignore by gogp tool.}
```go
// #GOGP_IGNORE_BEGIN 
     {ignore-content} 
// #GOGP_IGNORE_END
```
- **11/19 #gp-only**<br>
  {txt that will stay at .gp file only. Which will ignored at final .go file.}
```go
// #GOGP_GPONLY_BEGIN 
     {gp-only content} 
// #GOGP_GPONLY_END
```
- **12/19 #empty-line**<br>
  {empty line.}
```go
{empty-lines} 
```
- **13/19 #trim-empty-line**<br>
  {trim empty line}
```go
{empty-lines} 
{contents}
{empty-lines} 
```
- **14/19 #gpg-config**<br>
  {refer .gpg config}
```go
[//] #GOGP_GPGCFG(<GPGCFG>)
```
- **15/19 #once**<br>
  {code that will generate once during one .gp file processing.}
```go
// #GOGP_ONCE 
    {only generate once from a gp file} 
// #GOGP_END_ONCE 
```
- **16/19 #file-begin**<br>
  {file head of a fake .go file.}
```go
// #GOGP_FILE_BEGIN
```
- **17/19 #file-end**<br>
  {file tail of a fake .go file.}
```go
// #GOGP_FILE_END
```
- **18/19 #to-replace**<br>
  {literal that waiting to replacing.}
```go
<{to-replace}>
```
- **19/19 #condition**<br>
  {txt that for #if or #case condition field parser.}
```go
<key> || !<key> || <key> == xxx || <key> != xxx || <SwitchKeyValue> || !<SwitchKeyValue>
```
	
## More gogp details:

### 1. Working flow:
	
DummyGoFile  --(GPGFile[1])-->  gp_file  --(GPGFile[2])-->  real_go_files
	
	   In this flow, DummyGoFile and GPGFile are hand making, and the rests are products 
	of gogp tool.
	
#### 1.1 DummyGoFile
	    Sample: https://github.com/vipally/gogp/blob/master/examples/stack.go
		
	    This is a "normal" go file with WELL-DESIGNED structure.
	    Texts that matches 
	         "//#GOGP_IGNORE_BEGIN ... //#GOGP_IGNORE_END ...\n"
	case will be ingored by gogp tool when loading.
	    From line 3~14, we add some help info about this DummyGoFile, and that will
	not exists in products.
	    At line-6 (/*   //<----This line can be...), we setted a whole-file comment
	switch corresponding to line-89 (//*/).If add "//" to head of this line, this
	file comes to a "normal" go file, we can edit,compile,test, and of cause, use
	go-fmt tool to format this file.
	    After that, remove "//" from line-6. This file becomes a big-commented file.
	And will have noting for go-doc tool and no export-symbols.Of cause, this
	does nothing to do with the final products real-go files.
	    But there is one limit, we can not use "/* ... */" style comment in this file
	anywhere again.
	
	    Any more, from line 18~35, we defines some dummy types and methods.For making
	this file LEGAL.What we exactly need is the unique dummy identifiers (GOGPStackElem). 
	Which is similar to template parameter T in C++.
	
	    After that, we have a go-like file, but anywhere we want to be replacing with 
	has been set to a unique legal identifiers.
	
#### 1.2 GPGFile
	   GPGFile is an ini-format file, that defines key-value replacing cases from 
	source to the product.
	   "GOGP_IGNORExxxx" style sections will ignore by gogp tool.
	   "GOGP_REVERSExxxx" style sections are used as GPGFile[1](reverse) flow.
	Which is used to generate .gp file from DummyGoFile.
	   Reverse process replaces value(GOGPStackElem) with <key>(<STACK_ELEM>) in .gp file.
	   So .gp file is a normal-go-like file that exists some <xxx> format template 
	keys, which	need to be replaced with proper txt to generate real-go file.
	   Other styles of gpg sections are used as the last flow: generate go code 
	file from .gp file. It is a mechanical matches from keys to values.
	
	   Moreover, "GOGP_xxx" style keys are reserved by gogp tool which will never 
	be replacing with.
	   "GOGP_Name" is used to specify DummyGoFileName in the first flow, and specify 
	go-file-name-suffix in the second flow.
	   "GOGP_GpFilePath" is used to specify .gp file path in the second flow.
	
	   
	
