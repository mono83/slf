package slfcobra

import (
	"github.com/mono83/slf"
	"github.com/mono83/slf/recievers/writer"
	"github.com/mono83/slf/wd"
	"github.com/spf13/cobra"
)

var commonLoggingPredicate = func(e slf.Event) bool {
	return e.IsLog() && e.Type >= slf.TypeError
}
var verboseLoggingPredicate = func(e slf.Event) bool {
	return e.IsLog() && e.Type >= slf.TypeInfo
}
var allLoggingPredicate = func(e slf.Event) bool {
	return e.IsLog()
}

// Wrap method wraps provided cobra command, adding support for stdout logging
func Wrap(cmd *cobra.Command) *cobra.Command {
	if cmd.PersistentFlags().Lookup("verbose") == nil {
		cmd.PersistentFlags().BoolP("verbose", "v", false, "Display info-level logging and higher")
	}
	if cmd.PersistentFlags().Lookup("vv") == nil {
		cmd.PersistentFlags().Bool("vv", false, "Very verbose mode, trace and debug will be displayed")
	}
	if cmd.PersistentFlags().Lookup("quiet") == nil {
		cmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode, logging output will be suppressed")
	}
	if cmd.PersistentFlags().Lookup("no-ansi") == nil {
		cmd.PersistentFlags().Bool("no-ansi", false, "Disable ANSI coloring for logs")
	}

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		vv, _ := cmd.Flags().GetBool("vv")
		verbose, _ := cmd.Flags().GetBool("verbose")
		quiet, _ := cmd.Flags().GetBool("quiet")
		nocolor, _ := cmd.Flags().GetBool("no-ansi")
		// Enabling logger
		if !quiet {
			var predicate func(slf.Event) bool
			if vv {
				// Very verbose mode
				predicate = allLoggingPredicate
			} else if verbose {
				// Info+ logging
				predicate = verboseLoggingPredicate
			} else {
				// Error+ logging
				predicate = commonLoggingPredicate
			}

			wd.AddReceiver(slf.Filter(
				writer.New(writer.Options{
					NoColor: nocolor,
					Marker:  true,
				}),
				predicate,
			))
		}
	}

	return cmd
}
