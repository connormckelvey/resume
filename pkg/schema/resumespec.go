package schema

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ResumeSpec struct {
	PersonalInfo *PersonalInfoSpec `yaml:"personal"`
	Employment   []*EmploymentSpec `yaml:"employment"`
	Education    []*EducationSpec  `yaml:"education"`
	Experience   []*ExperienceSpec `yaml:"experience"`
	Projects     []*ProjectSpec    `yaml:"projects"`
}

type PersonalInfoSpec struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	Phone    string `yaml:"phone"`
	Email    *Link  `yaml:"email"`
	Website  *Link  `yaml:"website"`
	Github   *Link  `yaml:"github"`
}

type EmploymentSpec struct {
	Dates      *DateRange `yaml:"dates"`
	Title      string     `yaml:"title"`
	Company    string     `yaml:"company"`
	Location   string     `yaml:"location"`
	Summary    string     `yaml:"summary"`
	Technology []string   `yaml:"technology"`
	Notes      []string   `yaml:"notes"`
}

type EducationSpec struct {
	Dates       *DateRange `yaml:"dates"`
	Degree      string     `yaml:"degree"`
	Institution string     `yaml:"institution"`
	Location    string     `yaml:"location"`
	Notes       []string   `yaml:"notes"`
}

type ExperienceSpec struct {
	Title string   `yaml:"title"`
	Items []string `yaml:"items"`
}

type ProjectSpec struct {
	Title       string   `yaml:"title"`
	Link        *Link    `yaml:"link"`
	Description string   `yaml:"description"`
	Technology  []string `yaml:"technology"`
}

type Link struct {
	Text string `yaml:"text"`
	Href string `yaml:"href"`
}

type DateRange struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

func LoadResume(path string) (*ResumeSpec, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	resume := &ResumeSpec{}
	err = yaml.Unmarshal(f, resume)
	if err != nil {
		return nil, err
	}

	return resume, nil
}
