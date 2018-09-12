package lc

import (
	"encoding/binary"

	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
	load from file implementation
*/
type FileLoader struct {
	caches map[string]interface{}
}

func (loader *FileLoader) LoadImport(importName string) (interface{}, error) {
	if loader.caches != nil && loader.caches[importName] != nil {
		return loader.caches[importName], nil
	}
	var realPaths []*Resource
	for _, v := range compiler.lucyPaths {
		p := filepath.Join(v, "class", importName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realPaths = append(realPaths, &Resource{
				kind:     ResourceKindLucyPackage,
				realPath: p,
				name:     importName,
			})
		}
		p = filepath.Join(v, "class", importName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // class file
			realPaths = append(realPaths, &Resource{
				kind:     ResourceKindLucyClass,
				realPath: p,
				name:     importName,
			})
		}
	}
	for _, v := range compiler.ClassPaths {
		p := filepath.Join(v, importName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realPaths = append(realPaths, &Resource{
				kind:     ResourceKindJavaPackage,
				realPath: p,
				name:     importName,
			})
		}
		p = filepath.Join(v, importName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // directory is package
			realPaths = append(realPaths, &Resource{
				kind:     ResourceKindJavaClass,
				realPath: p,
				name:     importName,
			})
		}
	}
	if len(realPaths) == 0 {
		return nil, fmt.Errorf("resource '%v' not found", importName)
	}
	realPathMap := make(map[string][]*Resource)
	for _, v := range realPaths {
		_, ok := realPathMap[v.realPath]
		if ok {
			realPathMap[v.realPath] = append(realPathMap[v.realPath], v)
		} else {
			realPathMap[v.realPath] = []*Resource{v}
		}
	}
	if len(realPathMap) > 1 {
		errMsg := "not 1 resource named '" + importName + "' present:\n"
		for _, v := range realPathMap {
			switch v[0].kind {
			case ResourceKindJavaClass:
				errMsg += fmt.Sprintf("\t in '%s' is a java class\n", v[0].realPath)
			case ResourceKindJavaPackage:
				errMsg += fmt.Sprintf("\t in '%s' is a java package\n", v[0].realPath)
			case ResourceKindLucyClass:
				errMsg += fmt.Sprintf("\t in '%s' is a lucy class\n", v[0].realPath)
			case ResourceKindLucyPackage:
				errMsg += fmt.Sprintf("\t in '%s' is a lucy package\n", v[0].realPath)
			}
		}
		return nil, fmt.Errorf(errMsg)
	}
	if realPaths[0].kind == ResourceKindLucyClass {
		if filepath.Base(realPaths[0].realPath) == mainClassName {
			return nil, fmt.Errorf("%s is special class for global variable and other things", mainClassName)
		}
	}
	if realPaths[0].kind == ResourceKindJavaClass {
		class, err := loader.loadClass(realPaths[0])
		if class != nil {
			loader.caches[importName] = class
		}
		return class, err
	} else if realPaths[0].kind == ResourceKindLucyClass {
		t, err := loader.loadClass(realPaths[0])
		if t != nil {
			loader.caches[importName] = t
		}
		return t, err
	} else if realPaths[0].kind == ResourceKindJavaPackage {
		p, err := loader.loadJavaPackage(realPaths[0])
		if p != nil {
			loader.caches[importName] = p
		}
		return p, err
	} else { // lucy package
		p, err := loader.loadLucyPackage(realPaths[0])
		if p != nil {
			loader.caches[importName] = p
		}
		return p, err
	}
}

/*
	lucy and java have no difference
*/
func (loader *FileLoader) loadInterfaces(astClass *ast.Class, c *cg.Class) error {
	astClass.InterfaceNames = make([]*ast.NameWithPos, len(c.Interfaces))
	for k, v := range c.Interfaces {
		astClass.InterfaceNames[k] = &ast.NameWithPos{
			Name: string(c.ConstPool[v].Info),
		}
	}
	astClass.Interfaces = make([]*ast.Class, len(astClass.InterfaceNames))
	for k, c := range astClass.InterfaceNames {
		i := &ast.Class{}
		i.Name = c.Name
		i.NotImportedYet = true
		astClass.Interfaces[k] = i
	}
	return nil
}

func (loader *FileLoader) loadAsJava(c *cg.Class) (*ast.Class, error) {
	//name
	if t := c.AttributeGroupedByName.GetByName(cg.AttributeNameSignature); t != nil && len(t) > 0 {
		//TODO:: support signature???
	}
	astClass := &ast.Class{}
	{
		nameIndex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		astClass.Name = string(c.ConstPool[nameIndex].Info)
		if astClass.Name != ast.JavaRootClass {
			nameIndex = binary.BigEndian.Uint16(c.ConstPool[c.SuperClass].Info)
			astClass.SuperClassName = string(c.ConstPool[nameIndex].Info)
		}
	}
	loader.loadInterfaces(astClass, c)
	astClass.AccessFlags = c.AccessFlag
	astClass.IsJava = true // class compiled from java
	var err error
	astClass.Fields = make(map[string]*ast.ClassField)
	astClass.LoadFromOutSide = true
	for _, v := range c.Fields {
		f := &ast.ClassField{}
		f.LoadFromOutSide = true
		f.AccessFlags = v.AccessFlags
		f.JvmDescriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		f.Name = string(c.ConstPool[v.NameIndex].Info)
		_, f.Type, err = jvm.Descriptor.ParseType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		astClass.Fields[f.Name] = f
	}
	astClass.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range c.Methods {
		m := &ast.ClassMethod{}
		m.Function = &ast.Function{}
		m.LoadFromOutSide = true
		m.Function.Name = string(c.ConstPool[v.NameIndex].Info)
		m.Function.JvmDescriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		m.Function.AccessFlags = v.AccessFlags
		m.Function.Type, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		if t := v.AttributeGroupedByName.GetByName(cg.AttributeNameMethodParameters); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, m.Function)
		}
		if (v.AccessFlags & cg.ACC_METHOD_VARARGS) != 0 {
			m.Function.Type.VArgs = m.Function.Type.ParameterList[len(m.Function.Type.ParameterList)-1]
			if m.Function.Type.VArgs.Type.Type != ast.VariableTypeJavaArray {
				panic("variable args is not array")
			}
			m.Function.Type.VArgs.Type.IsVArgs = true
			m.Function.Type.ParameterList = m.Function.Type.ParameterList[:len(m.Function.Type.ParameterList)-1]
		}
		if astClass.Methods[m.Function.Name] == nil {
			astClass.Methods[m.Function.Name] = []*ast.ClassMethod{m}
		} else {
			astClass.Methods[m.Function.Name] = append(astClass.Methods[m.Function.Name], m)
		}
	}
	return astClass, nil
}

