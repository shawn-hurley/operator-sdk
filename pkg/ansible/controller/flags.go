package controller

import (
	"flag"

	"github.com/spf13/pflag"
)

func AddPflags(fs *pflag.FlagSet) *Options {
	return &Options{}
}

func AddGoFlags(fs *flag.FlagSet) *Options {
	o := &Options{}
	i := o.LoggingLevel
	fs.Int("runner-loglevel", int(i), "The level for the ansible operator event logging, 0 for tasks, 1 for everyting and >1 for nothing")
	return o
}
