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
		Prefix: "+gobatis:",
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
