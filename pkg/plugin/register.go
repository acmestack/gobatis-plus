// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

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