func (loader *FileLoader) loadAsLucy(c *cg.Class) (*ast.Class, error) {
	if c.IsSynthetic() {
		return nil, nil
	}
	//name
	astClass := &ast.Class{}
	{
		nameIndex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		astClass.Name = string(c.ConstPool[nameIndex].Info)
		if astClass.Name != ast.JavaRootClass {
			nameIndex = binary.BigEndian.Uint16(c.ConstPool[c.SuperClass].Info)
			astClass.SuperClassName = string(c.ConstPool[nameIndex].Info)
		}
	}
	loader.loadInterfaces(astClass, c)
	astClass.AccessFlags = c.AccessFlag
	astClass.LoadFromOutSide = true
	var err error
	astClass.Fields = make(map[string]*ast.ClassField)
	for _, v := range c.Fields {
		f := &ast.ClassField{}
		f.Name = string(c.ConstPool[v.NameIndex].Info)
		f.JvmDescriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		f.LoadFromOutSide = true
		_, f.Type, err = jvm.Descriptor.ParseType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		if t := v.AttributeGroupedByName.GetByName(cg.AttributeNameLucyFieldDescriptor); t != nil && len(t) > 0 {
			d := &cg.AttributeLucyFieldDescriptor{}
			d.FromBs(c, t[0].Info)
			_, f.Type, err = jvm.LucyFieldSignatureParser.Decode([]byte(d.Descriptor))
			if err != nil {
				return nil, err
			}
			if f.Type.Type == ast.VariableTypeFunction && d.MethodAccessFlag&cg.ACC_METHOD_VARARGS != 0 {
				if f.Type.FunctionType.ParameterList[len(f.Type.FunctionType.ParameterList)-1].Type.Type != ast.VariableTypeJavaArray {
					panic("not a java array")
				}
				f.Type.FunctionType.VArgs = f.Type.FunctionType.ParameterList[len(f.Type.FunctionType.ParameterList)-1]
				f.Type.FunctionType.VArgs.Type.IsVArgs = true
				f.Type.FunctionType.ParameterList = f.Type.FunctionType.ParameterList[:len(f.Type.FunctionType.ParameterList)-1]
			}
		}
		if f.Type.Type == ast.VariableTypeEnum {
			loadEnumForVariableType(f.Type)
		}
		f.AccessFlags = v.AccessFlags
		astClass.Fields[f.Name] = f
	}
	astClass.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range c.Methods {
		m := &ast.ClassMethod{}
		m.Function = &ast.Function{}
		m.Function.Name = string(c.ConstPool[v.NameIndex].Info)
		m.Function.Type, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		m.Function.AccessFlags = v.AccessFlags
		m.LoadFromOutSide = true
		m.Function.JvmDescriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		if t := v.AttributeGroupedByName.GetByName(cg.AttributeNameLucyMethodDescriptor); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			_, err = jvm.LucyMethodSignatureParser.Decode(&m.Function.Type, c.ConstPool[index].Info)
			if err != nil {
				return nil, err
			}
		}
		if t := v.AttributeGroupedByName.GetByName(cg.AttributeNameLucyDefaultParameters); t != nil && len(t) > 0 {
			dp := &cg.AttributeDefaultParameters{}
			dp.FromBytes(t[0].Info)
			jvm.DefaultValueParser.Decode(c, m.Function, dp)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.AttributeNameMethodParameters); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, m.Function)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.AttributeNameLucyReturnListNames); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, m.Function)
		}
		err = loadEnumForFunction(m.Function)
		if err != nil {
			return nil, err
		}
		if (v.AccessFlags & cg.ACC_METHOD_VARARGS) != 0 {
			m.Function.Type.VArgs = m.Function.Type.ParameterList[len(m.Function.Type.ParameterList)-1]
			if m.Function.Type.VArgs.Type.Type != ast.VariableTypeJavaArray {
				panic("variable args is not array")
			}
			m.Function.Type.VArgs.Type.IsVArgs = true
			m.Function.Type.ParameterList = m.Function.Type.ParameterList[:len(m.Function.Type.ParameterList)-1]
		}
		if astClass.Methods[m.Function.Name] == nil {
			astClass.Methods[m.Function.Name] = []*ast.ClassMethod{m}
		} else {
			astClass.Methods[m.Function.Name] = append(astClass.Methods[m.Function.Name], m)
		}
	}
	return astClass, nil
}

