package rep

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prvst/philosopher/lib/bio"
	"github.com/prvst/philosopher/lib/cla"
	"github.com/prvst/philosopher/lib/id"
	"github.com/prvst/philosopher/lib/mod"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/uti"
	"github.com/sirupsen/logrus"
)

// AssembleIonReport reports consist on ion reporting
func (e *Evidence) AssembleIonReport(ion id.PepIDList, decoyTag string) error {

	var list IonEvidenceList
	var psmPtMap = make(map[string][]string)
	var psmIonMap = make(map[string][]string)
	var bestProb = make(map[string]float64)
	var err error

	var ionMods = make(map[string][]mod.Modification)

	// collapse all psm to protein based on Peptide-level identifications
	for _, i := range e.PSM {

		psmIonMap[i.IonForm] = append(psmIonMap[i.IonForm], i.Spectrum)
		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.Protein)

		if i.Probability > bestProb[i.IonForm] {
			bestProb[i.IonForm] = i.Probability
		}

		for _, j := range i.Modifications.Index {
			ionMods[i.IonForm] = append(ionMods[i.IonForm], j)
		}

	}

	for _, i := range ion {
		var pr IonEvidence

		pr.IonForm = fmt.Sprintf("%s#%d#%.4f", i.Peptide, i.AssumedCharge, i.CalcNeutralPepMass)

		pr.Spectra = make(map[string]int)
		pr.MappedProteins = make(map[string]int)
		pr.Modifications.Index = make(map[string]mod.Modification)

		v, ok := psmIonMap[pr.IonForm]
		if ok {
			for _, j := range v {
				pr.Spectra[j]++
			}
		}

		pr.Sequence = i.Peptide
		pr.ModifiedSequence = i.ModifiedPeptide
		pr.MZ = uti.Round(((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)), 5, 4)
		pr.ChargeState = i.AssumedCharge
		pr.PeptideMass = i.CalcNeutralPepMass
		pr.PrecursorNeutralMass = i.PrecursorNeutralMass
		pr.Expectation = i.Expectation
		pr.Protein = i.Protein
		pr.MappedProteins[i.Protein] = 0
		pr.Modifications = i.Modifications
		pr.Probability = bestProb[pr.IonForm]

		// get he list of indi proteins from pepXML data
		v, ok = psmPtMap[i.Spectrum]
		if ok {
			for _, j := range v {
				pr.MappedProteins[j] = 0
			}
		}

		mods, ok := ionMods[pr.IonForm]
		if ok {
			for _, j := range mods {
				_, okMod := pr.Modifications.Index[j.Index]
				if !okMod {
					pr.Modifications.Index[j.Index] = j
				}
			}
		}

		// is this bservation a decoy ?
		if cla.IsDecoyPSM(i, decoyTag) {
			pr.IsDecoy = true
		}

		list = append(list, pr)
	}

	sort.Sort(list)
	e.Ions = list

	return err
}

