package pro

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	"github.com/rogpeppe/go-charset/charset"
	// anon charset
	_ "github.com/rogpeppe/go-charset/data"
)

// XML is the root tag
type XML struct {
	Name           string
	ProteinSummary ProteinSummary
}

// ProteinSummary tag is the root level
type ProteinSummary struct {
	XMLName              xml.Name             `xml:"protein_summary"`
	ProteinSummaryHeader ProteinSummaryHeader `xml:"protein_summary_header"`
	ProteinGroup         []ProteinGroup       `xml:"protein_group"`
}

// ProteinSummaryHeader tag
type ProteinSummaryHeader struct {
	XMLName                     xml.Name       `xml:"protein_summary_header"`
	ReferenceDatabase           []byte         `xml:"reference_database,attr"`
	ResidueSubstitutionList     []byte         `xml:"residue_substitution_list,attr"`
	MinPeptideProbability       float32        `xml:"min_peptide_probability,attr"`
	MinPeptideWeight            float32        `xml:"min_peptide_weight,attr"`
	NumPredictedCorrectProteins float32        `xml:"num_predicted_correct_prots,attr"`
	NumInput1Spectra            uint32         `xml:"num_input_1_spectra,attr"`
	NumInput2Spectra            uint32         `xml:"num_input_2_spectra,attr"`
	NumInput3Spectra            uint32         `xml:"num_input_3_spectra,attr"`
	NumInput4Spectra            uint32         `xml:"num_input_4_spectra,attr"`
	NumInput5Spectra            uint32         `xml:"num_input_5_spectra,attr"`
	TotalNumberSpectrumIDs      float32        `xml:"total_no_spectrum_ids,attr"`
	SampleEnzyme                []byte         `xml:"sample_enzyme,attr"`
	ProgramDetails              ProgramDetails `xml:"program_details"`
}

// ProgramDetails tag
type ProgramDetails struct {
	XMLName               xml.Name              `xml:"program_details"`
	Analysis              []byte                `xml:"analysis,attr"`
	Time                  []byte                `xml:"time,attr"`
	Version               []byte                `xml:"version,attr"`
	ProteinProphetDetails ProteinProphetDetails `xml:"proteinprophet_details"`
}

// ProteinProphetDetails tag
type ProteinProphetDetails struct {
	XMLName               xml.Name `xml:"proteinprophet_details"`
	OccamFlag             []byte   `xml:"occam_flag,attr"`
	GroupsFlag            []byte   `xml:"groups_flag,attr"`
	DegenFlag             []byte   `xml:"degen_flag,attr"`
	NSPFlag               []byte   `xml:"nsp_flag,attr"`
	FPKMFlag              []byte   `xml:"fpkm_flag,attr"`
	InitialPeptideWtIters []byte   `xml:"initial_peptide_wt_iters,attr"`
	NspDistributionIters  []byte   `xml:"nsp_distribution_iters,attr"`
	FinalPeptideWtIters   []byte   `xml:"final_peptide_wt_iters,attr"`
	RunOptions            []byte   `xml:"run_options,attr"`
}

// ProteinGroup tag
type ProteinGroup struct {
	XMLName     xml.Name  `xml:"protein_group"`
	GroupNumber uint32    `xml:"group_number,attr"`
	Probability float64   `xml:"probability,attr"`
	Protein     []Protein `xml:"protein"`
}

// Protein tag
type Protein struct {
	XMLName                         xml.Name                   `xml:"protein"`
	ProteinName                     []byte                     `xml:"protein_name,attr"`
	NumberIndistinguishableProteins int16                      `xml:"n_indistinguishable_proteins,attr"`
	Probability                     float64                    `xml:"probability,attr"`
	PercentCoverage                 float32                    `xml:"percent_coverage,attr"`
	UniqueStrippedPeptides          []byte                     `xml:"unique_stripped_peptides,attr"`
	GroupSiblingID                  []byte                     `xml:"group_sibling_id,attr"`
	TotalNumberPeptides             int                        `xml:"total_number_peptides,attr"`
	TotalNumberIndPeptides          int                        `xml:"total_number_distinct_peptides,attr"`
	PctSpectrumIDs                  float32                    `xml:"pct_spectrum_ids,attr"`
	Parameter                       Parameter                  `xml:"parameter"`
	Annotation                      Annotation                 `xml:"annotation"`
	IndistinguishableProtein        []IndistinguishableProtein `xml:"indistinguishable_protein"`
	Peptide                         []Peptide                  `xml:"peptide"`
	TopPepProb                      float64
	//Confidence                      float64                    `xml:"confidence,attr"`
}