func (loader *FileLoader) loadLucyEnum(c *cg.Class) (*ast.Enum, error) {
	e := &ast.Enum{}
	{
		nameIndex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		e.Name = string(c.ConstPool[nameIndex].Info)
	}
	e.AccessFlags = c.AccessFlag
	for _, v := range c.Fields {
		en := &ast.EnumName{}
		name := string(c.ConstPool[v.NameIndex].Info)
		en.Name = name
		en.Enum = e
		constValue := v.AttributeGroupedByName[cg.AttributeNameConstValue][0] // must have this attribute
		en.Value = int32(binary.BigEndian.Uint32(c.ConstPool[binary.BigEndian.Uint16(constValue.Info)].Info))
		e.Enums = append(e.Enums, en)
	}
	return e, nil
}

func (loader *FileLoader) loadLucyMainClass(pack *ast.Package, c *cg.Class) error {
	var err error
	mainClassName := &cg.ClassHighLevel{}
	mainClassName.Name = pack.Name + "/main"
	pack.Block.Variables = make(map[string]*ast.Variable)
	pack.Block.Constants = make(map[string]*ast.Constant)
	pack.Block.Functions = make(map[string]*ast.Function)
	for _, f := range c.Fields {
		name := string(c.ConstPool[f.NameIndex].Info)
		constValue := f.AttributeGroupedByName.GetByName(cg.AttributeNameConstValue)
		if len(constValue) > 1 {
			return fmt.Errorf("constant value length greater than  1 at class 'main'  field '%s'", name)
		}
		_, typ, err := jvm.Descriptor.ParseType(c.ConstPool[f.DescriptorIndex].Info)
		if err != nil {
			return err
		}
		if len(f.AttributeGroupedByName.GetByName(cg.AttributeNameLucyConst)) > 0 {
			//const
			cos := &ast.Constant{}
			cos.Name = name
			cos.AccessFlags = f.AccessFlags
			cos.Type = typ
			_, cos.Type, err = jvm.Descriptor.ParseType(c.ConstPool[f.DescriptorIndex].Info)
			if err != nil {
				return err
			}
			valueIndex := binary.BigEndian.Uint16(constValue[0].Info)
			switch cos.Type.Type {
			case ast.VariableTypeBool:
				cos.Value = binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info) != 0
			case ast.VariableTypeByte:
				cos.Value = byte(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VariableTypeShort:
				fallthrough
			case ast.VariableTypeInt:
				cos.Value = int32(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VariableTypeLong:
				cos.Value = int64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VariableTypeFloat:
				cos.Value = float32(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VariableTypeDouble:
				cos.Value = float64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VariableTypeString:
				valueIndex = binary.BigEndian.Uint16(c.ConstPool[valueIndex].Info) // const_string_info
				cos.Value = string(c.ConstPool[valueIndex].Info)                   // utf 8
			}
			pack.Block.Constants[name] = cos
		} else {
			//global vars
			vd := &ast.Variable{}
			vd.Name = name
			vd.AccessFlags = f.AccessFlags
			vd.JvmDescriptor = string(c.ConstPool[f.DescriptorIndex].Info)
			vd.Type = typ
			vd.IsGlobal = true
			pack.Block.Variables[name] = vd
			if t := f.AttributeGroupedByName.GetByName(cg.AttributeNameLucyFieldDescriptor); t != nil && len(t) > 0 {
				d := &cg.AttributeLucyFieldDescriptor{}
				d.FromBs(c, t[0].Info)
				_, vd.Type, err = jvm.LucyFieldSignatureParser.Decode([]byte(d.Descriptor))
				if err != nil {
					return err
				}
				if vd.Type.Type == ast.VariableTypeFunction && d.MethodAccessFlag&cg.ACC_METHOD_VARARGS != 0 {
					if vd.Type.FunctionType.ParameterList[len(vd.Type.FunctionType.ParameterList)-1].Type.Type !=
						ast.VariableTypeJavaArray {
						panic("not a java array")
					}
					vd.Type.FunctionType.VArgs = vd.Type.FunctionType.ParameterList[len(vd.Type.FunctionType.ParameterList)-1]
					vd.Type.FunctionType.VArgs.Type.IsVArgs = true
					vd.Type.FunctionType.ParameterList = vd.Type.FunctionType.ParameterList[:len(vd.Type.FunctionType.ParameterList)-1]
				}
			}
			if typ.Type == ast.VariableTypeEnum {
				loadEnumForVariableType(typ)
			}
		}
	}

	for _, m := range c.Methods {
		if t := m.AttributeGroupedByName.GetByName(cg.AttributeNameLucyTriggerPackageInit); t != nil && len(t) > 0 {
			pack.TriggerPackageInitMethodName = string(c.ConstPool[m.NameIndex].Info)
			continue
		}
		name := string(c.ConstPool[m.NameIndex].Info)
		if name == ast.MainFunctionName {
			// this is main function
			continue
		}
		function := &ast.Function{}
		function.Name = name
		function.AccessFlags = m.AccessFlags
		function.JvmDescriptor = string(c.ConstPool[m.DescriptorIndex].Info)
		function.Type, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[m.DescriptorIndex].Info)
		if err != nil {
			return err
		}
		if t := m.AttributeGroupedByName.GetByName(cg.AttributeNameLucyMethodDescriptor); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			_, err = jvm.LucyMethodSignatureParser.Decode(&function.Type, c.ConstPool[index].Info)
			if err != nil {
				return err
			}
		}
		err = loadEnumForFunction(function)
		if err != nil {
			return err
		}
		if t := m.AttributeGroupedByName.GetByName(cg.AttributeNameMethodParameters); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, function)
		}
		if t := m.AttributeGroupedByName.GetByName(cg.AttributeNameLucyReturnListNames); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, function)
		}
		if t := m.AttributeGroupedByName.GetByName(cg.AttributeNameLucyDefaultParameters); t != nil && len(t) > 0 {
			dp := &cg.AttributeDefaultParameters{}
			dp.FromBytes(t[0].Info)
			jvm.DefaultValueParser.Decode(c, function, dp)
		}

		if (function.AccessFlags & cg.ACC_METHOD_VARARGS) != 0 {
			function.Type.VArgs = function.Type.ParameterList[len(function.Type.ParameterList)-1]
			if function.Type.VArgs.Type.Type != ast.VariableTypeJavaArray {
				panic("variable args is not array")
			}
			function.Type.VArgs.Type.IsVArgs = true
			function.Type.ParameterList = function.Type.ParameterList[:len(function.Type.ParameterList)-1]
		}

		function.Entrance = &cg.MethodHighLevel{}
		function.Entrance.Name = function.Name
		function.Entrance.Class = mainClassName
		function.Entrance.Descriptor = function.JvmDescriptor
		function.IsGlobal = true
		pack.Block.Functions[name] = function
	}

	if pack.Block.TypeAliases == nil {
		pack.Block.TypeAliases = make(map[string]*ast.Type)
	}
	for _, v := range c.AttributeGroupedByName.GetByName(cg.AttributeNameLucyTypeAlias) {
		index := binary.BigEndian.Uint16(v.Info)
		name, typ, err := jvm.LucyTypeAliasParser.Decode(c.ConstPool[index].Info)
		if err != nil {
			return err
		}
		typ.Alias = name
		pack.Block.TypeAliases[name] = typ
		if typ.Type == ast.VariableTypeEnum {
			err = loadEnumForVariableType(typ)
			if err != nil {
				return err
			}
		}
	}
	for _, v := range c.AttributeGroupedByName.GetByName(cg.AttributeNameLucyTemplateFunction) {
		attr := &cg.AttributeTemplateFunction{}
		attr.FromBytes(c, v.Info)
		f, es := ast.ParseFunctionHandler([]byte(attr.Code), &ast.Pos{
			Filename:    attr.Filename,
			StartLine:   int(attr.StartLine),
			StartColumn: int(attr.StartColumn),
		})
		if len(es) > 0 { // looks impossible
			return es[0]
		}
		f.AccessFlags = attr.AccessFlag
		f.TemplateFunction = &ast.TemplateFunction{}
		pack.Block.Functions[attr.Name] = f
	}
	return nil
}