// PeptideIonReport reports consist on ion reporting
func (e *Evidence) PeptideIonReport(hasDecoys bool) {

	output := fmt.Sprintf("%s%sion.tsv", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide Sequence\tModified Sequence\tM/Z\tCharge\tExperimental Mass\tProbability\tExpectation\tSpectral Count\tUnmodified Observations\tModified Observations\tIntensity\tAssigned Modifications\tObserved Modifications\tIntensity\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Proteins\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide ion report header")
	}

	// building the printing set tat may or not contain decoys
	var printSet IonEvidenceList
	for _, i := range e.Ions {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	// peptides with no mapped poteins are related to contaminants
	// and reverse sequences. They are dificult to clean because
	// in some cases they are shared between a match decoy and a target,
	// so they stay on the lists but cannot be mapped back to the
	// original proteins. These cases should be rare to find.
	for _, i := range printSet {

		if len(i.MappedProteins) > 0 {

			assL, obs := getModsList(i.Modifications.Index)

			var mappedProteins []string
			for j := range i.MappedProteins {
				if j != i.Protein {
					mappedProteins = append(mappedProteins, j)
				}
			}

			sort.Strings(mappedProteins)
			sort.Strings(assL)
			sort.Strings(obs)

			line := fmt.Sprintf("%s\t%s\t%.4f\t%d\t%.4f\t%.4f\t%.4f\t%d\t%.4f\t%s\t%s\t%.4f\t%s\t%s\t%s\t%s\t%s\t%s\n",
				i.Sequence,
				i.ModifiedSequence,
				i.MZ,
				i.ChargeState,
				i.PeptideMass,
				i.Probability,
				i.Expectation,
				len(i.Spectra),
				i.Intensity,
				strings.Join(assL, ", "),
				strings.Join(obs, ", "),
				i.Intensity,
				i.Protein,
				i.ProteinID,
				i.EntryName,
				i.GeneName,
				i.ProteinDescription,
				strings.Join(mappedProteins, ","),
			)
			_, err = io.WriteString(file, line)
			if err != nil {
				logrus.Fatal("Cannot print PSM to file")
			}
			//}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PeptideIonTMTReport reports the ion table with TMT quantification
func (e *Evidence) PeptideIonTMTReport(labels map[string]string, hasDecoys bool) {

	output := fmt.Sprintf("%s%sion.tsv", sys.MetaDir(), string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	header := "Peptide Sequence\tModified Sequence\tM/Z\tCharge\tExperimental Mass\tProbability\tExpectation\tSpectral Count\tUnmodified Observations\tModified Observations\tIntensity\tAssigned Modifications\tObserved Modifications\tIntensity\tProtein\tProtein ID\tEntry Name\tGene\tProtein Description\tMapped Proteins\t126 Abundance\t127N Abundance\t127C Abundance\t128N Abundance\t128C Abundance\t129N Abundance\t129C Abundance\t130N Abundance\t130C Abundance\t131N Abundance\t131C Abundance\n"

	if len(labels) > 0 {
		for k, v := range labels {
			header = strings.Replace(header, k, v, -1)
		}
	}

	_, err = io.WriteString(file, header)
	if err != nil {
		logrus.Fatal("Cannot create peptide ion report header")
	}

	// building the printing set tat may or not contain decoys
	var printSet IonEvidenceList
	for _, i := range e.Ions {
		if hasDecoys == false {
			if i.IsDecoy == false {
				printSet = append(printSet, i)
			}
		} else {
			printSet = append(printSet, i)
		}
	}

	// peptides with no mapped poteins are related to contaminants
	// and reverse sequences. They are dificult to clean because
	// in some cases they are shared between a match decoy and a target,
	// so they stay on the lists but cannot be mapped back to the
	// original proteins. These cases should be rare to find.
	for _, i := range printSet {

		if len(i.MappedProteins) > 0 {

			assL, obs := getModsList(i.Modifications.Index)

			var mappedProteins []string
			for j := range i.MappedProteins {
				if j != i.Protein {
					mappedProteins = append(mappedProteins, j)
				}
			}

			sort.Strings(mappedProteins)
			sort.Strings(assL)
			sort.Strings(obs)

			line := fmt.Sprintf("%s\t%s\t%.4f\t%d\t%.4f\t%.4f\t%.4f\t%d\t%.4f\t%s\t%s\t%.4f\t%s\t%s\t%s\t%s\t%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
				i.Sequence,
				i.ModifiedSequence,
				i.MZ,
				i.ChargeState,
				i.PeptideMass,
				i.Probability,
				i.Expectation,
				len(i.Spectra),
				i.Intensity,
				strings.Join(assL, ", "),
				strings.Join(obs, ", "),
				i.Intensity,
				i.Protein,
				i.ProteinID,
				i.EntryName,
				i.GeneName,
				i.ProteinDescription,
				strings.Join(mappedProteins, ","),
				i.Labels.Channel1.Intensity,
				i.Labels.Channel2.Intensity,
				i.Labels.Channel3.Intensity,
				i.Labels.Channel4.Intensity,
				i.Labels.Channel5.Intensity,
				i.Labels.Channel6.Intensity,
				i.Labels.Channel7.Intensity,
				i.Labels.Channel8.Intensity,
				i.Labels.Channel9.Intensity,
				i.Labels.Channel10.Intensity,
				i.Labels.Channel11.Intensity,
			)
			_, err = io.WriteString(file, line)
			if err != nil {
				logrus.Fatal("Cannot print PSM to file")
			}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}