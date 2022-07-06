/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
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

import "k8s.io/gengo/types"

var (
	gPlugins = []Plugin{
		&mapperPlugin{},
		&dataPlugin{},
	}
)

func FindPlugin(t *types.Type) Plugin {
	for i := len(gPlugins) - 1; i >= 0; i-- {
		v := gPlugins[i]
		if v.CouldHandle(t) {
			return v
		}
	}
	return nil
}

func RegisterPlugin(plugin Plugin) {
	gPlugins = append(gPlugins, plugin)
}
