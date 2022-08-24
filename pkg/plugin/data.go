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

package plugin

import (
	"io"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
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
