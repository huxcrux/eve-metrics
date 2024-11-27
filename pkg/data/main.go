package data

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Data struct {
	AllianceMembers       []AllianceMember      `yaml:"allianceMembers"`
	AllianceContacts      []AllianceContacts    `yaml:"allianceContacts"`
	CorporationContacts   []CorporationContacts `yaml:"corporationContacts"`
	CompletedIndustryJobs []int32               `yaml:"completedIndustryJobs"`
}

type AllianceMember struct {
	AllianceID int32   `yaml:"allianceID"`
	Members    []int32 `yaml:"members"`
}

type AllianceContacts struct {
	AllianceID int32             `yaml:"allianceID"`
	Contacts   []AllianceContact `yaml:"contacts"`
}

type AllianceContact struct {
	ContactId   int32   `json:"contact_id,omitempty"`   /* contact_id integer */
	ContactType string  `json:"contact_type,omitempty"` /* contact_type string */
	LabelIds    []int64 `json:"label_ids,omitempty"`    /* label_ids array */
	Standing    float32 `json:"standing,omitempty"`     /* Standing of the contact */
}

type CorporationContacts struct {
	CorporationID int32                `yaml:"corporationID"`
	Contacts      []CorporationContact `yaml:"contacts"`
}

type Contact struct {
	ID       int32   `yaml:"id"`
	Standing float64 `yaml:"standing"`
	LabelIds []int32 `yaml:"labelIds"`
}

type CorporationContact struct {
	ID          int32   `yaml:"id"`
	ContactType string  `yaml:"contact_type,omitempty"`
	Standing    float32 `yaml:"standing"`
	LabelIds    []int64 `yaml:"labelIds,omitempty"`
}

var (
	filePath = "data.yml"
)

func initDataFile() {

	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// Data file exists
		return
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	fmt.Println("Data file created")
}

func ReadData() Data {

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return Data{}
	}

	// Read the YAML file
	datafile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the YAML content into the Config struct
	var data Data
	err = yaml.Unmarshal(datafile, &data)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	return data
}

func WriteData(datadile Data) {

	// Marshal the Config struct into YAML
	data, err := yaml.Marshal(&datadile)
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	initDataFile()

	// Write the YAML content to the file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Println("Data file written")
}
