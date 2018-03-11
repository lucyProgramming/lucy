package signature

//
//import (
//	"github.com/756445638/lucy/src/cmd/compile/ast"
//)
//
//type Sinature struct {
//}
//
//type ClassTypeSignature struct {
//}
//
//type ArrayTypeSignature struct {
//}
//
//type TypeVariableSignature struct {
//}
//
//type ReferenceTypeSignature struct {
//	Typ *ast.VariableType
//}
//
//func (s *Sinature) parseIdentifier(bs []byte) (string, []byte) {
//	b := []byte{}
//	length := len(bs)
//	var i int
//	for i = 0; i < length; i++ {
//		if (bs[i] >= 'A' && bs[i] <= 'Z') ||
//			(bs[i] >= 'a' && bs[i] <= 'z') ||
//			bs[i] == '_' {
//			b = append(b, bs[i])
//		} else {
//			break
//		}
//	}
//	ret := string(b)
//	bs = bs[i:]
//	return ret, bs
//}
//
//func (s *Sinature) parseReferenceTypeSignature(bs []byte) (*ast.VariableType, []byte) {
//	switch bs[0] {
//	case 'L':
//		return s.parserClassTypeSignature(bs)
//	case '[':
//		return s.parserArrayTypeSignature(bs)
//	case 'T':
//		panic("111")
//	}
//	return nil, nil
//}
//
//func (s *Sinature) parserArrayTypeSignature(bs []byte) (*ast.VariableType, []byte) {
//	// skip [
//	bs = bs[1:]
//	ret := &ast.VariableType{}
//	ret.Typ = ast.VARIABLE_TYPE_JAVA_ARRAY
//	ret.ArrayType, bs = s.parseJvaTypeSignature(bs)
//	return ret, bs
//}
//func (s *Sinature) parserClassTypeSignature(bs []byte) (*ast.VariableType, []byte) {
//	bs = bs[1:] // skip L
//
//	return nil, bs
//}
//
//func (s *Sinature) parseTypedArguments(bs []byte) (*ast.TypedParameters, []byte) {
//	bs = bs[1:]
//	ret := &ast.TypedParameters{}
//	var identifier string
//	for bs[0] != '>' { // '>' is end of typed parameter
//		identifier, bs = s.parseIdentifier(bs)
//		firstOne := true
//		var t *ast.VariableType
//		for bs[0] == ':' {
//			bs = bs[1:] // skip :
//			t, bs = s.parseReferenceTypeSignature(bs)
//			if firstOne {
//				firstOne = false
//			} else {
//			}
//		}
//		p := &ast.TypedParameter{}
//		p.Name = identifier
//		p.Typ = t
//		ret.Parameters = append(ret.Parameters, p)
//	}
//	return nil, nil
//}
//
//func (s *Sinature) parseTypedArgument(bs []byte) (*ast.TypedParameters, []byte) {
//	if bs[0] == '+' || bs[0] == '-' {
//		bs = bs[1:] // TODO::WildcardIndicator I don`t know what`s it means
//	}
//	return nil, nil
//}
//
//func (s *Sinature) parseJvaTypeSignature(bs []byte) (*ast.VariableType, []byte) {
//	ret := &ast.VariableType{}
//	switch bs[0] {
//	case 'B':
//		ret.Typ = ast.VARIABLE_TYPE_BYTE
//		bs = bs[1:]
//		return ret, bs
//	case 'C':
//		ret.Typ = ast.VARIABLE_TYPE_SHORT
//		bs = bs[1:]
//		return ret, bs
//	case 'D':
//		ret.Typ = ast.VARIABLE_TYPE_DOUBLE
//		bs = bs[1:]
//		return ret, bs
//	case 'F':
//		ret.Typ = ast.VARIABLE_TYPE_FLOAT
//		bs = bs[1:]
//		return ret, bs
//	case 'I':
//		ret.Typ = ast.VARIABLE_TYPE_INT
//		bs = bs[1:]
//		return ret, bs
//	case 'J':
//		ret.Typ = ast.VARIABLE_TYPE_LONG
//		bs = bs[1:]
//		return ret, bs
//	case 'S':
//		ret.Typ = ast.VARIABLE_TYPE_SHORT
//		bs = bs[1:]
//		return ret, bs
//	case 'Z':
//		ret.Typ = ast.VARIABLE_TYPE_SHORT
//		bs = bs[1:]
//		return ret, bs
//	default:
//
//	}
//	return nil, nil
//}
