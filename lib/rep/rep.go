package rep

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/cmsl/bio"
	"github.com/prvst/cmsl/err"
	"github.com/prvst/cmsl/utils"
	"github.com/prvst/philosopher/lib/clas"
	"github.com/prvst/philosopher/lib/data"
	"github.com/prvst/philosopher/lib/meta"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/prvst/philosopher/lib/tmt"
	"github.com/prvst/philosopher/lib/uni"
	"github.com/prvst/philosopher/lib/xml"
)

// Evidence ...
type Evidence struct {
	meta.Data
	PSM           PSMEvidenceList
	Ions          IonEvidenceList
	Peptides      PeptideEvidenceList
	Proteins      ProteinEvidenceList
	Mods          Modifications
	Modifications ModificationEvidence
	Combined      CombinedEvidenceList
}

// Modifications ...
type Modifications struct {
	DefinedModMassDiff  map[float64]float64
	DefinedModAminoAcid map[float64]string
}

// PSMEvidence struct
type PSMEvidence struct {
	Index                 uint32
	Spectrum              string
	Scan                  int
	Peptide               string
	Protein               string
	ProteinID             string
	GeneName              string
	ModifiedPeptide       string
	AlternativeProteins   []string
	ModPositions          []string
	AssignedModMasses     []float64
	AssignedMassDiffs     []float64
	AssignedAminoAcid     []string
	AssignedModifications map[string]uint16
	ObservedModifications map[string]uint16
	AssumedCharge         uint8
	HitRank               uint8
	PrecursorNeutralMass  float64
	PrecursorExpMass      float64
	RetentionTime         float64
	CalcNeutralPepMass    float64
	RawMassdiff           float64
	Massdiff              float64
	LocalizedMassDiff     string
	Probability           float64
	Expectation           float64
	Xcorr                 float64
	DeltaCN               float64
	DeltaCNStar           float64
	SPScore               float64
	SPRank                float64
	Hyperscore            float64
	Nextscore             float64
	DiscriminantValue     float64
	Intensity             float64
	Purity                float64
	IsUnique              bool
	IsURazor              bool
	Labels                tmt.Labels
}

// PSMEvidenceList ...
type PSMEvidenceList []PSMEvidence

func (a PSMEvidenceList) Len() int           { return len(a) }
func (a PSMEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PSMEvidenceList) Less(i, j int) bool { return a[i].Spectrum < a[j].Spectrum }

// IonEvidence groups all valid info about peptide ions for reports
type IonEvidence struct {
	Sequence                string
	ModifiedSequence        string
	AssignedModifications   map[string]uint16
	ObservedModifications   map[string]uint16
	RetentionTime           string
	Spectra                 map[string]int
	MappedProteins          map[string]uint8
	ChargeState             uint8
	Spc                     int
	MZ                      float64
	PeptideMass             float64
	PrecursorNeutralMass    float64
	IsNondegenerateEvidence bool
	Weight                  float64
	GroupWeight             float64
	Intensity               float64
	Probability             float64
	Expectation             float64
	IsURazor                bool
	Labels                  tmt.Labels
	SummedLabelIntensity    float64
	ModifiedObservations    int
	UnModifiedObservations  int
}

// IonEvidenceList ...
type IonEvidenceList []IonEvidence

