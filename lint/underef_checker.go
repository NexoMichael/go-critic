package lint

//! Detects dereference expressions that can be omitted.
//
// @Before:
// (*k).field = 5
// _ := (*a)[5] // only if a is array
//
// @After:
// k.field = 5
// _ := a[5]

import (
	"go/ast"
	"go/types"
)

func init() {
	addChecker(&underefChecker{})
}

type underefChecker struct {
	checkerBase
}

func (c *underefChecker) VisitExpr(expr ast.Expr) {
	switch n := expr.(type) {
	case *ast.SelectorExpr:
		expr, ok := n.X.(*ast.ParenExpr)
		if !ok {
			return
		}
		if expr, ok := expr.X.(*ast.StarExpr); ok {
			if c.checkStarExpr(expr) {
				c.warnSelect(n)
			}
		}
	case *ast.IndexExpr:
		expr, ok := n.X.(*ast.ParenExpr)
		if !ok {
			return
		}
		if expr, ok := expr.X.(*ast.StarExpr); ok {
			if !c.checkStarExpr(expr) {
				return
			}
			if c.checkArray(expr) {
				c.warnArray(n)
			}
		}
	}
}

func (c *underefChecker) warnSelect(expr *ast.SelectorExpr) {
	// TODO: add () to function output.
	c.ctx.Warn(expr, "could simplify %s to %s.%s",
		expr,
		expr.X.(*ast.ParenExpr).X.(*ast.StarExpr).X,
		expr.Sel.Name)
}

func (c *underefChecker) warnArray(expr *ast.IndexExpr) {
	c.ctx.Warn(expr, "could simplify %s to %s[%s]",
		expr,
		expr.X.(*ast.ParenExpr).X.(*ast.StarExpr).X,
		expr.Index)
}

// checkStarExpr checks if ast.StarExpr could be simplified.
func (c *underefChecker) checkStarExpr(expr *ast.StarExpr) bool {
	// Checks if expr is pointer type.
	typ, ok := c.ctx.typesInfo.TypeOf(expr.X).(*types.Pointer)
	if !ok {
		return false
	}

	// Checks if dereference of typ is pointer.
	if _, ok := typ.Elem().(*types.Pointer); ok {
		return false
	}

	return true
}

func (c *underefChecker) checkArray(expr *ast.StarExpr) bool {
	typ, ok := c.ctx.typesInfo.TypeOf(expr.X).(*types.Pointer)
	if !ok {
		return false
	}
	_, ok = typ.Elem().(*types.Array)
	return ok
}
