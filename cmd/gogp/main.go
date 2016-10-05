//this is a cmdline interface of gogp tool.
package main

import (
	"github.com/vipally/cmdline"
	"github.com/vipally/cpright"
	"github.com/vipally/gogp"
)

func main() {
	var (
		filePath    = ""
		codeExt     = ""
		reverseWork = false
		forceUpdate = false
		exit_code   = 0
	)

	cmdline.Version(gogp.Version())
	cmdline.CopyRight(cpright.CopyRight())

	cmdline.Summary("Tool <thiscmd> is a generic-programming solution for golang or any other languages.")
	cmdline.Details(`1. .gpg files
	  An ini file in fact.It's used to define generic parameters's replacing relation.
	  Corresponding .gp file may with the same path and name.
	  But we can redirect it by key "GOGP_GpFilePath".
	  Section "GOGP_REVERSE" is defined for ReverseWork to generate .gp file from .go file.
	  So normal work mode will not generate go code file for this section.

	2. .gp files
	  A go-like file, but exists some <xxx> format keys,
	  that need to be replaced with which defined in .gpg file.

	3. .go files
	  gogp tool auto-generated .go files are exactly normal go code files.
	  But never modify it manually, you can see this warning at the first line in every file.
	  Auto work on GoPath is recmmended.
	  gogp tool will deep travel the path to find all gpg files for processing.
	  If the generated go code file's body has no changes, this file will not be updated.
	  So run gogp tool any times on GoPath is harmless, unless there are indeed changes.
	  So any manually modification will be restored by this tool.
	  Take care of that.
	
	usage example:
	  gogp
	  gogp gopath`)

	cmdline.StringVar(&filePath, "", "filePath", filePath, false, "Path that gogp will work. GoPath and WorkPath is allowed.")
	cmdline.BoolVar(&reverseWork, "r", "reverse", reverseWork, false,
		`Reverse work, this mode is used to gen .gp file from a real-go file.
		If set this flag, the filePath flag must be a .gpg file path related to GoPath.`)
	cmdline.BoolVar(&forceUpdate, "f", "forceUpdate", forceUpdate, false, "Force update all products.")
	cmdline.StringVar(&codeExt, "e", "codeExt", codeExt, false, "Code file ext name. [.go] is default. [.gp] and [.gpg] is not allowed.")
	cmdline.AnotherName("ext", "e")
	cmdline.AnotherName("reverse", "r")
	cmdline.AnotherName("force", "f")
	cmdline.Parse()

	gogp.ForceUpdate(forceUpdate)
	gogp.CodeExtName(codeExt)
	if reverseWork {
		gogp.ReverseWork(filePath)
	} else {
		gogp.Work(filePath)
	}

	cmdline.Exit(exit_code)
}
