package astutil

import (
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gosrob/autumn/internal/util/pkginfo"
	"golang.org/x/tools/go/packages"
)

func LoadFileAST(path string) (*ast.File, error) {
	return load(path)
}

func FindTopNodeByName(f *ast.File, name string) (ast.Node, []ast.Node, bool) {
	return findHighLevelNode(f, name)
}

func Comment(n ast.Node) (string, bool) {
	return extractComment(n)
}

func Walk(f *ast.File, visitor func(node ast.Node, parents []ast.Node) bool) {
	walk(f, visitor)
}

// ====
//
//
//

// 充当一个简易的缓存结构
type CheckerCache struct {
	mu sync.RWMutex
	// 存储检查结果的映射，key为类型和接口的标识符，value为检查结果
	cache map[string]bool
}

var cache = CheckerCache{
	mu:    sync.RWMutex{},
	cache: map[string]bool{},
}

var pkgsCache []*packages.Package

func NewCheckerCache() *CheckerCache {
	return &CheckerCache{
		cache: make(map[string]bool),
	}
}

func (cc *CheckerCache) get(key string) (bool, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	result, exists := cc.cache[key]
	return result, exists
}

func (cc *CheckerCache) set(key string, value bool) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.cache[key] = value
}

// 提取包名和类型名
func extractPackageNameFromFullIdentity(fullIdentity string) (string, string) {
	lastSlash := strings.LastIndex(fullIdentity, ".")
	if lastSlash == -1 {
		return "", ""
	}
	pkgName := fullIdentity[:lastSlash]
	typeName := fullIdentity[lastSlash+1:]
	return pkgName, typeName
}

// 收集目录中的所有Go文件，过滤掉测试文件
func collectGoFilesInDirectory(dir string) ([]string, error) {
	var goFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 过滤条件: 不是目录且是Go文件，但不是测试文件
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	return goFiles, err
}

// 主要的函数：检查类型是否实现了接口，使用缓存
// 使用 packages 包进行类型检查
func CheckIfTypeImplementsInterfaceWithCache(typeFullIdentity, interfaceFullIdentity string) (bool, error) {
	dir := pkginfo.GetFullPackage("").Module.Dir
	cacheKey := fmt.Sprintf("%s implements %s", typeFullIdentity, interfaceFullIdentity)
	// 检查缓存
	if result, exists := cache.get(cacheKey); exists {
		return result, nil
	}

	// 提取包名和类型名
	pkgType, typeName := extractPackageNameFromFullIdentity(typeFullIdentity)
	pkgInterface, interfaceName := extractPackageNameFromFullIdentity(interfaceFullIdentity)

	// 使用 packages 包加载目录中的所有Go文件
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Dir:   dir,
		Tests: true,
	}
	if pkgsCache == nil {
		pkgs, err := packages.Load(cfg, "./...")
		if err != nil || len(pkgs) == 0 {
			return false, fmt.Errorf("failed to load packages: %s", err)
		}
		pkgsCache = pkgs

	}
	pkgs := pkgsCache

	// // 构建 AST 和 type 信息
	// pkgMap := make(map[string]*packages.Package)
	// for _, pkg := range pkgs {
	// 	pkgMap[pkg.PkgPath] = pkg
	// 	for range pkg.Syntax {
	// 		// 使用 types 信息分析
	// 		typesInfo := pkg.TypesInfo
	// 		for _, obj := range typesInfo.Defs {
	// 			// if obj != nil {
	// 			// 	logger.Logger.Debugf("struct name %s vs checkName %s, struct pkgPath %s vs checkPath %s, ", obj.Name(), typeName, pkg.PkgPath, pkgType)
	// 			// }
	// 			// 检查类型
	// 			if obj != nil && obj.Name() == typeName && pkg.PkgPath == pkgType {
	// 				typ := obj.Type()
	// 				for _, obj := range typesInfo.Defs {
	// 					// 检查接口
	// 					if obj == nil || obj.Type() == nil || obj.Pkg() == nil {
	// 						continue
	// 					}
	//
	// 					if obj != nil {
	// 						name := obj.Name()
	// 						name = name
	// 						logger.Logger.Debugf("name %s interfaceName %s, path %s pkgInterface %s", name, interfaceName, obj.Pkg().Path(), pkgInterface)
	// 					}
	// 					if obj.Name() == interfaceName && obj.Pkg().Path() == pkgInterface {
	// 						if obj != nil {
	// 							itype, ok := obj.Type().Underlying().(*types.Interface)
	// 							if ok {
	// 								logger.Logger.Debugf("interface name %s vs checkName %s, interface pkgPath %s vs checkPath %s, itype %s", obj.Name(), interfaceName, pkg.PkgPath, pkgInterface, itype)
	// 							}
	// 						}
	// 						if itype, ok := obj.Type().Underlying().(*types.Interface); ok && types.Implements(types.NewPointer(typ), itype) {
	// 							cache.set(cacheKey, true)
	// 							return true, nil
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	/// ===

	var typ types.Type

	var iface *types.Interface
	foundStruct := false
	for _, pkg := range pkgs {
		for range pkg.Syntax {
			if foundStruct {
				break
			}
			typesInfo := pkg.TypesInfo
			for _, obj := range typesInfo.Defs {
				if obj == nil {
					continue
				}
				if obj.Pkg().Path() == pkgType && obj.Name() == typeName {
					typ = obj.Type()
					foundStruct = true
					break
				}
			}
		}
		if foundStruct {
			break
		}
	}

	foundInterface := false
	for _, pkg := range pkgs {
		if foundInterface {
			break
		}
		for range pkg.Syntax {
			if foundInterface {
				break
			}
			typesInfo := pkg.TypesInfo
			for _, obj := range typesInfo.Defs {
				if obj == nil {
					continue
				}
				if obj.Pkg().Path() == pkgInterface && obj.Name() == interfaceName {
					for _, ifaceObj := range typesInfo.Defs {
						if ifaceObj == nil {
							continue
						}
						// logger.Logger.Infof("path %s", ifaceObj.Pkg().Path())
						if ifaceObj.Pkg().Path() == pkgInterface && ifaceObj.Name() == interfaceName {
							if ifaces, ok := ifaceObj.Type().Underlying().(*types.Interface); ok {
								iface = ifaces
								foundInterface = true
								break
							}
						}
					}
				}
			}
		}
	}

	if iface != nil && typ != nil {
		if types.Implements(types.NewPointer(typ), iface) {
			cache.set(cacheKey, true)
			return true, nil
		}
	}

	// ===

	cache.set(cacheKey, false)
	return false, nil
}
