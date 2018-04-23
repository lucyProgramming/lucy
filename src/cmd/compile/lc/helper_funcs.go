package lc

// import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

// func loadClass(v *ast.VariableType) error {
// 	switch v.Typ {
// 	case ast.VARIABLE_TYPE_OBJECT:
// 		_, c, err := loader.LoadName(v.Class.Name)
// 		if err != nil {
// 			return err
// 		}
// 		v.Class = c.(*ast.Class)
// 	case ast.VARIABLE_TYPE_MAP:
// 		err := loadClass(v.Map.K)
// 		if err != nil {
// 			return err
// 		}
// 		return loadClass(v.Map.V)
// 	case ast.VARIABLE_TYPE_ARRAY:
// 		return loadClass(v.ArrayType)
// 	case ast.VARIABLE_TYPE_JAVA_ARRAY:
// 		return loadClass(v.ArrayType)
// 	}
// 	return nil
// }