func (a IonEvidenceList) Len() int           { return len(a) }
func (a IonEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IonEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// PeptideEvidence groups all valid info about peptide ions for reports
type PeptideEvidence struct {
	Sequence               string
	ChargeState            map[uint8]uint8
	Spc                    int
	Intensity              float64
	ModifiedObservations   int
	UnModifiedObservations int
}

// PeptideEvidenceList ...
type PeptideEvidenceList []PeptideEvidence

func (a PeptideEvidenceList) Len() int           { return len(a) }
func (a PeptideEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PeptideEvidenceList) Less(i, j int) bool { return a[i].Sequence < a[j].Sequence }

// ProteinEvidence ...
type ProteinEvidence struct {
	OriginalHeader         string
	ProteinName            string
	ProteinGroup           uint32
	ProteinSubGroup        string
	ProteinID              string
	EntryName              string
	Description            string
	Organism               string
	Length                 int
	Coverage               float32
	GeneNames              string
	ProteinExistence       string
	Sequence               string
	SupportingSpectra      map[string]int
	IndiProtein            map[string]uint8
	UniqueStrippedPeptides int
	// TotalNumRazorPeptides        int
	// TotalNumPeptideIons          int
	// NumURazorPeptideIons         int // Unique + razor
	TotalPeptideIons             map[string]IonEvidence
	UniquePeptideIons            map[string]IonEvidence
	URazorPeptideIons            map[string]IonEvidence // Unique + razor
	TotalSpC                     int
	UniqueSpC                    int
	URazorSpC                    int // Unique + razor
	TotalIntensity               float64
	UniqueIntensity              float64
	URazorIntensity              float64 // Unique + razor
	Probability                  float64
	TopPepProb                   float64
	IsDecoy                      bool
	IsContaminant                bool
	URazorModifiedObservations   int
	URazorUnModifiedObservations int
	URazorAssignedModifications  map[string]uint16
	URazorObservedModifications  map[string]uint16
	TotalLabels                  tmt.Labels
	UniqueLabels                 tmt.Labels
	URazorLabels                 tmt.Labels // Unique + razor
}

// ProteinEvidenceList list
type ProteinEvidenceList []ProteinEvidence

func (a ProteinEvidenceList) Len() int           { return len(a) }
func (a ProteinEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProteinEvidenceList) Less(i, j int) bool { return a[i].ProteinGroup < a[j].ProteinGroup }

// CombinedEvidence represents all combined proteins detected
type CombinedEvidence struct {
	GroupNumber            uint32
	SiblingID              string
	ProteinName            string
	ProteinID              string
	IndiProtein            []string
	EntryName              string
	GeneNames              string
	Description            string
	Length                 int
	Names                  []string
	UniqueStrippedPeptides int
	TotalIons              int
	SupportingSpectra      map[string]string
	ProteinProbability     float64
	TopPepProb             float64
	PeptideIons            []xml.PeptideIonIdentification
	TotalSpc               map[string]int
	UniqueSpc              map[string]int
	UrazorSpc              map[string]int
	TotalIntensity         map[string]float64
	UniqueIntensity        map[string]float64
	UrazorIntensity        map[string]float64
	//UniquePeptideIons      []xml.PeptideIonIdentification
	// TotalPeptideIonStrings  map[string]int
	// UniquePeptideIonStrings map[string]int
	// RazorPeptideIonStrings  map[string]int
	// TotalPeptideIntensity   map[string]float64
	// UniquePeptideIntensity  map[string]float64
	// RazorPeptideIntensity   map[string]float64
	// TotalPeptideIons        int
	// UniquePeptideIons       int
	// SharedPeptideIons       int
	// RazorPeptideIons        int
}

// CombinedEvidenceList ...
type CombinedEvidenceList []CombinedEvidence

func (a CombinedEvidenceList) Len() int           { return len(a) }
func (a CombinedEvidenceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CombinedEvidenceList) Less(i, j int) bool { return a[i].GroupNumber < a[j].GroupNumber }

// ModificationEvidence represents the list of modifications and the mod bins
type ModificationEvidence struct {
	MassBins []MassBin
}

// MassBin represents each bin from the mass distribution
type MassBin struct {
	LowerMass     float64
	HigherRight   float64
	MassCenter    float64
	AverageMass   float64
	CorrectedMass float64
	Modifications []string
	AssignedMods  PSMEvidenceList
	ObservedMods  PSMEvidenceList
}

// New constructor
func New() Evidence {

	var o Evidence
	var m meta.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp
	o.OS = m.OS
	o.Arch = m.Arch

	return o
}

// AssemblePSMReport ...
func (e *Evidence) AssemblePSMReport(pep xml.PepIDList, decoyTag string) error {

	var list PSMEvidenceList

	// collect database information
	var dtb data.Base
	dtb.Restore()

	var genes = make(map[string]string)
	var ptid = make(map[string]string)
	for _, j := range dtb.Records {
		genes[j.PartHeader] = j.GeneNames
		ptid[j.PartHeader] = j.ID
	}

	for _, i := range pep {
		if !clas.IsDecoyPSM(i, decoyTag) {

			var p PSMEvidence

			p.Index = i.Index
			p.Spectrum = i.Spectrum
			p.Scan = i.Scan
			p.Peptide = i.Peptide
			p.Protein = i.Protein
			p.ModifiedPeptide = i.ModifiedPeptide
			p.AlternativeProteins = i.AlternativeProteins
			//p.AlternativeTargetProteins = i.AlternativeTargetProteins
			p.ModPositions = i.ModPositions
			p.AssignedModMasses = i.AssignedModMasses
			p.AssignedMassDiffs = i.AssignedMassDiffs
			p.AssignedAminoAcid = i.AssignedAminoAcid
			p.AssumedCharge = i.AssumedCharge
			p.HitRank = i.HitRank
			p.PrecursorNeutralMass = i.PrecursorNeutralMass
			p.PrecursorExpMass = i.PrecursorExpMass
			p.RetentionTime = i.RetentionTime
			p.CalcNeutralPepMass = i.CalcNeutralPepMass
			p.RawMassdiff = i.RawMassDiff
			p.Massdiff = i.Massdiff
			p.LocalizedMassDiff = i.LocalizedMassDiff
			p.Probability = i.Probability
			p.Expectation = i.Expectation
			p.Xcorr = i.Xcorr
			p.DeltaCN = i.DeltaCN
			p.SPRank = i.SPRank
			p.Hyperscore = i.Hyperscore
			p.Nextscore = i.Nextscore
			p.DiscriminantValue = i.DiscriminantValue
			p.Intensity = i.Intensity
			p.AssignedModifications = make(map[string]uint16)
			p.ObservedModifications = make(map[string]uint16)

			// TODO find a better way to map gene names to the psm
			gn, ok := genes[i.Protein]
			if ok {
				p.GeneName = gn
			}

			id, ok := ptid[i.Protein]
			if ok {
				p.ProteinID = id
			}

			list = append(list, p)
		}
	}

	sort.Sort(list)
	e.PSM = list

	return nil
}

// PSMReport report all psms from study that passed the FDR filter
func (e *Evidence) PSMReport() {

	output := fmt.Sprintf("%s%spsm.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tAssigned Modifications\tOberved Modifications\tObserved Mass Localization\tIs Unique\tIs Razor\tMapped Proteins\tProtein\tAlternative Proteins\n")
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	for _, i := range e.PSM {

		var assL []string

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "n" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] != "n" && i.AssignedAminoAcid[j] != "c" {
				loc := fmt.Sprintf("%s%s(%.4f)", i.ModPositions[j], i.AssignedAminoAcid[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "c" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		var obs []string
		for j := range i.ObservedModifications {
			obs = append(obs, j)
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t%s\t%s\t%t\t%t\t%d\t%s\t%s\n",
			i.Spectrum,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.RawMassdiff,
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Xcorr,
			i.DeltaCN,
			i.DeltaCNStar,
			i.SPScore,
			i.SPRank,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedMassDiff,
			i.IsUnique,
			i.IsURazor,
			len(i.AlternativeProteins)+1, // Mapped Proteins //len(i.AlternativeTargetProteins)+1,
			i.Protein,
			strings.Join(i.AlternativeProteins, ", "), // strings.Join(i.AlternativeTargetProteins, ", "),
		)
		_, err = io.WriteString(file, line)
		if err != nil {
			logrus.Fatal("Cannot print PSM to file")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PSMQuantReport report all psms with TMT labels from study that passed the FDR filter
func (e *Evidence) PSMQuantReport() {

	output := fmt.Sprintf("%s%spsm.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, "Spectrum\tPeptide\tModified Peptide\tCharge\tRetention\tCalculated M/Z\tObserved M/Z\tOriginal Delta Mass\tAdjusted Delta Mass\tExperimental Mass\tPeptide Mass\tXCorr\tDeltaCN\tDeltaCNStar\tSPScore\tSPRank\tExpectation\tHyperscore\tNextscore\tPeptideProphet Probability\tIntensity\tIs Unique\tIs Razor\tAssigned Modifications\tOberved Modifications\tObserved Mass Localization\tMapped Proteins\tGene Name\tProtein\tAlternative Proteins\tPurity\tRaw Channel 1\tRaw Channel 2\tRaw Channel 3\tRaw Channel 4\tRaw Channel 5\tRaw Channel 6\tRaw Channel 7\tRaw Channel 8\tRaw Channel 9\tRaw Channel 10\n")
	if err != nil {
		logrus.Fatal("Cannot print PSM to file")
	}

	for _, i := range e.PSM {

		var assL []string

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "n" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] != "n" && i.AssignedAminoAcid[j] != "c" {
				loc := fmt.Sprintf("%s%s(%.4f)", i.ModPositions[j], i.AssignedAminoAcid[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		for j := 0; j <= len(i.ModPositions)-1; j++ {
			if i.AssignedMassDiffs[j] != 0 && i.AssignedAminoAcid[j] == "c" {
				loc := fmt.Sprintf("%s(%.4f)", i.ModPositions[j], i.AssignedMassDiffs[j])
				assL = append(assL, loc)
			}
		}

		var obs []string
		for j := range i.ObservedModifications {
			obs = append(obs, j)
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%d\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%e\t%.4f\t%.4f\t%.4f\t%.4f\t%t\t%t\t%s\t%s\t%s\t%d\t%s\t%s\t%s\t%.2f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
			i.Spectrum,
			i.Peptide,
			i.ModifiedPeptide,
			i.AssumedCharge,
			i.RetentionTime,
			((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			((i.PrecursorNeutralMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)),
			i.RawMassdiff,
			i.Massdiff,
			i.PrecursorNeutralMass,
			i.CalcNeutralPepMass,
			i.Xcorr,
			i.DeltaCN,
			i.DeltaCNStar,
			i.SPScore,
			i.SPRank,
			i.Expectation,
			i.Hyperscore,
			i.Nextscore,
			i.Probability,
			i.Intensity,
			i.IsUnique,
			i.IsURazor,
			strings.Join(assL, ", "),
			strings.Join(obs, ", "),
			i.LocalizedMassDiff,
			len(i.AlternativeProteins)+1, //len(i.AlternativeTargetProteins)+1,
			i.GeneName,
			i.Protein,
			strings.Join(i.AlternativeProteins, ", "), // strings.Join(i.AlternativeTargetProteins, ", "),
			i.Purity,
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
		)
		_, err = io.WriteString(file, line)
		if err != nil {
			logrus.Fatal("Cannot print PSM to file")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// AssembleIonReport reports consist on ion reporting
func (e *Evidence) AssembleIonReport(ion xml.PepIDList, decoyTag string) error {

	var list IonEvidenceList
	var psmPtMap = make(map[string][]string)
	var psmIonMap = make(map[string][]string)
	var assignedModMap = make(map[string][]string)
	var observedModMap = make(map[string][]string)
	var err error

	// collapse all psm to protein based on Peptide-level identifications
	for _, i := range e.PSM {

		var ion string
		if len(i.ModifiedPeptide) > 0 {
			ion = fmt.Sprintf("%s#%d", i.ModifiedPeptide, i.AssumedCharge)
		} else {
			ion = fmt.Sprintf("%s#%d", i.Peptide, i.AssumedCharge)
		}

		psmIonMap[ion] = append(psmIonMap[ion], i.Spectrum)

		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.Protein)
		psmPtMap[i.Spectrum] = append(psmPtMap[i.Spectrum], i.AlternativeProteins...)

		// get the list of all assigned modifications
		if len(i.AssignedModifications) > 0 {
			for k := range i.AssignedModifications {
				assignedModMap[i.Spectrum] = append(assignedModMap[i.Spectrum], k)
			}
		}

		// get the list of all observed modifications
		if len(i.ObservedModifications) > 0 {
			for k := range i.ObservedModifications {
				observedModMap[i.Spectrum] = append(observedModMap[i.Spectrum], k)
			}
		}
	}

	for _, i := range ion {
		if !clas.IsDecoyPSM(i, decoyTag) {

			var key string
			if len(i.ModifiedPeptide) > 0 {
				key = fmt.Sprintf("%s#%d", i.ModifiedPeptide, i.AssumedCharge)
			} else {
				key = fmt.Sprintf("%s#%d", i.Peptide, i.AssumedCharge)
			}

			var pr IonEvidence

			pr.Spectra = make(map[string]int)
			pr.MappedProteins = make(map[string]uint8)
			pr.ObservedModifications = make(map[string]uint16)
			pr.AssignedModifications = make(map[string]uint16)

			v, ok := psmIonMap[key]
			if ok {
				for _, j := range v {
					pr.Spectra[j]++
				}
			}

			pr.Sequence = i.Peptide
			pr.ModifiedSequence = i.ModifiedPeptide
			pr.MZ = utils.Round(((i.CalcNeutralPepMass + (float64(i.AssumedCharge) * bio.Proton)) / float64(i.AssumedCharge)), 5, 4)
			pr.ChargeState = i.AssumedCharge
			pr.PeptideMass = i.CalcNeutralPepMass
			pr.PrecursorNeutralMass = i.PrecursorNeutralMass
			pr.Probability = i.Probability
			pr.Expectation = i.Expectation

			// get he list of indi proteins from pepXML data
			v, ok = psmPtMap[i.Spectrum]
			if ok {
				for _, j := range v {
					pr.MappedProteins[j] = 0
				}
			}

			va, oka := assignedModMap[i.Spectrum]
			if oka {
				for _, j := range va {
					pr.AssignedModifications[j] = 0
				}
			}

			vo, oko := observedModMap[i.Spectrum]
			if oko {
				for _, j := range vo {
					pr.ObservedModifications[j] = 0
				}
			} else {
				pr.UnModifiedObservations++
			}

			list = append(list, pr)
		}
	}

	sort.Sort(list)
	e.Ions = list

	return err
}

// UpdateIonStatus pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
func (e *Evidence) UpdateIonStatus() {

	var uniqueMap = make(map[string]bool)
	var urazorMap = make(map[string]bool)

	for _, i := range e.Proteins {

		for _, j := range i.UniquePeptideIons {
			var ion string
			if len(j.ModifiedSequence) > 0 {
				ion = fmt.Sprintf("%s#%d", j.ModifiedSequence, j.ChargeState)
			} else {
				ion = fmt.Sprintf("%s#%d", j.Sequence, j.ChargeState)
			}
			//key := fmt.Sprintf("%s#%s", i.ProteinName, ion)
			key := fmt.Sprintf("%s", ion)
			uniqueMap[key] = true
		}

		for _, j := range i.URazorPeptideIons {
			var ion string
			if len(j.ModifiedSequence) > 0 {
				ion = fmt.Sprintf("%s#%d", j.ModifiedSequence, j.ChargeState)
			} else {
				ion = fmt.Sprintf("%s#%d", j.Sequence, j.ChargeState)
			}
			//key := fmt.Sprintf("%s#%s", i.ProteinName, ion)
			key := fmt.Sprintf("%s", ion)
			urazorMap[key] = true
		}

	}

	for i := range e.PSM {

		var ion string
		if len(e.PSM[i].ModifiedPeptide) > 0 {
			ion = fmt.Sprintf("%s#%d", e.PSM[i].ModifiedPeptide, e.PSM[i].AssumedCharge)
		} else {
			ion = fmt.Sprintf("%s#%d", e.PSM[i].Peptide, e.PSM[i].AssumedCharge)
		}

		//key := fmt.Sprintf("%s#%s", e.PSM[i].Protein, ion)
		key := fmt.Sprintf("%s", ion)

		_, uOK := uniqueMap[key]
		if uOK {
			e.PSM[i].IsUnique = true
		}

		_, rOK := urazorMap[key]
		if rOK {
			e.PSM[i].IsURazor = true
		}

	}

	for i := range e.Ions {

		var ion string
		if len(e.PSM[i].ModifiedPeptide) > 0 {
			ion = fmt.Sprintf("%s#%d", e.PSM[i].ModifiedPeptide, e.PSM[i].AssumedCharge)
		} else {
			ion = fmt.Sprintf("%s#%d", e.PSM[i].Peptide, e.PSM[i].AssumedCharge)
		}

		//key := fmt.Sprintf("%s#%s", e.PSM[i].Protein, ion)
		key := fmt.Sprintf("%s", ion)

		_, uOK := uniqueMap[key]
		if uOK {
			e.Ions[i].IsNondegenerateEvidence = true
		}

		_, rOK := urazorMap[key]
		if rOK {
			e.PSM[i].IsURazor = true
		}

	}

	return
}

// UpdateIonModCount counts how many times each ion is observed modified and not modified
func (e *Evidence) UpdateIonModCount() {

	// recreate the ion list from the main report object
	var AllIons = make(map[string]int)
	var ModIons = make(map[string]int)
	var UnModIons = make(map[string]int)

	for _, i := range e.Ions {
		var ion string
		if len(i.ModifiedSequence) > 0 {
			ion = fmt.Sprintf("%s#%d", i.ModifiedSequence, i.ChargeState)
		} else {
			ion = fmt.Sprintf("%s#%d", i.Sequence, i.ChargeState)
		}
		AllIons[ion] = 0
		ModIons[ion] = 0
		UnModIons[ion] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range e.PSM {
		var psmIon string
		if len(i.ModifiedPeptide) > 0 {
			psmIon = fmt.Sprintf("%s#%d", i.ModifiedPeptide, i.AssumedCharge)
		} else {
			psmIon = fmt.Sprintf("%s#%d", i.Peptide, i.AssumedCharge)
		}

		// check the map
		_, ok := AllIons[psmIon]
		if ok {

			if i.Massdiff >= -0.99 && i.Massdiff <= 0.99 {
				UnModIons[psmIon]++
			} else {
				ModIons[psmIon]++
			}

		}
	}

	for i := range e.Ions {
		var ion string
		if len(e.Ions[i].ModifiedSequence) > 0 {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
		} else {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
		}

		v1, ok1 := UnModIons[ion]
		if ok1 {
			e.Ions[i].UnModifiedObservations = v1
		}

		v2, ok2 := ModIons[ion]
		if ok2 {
			e.Ions[i].ModifiedObservations = v2
		}

	}

	return
}

// UpdateIndistinguishableProteinLists pushes back to ion and psm evideces the uniqueness and razorness status of each peptide and ion
func (e *Evidence) UpdateIndistinguishableProteinLists() {

	var ptNetMap = make(map[string][]string)

	for _, i := range e.Proteins {

		var list []string
		for k := range i.IndiProtein {
			list = append(list, k)
		}

		ptNetMap[i.ProteinID] = list

	}

	for i := range e.PSM {
		v, ok := ptNetMap[e.PSM[i].ProteinID]
		if ok {
			e.PSM[i].AlternativeProteins = v
		}
	}

	// for i := range e.Ions {
	// 	v, ok := ptNetMap[e.Ions[i].]
	// 	if ok {
	// 		e.PSM[i].AlternativeProteins = v
	// 	}
	// }

	return
}

// UpdateIonAssignedAndObservedMods collects all Assigned and Observed modifications from
// individual PSM and assign them to ions
func (e *Evidence) UpdateIonAssignedAndObservedMods() {

	//var list IonEvidenceList

	for i := range e.Ions {
		var ion string
		if len(e.Ions[i].ModifiedSequence) > 0 {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].ModifiedSequence, e.Ions[i].ChargeState)
		} else {
			ion = fmt.Sprintf("%s#%d", e.Ions[i].Sequence, e.Ions[i].ChargeState)
		}

		for _, j := range e.PSM {
			var psmIon string
			if len(j.ModifiedPeptide) > 0 {
				psmIon = fmt.Sprintf("%s#%d", j.ModifiedPeptide, j.AssumedCharge)
			} else {
				psmIon = fmt.Sprintf("%s#%d", j.Peptide, j.AssumedCharge)
			}

			if ion == psmIon {
				for k := range j.AssignedModifications {
					e.Ions[i].AssignedModifications[k]++
				}
				for k := range j.ObservedModifications {
					e.Ions[i].ObservedModifications[k]++
				}

				break
			}

		}

		//list = append(list, e.Ions[i])
	}

	//e.Ions = list

	return
}

// PeptideIonReport reports consist on ion reporting
func (e *Evidence) PeptideIonReport() {

	output := fmt.Sprintf("%s%sion.tsv", e.Temp, string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide Sequence\tModified Sequence\tM/Z\tCharge\tExperimental Mass\tProbability\tExpectation\tSpectral Count\tUnmodified Observations\tModified Observations\tIntensity\tAssigned Modifications\tObserved Modifications\tIntensity\tMapped Proteins\tProtein IDs\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide ion report header")
	}

	// peptides with no mapped poteins are related to contaminants
	// and reverse sequences. They are dificult to clean because
	// in some cases they are shared between a match decoy and a target,
	// so they stay on the lists but cannot be mapped back to the
	// original proteins. These cases should be rare to find.
	for _, i := range e.Ions {

		var pts []string
		//var ipts []string

		if len(i.MappedProteins) > 0 {

			if len(e.Proteins) > 1 {

				for k := range i.MappedProteins {
					pts = append(pts, k)
				}

				var amods []string
				for j := range i.AssignedModifications {
					amods = append(amods, j)
				}

				var omods []string
				for j := range i.ObservedModifications {
					omods = append(omods, j)
				}

				line := fmt.Sprintf("%s\t%s\t%.4f\t%d\t%.4f\t%.4f\t%.4f\t%d\t%d\t%d\t%.4f\t%s\t%s\t%.4f\t%d\t%s\n",
					i.Sequence,
					i.ModifiedSequence,
					i.MZ,
					i.ChargeState,
					i.PeptideMass,
					i.Probability,
					i.Expectation,
					i.Spc,
					i.UnModifiedObservations,
					i.ModifiedObservations,
					i.Intensity,
					strings.Join(amods, ", "),
					strings.Join(omods, ", "),
					i.Intensity,
					len(i.MappedProteins),
					strings.Join(pts, ", "),
				)
				_, err = io.WriteString(file, line)
				if err != nil {
					logrus.Fatal("Cannot print PSM to file")
				}
			}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// AssemblePeptideReport reports consist on ion reporting
func (e *Evidence) AssemblePeptideReport(pep xml.PepIDList, decoyTag string) error {

	var list PeptideEvidenceList
	var pepSeqMap = make(map[string]uint8)
	var pepCSMap = make(map[string][]uint8)
	var pepSpc = make(map[string]int)
	var pepInt = make(map[string]float64)
	var err error

	for _, i := range pep {
		if !clas.IsDecoyPSM(i, decoyTag) {
			pepSeqMap[i.Peptide] = 0
			pepSpc[i.Peptide] = 0
			pepInt[i.Peptide] = 0
		}
	}

	// TODO review this method, Intensity quant is not working
	for _, i := range e.PSM {
		_, ok := pepSeqMap[i.Peptide]
		if ok {
			pepCSMap[i.Peptide] = append(pepCSMap[i.Peptide], i.AssumedCharge)
			pepSpc[i.Peptide]++
			if i.Intensity > pepInt[i.Peptide] {
				pepInt[i.Peptide] = i.Intensity
			}
		}
	}

	for k := range pepSeqMap {

		var pep PeptideEvidence
		pep.ChargeState = make(map[uint8]uint8)
		pep.Sequence = k

		for _, i := range pepCSMap[k] {
			pep.ChargeState[i] = 0
		}
		pep.Spc = pepSpc[k]
		pep.Intensity = pepInt[k]

		list = append(list, pep)
	}

	sort.Sort(list)
	e.Peptides = list

	return err
}

// PeptideReport reports consist on ion reporting
func (e *Evidence) PeptideReport() {

	output := fmt.Sprintf("%s%speptide.tsv", e.Temp, string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create peptide output file")
	}
	defer file.Close()

	_, err = io.WriteString(file, "Peptide\tCharges\tSpectral Count\tUnmodified Observations\tModified Observations\n")
	if err != nil {
		logrus.Fatal("Cannot create peptide report header")
	}

	for _, i := range e.Peptides {

		var cs []string
		for j := range i.ChargeState {
			cs = append(cs, strconv.Itoa(int(j)))
		}
		sort.Strings(cs)

		line := fmt.Sprintf("%s\t%s\t%d\t%d\t%d\n",
			i.Sequence,
			strings.Join(cs, ", "),
			i.Spc,
			i.UnModifiedObservations,
			i.ModifiedObservations,
			//i.Intensity,
		)
		_, err = io.WriteString(file, line)
		if err != nil {
			logrus.Fatal("Cannot print PSM to file")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// ProteinFastaReport saves to disk a filtered FASTA file with FDR aproved proteins
func (e *Evidence) ProteinFastaReport() error {

	output := fmt.Sprintf("%s%sproteins.fas", e.Temp, string(filepath.Separator))

	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Could not create output file")
	}
	defer file.Close()

	for _, i := range e.Proteins {
		header := i.OriginalHeader
		line := ">" + header + "\n" + i.Sequence + "\n"
		_, err = io.WriteString(file, line)
		if err != nil {
			return errors.New("Cannot print PSM to file")
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return nil
}

// AssembleProteinReport ...
func (e *Evidence) AssembleProteinReport(pro xml.ProtIDList, decoyTag string) error {

	var list ProteinEvidenceList
	var err error

	var evidenceIons = make(map[string]IonEvidence)
	for _, i := range e.Ions {
		var ion string
		if len(i.ModifiedSequence) > 0 {
			ion = fmt.Sprintf("%s#%d#%.4f", i.ModifiedSequence, i.ChargeState, i.PeptideMass)
		} else {
			ion = fmt.Sprintf("%s#%d#%.4f", i.Sequence, i.ChargeState, i.PeptideMass)
		}
		evidenceIons[ion] = i
	}

	for _, i := range pro {
		if !strings.HasPrefix(i.ProteinName, decoyTag) {

			var rep ProteinEvidence

			rep.SupportingSpectra = make(map[string]int)
			rep.TotalPeptideIons = make(map[string]IonEvidence)
			rep.UniquePeptideIons = make(map[string]IonEvidence)
			rep.URazorPeptideIons = make(map[string]IonEvidence)
			rep.IndiProtein = make(map[string]uint8)
			rep.URazorAssignedModifications = make(map[string]uint16)
			rep.URazorObservedModifications = make(map[string]uint16)

			rep.ProteinName = i.ProteinName
			rep.ProteinGroup = i.GroupNumber
			rep.ProteinSubGroup = i.GroupSiblingID
			rep.Length = i.Length
			rep.Coverage = i.PercentCoverage
			rep.UniqueStrippedPeptides = len(i.UniqueStrippedPeptides)
			rep.Probability = i.Probability
			rep.TopPepProb = i.TopPepProb

			if strings.Contains(i.ProteinName, decoyTag) {
				rep.IsDecoy = true
			} else {
				rep.IsDecoy = false
			}

			for j := range i.IndistinguishableProtein {
				rep.IndiProtein[i.IndistinguishableProtein[j]] = 0
			}

			for _, k := range i.PeptideIons {

				// var ion string
				// if len(k.ModifiedPeptide) > 0 {
				// 	ion = fmt.Sprintf("%s#%d#%.4f", k.ModifiedPeptide, k.Charge, k.CalcNeutralPepMass)
				// } else {
				// 	ion = fmt.Sprintf("%s#%d#%.4f", k.PeptideSequence, k.Charge, k.CalcNeutralPepMass)
				// }

				ion := fmt.Sprintf("%s#%d#%.4f", k.PeptideSequence, k.Charge, k.CalcNeutralPepMass)

				v, ok := evidenceIons[ion]
				if ok {

					for spec := range v.Spectra {
						rep.SupportingSpectra[spec]++
					}

					v.MappedProteins = nil

					ref := v
					ref.Weight = k.Weight
					ref.GroupWeight = k.GroupWeight
					ref.IsNondegenerateEvidence = k.IsNondegenerateEvidence
					if k.Razor == 1 {
						ref.IsURazor = true
					}
					evidenceIons[ion] = ref
				}

				rep.TotalPeptideIons[ion] = evidenceIons[ion]

				// if k.IsUnique == true && k.Razor == 1 {
				// 	rep.UniquePeptideIons[ion] = evidenceIons[ion]
				// 	rep.URazorPeptideIons[ion] = evidenceIons[ion]
				//
				// 	rep.URazorUnModifiedObservations += evidenceIons[ion].UnModifiedObservations
				// 	rep.URazorModifiedObservations += evidenceIons[ion].ModifiedObservations
				//
				// 	for key, value := range evidenceIons[ion].AssignedModifications {
				// 		rep.URazorAssignedModifications[key] += value
				// 	}
				//
				// 	for key, value := range evidenceIons[ion].ObservedModifications {
				// 		rep.URazorObservedModifications[key] += value
				// 	}
				// } else if k.IsUnique == false && k.Razor == 1 {
				// 	rep.URazorUnModifiedObservations += evidenceIons[ion].UnModifiedObservations
				// 	rep.URazorModifiedObservations += evidenceIons[ion].ModifiedObservations
				//
				// 	rep.URazorPeptideIons[ion] = evidenceIons[ion]
				//
				// 	for key, value := range evidenceIons[ion].AssignedModifications {
				// 		rep.URazorAssignedModifications[key] += value
				// 	}
				//
				// 	for key, value := range evidenceIons[ion].ObservedModifications {
				// 		rep.URazorObservedModifications[key] += value
				// 	}
				// }

				if k.IsUnique == true {
					rep.UniquePeptideIons[ion] = evidenceIons[ion]
					rep.URazorPeptideIons[ion] = evidenceIons[ion]
					// rep.NumURazorPeptideIons++
					// rep.TotalNumRazorPeptides++

					rep.URazorUnModifiedObservations += evidenceIons[ion].UnModifiedObservations
					rep.URazorModifiedObservations += evidenceIons[ion].ModifiedObservations

					for key, value := range evidenceIons[ion].AssignedModifications {
						rep.URazorAssignedModifications[key] += value
					}

					for key, value := range evidenceIons[ion].ObservedModifications {
						rep.URazorObservedModifications[key] += value
					}
				}

				if k.Razor == 1 {
					rep.URazorUnModifiedObservations += evidenceIons[ion].UnModifiedObservations
					rep.URazorModifiedObservations += evidenceIons[ion].ModifiedObservations

					rep.URazorPeptideIons[ion] = evidenceIons[ion]
					// rep.TotalNumRazorPeptides++

					for key, value := range evidenceIons[ion].AssignedModifications {
						rep.URazorAssignedModifications[key] += value
					}

					for key, value := range evidenceIons[ion].ObservedModifications {
						rep.URazorObservedModifications[key] += value
					}
				}

			}

			list = append(list, rep)
		}
	}

	var dtb data.Base
	dtb.Restore()

	if len(dtb.Records) < 1 {
		return errors.New("Cant locate database data")
	}

	for i := range list {
		for _, j := range dtb.Records {
			// fix the name sand headers and pull database information into proteinreport
			if strings.Contains(j.OriginalHeader, list[i].ProteinName) {
				if (j.IsDecoy == true && list[i].IsDecoy == true) || (j.IsDecoy == false && list[i].IsDecoy == false) {
					list[i].OriginalHeader = j.OriginalHeader
					list[i].ProteinID = j.ID
					list[i].EntryName = j.EntryName
					list[i].ProteinExistence = j.ProteinExistence
					list[i].GeneNames = j.GeneNames
					list[i].Sequence = j.Sequence
					list[i].ProteinName = j.ProteinName
					list[i].Organism = j.Organism

					// uniprot entries have the description on ProteinName
					if len(j.Description) < 1 {
						list[i].Description = j.ProteinName
					} else {
						list[i].Description = j.Description
					}

					break
				}
			}
		}
	}

	sort.Sort(list)
	e.Proteins = list

	return err
}

// ProteinReport ...
func (e *Evidence) ProteinReport() {

	// create result file
	output := fmt.Sprintf("%s%sreport.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tOrganism\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tStripped Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tRazor Spectral Count\tTotal Intensity\tUnique Intensity\tRazor Intensity\tRazor Assigned Modifications\tRazor Observed Modifications\tIndistinguishable Proteins\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Proteins {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		var amods []string
		if len(i.URazorAssignedModifications) > 0 {
			for j := range i.URazorAssignedModifications {
				amods = append(amods, j)
			}
		}

		var omods []string
		if len(i.URazorObservedModifications) > 0 {
			for j := range i.URazorObservedModifications {
				omods = append(omods, j)
			}
		}

		// proteins with almost no evidences, and completely shared with decoys are eliminated from the analysis,
		// in most cases proteins with one small peptide shared with a decoy
		//if len(i.TotalPeptideIons) > 0 {

		line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%6.f\t%s\t%s\t%s\t",
			i.ProteinGroup,           // Group
			i.ProteinSubGroup,        // SubGroup
			i.ProteinID,              // Protein ID
			i.EntryName,              // Entry Name
			i.Length,                 // Length
			i.Coverage,               // Percent Coverage
			i.Organism,               // Organism
			i.Description,            // Description
			i.ProteinExistence,       // Protein Existence
			i.GeneNames,              // Genes
			i.Probability,            // Protein Probability
			i.TopPepProb,             // Top Peptide Probability
			i.UniqueStrippedPeptides, // Stripped Peptides
			// i.TotalNumPeptideIons,    // Total Peptide Ions
			// i.NumURazorPeptideIons,   // Unique Peptide Ions
			len(i.TotalPeptideIons),
			len(i.UniquePeptideIons),
			i.TotalSpC,  // Total Spectral Count
			i.UniqueSpC, // Unique Spectral Count
			i.URazorSpC, // Razor Spectral Count
			//i.URazorUnModifiedObservations, // Unmodified Occurrences
			//i.URazorModifiedObservations,   // Modified Occurrences
			i.TotalIntensity,          // Total Intensity
			i.UniqueIntensity,         // Unique Intensity
			i.URazorIntensity,         // Razor Intensity
			strings.Join(amods, ", "), // Razor Assigned Modifications
			strings.Join(omods, ", "), // Razor Observed Modifications
			strings.Join(ip, ", "),    // Indistinguishable Proteins
		)

		line += "\n"
		n, err := io.WriteString(file, line)
		if err != nil {
			logrus.Fatal(n, err)
		}
		//}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// ProteinQuantReport ...
func (e *Evidence) ProteinQuantReport() {

	// create result file
	output := fmt.Sprintf("%s%sreport.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	//Total Channel 1\tTotal Channel 2\tTotal Channel 3\t Total Channel 4\t Total Channel 5\t Total Channel 6\t Total Channel 7\tTotal Channel 8\tTotal Channel 9\tTotal Channel 10\tUnique Channel 1\tUnique Channel 2\tUnique Channel 3\tUnique Channel 4\tUnique Channel 5\tUnique Channel 6\tUnique Channel 7\tUnique Channel 8\tUnique Channel 9\tUnique Channel 10\n")
	line := fmt.Sprintf("Group\tSubGroup\tProtein ID\tEntry Name\tLength\tPercent Coverage\tDescription\tProtein Existence\tGenes\tProtein Probability\tTop Peptide Probability\tUnique Stripped Peptides\tRazor Peptides\tTotal Peptide Ions\tUnique Peptide Ions\tTotal Spectral Count\tUnique Spectral Count\tTotal Intensity\tUnique Intensity\tTotal Raw Channel 1\tTotal Raw Channel 2\tTotal Raw Channel 3\t Total Raw Channel 4\t Total Raw Channel 5\t Total Raw Channel 6\t Total Raw Channel 7\tTotal Raw Channel 8\tTotal Raw Channel 9\tTotal Raw Channel 10\tUnique Raw Channel 1\tUnique Raw Channel 2\tUnique Raw Channel 3\tUnique Raw Channel 4\tUnique Raw Channel 5\tUnique Raw Channel 6\tUnique Raw Channel 7\tUnique Raw Channel 8\tUnique Raw Channel 9\tUnique Raw Channel 10\tRazor Raw Channel 1\tRazor Raw Channel 2\tRazor Raw Channel 3\tRazor Raw Channel 4\tRazor Raw Channel 5\tRazor Raw Channel 6\tRazor Raw Channel 7\tRazor Raw Channel 8\tRazor Raw Channel 9\tRazor Raw Channel 10\tTotal Channel 1\tTotal Channel 2\tTotal Channel 3\t Total Channel 4\t Total Channel 5\t Total Channel 6\t Total Channel 7\tTotal Channel 8\tTotal Channel 9\tTotal Channel 10\tUnique Channel 1\tUnique Channel 2\tUnique Channel 3\tUnique Channel 4\tUnique Channel 5\tUnique Channel 6\tUnique Channel 7\tUnique Channel 8\tUnique Channel 9\tUnique Channel 10\tRazor Channel 1\tRazor Channel 2\tRazor Channel 3\tRazor Channel 4\tRazor Channel 5\tRazor Channel 6\tRazor Channel 7\tRazor Channel 8\tRazor Channel 9\tRazor Channel 10\tIndistinguishable Proteins\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Proteins {

		var ip []string
		for k := range i.IndiProtein {
			ip = append(ip, k)
		}

		//%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t
		if len(i.TotalPeptideIons) > 0 {
			line = fmt.Sprintf("%d\t%s\t%s\t%s\t%d\t%.2f\t%s\t%s\t%s\t%.4f\t%.4f\t%d\t%d\t%d\t%d\t%d\t%d\t%6.f\t%6.f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%s\t",
				i.ProteinGroup,           // Group
				i.ProteinSubGroup,        // SubGroup
				i.ProteinID,              // Protein ID
				i.EntryName,              // EntryName
				i.Length,                 // Length
				i.Coverage,               // Percent Coverage
				i.Description,            // Description
				i.ProteinExistence,       // Protein Existence
				i.GeneNames,              // Genes
				i.Probability,            // Protein Probability
				i.TopPepProb,             // Top peptide Probability
				i.UniqueStrippedPeptides, // Unique Stripped Peptides
				len(i.URazorPeptideIons), // Razor peptides
				len(i.TotalPeptideIons),  // Total peptide Ions
				len(i.UniquePeptideIons), // Unique Peptide Ions
				// len(i.ur),
				// len(i.TotalPeptideIons),
				// len(i.URazorPeptideIons),
				i.TotalSpC,        // Total Spectral Count
				i.UniqueSpC,       // Unique Spectral Count
				i.TotalIntensity,  // Total Intensity
				i.UniqueIntensity, // Unique Intensity
				i.TotalLabels.Channel1.NormIntensity,
				i.TotalLabels.Channel2.NormIntensity,
				i.TotalLabels.Channel3.NormIntensity,
				i.TotalLabels.Channel4.NormIntensity,
				i.TotalLabels.Channel5.NormIntensity,
				i.TotalLabels.Channel6.NormIntensity,
				i.TotalLabels.Channel7.NormIntensity,
				i.TotalLabels.Channel8.NormIntensity,
				i.TotalLabels.Channel9.NormIntensity,
				i.TotalLabels.Channel10.NormIntensity,
				i.UniqueLabels.Channel1.NormIntensity,
				i.UniqueLabels.Channel2.NormIntensity,
				i.UniqueLabels.Channel3.NormIntensity,
				i.UniqueLabels.Channel4.NormIntensity,
				i.UniqueLabels.Channel5.NormIntensity,
				i.UniqueLabels.Channel6.NormIntensity,
				i.UniqueLabels.Channel7.NormIntensity,
				i.UniqueLabels.Channel8.NormIntensity,
				i.UniqueLabels.Channel9.NormIntensity,
				i.UniqueLabels.Channel10.NormIntensity,
				i.URazorLabels.Channel1.NormIntensity,
				i.URazorLabels.Channel2.NormIntensity,
				i.URazorLabels.Channel3.NormIntensity,
				i.URazorLabels.Channel4.NormIntensity,
				i.URazorLabels.Channel5.NormIntensity,
				i.URazorLabels.Channel6.NormIntensity,
				i.URazorLabels.Channel7.NormIntensity,
				i.URazorLabels.Channel8.NormIntensity,
				i.URazorLabels.Channel9.NormIntensity,
				i.URazorLabels.Channel10.NormIntensity,
				i.TotalLabels.Channel1.RatioIntensity,
				i.TotalLabels.Channel2.RatioIntensity,
				i.TotalLabels.Channel3.RatioIntensity,
				i.TotalLabels.Channel4.RatioIntensity,
				i.TotalLabels.Channel5.RatioIntensity,
				i.TotalLabels.Channel6.RatioIntensity,
				i.TotalLabels.Channel7.RatioIntensity,
				i.TotalLabels.Channel8.RatioIntensity,
				i.TotalLabels.Channel9.RatioIntensity,
				i.TotalLabels.Channel10.RatioIntensity,
				i.UniqueLabels.Channel1.RatioIntensity,
				i.UniqueLabels.Channel2.RatioIntensity,
				i.UniqueLabels.Channel3.RatioIntensity,
				i.UniqueLabels.Channel4.RatioIntensity,
				i.UniqueLabels.Channel5.RatioIntensity,
				i.UniqueLabels.Channel6.RatioIntensity,
				i.UniqueLabels.Channel7.RatioIntensity,
				i.UniqueLabels.Channel8.RatioIntensity,
				i.UniqueLabels.Channel9.RatioIntensity,
				i.UniqueLabels.Channel10.RatioIntensity,
				i.URazorLabels.Channel1.RatioIntensity,
				i.URazorLabels.Channel2.RatioIntensity,
				i.URazorLabels.Channel3.RatioIntensity,
				i.URazorLabels.Channel4.RatioIntensity,
				i.URazorLabels.Channel5.RatioIntensity,
				i.URazorLabels.Channel6.RatioIntensity,
				i.URazorLabels.Channel7.RatioIntensity,
				i.URazorLabels.Channel8.RatioIntensity,
				i.URazorLabels.Channel9.RatioIntensity,
				i.URazorLabels.Channel10.RatioIntensity,
				strings.Join(ip, ", "))

			line += "\n"
			n, err := io.WriteString(file, line)
			if err != nil {
				logrus.Fatal(n, err)
			}
		}
	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// AssembleModificationReport cretaes the modifications lists
func (e *Evidence) AssembleModificationReport() error {

	var modEvi ModificationEvidence

	var massWindow = float64(0.5)
	var binsize = float64(0.1)
	var amplitude = float64(500)

	var bins []MassBin

	nBins := (amplitude*(1/binsize) + 1) * 2
	for i := 0; i <= int(nBins); i++ {
		var b MassBin

		b.LowerMass = -(amplitude) - (massWindow * binsize) + (float64(i) * binsize)
		b.LowerMass = utils.Round(b.LowerMass, 5, 4)

		b.HigherRight = -(amplitude) + (massWindow * binsize) + (float64(i) * binsize)
		b.HigherRight = utils.Round(b.HigherRight, 5, 4)

		b.MassCenter = -(amplitude) + (float64(i) * binsize)
		b.MassCenter = utils.Round(b.MassCenter, 5, 4)

		bins = append(bins, b)
	}

	// calculate the total number of PSMs per cluster
	for i := range e.PSM {

		// the checklist will not allow the same PSM to be added multiple times to the
		// same bin in case multiple identical mods are present in te sequence
		var assignChecklist = make(map[float64]uint8)
		var obsChecklist = make(map[float64]uint8)

		for j := range bins {

			// for assigned mods
			// 0 here means something that doest not map to the pepXML header
			// like multiple mods on n-term
			for _, l := range e.PSM[i].AssignedMassDiffs {

				if l > bins[j].LowerMass && l <= bins[j].HigherRight && l != 0 {
					_, ok := assignChecklist[l]
					if !ok {
						bins[j].AssignedMods = append(bins[j].AssignedMods, e.PSM[i])
						assignChecklist[l] = 0
					}
				}
			}

			// for delta masses
			if e.PSM[i].Massdiff > bins[j].LowerMass && e.PSM[i].Massdiff <= bins[j].HigherRight {
				_, ok := obsChecklist[e.PSM[i].Massdiff]
				if !ok {
					bins[j].ObservedMods = append(bins[j].ObservedMods, e.PSM[i])
					obsChecklist[e.PSM[i].Massdiff] = 0
				}
			}

		}
	}

	// calculate average mass for each cluster
	var zeroBinMassDeviation float64
	for i := range bins {
		pep := bins[i].ObservedMods
		total := 0.0
		for j := range pep {
			total += pep[j].Massdiff
		}
		if len(bins[i].ObservedMods) > 0 {
			bins[i].AverageMass = (float64(total) / float64(len(pep)))
		} else {
			bins[i].AverageMass = 0
		}
		if bins[i].MassCenter == 0 {
			zeroBinMassDeviation = bins[i].AverageMass
		}

		bins[i].AverageMass = utils.Round(bins[i].AverageMass, 5, 4)
	}

	// correcting mass values based on Bin 0 average mass
	for i := range bins {
		if len(bins[i].ObservedMods) > 0 {
			if bins[i].AverageMass > 0 {
				bins[i].CorrectedMass = (bins[i].AverageMass - zeroBinMassDeviation)
			} else {
				bins[i].CorrectedMass = (bins[i].AverageMass + zeroBinMassDeviation)
			}
		} else {
			bins[i].CorrectedMass = bins[i].MassCenter
		}
		bins[i].CorrectedMass = utils.Round(bins[i].CorrectedMass, 5, 4)
	}

	//e.Modifications = modEvi
	//e.Modifications.MassBins = bins

	modEvi.MassBins = bins
	e.Modifications = modEvi

	return nil
}

// MapMassDiffToUniMod maps PSMs to modifications based on their mass shifts
func (e *Evidence) MapMassDiffToUniMod() *err.Error {

	// 10 ppm
	var tolerance = 0.01

	u := uni.New()
	u.ProcessUniMOD()

	for _, i := range u.Modifications {

		for j := range e.PSM {

			// for fixed and variable modifications
			for k := range e.PSM[j].AssignedMassDiffs {
				if e.PSM[j].AssignedMassDiffs[k] >= (i.MonoMass-tolerance) && e.PSM[j].AssignedMassDiffs[k] <= (i.MonoMass+tolerance) {
					if !strings.Contains(i.Description, "substitution") {
						fullname := fmt.Sprintf("%.4f:%s (%s)", i.MonoMass, i.Title, i.Description)
						e.PSM[j].AssignedModifications[fullname] = 0
					}
				}
			}

			// for delta masses
			if e.PSM[j].Massdiff >= (i.MonoMass-tolerance) && e.PSM[j].Massdiff <= (i.MonoMass+tolerance) {
				fullName := fmt.Sprintf("%.4f:%s (%s)", i.MonoMass, i.Title, i.Description)
				_, ok := e.PSM[j].AssignedModifications[fullName]
				if !ok {
					e.PSM[j].ObservedModifications[fullName] = 0
				}
			}

		}
	}

	for j := range e.PSM {
		if e.PSM[j].Massdiff != 0 && len(e.PSM[j].ObservedModifications) == 0 {
			e.PSM[j].ObservedModifications["Unknown"] = 0
		}
	}

	return nil
}

// UpdatePeptideModCount counts how many times each peptide is observed modified and not modified
func (e *Evidence) UpdatePeptideModCount() {

	// recreate the ion list from the main report object
	var all = make(map[string]int)
	var mod = make(map[string]int)
	var unmod = make(map[string]int)

	for _, i := range e.Peptides {
		all[i.Sequence] = 0
		mod[i.Sequence] = 0
		unmod[i.Sequence] = 0
	}

	// range over PSMs looking for modified and not modified evidences
	// if they exist on the ions map, get the numbers
	for _, i := range e.PSM {

		_, ok := all[i.Peptide]
		if ok {

			if i.Massdiff >= -0.99 && i.Massdiff <= 0.99 {
				unmod[i.Peptide]++
			} else {
				mod[i.Peptide]++
			}

		}
	}

	for i := range e.Peptides {

		v1, ok1 := unmod[e.Peptides[i].Sequence]
		if ok1 {
			e.Peptides[i].UnModifiedObservations = v1
		}

		v2, ok2 := mod[e.Peptides[i].Sequence]
		if ok2 {
			e.Peptides[i].ModifiedObservations = v2
		}

	}

	return
}

// ModificationReport ...
func (e *Evidence) ModificationReport() {

	// create result file
	output := fmt.Sprintf("%s%smodifications.tsv", e.Temp, string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		logrus.Fatal("Cannot create report file:", err)
	}
	defer file.Close()

	line := fmt.Sprintf("Mass Bin\tPSMs with Assigned Modifications\tPSMs with Observed Modifications\n")

	n, err := io.WriteString(file, line)
	if err != nil {
		logrus.Fatal(n, err)
	}

	for _, i := range e.Modifications.MassBins {

		line = fmt.Sprintf("%.4f\t%d\t%d",
			i.CorrectedMass,
			len(i.AssignedMods),
			len(i.ObservedMods),
		)

		line += "\n"
		n, err := io.WriteString(file, line)
		if err != nil {
			logrus.Fatal(n, err)
		}

	}

	// copy to work directory
	sys.CopyFile(output, filepath.Base(output))

	return
}

// PlotMassHist plots the delta mass histogram
func (e *Evidence) PlotMassHist() error {

	outfile := fmt.Sprintf("%s%sdelta-mass.html", e.Temp, string(filepath.Separator))

	file, err := os.Create(outfile)
	if err != nil {
		return errors.New("Could not create output for delta mass binning")
	}
	defer file.Close()

	var xvar []string
	var y1var []string
	var y2var []string

	for _, i := range e.Modifications.MassBins {
		xel := fmt.Sprintf("'%.2f',", i.MassCenter)
		xvar = append(xvar, xel)
		y1el := fmt.Sprintf("'%d',", len(i.AssignedMods))
		y1var = append(y1var, y1el)
		y2el := fmt.Sprintf("'%d',", len(i.ObservedMods))
		y2var = append(y2var, y2el)
	}

	xAxis := fmt.Sprintf("	  x: %s,", xvar)
	AssAxis := fmt.Sprintf("	  y: %s,", y1var)
	ObsAxis := fmt.Sprintf("	  y: %s,", y2var)

	io.WriteString(file, "<head>\n")
	io.WriteString(file, "  <script src=\"https://cdn.plot.ly/plotly-latest.min.js\"></script>\n")
	io.WriteString(file, "</head>\n")
	io.WriteString(file, "<body>\n")
	io.WriteString(file, "<div id=\"myDiv\" style=\"width: 1024px; height: 768px;\"></div>\n")
	io.WriteString(file, "<script>\n")
	io.WriteString(file, "var trace1 = {")
	io.WriteString(file, xAxis)
	io.WriteString(file, ObsAxis)
	io.WriteString(file, "name: 'Observed',")
	io.WriteString(file, "type: 'bar',")
	io.WriteString(file, "};")
	io.WriteString(file, "var trace2 = {")
	io.WriteString(file, xAxis)
	io.WriteString(file, AssAxis)
	io.WriteString(file, "name: 'Assigned',")
	io.WriteString(file, "type: 'bar',")
	io.WriteString(file, "};")
	io.WriteString(file, "var data = [trace1, trace2];\n")
	io.WriteString(file, "var layout = {barmode: 'stack', title: 'Distribution of Mass Modifications', xaxis: {title: 'mass bins'}, yaxis: {title: '# PSMs'}};\n")
	io.WriteString(file, "Plotly.newPlot('myDiv', data, layout);\n")
	io.WriteString(file, "</script>\n")
	io.WriteString(file, "</body>")

	if err != nil {
		logrus.Warning("There was an error trying to plot the mass distribution")
	}

	// copy to work directory
	sys.CopyFile(outfile, filepath.Base(outfile))

	return nil
}
