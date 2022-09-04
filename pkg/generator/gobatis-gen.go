/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

import (
	"fmt"
	"io"
	"k8s.io/klog/v2"
	"log"
	"path/filepath"
	"strings"

	"github.com/acmestack/gobatis-plus/pkg/plugin"
	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

var (
	gobatisImports = []string{"github.com/acmestack/gobatis"}
)

type gobatisAnnotion struct {
	rawTypeName string
	body        string
}

type gobatisAnnotions map[string]*gobatisAnnotion

type gobatisGen struct {
	name      string
	prefix    string
	targetPkg string
	pkg       *types.Package
	imports   namer.ImportTracker
}

func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPrivateNamer(0, ""),
		"raw":    namer.NewRawNamer("", nil),
	}
}

func DefaultNameSystem() string {
	return "public"
}

func checkEnable(annotation string, comments []string) bool {
	key := annotation + "enable"
	for _, c := range comments {
		if strings.HasPrefix(c, key) {
			return true
		}
	}
	return false
}

func GenPackages(ctx *generator.Context, args *args.GeneratorArgs) generator.Packages {
	inputs := sets.NewString(ctx.Inputs...)
	pkgs := generator.Packages{}
	annotation := args.CustomArgs.(fmt.Stringer).String()

	log.Println(args)
	boilerplate, err := args.LoadGoBoilerplate()
	if err != nil {
		klog.Warningf("LoadGoBoilerplate failed: %v. ", err)
		boilerplate = nil
	}
	header := []byte(fmt.Sprintf("// +build !%s\n\n", args.GeneratedBuildTag))
	if boilerplate != nil {
		header = append(header, boilerplate...)
	}

	for in := range inputs {
		klog.V(5).Infof("Parsing pkg %s\n", in)
		pkg := ctx.Universe[in]
		if pkg == nil {
			continue
		}
		for _, i := range pkg.Imports {
			ctx.AddDirectory(i.Path)
		}

		if !checkEnable(annotation, pkg.Comments) {
			continue
		}

		klog.V(5).Infof("Generating package %s...\n", in)

		pkgs = append(pkgs, &generator.DefaultPackage{
			PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
			PackagePath: pkg.Path,
			HeaderText:  header,
			GeneratorFunc: func(context *generator.Context) []generator.Generator {
				return []generator.Generator{
					NewGobatisGenerator(args.OutputFileBaseName, annotation, pkg),
				}
			},
			FilterFunc: func(context *generator.Context, i *types.Type) bool {
				return i.Name.Package == pkg.Path
			},
		})
	}
	return pkgs
}

func NewGobatisGenerator(name, prefix string, pkg *types.Package) *gobatisGen {
	ret := &gobatisGen{
		name:    name,
		prefix:  prefix,
		pkg:     pkg,
		imports: generator.NewImportTracker(),
	}

	return ret
}

// The name of this generator. Will be included in generated comments.
func (g *gobatisGen) Name() string {
	return g.name
}

// Filter should return true if this generator cares about this type.
// (otherwise, GenerateType will not be called.)
//
// Filter is called before any of the generator's other functions;
// subsequent calls will get a context with only the types that passed
// this filter.
func (g *gobatisGen) Filter(ctx *generator.Context, t *types.Type) bool {
	return true
}

// If this generator needs special namers, return them here. These will
// override the original namers in the context if there is a collision.
// You may return nil if you don't need special names. These names will
// be available in the context passed to the rest of the generator's
// functions.
//
// A use case for this is to return a namer that tracks imports.
func (g *gobatisGen) Namers(ctx *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.targetPkg, g.imports),
	}
}

// Init should write an init function, and any other content that's not
// generated per-type. (It's not intended for generator specific
// initialization! Do that when your Package constructs the
// Generators.)
func (g *gobatisGen) Init(ctx *generator.Context, w io.Writer) error {
	return nil
}

// Finalize should write finish up functions, and any other content that's not
// generated per-type.
func (g *gobatisGen) Finalize(ctx *generator.Context, w io.Writer) error {
	return nil
}

// PackageVars should emit an array of variable lines. They will be
// placed in a var ( ... ) block. There's no need to include a leading
// \t or trailing \n.
func (g *gobatisGen) PackageVars(ctx *generator.Context) []string {
	return nil
}

// PackageConsts should emit an array of constant lines. They will be
// placed in a const ( ... ) block. There's no need to include a leading
// \t or trailing \n.
func (g *gobatisGen) PackageConsts(ctx *generator.Context) []string {
	return nil
}

// GenerateType should emit the code for a particular type.
func (g *gobatisGen) GenerateType(ctx *generator.Context, t *types.Type, w io.Writer) error {
	p := plugin.FindPlugin(t)
	if p == nil {
		return fmt.Errorf("Cannot handle type: %s. ", t.String())
	}
	return p.Generate(ctx, w, t)
}

func parseAnnotations(annotation string, t *types.Type) gobatisAnnotions {
	//ret := gobatisAnnotions{}
	//t.CommentLines
	return nil
}

// Imports should return a list of necessary imports. They will be
// formatted correctly. You do not need to include quotation marks,
// return only the package name; alternatively, you can also return
// imports in the format `name "path/to/pkg"`. Imports will be called
// after Init, PackageVars, PackageConsts, and GenerateType, to allow
// you to keep track of what imports you actually need.
func (g *gobatisGen) Imports(ctx *generator.Context) []string {
	imports := g.imports.ImportLines()
	imports = append(imports, gobatisImports...)
	return imports
}

// Preferred file name of this generator, not including a path. It is
// allowed for multiple generators to use the same filename, but it's
// up to you to make sure they don't have colliding import names.
// TODO: provide per-file import tracking, removing the requirement
// that generators coordinate..
func (g *gobatisGen) Filename() string {
	return g.name + ".go"
}

// A registered file type in the context to generate this file with. If
// the FileType is not found in the context, execution will stop.
func (g *gobatisGen) FileType() string {
	return generator.GolangFileType
}