func (loader *FileLoader) loadLucyPackage(r *Resource) (*ast.Package, error) {
	fis, err := ioutil.ReadDir(r.realPath)
	if err != nil {
		return nil, err
	}
	fisM := make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") {
			fisM[v.Name()] = v
		}
	}

	_, ok := fisM[mainClassName]
	if ok == false {
		return nil, fmt.Errorf("main class not found")
	}
	bs, err := ioutil.ReadFile(filepath.Join(r.realPath, mainClassName))
	if err != nil {
		return nil, fmt.Errorf("read main.class failed,err:%v", err)
	}
	c, err := (&ClassDecoder{}).decode(bs)
	if err != nil {
		return nil, fmt.Errorf("decode main class failed,err:%v", err)
	}
	p := &ast.Package{}
	p.Name = r.name
	err = loader.loadLucyMainClass(p, c)
	if err != nil {
		return nil, fmt.Errorf("parse main class failed,err:%v", err)
	}

	delete(fisM, mainClassName)
	mkEnums := func(e *ast.Enum) {
		if p.Block.Enums == nil {
			p.Block.Enums = make(map[string]*ast.Enum)
		}
		if p.Block.EnumNames == nil {
			p.Block.EnumNames = make(map[string]*ast.EnumName)
		}
		p.Block.Enums[filepath.Base(e.Name)] = e
		for _, v := range e.Enums {
			p.Block.EnumNames[v.Name] = v
		}
	}
	for _, v := range fisM {

		bs, err := ioutil.ReadFile(filepath.Join(r.realPath, v.Name()))
		if err != nil {
			return p, fmt.Errorf("read class failed,err:%v", err)
		}
		c, err := (&ClassDecoder{}).decode(bs)
		if err != nil {
			return nil, fmt.Errorf("decode class failed,err:%v", err)
		}
		if len(c.AttributeGroupedByName.GetByName(cg.AttributeNameLucyEnum)) > 0 {
			e, err := loader.loadLucyEnum(c)
			if err != nil {
				return nil, err
			}
			mkEnums(e)
			continue
		}
		class, err := loader.loadAsLucy(c)
		if err != nil {
			return nil, fmt.Errorf("decode class failed,err:%v", err)
		}
		if p.Block.Classes == nil {
			p.Block.Classes = make(map[string]*ast.Class)
		}
		p.Block.Classes[filepath.Base(class.Name)] = class
	}
	return p, nil
}

