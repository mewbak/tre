package compiler

import (
	"github.com/zegl/tre/compiler/compiler/strings"
	"github.com/zegl/tre/compiler/compiler/types"
	"github.com/zegl/tre/compiler/compiler/value"
	"github.com/zegl/tre/compiler/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
)

func (c *Compiler) compileConstantNode(v parser.ConstantNode) value.Value {
	switch v.Type {
	case parser.NUMBER:
		var intType types.Type = i64

		// Use context to detect which integer type we should create
		if len(c.contextAssignDest) > 0 {
			intType = c.contextAssignDest[len(c.contextAssignDest)-1].Type
		}

		return value.Value{
			Value:      constant.NewInt(v.Value, intType.LLVM()),
			Type:       i64,
			IsVariable: false,
		}

	case parser.STRING:
		var constString *ir.Global

		// Reuse the *ir.Global if it has already been created
		if reusedConst, ok := c.stringConstants[v.ValueStr]; ok {
			constString = reusedConst
		} else {
			constString = c.module.NewGlobalDef(strings.NextStringName(), strings.Constant(v.ValueStr))
			constString.IsConst = true
			c.stringConstants[v.ValueStr] = constString
		}

		alloc := c.contextBlock.NewAlloca(typeConvertMap["string"].LLVM())

		// Save length of the string
		lenItem := c.contextBlock.NewGetElementPtr(alloc, constant.NewInt(0, i32.LLVM()), constant.NewInt(0, i32.LLVM()))
		c.contextBlock.NewStore(constant.NewInt(int64(len(v.ValueStr)), i64.LLVM()), lenItem)

		// Save i8* version of string
		strItem := c.contextBlock.NewGetElementPtr(alloc, constant.NewInt(0, i32.LLVM()), constant.NewInt(1, i32.LLVM()))
		c.contextBlock.NewStore(strings.Toi8Ptr(c.contextBlock, constString), strItem)

		return value.Value{
			Value:      c.contextBlock.NewLoad(alloc),
			Type:       types.String,
			IsVariable: false,
		}

	case parser.BOOL:
		return value.Value{
			Value:      constant.NewInt(v.Value, i1.LLVM()),
			Type:       i64,
			IsVariable: false,
		}

	default:
		panic("Unknown constant Type")
	}
}
