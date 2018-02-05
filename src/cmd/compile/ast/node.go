package ast

//代表语法数的一个节点
type Node struct {
	//Pos  Pos
	Data interface{} //class defination or varialbe Defination
}

//type Tops []*Node //语法树顶层结构

type ConvertTops2Package struct {
	Name    []string //package name
	Blocks  []*Block
	Funcs   []*Function
	Classes []*Class
	Enums   []*Enum
	Vars    []*VariableDefinition
	Consts  []*Const
	Import  []*Imports
}

func (c *ConvertTops2Package) ConvertTops2Package(t []*Node) (p *Package, redeclareErrors []*RedeclareError, errs []error) {
	errs = make([]error, 0)
	p = &Package{}
	p.Files = make(map[string]*File)
	c.Name = []string{}
	c.Blocks = []*Block{}
	c.Funcs = make([]*Function, 0)
	c.Classes = make([]*Class, 0)
	c.Enums = make([]*Enum, 0)
	c.Vars = make([]*VariableDefinition, 0)
	c.Consts = make([]*Const, 0)
	expressions := []*Expression{}
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			t := v.Data.(*Block)
			c.Blocks = append(c.Blocks, t)
		case *Function:
			t := v.Data.(*Function)
			c.Funcs = append(c.Funcs, t)
		case *Enum:
			t := v.Data.(*Enum)
			c.Enums = append(c.Enums, t)
		case *Class:
			t := v.Data.(*Class)
			c.Classes = append(c.Classes, t)
			//		case *VariableDefinition:
			//			t := v.Data.(*VariableDefinition)
			//			c.Vars = append(c.Vars, t)
		case *Const:
			t := v.Data.(*Const)
			c.Consts = append(c.Consts, t)
		case *Imports:
			i := v.Data.(*Imports)
			if p.Files[i.Pos.Filename] == nil {
				p.Files[i.Pos.Filename] = &File{Imports: make(map[string]*Imports)}
			}
			p.Files[i.Pos.Filename].Imports[i.Name] = i
		case *Expression: // a,b = f();
			t := v.Data.(*Expression)
			expressions = append(expressions, t)
		default:
			panic("tops have unkown type")
		}
	}
	errs = append(errs, checkEnum(c.Enums)...)
	redeclareErrors = c.redeclareErrors()
	p.Block.Consts = make(map[string]*Const)
	for _, v := range c.Consts {
		p.Block.insert(v.Name, v.Pos, v)
	}
	p.Block.Vars = make(map[string]*VariableDefinition)
	//	for _, v := range c.Vars {
	//		p.Block.Vars[v.Name] = v
	//		v.IsGlobal = true
	//	}
	p.Block.Funcs = make(map[string]*Function)
	for _, v := range c.Funcs {
		v.MkVariableType()
		err := p.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
		v.IsGlobal = true
	}
	p.Block.Classes = make(map[string]*Class)
	for _, v := range c.Classes {
		v.mkVariableType()
		p.Block.Classes[v.Name] = v
	}
	p.Block.Enums = make(map[string]*Enum)
	p.Block.EnumNames = make(map[string]*EnumName)
	for _, v := range c.Enums {
		v.mkVariableType()
		p.Block.Enums[v.Name] = v
		for _, vv := range v.Names {
			p.Block.EnumNames[vv.Name] = vv
		}
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
		c.Blocks = append([]*Block{b}, c.Blocks...)
	}
	p.mkInitFunctions(c.Blocks)
	return
}

func (p *ConvertTops2Package) redeclareErrors() []*RedeclareError {
	ret := []*RedeclareError{}
	m := make(map[string][]interface{})
	//eums
	for _, v := range p.Enums {
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
	for _, v := range p.Consts {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//vars
	for _, v := range p.Vars {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//funcs
	for _, v := range p.Funcs {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//classes
	for _, v := range p.Classes {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	for k, v := range m {
		if len(v) == 1 { //very good
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
			case *ExpressionDeclareVariable:
				panic("1")
			default:
				panic("make error")
			}
		}
		ret = append(ret, r)
	}
	return ret
}
