package ast

//代表语法数的一个节点
type Node struct {
	//Pos  Pos
	Data interface{} //class defination or varialbe Defination
}

//type Tops []*Node //语法树顶层结构

type ConvertTops2Package struct {
	Name      []string //package name
	Blocks    []*Block
	Funcs     []*Function
	Classes   []*Class
	Enums     []*Enum
	Vars      []*VariableDefinition
	Consts    []*Const
	Import    []*Import
	TypeAlias []*ExpressionTypeAlias
}

func (convertor *ConvertTops2Package) ConvertTops2Package(t []*Node) (Pack *Package, redeclareErrors []*RedeclareError, errs []error) {
	errs = make([]error, 0)
	Pack = &Package{}
	Pack.Files = make(map[string]*File)
	convertor.Name = []string{}
	convertor.Blocks = []*Block{}
	convertor.Funcs = make([]*Function, 0)
	convertor.Classes = make([]*Class, 0)
	convertor.Enums = make([]*Enum, 0)
	convertor.Vars = make([]*VariableDefinition, 0)
	convertor.Consts = make([]*Const, 0)
	expressions := []*Expression{}
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			t := v.Data.(*Block)
			convertor.Blocks = append(convertor.Blocks, t)
		case *Function:
			t := v.Data.(*Function)
			convertor.Funcs = append(convertor.Funcs, t)
		case *Enum:
			t := v.Data.(*Enum)
			convertor.Enums = append(convertor.Enums, t)
		case *Class:
			t := v.Data.(*Class)
			convertor.Classes = append(convertor.Classes, t)
		case *Const:
			t := v.Data.(*Const)
			convertor.Consts = append(convertor.Consts, t)
		case *Import:
			i := v.Data.(*Import)
			if Pack.Files[i.Pos.Filename] == nil {
				Pack.Files[i.Pos.Filename] = &File{Imports: make(map[string]*Import)}
			}
			Pack.Files[i.Pos.Filename].Imports[i.AccessName] = i
		case *Expression: // a,b = f();
			t := v.Data.(*Expression)
			expressions = append(expressions, t)
		case *ExpressionTypeAlias:
			t := v.Data.(*ExpressionTypeAlias)
			convertor.TypeAlias = append(convertor.TypeAlias, t)
		default:
			panic("tops have unkown type")
		}
	}

	errs = append(errs, checkEnum(convertor.Enums)...)
	redeclareErrors = convertor.redeclareErrors()
	Pack.Block.Consts = make(map[string]*Const)
	for _, v := range convertor.Consts {
		Pack.Block.insert(v.Name, v.Pos, v)
	}
	Pack.Block.Vars = make(map[string]*VariableDefinition)
	Pack.Block.Funcs = make(map[string]*Function)
	for _, v := range convertor.Funcs {
		err := Pack.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
		v.IsGlobal = true
	}
	Pack.Block.Classes = make(map[string]*Class)
	for _, v := range convertor.Classes {
		Pack.Block.Classes[v.Name] = v
	}
	Pack.Block.Enums = make(map[string]*Enum)
	Pack.Block.EnumNames = make(map[string]*EnumName)
	for _, v := range convertor.Enums {
		Pack.Block.Enums[v.Name] = v
		for _, vv := range v.Names {
			Pack.Block.EnumNames[vv.Name] = vv
		}
	}
	//after class inserted,then resolve type
	for _, v := range convertor.TypeAlias {
		err := v.Typ.resolve(&Pack.Block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		Pack.Block.Types[v.Name] = v.Typ
	}
	if len(expressions) > 0 {
		s := make([]*Statement, len(expressions))
		for k, v := range expressions {
			s[k] = &Statement{
				Typ:        STATEMENT_TYPE_EXPRESSION,
				Expression: v,
			}
		}
		b := &Block{}
		b.Statements = s
		b.isGlobalVariableDefinition = true
		convertor.Blocks = append([]*Block{b}, convertor.Blocks...)
	}
	Pack.mkInitFunctions(convertor.Blocks)
	return
}

func (convertor *ConvertTops2Package) redeclareErrors() []*RedeclareError {
	ret := []*RedeclareError{}
	m := make(map[string][]interface{})
	//eums
	for _, v := range convertor.Enums {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
		for _, vv := range v.Names {
			if _, ok := m[vv.Name]; ok {
				m[vv.Name] = append(m[vv.Name], vv)
			} else {
				m[vv.Name] = []interface{}{vv}
			}
		}
	}
	//const
	for _, v := range convertor.Consts {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//vars
	for _, v := range convertor.Vars {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//funcs
	for _, v := range convertor.Funcs {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//classes
	for _, v := range convertor.Classes {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	// type alias
	for _, v := range convertor.TypeAlias {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}

	for k, v := range m {
		if len(v) == 1 || len(v) == 0 { //very good  , 0 looks impossible
			continue
		}
		r := &RedeclareError{}
		r.Name = k
		r.Pos = make([]*Pos, 0)
		for _, vv := range v {
			switch vv.(type) {
			case *Const:
				t := vv.(*Const)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "const"
			case *Enum:
				t := vv.(*EnumName)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "enum"
			case *Function:
				t := vv.(*Function)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "function"
			case *Class:
				t := vv.(*Class)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "class"
			case *ExpressionTypeAlias:
				t := vv.(*ExpressionTypeAlias)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "type alias"
			default:
				panic("make error")
			}
		}
		ret = append(ret, r)
	}
	return ret
}
