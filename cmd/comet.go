package cmd

import (
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher-source/lib/ext/comet"
	"github.com/prvst/philosopher-source/lib/meta"
	"github.com/prvst/philosopher-source/lib/sys"

	"github.com/spf13/cobra"
)

var cmt comet.Comet

// cometCmd represents the comet command
var cometCmd = &cobra.Command{
	Use:   "comet",
	Short: "MS/MS database search",
	//Long:  "Peptide Spectrum Matching using the Comet algorithm\nComet release 2016.01.rev 2",
	Run: func(cmd *cobra.Command, args []string) {

		var m meta.Data
		m.Restore(sys.Meta())
		if len(m.UUID) < 1 && len(m.Home) < 1 {
			logrus.Fatal("Workspace not found. Run 'philosopher init' to create a workspace")
		}

		if len(cmt.Param) < 1 {
			logrus.Fatal("No parameter file found. Run 'comet --help' for more information")
		}

		if cmt.Print == false && len(args) < 1 {
			logrus.Fatal("Missing parameter file or data file for analysis")
		}

		// deploy the binaries
		cmt.Deploy()

		if cmt.Print == true {
			sys.CopyFile(cmt.DefaultParam, filepath.Base(cmt.DefaultParam))
			return
		}

		// var binFile []byte
		// binFile, err := ioutil.ReadFile(cmt.DefaultParam)
		// if err != nil {
		// 	logrus.Fatal(err)
		// }

		//m.Experimental.CometParam = binFile

		// run
		cmt.Run(args)

		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	cmt = comet.New()

	cometCmd.Flags().BoolVarP(&cmt.Print, "print", "", false, "print a comet.params file")
	cometCmd.Flags().StringVarP(&cmt.Param, "param", "", "comet.params.txt", "comet parameter file")

	RootCmd.AddCommand(cometCmd)
}
