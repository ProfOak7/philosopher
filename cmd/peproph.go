package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/ext/peptideprophet"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/spf13/cobra"
)

// peprophCmd represents the peproph command
var peprophCmd = &cobra.Command{
	Use:   "peptideprophet",
	Short: "Peptide assignment validation",
	//Long:  "Statistical validation of peptide assignments for MS/MS Proteomics data\nPeptidProphet v5.0",
	Run: func(cmd *cobra.Command, args []string) {

		if len(m.UUID) < 1 && len(m.Home) < 1 {
			e := &err.Error{Type: err.WorkspaceNotFound, Class: err.FATA}
			logrus.Fatal(e.Error())
		}

		logrus.Info("Executing PeptideProphet")

		peptideprophet.Run(m, args)
		m.Serialize()

		logrus.Info("Done")
		return
	},
}

func init() {

	if len(os.Args) > 1 && os.Args[1] == "peptideprophet" {

		m.Restore(sys.Meta())

		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Exclude, "exclude", "", false, "exclude deltaCn*, Mascot*, and Comet* results from results (default Penalize * results)")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Leave, "leave", "", false, "leave alone deltaCn*, Mascot*, and Comet* results from results (default Penalize * results)")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Icat, "icat", "", false, "apply ICAT model (default Autodetect ICAT)")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Noicat, "noicat", "", false, "do no apply ICAT model (default Autodetect ICAT)")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Zero, "zero", "", false, "report results with minimum probability 0")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Accmass, "accmass", "", false, "use Accurate Mass model binning")
		peprophCmd.Flags().IntVarP(&m.PeptideProphet.Clevel, "clevel", "", 0, "set Conservative Level in neg_stdev from the neg_mean, low numbers are less conservative, high numbers are more conservative")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Ppm, "ppm", "", false, "use PPM mass error instead of Daltons for mass modeling")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Nomass, "nomass", "", false, "disable mass model")
		peprophCmd.Flags().Float64VarP(&m.PeptideProphet.Masswidth, "masswidth", "", 5.0, "model mass width")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Pi, "pi", "", false, "enable peptide pI model")
		peprophCmd.Flags().IntVarP(&m.PeptideProphet.Minpintt, "minpintt", "", 2, "minimum number of NTT in a peptide used for positive pI model")
		peprophCmd.Flags().Float64VarP(&m.PeptideProphet.Minpiprob, "minpiprob", "", 0.9, "minimum probability after first pass of a peptide used for positive pI model")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Rt, "rt", "", false, "enable peptide RT model")
		peprophCmd.Flags().Float64VarP(&m.PeptideProphet.Minrtprob, "minrtprob", "", 0.9, "minimum probability after first pass of a peptide used for positive RT model")
		peprophCmd.Flags().IntVarP(&m.PeptideProphet.Minrtntt, "minrtntt", "", 2, "minimum number of NTT in a peptide used for positive RT model")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Glyc, "glyc", "", false, "enable peptide Glyco motif model")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Phospho, "phospho", "", false, "enable peptide Phospho motif model")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Maldi, "maldi", "", false, "enable MALDI mode")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Instrwarn, "instrwarn", "", false, "warn and continue if combined data was generated by different instrument models")
		peprophCmd.Flags().Float64VarP(&m.PeptideProphet.Minprob, "minprob", "", 0.05, "report results with minimum probability")
		peprophCmd.Flags().StringVarP(&m.PeptideProphet.Decoy, "decoy", "", "", "semi-supervised mode, protein name prefix to identify Decoy entries")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Decoyprobs, "decoyprobs", "", false, "compute possible non-zero probabilities for Decoy entries on the last iteration")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Nontt, "nontt", "", false, "disable NTT enzymatic termini model")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Nonmc, "nonmc", "", false, "disable NMC missed cleavage model")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Expectscore, "expectscore", "", false, "use expectation value as the only contributor to the f-value for modeling")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Nonparam, "nonparam", "", false, "use semi-parametric modeling, must be used in conjunction with --decoy option")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Neggamma, "neggamma", "", false, "use Gamma distribution to model the negative hits")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Forcedistr, "forcedistr", "", false, "bypass quality control checks, report model despite bad modelling")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Optimizefval, "optimizefval", "", false, "(SpectraST only) optimize f-value function f(dot,delta) using PCA")
		peprophCmd.Flags().IntVarP(&m.PeptideProphet.MinPepLen, "minpeplen", "", 7, "minimum peptide length not rejected")
		peprophCmd.Flags().StringVarP(&m.PeptideProphet.Output, "output", "", "interact", "Output name prefix")
		peprophCmd.Flags().BoolVarP(&m.PeptideProphet.Combine, "combine", "", false, "combine the results from PeptideProphet into a single result file")
		peprophCmd.Flags().StringVarP(&m.PeptideProphet.Database, "database", "", "", "path to the database")
	}

	RootCmd.AddCommand(peprophCmd)
}
