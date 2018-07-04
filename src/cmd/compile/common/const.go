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
	MapMethodKeyExists = "keyExists"
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
)

const (
	MainClassName = "main.class"
)
