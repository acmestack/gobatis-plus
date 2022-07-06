// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package plugin

import (
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
	"strings"
)

type dataPlugin struct {
}

func (p *dataPlugin) Annotation() string {
	return "+gobatis:data"
}

func (p *dataPlugin) CouldHandle(t *types.Type) bool {
	if t.Kind == types.Struct {
		for _, c := range t.CommentLines {
			i := strings.Index(c, p.Annotation())
			if i != -1 {
				return true
			}
		}
	}
	return false
}

func (p *dataPlugin) Generate(ctx *generator.Context, w io.Writer, t *types.Type) (err error) {
	//sw := generator.NewSnippetWriter(w, ctx, delimiterLeft, delimiterRight)
	return nil
}
