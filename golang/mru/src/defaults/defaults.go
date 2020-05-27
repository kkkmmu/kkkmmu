package defaults

import (
	"device"
	//"emulator"
	"regexp"
)

const (
	CaseMark     = "(!@#$^&)"
	GroupMark    = "($^&*@!)"
	SubGroupMark = "(*&&^$!)"
	FeatureMark  = "(#$@#$!)"

	TopologySuffix    = ".to"
	TopologyFileMark  = "!@#tp"
	ScriptSuffix      = ".sc"
	ScriptFileMark    = "!@#sc"
	PreSuffix         = ".pr"
	PreFileMark       = "!@#pr"
	PostSuffix        = ".po"
	PostFileMark      = "!@#po"
	ConfigureSuffix   = ".co"
	ConfigureFileMark = "!@#co"
	DevMark           = "[@DEVS@]"
	IfpMark           = "[@IFPS@]"
	OtherMark         = "[@OTHER@]"
	DUTMark           = "[@DUTS@]"
	LinkMark          = "[@LINKS@]"
	TesterMark        = "[@TESTER@]"
	EmulatorMark      = "[@EMULATOR@]"
	FunctionMark      = "Function"
	AssertionMark     = "Assertion"
	CommentMark       = "###"
)

var Devices = map[string]device.Switch{
	"M3000": device.M3000,
	"M2400": device.M3000,
	"M3500": device.M3000,
}

/*
var Emulators = map[string]emulator.Tester{
	"N2X": emulator.EN2X,
}
*/

var Parsers = map[string]*regexp.Regexp{
	LinkMark:      regexp.MustCompile(`\[@LINKS@\](?P<links>[[:graph:][:space:]]+)\[@LINKS@\]`),
	TesterMark:    regexp.MustCompile(`\[@TESTERS@\](?P<testers>[[:graph:][:space:]]+)\[@TESTERS@\]`),
	EmulatorMark:  regexp.MustCompile(`\[@EMULATOR@\](?P<emu>[[:graph:][:space:]]+)\[@EMULATOR@\]`),
	DUTMark:       regexp.MustCompile(`\[@DUTS@\](?P<duts>[[:graph:][:space:]]+)\[@DUTS@\]`),
	IfpMark:       regexp.MustCompile(`\[@IFPS@\](?P<ifps>[[:graph:][:space:]]+)\[@IFPS@\]`),
	DevMark:       regexp.MustCompile(`\[@DEVS@\](?P<devs>[[:graph:][:space:]]+)\[@DEVS@\]`),
	OtherMark:     regexp.MustCompile(`\[@OTHER@\](?P<devs>[[:graph:][:space:]]+)\[@OTHER@\]`),
	FunctionMark:  regexp.MustCompile(`^[[:space:]]*(?P<dut>[[:word:]]+)\.(?P<function>[[:word:]]+)[[:space:]]*\((?P<parameters>[[:word:][:space:],]*)\)`),
	AssertionMark: regexp.MustCompile(`ASSERT[[:space:]]*\([[:space:]]*(?P<dut>[[:word:]]+)\.(?P<function>[[:word:]]+)[[:space:]]*\((?P<parameters>[[:word:][:space:],]*)\)[[:space:]]*,[[:space:]]*(?P<result>[[:word:]\"]+)[[:space:]]*\)`),
}

func GetTopologyFile(name string) string {
	return "asset/cases/" + name + "/" + name + TopologySuffix
}

func GetScriptFile(name string) string {
	return "asset/cases/" + name + "/" + name + ScriptSuffix
}

func GetPreRequisiteFile(name string) string {
	return "asset/cases/" + name + "/" + name + PreSuffix
}

func GetPostRequisiteFile(name string) string {
	return "asset/cases/" + name + "/" + name + PostSuffix
}

func GetConfigureFile(name string) string {
	return "asset/cases/" + name + "/" + name + ConfigureSuffix
}
