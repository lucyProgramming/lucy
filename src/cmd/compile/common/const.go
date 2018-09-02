package common

const (
	//array
	ArrayMethodSize      = "size"
	ArrayMethodCap       = "cap"   // for debug
	ArrayMethodStart     = "start" // for debug
	ArrayMethodEnd       = "end"   // for debug
	ArrayMethodAppend    = "append"
	ArrayMethodAppendAll = "appendAll"
	// map
	MapMethodRemove    = "remove"
	MapMethodRemoveAll = "removeAll"
	MapMethodKeyExist  = "keyExist"
	MapMethodSize      = "size"
	// buildIn function
	BuildInFunctionPanic        = "panic"
	BuildInFunctionCatch        = "catch"
	BuildInFunctionPrint        = "print"
	BuildInFunctionPrintf       = "printf"
	BuildInFunctionSprintf      = "sprintf"
	BuildInFunctionLen          = "len"
	BuildInFunctionMonitorEnter = "monitorEnter"
	BuildInFunctionMonitorExit  = "monitorExit"
	BuildInFunctionBlockHole    = "blackHole"
	BuildInFunctionAssert       = "assert"
	//BuildInFunctionTypeOf       = "typeOf"
)

const (
	MainClassName = "main.class"
)
