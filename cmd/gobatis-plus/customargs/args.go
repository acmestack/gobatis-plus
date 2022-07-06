// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package customargs

import (
	"fmt"
	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
)

type gobatisArgs struct {
	Prefix string
}

func NewDefault() (*args.GeneratorArgs, *gobatisArgs) {
	args := args.Default().WithoutDefaultFlagParsing()
	cusArgs := &gobatisArgs{
		Prefix: "gobatis",
	}
	args.CustomArgs = cusArgs
	args.OutputFileBaseName = "zz_generated"
	return args, cusArgs
}

func (arg *gobatisArgs) String() string {
	return arg.Prefix
}

func (arg *gobatisArgs) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&arg.Prefix, "annotation", arg.Prefix, "Annotation name")
}

func Validate(args *args.GeneratorArgs) error {
	_ = args.CustomArgs.(*gobatisArgs)

	if len(args.OutputFileBaseName) == 0 {
		return fmt.Errorf("Output file base name cannot be empty. ")
	}
	if len(args.InputDirs) == 0 {
		return fmt.Errorf("Input directory cannot be empty. ")
	}
	return nil
}