// Parameter tag
type Parameter struct {
	XMLName xml.Name `xml:"parameter"`
	Name    []byte   `xml:"name,attr"`
	Value   int      `xml:"value,attr"`
}

// Annotation tag
type Annotation struct {
	XMLName            xml.Name `xml:"annotation"`
	ProteinDescription []byte   `xml:"protein_description,attr"`
}

// IndistinguishableProtein tag
type IndistinguishableProtein struct {
	XMLName     xml.Name   `xml:"indistinguishable_protein"`
	ProteinName string     `xml:"protein_name,attr"`
	Annotation  Annotation `xml:"annotation"`
}

// Peptide tag
type Peptide struct {
	XMLName                  xml.Name                   `xml:"peptide"`
	PeptideSequence          []byte                     `xml:"peptide_sequence,attr"`
	Charge                   uint8                      `xml:"charge,attr"`
	InitialProbability       float64                    `xml:"initial_probability,attr"`
	NSPAdjustedPprobability  float32                    `xml:"nsp_adjusted_probability,attr"`
	FPKMAdjustedProbability  float32                    `xml:"fpkm_adjusted_probability,attr"`
	Weight                   float64                    `xml:"weight,attr"`
	GroupWeight              float64                    `xml:"group_weight,attr"`
	IsNondegenerateEvidence  []byte                     `xml:"is_nondegenerate_evidence,attr"`
	NEnzymaticTermini        uint8                      `xml:"n_enzymatic_termini,attr"`
	NSiblingPeptides         float32                    `xml:"n_sibling_peptides,attr"`
	NSiblingPeptidesBin      float32                    `xml:"n_sibling_peptides_bin,attr"`
	NIstances                int                        `xml:"n_instances,attr"`
	ExpTotInstances          float32                    `xml:"exp_tot_instances,attr"`
	IsContributingEvidence   []byte                     `xml:"is_contributing_evidence,attr"`
	CalcNeutralPepMass       float64                    `xml:"calc_neutral_pep_mass,attr"`
	ModificationInfo         ModificationInfo           `xml:"modification_info"`
	PeptideParentProtein     []PeptideParentProtein     `xml:"peptide_parent_protein"`
	IndistinguishablePeptide []IndistinguishablePeptide `xml:"indistinguishable_peptide"`
	//MaxFPKM                  float64                    `xml:"max_fpkm,attr"`
	//FPKMBin                  int                        `xml:"fpkm_bin,attr"`
}

// PeptideParentProtein tag
type PeptideParentProtein struct {
	XMLName     xml.Name `xml:"peptide_parent_protein"`
	ProteinName []byte   `xml:"protein_name,attr"`
}

// IndistinguishablePeptide tag
type IndistinguishablePeptide struct {
	XMLName            xml.Name `xml:"indistinguishable_peptide"`
	PeptideSequence    []byte   `xml:"peptide_sequence,attr"`
	Charge             uint8    `xml:"charge,attr"`
	CalcNeutralPepMass float32  `xml:"calc_neutral_pep_mass,attr"`
}

// ModificationInfo tag
type ModificationInfo struct {
	XMLName          xml.Name         `xml:"modification_info"`
	ModifiedPeptide  []byte           `xml:"modified_peptide,attr"`
	ModAminoacidMass ModAminoacidMass `xml:"mod_aminoacid_mass"`
}

// ModAminoacidMass tag
type ModAminoacidMass struct {
	XMLName  xml.Name `xml:"mod_aminoacid_mass"`
	Position uint8    `xml:"position,attr"`
	Mass     float32  `xml:"mass,attr"`
}

// Parse is the main function for parsing pepxml data
func (p *XML) Parse(f string) error {

	xmlFile, e := os.Open(f)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: filepath.Base(f)}
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var ps ProteinSummary

	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader

	if e = decoder.Decode(&ps); e != nil {
		return &err.Error{Type: err.CannotParseXML, Class: err.FATA, Argument: filepath.Base(f)}
	}

	p.ProteinSummary = ps
	p.Name = filepath.Base(f)

	return nil
}
