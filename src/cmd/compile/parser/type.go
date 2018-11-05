package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (this *Parser) parseType() (*ast.Type, error) {
	var err error
	var ret *ast.Type
	pos := this.mkPos()
	switch this.token.Type {
	case lex.TokenLb:
		this.Next(lfIsToken)
		this.unExpectNewLineAndSkip()
		if this.token.Type != lex.TokenRb {
			// [ and ] not match
			err = fmt.Errorf("%s '[' and ']' not match", this.errMsgPrefix())
			this.errs = append(this.errs, err)
			return nil, err
		}
		//lookahead
		this.Next(lfIsToken) //skip ]
		this.unExpectNewLineAndSkip()
		array, err := this.parseType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{}
		ret.Pos = pos
		ret.Type = ast.VariableTypeArray
		ret.Array = array
	case lex.TokenBool:
		ret = &ast.Type{
			Type: ast.VariableTypeBool,
			Pos:  pos,
		}
		this.Next(lfIsToken)
	case lex.TokenByte:
		ret = &ast.Type{
			Type: ast.VariableTypeByte,
			Pos:  pos,
		}
		this.Next(lfIsToken)

	case lex.TokenShort:
		ret = &ast.Type{
			Type: ast.VariableTypeShort,
			Pos:  pos,
		}
		this.Next(lfIsToken)
	case lex.TokenChar:
		ret = &ast.Type{
			Type: ast.VariableTypeChar,
			Pos:  pos,
		}
		this.Next(lfIsToken)
	case lex.TokenInt:
		ret = &ast.Type{
			Type: ast.VariableTypeInt,
			Pos:  pos,
		}
		this.Next(lfIsToken)
	case lex.TokenFloat:
		ret = &ast.Type{
			Type: ast.VariableTypeFloat,
			Pos:  pos,
		}
		this.Next(lfIsToken)

	case lex.TokenDouble:
		ret = &ast.Type{
			Type: ast.VariableTypeDouble,
			Pos:  pos,
		}
		this.Next(lfIsToken)
	case lex.TokenLong:
		ret = &ast.Type{
			Type: ast.VariableTypeLong,
			Pos:  pos,
		}
		this.Next(lfIsToken)
	case lex.TokenString:
		ret = &ast.Type{
			Type: ast.VariableTypeString,
			Pos:  pos,
		}
		this.Next(lfIsToken)
	case lex.TokenIdentifier:
		ret, err = this.parseIdentifierType()
		if err != nil {
			this.errs = append(this.errs, err)
		}
	case lex.TokenMap:
		this.Next(lfNotToken) // skip map key word
		if this.token.Type != lex.TokenLc {
			return nil, fmt.Errorf("%s expect '{',but '%s'",
				this.errMsgPrefix(), this.token.Description)
		}
		this.Next(lfNotToken) // skip {
		var k, v *ast.Type
		k, err = this.parseType()
		if err != nil {
			return nil, err
		}
		this.ifTokenIsLfThenSkip()
		if this.token.Type != lex.TokenArrow {
			return nil, fmt.Errorf("%s expect '->',but '%s'",
				this.errMsgPrefix(), this.token.Description)
		}
		this.Next(lfNotToken) // skip ->
		v, err := this.parseType()
		if err != nil {
			return nil, err
		}
		this.ifTokenIsLfThenSkip()
		if this.token.Type != lex.TokenRc {
			return nil, fmt.Errorf("%s expect '}',but '%s'",
				this.errMsgPrefix(), this.token.Description)
		}
		this.Next(lfIsToken)
		m := &ast.Map{
			K: k,
			V: v,
		}
		ret = &ast.Type{
			Type: ast.VariableTypeMap,
			Map:  m,
			Pos:  pos,
		}
	case lex.TokenFn:
		this.Next(lfIsToken)
		ft, err := this.parseFunctionType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{
			Type:         ast.VariableTypeFunction,
			Pos:          pos,
			FunctionType: &ft,
		}
	case lex.TokenGlobal:
		this.Next(lfIsToken)
		this.unExpectNewLineAndSkip()
		if this.token.Type != lex.TokenSelection {
			return nil, fmt.Errorf("%s expect '.' , but '%s'",
				this.errMsgPrefix(), this.token.Description)
		}
		this.Next(lfNotToken)
		if this.token.Type != lex.TokenIdentifier {
			this.errs = append(this.errs, fmt.Errorf("%s expect identifier , but '%s'",
				this.errMsgPrefix(), this.token.Description))
		} else {
			ret = &ast.Type{
				Type: ast.VariableTypeGlobal,
				Pos:  pos,
				Name: this.token.Data.(string),
			}
			this.Next(lfIsToken)
		}
	default:
		err := fmt.Errorf("%s unkown begining '%s' token for a type",
			this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		return nil, err
	}
	if err != nil {
		this.errs = append(this.errs, err)
		return nil, err
	}
	if this.token.Type == lex.TokenVArgs {
		newRet := &ast.Type{
			Pos:            this.mkPos(),
			Type:           ast.VariableTypeJavaArray,
			Array:          ret,
			IsVariableArgs: true,
		}
		this.Next(lfIsToken) // skip ...
		ret = newRet
		return ret, nil
	}
	for this.token.Type == lex.TokenLb { // int [
		this.Next(lfIsToken) // skip [
		this.unExpectNewLineAndSkip()
		if this.token.Type != lex.TokenRb {
			err = fmt.Errorf("%s '[' and ']' not match", this.errMsgPrefix())
			this.errs = append(this.errs, err)
			return ret, err
		}
		newRet := &ast.Type{
			Pos:   this.mkPos(),
			Type:  ast.VariableTypeJavaArray,
			Array: ret,
		}
		ret = newRet
		this.Next(lfIsToken) // skip ]
	}
	return ret, err
}

/*
	valid begin token of a type
*/
func (this *Parser) isValidTypeBegin() bool {
	return this.token.Type == lex.TokenLb ||
		this.token.Type == lex.TokenBool ||
		this.token.Type == lex.TokenByte ||
		this.token.Type == lex.TokenShort ||
		this.token.Type == lex.TokenChar ||
		this.token.Type == lex.TokenInt ||
		this.token.Type == lex.TokenFloat ||
		this.token.Type == lex.TokenDouble ||
		this.token.Type == lex.TokenLong ||
		this.token.Type == lex.TokenString ||
		this.token.Type == lex.TokenMap ||
		this.token.Type == lex.TokenIdentifier ||
		this.token.Type == lex.TokenFn
}

func (this *Parser) parseIdentifierType() (*ast.Type, error) {
	name := this.token.Data.(string)
	ret := &ast.Type{
		Pos:  this.mkPos(),
		Type: ast.VariableTypeName,
	}
	this.Next(lfIsToken) // skip name identifier
	for this.token.Type == lex.TokenSelection {
		this.Next(lfNotToken) // skip .
		if this.token.Type != lex.TokenIdentifier {
			return nil, fmt.Errorf("%s not a identifier after dot",
				this.errMsgPrefix())
		}
		name += "." + this.token.Data.(string)
		ret.Pos = this.mkPos() //  override pos
		this.Next(lfIsToken)   // skip identifier
	}
	ret.Name = name
	return ret, nil
}

func (this *Parser) parseTypes(endTokens ...lex.TokenKind) ([]*ast.Type, error) {
	ret := []*ast.Type{}
	for this.token.Type != lex.TokenEof {
		t, err := this.parseType()
		if err != nil {
			return ret, err
		}
		ret = append(ret, t)
		if this.token.Type != lex.TokenComma {
			if this.isValidTypeBegin() {
				this.errs = append(this.errs, fmt.Errorf("%s missing comma",
					this.errMsgPrefix()))
				continue
			}
			break
		}
		this.Next(lfNotToken) // skip ,
		for _, v := range endTokens {
			if v == this.token.Type {
				this.errs = append(this.errs, fmt.Errorf("%s extra comma", this.errMsgPrefix()))
				goto end
			}
		}
	}
end:
	return ret, nil
}
