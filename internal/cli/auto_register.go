package cli

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const modulesFilePath = "internal/modules/modules.go"

func AutoRegisterModule(module string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, modulesFilePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse modules.go: %w", err)
	}

	infraAlias := module + "Infra"
	usecaseAlias := module + "Usecase"

	infraImportPath := fmt.Sprintf(
		"github.com/yuristian/go-api/internal/modules/%s/infrastructure",
		module,
	)
	usecaseImportPath := fmt.Sprintf(
		"github.com/yuristian/go-api/internal/modules/%s/usecase",
		module,
	)

	addImportIfNotExists(file, infraImportPath, infraAlias)
	addImportIfNotExists(file, usecaseImportPath, usecaseAlias)

	if err := addRegisterBlockIfNotExists(file, module); err != nil {
		return err
	}

	out, err := os.Create(modulesFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	return format.Node(out, fset, file)
}

func RemoveModuleFromRegistry(module string) error {
	path := "internal/modules/modules.go"

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	removeImports(file, module)
	removeRegisterBlock(file, module)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	return format.Node(out, fset, file)
}

/* =========================
   IMPORT HANDLING
   ========================= */

func addImportIfNotExists(file *ast.File, path, alias string) {
	// 1. cek apakah import sudah ada
	for _, imp := range file.Imports {
		if strings.Trim(imp.Path.Value, `"`) == path {
			return
		}
	}

	spec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + path + `"`,
		},
	}

	if alias != "" {
		spec.Name = ast.NewIdent(alias)
	}

	// 2. cari import declaration
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}

		gen.Specs = append(gen.Specs, spec)
		file.Imports = append(file.Imports, spec)
		return
	}

	// 3. jika belum ada import block → buat baru
	gen := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{spec},
	}

	file.Decls = append([]ast.Decl{gen}, file.Decls...)
	file.Imports = append(file.Imports, spec)
}

func removeImports(file *ast.File, module string) {
	var newDecls []ast.Decl

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			newDecls = append(newDecls, decl)
			continue
		}

		var specs []ast.Spec
		for _, spec := range gen.Specs {
			imp := spec.(*ast.ImportSpec)
			if !strings.Contains(imp.Path.Value, "/modules/"+module+"/") {
				specs = append(specs, spec)
			}
		}

		if len(specs) > 0 {
			gen.Specs = specs
			newDecls = append(newDecls, gen)
		}
	}

	file.Decls = newDecls
}

/* =========================
   REGISTER BLOCK HANDLING
   ========================= */

func addRegisterBlockIfNotExists(file *ast.File, module string) error {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "RegisterAllModules" {
			continue
		}

		// idempotency check
		for _, stmt := range fn.Body.List {
			if containsModule(stmt, module) {
				return nil
			}
		}

		block := buildRegisterBlock(module)
		fn.Body.List = append(fn.Body.List, block...)
		return nil
	}

	return fmt.Errorf("RegisterAllModules function not found")
}

func buildRegisterBlock(module string) []ast.Stmt {
	cap := capitalize(module)

	repoVar := module + "Repo"
	ucVar := module + "UC"

	infraAlias := module + "Infra"
	usecaseAlias := module + "Usecase"

	return []ast.Stmt{
		// repo := moduleInfra.NewXGormRepository(gormDB)
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(repoVar)},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent(infraAlias),
						Sel: ast.NewIdent("New" + cap + "GormRepository"),
					},
					Args: []ast.Expr{
						ast.NewIdent("gormDB"),
					},
				},
			},
		},
		// uc := moduleUsecase.NewXUsecase(repo, jwtManager)
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent(ucVar)},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent(usecaseAlias),
						Sel: ast.NewIdent("New" + cap + "Usecase"),
					},
					Args: []ast.Expr{
						ast.NewIdent(repoVar),
						// ast.NewIdent("jwtManager"),
					},
				},
			},
		},
		// moduleInfra.RegisterRoutes(rg, uc)
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(infraAlias),
					Sel: ast.NewIdent("RegisterRoutes"),
				},
				Args: []ast.Expr{
					ast.NewIdent("rg"),
					ast.NewIdent(ucVar),
				},
			},
		},
	}
}

func containsModule(stmt ast.Stmt, module string) bool {
	return strings.Contains(fmt.Sprintf("%#v", stmt), module)
}

func removeRegisterBlock(file *ast.File, module string) {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "RegisterAllModules" {
			continue
		}

		var newStmts []ast.Stmt

		for _, stmt := range fn.Body.List {
			if isModuleConstructorStmt(stmt, module) {
				continue
			}

			if isModuleRegisterRoutesStmt(stmt, module) {
				continue
			}

			newStmts = append(newStmts, stmt)
		}

		fn.Body.List = newStmts
	}
}

/* =========================
   UTIL
   ========================= */

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// func astToString(node ast.Node) string {
// 	var buf bytes.Buffer
// 	_ = printer.Fprint(&buf, token.NewFileSet(), node)
// 	return buf.String()
// }

func isModuleConstructorStmt(stmt ast.Stmt, module string) bool {
	assign, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return false
	}

	for _, rhs := range assign.Rhs {
		call, ok := rhs.(*ast.CallExpr)
		if !ok {
			continue
		}

		switch fun := call.Fun.(type) {

		// Case: NewProduct(...)
		case *ast.Ident:
			if fun.Name == "New"+capitalize(module) {
				return true
			}

		// Case: productUsecase.NewProductUsecase(...)
		//       productInfra.NewProductGormRepository(...)
		case *ast.SelectorExpr:
			if strings.HasPrefix(fun.Sel.Name, "New"+capitalize(module)) {
				return true
			}
		}
	}

	return false
}

func isModuleRegisterRoutesStmt(stmt ast.Stmt, module string) bool {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}

	call, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkg, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	expectedPkg := module + "Infra"

	return pkg.Name == expectedPkg &&
		selector.Sel.Name == "RegisterRoutes"
}
