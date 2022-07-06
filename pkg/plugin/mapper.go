package plugin

import (
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
	"strings"
)

type mapperPlugin struct {
}

func (p *mapperPlugin) Annotation() string {
	return "+gobatis:mapper"
}

func (p *mapperPlugin) CouldHandle(t *types.Type) bool {
	if t.Kind == types.Interface {
		for _, c := range t.CommentLines {
			i := strings.Index(c, p.Annotation())
			if i != -1 {
				return true
			}
		}
	}
	return false
}

func (p *mapperPlugin) Generate(ctx *generator.Context, w io.Writer, t *types.Type) (err error) {
	//sw := generator.NewSnippetWriter(w, ctx, delimiterLeft, delimiterRight)
	return nil
}
