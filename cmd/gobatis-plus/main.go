// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package main

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/xfali/gobatis-plus/cmd/gobatis-plus/customargs"
	"github.com/xfali/gobatis-plus/pkg/generator"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	args, cusArgs := customargs.NewDefault()

	args.AddFlags(pflag.CommandLine)
	cusArgs.AddFlags(pflag.CommandLine)

	_ = flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := customargs.Validate(args); err != nil {
		klog.Fatalln(err)
	}

	if err := args.Execute(generator.NameSystems(), generator.DefaultNameSystem(), generator.GenPackages); err != nil {
		klog.Fatalln(err)
	}
	klog.V(2).Info("Success")
}
