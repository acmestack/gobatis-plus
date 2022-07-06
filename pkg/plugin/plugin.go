// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package plugin

import (
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

type Plugin interface {
	Annotation() string
	CouldHandle(t *types.Type) bool
	Generate(ctx *generator.Context, w io.Writer, t *types.Type) (err error)
}