func (loader *FileLoader) loadJavaPackage(r *Resource) (*ast.Package, error) {
	fis, err := ioutil.ReadDir(r.realPath)
	if err != nil {
		return nil, err
	}
	ret := &ast.Package{}
	ret.Block.Classes = make(map[string]*ast.Class)
	for _, v := range fis {
		var rr Resource
		rr.kind = ResourceKindJavaClass
		if strings.HasSuffix(v.Name(), ".class") == false || strings.Contains(v.Name(), "$") {
			continue
		}
		rr.realPath = filepath.Join(r.realPath, v.Name())
		class, err := loader.loadClass(&rr)
		if err != nil {
			return nil, err
		}
		if c, ok := class.(*ast.Class); ok && class != nil {
			ret.Block.Classes[filepath.Base(c.Name)] = c
		}
	}
	return ret, nil
}

func (loader *FileLoader) loadClass(r *Resource) (interface{}, error) {
	bs, err := ioutil.ReadFile(r.realPath)
	if err != nil {
		return nil, err
	}
	c, err := (&ClassDecoder{}).decode(bs)
	if r.kind == ResourceKindLucyClass {
		if t := c.AttributeGroupedByName[cg.AttributeNameLucyEnum]; len(t) > 0 {
			return loader.loadLucyEnum(c)
		} else {
			return loader.loadAsLucy(c)
		}
	}
	return loader.loadAsJava(c)
}
