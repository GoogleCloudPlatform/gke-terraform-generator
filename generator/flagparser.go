package generator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Vars
type Vars struct {
	ProjectName      string
	ProjectID        string
	Region           string
	Zones            string
	ServiceAccountID string

	Private  bool
	UseIstio bool
}

// NewVars creates creates a pointer to a new
// instance of Vars with default values set
func NewVars() *Vars {
	return &Vars{}
}

// Set is used to set the values of Vars via cli key=value input
func (v *Vars) Set(s string) error {
	kvR := strings.Split(s, ",")
	for _, val := range kvR {
		kv := strings.Split(val, "=")

		if len(kv) < 2 {
			return errors.New("input is invalid")
		}

		switch kv[0] {
		case "project_name":
			v.ProjectName = kv[1]
		case "project_id":
			v.ProjectID = kv[1]
		case "region":
			v.Region = kv[1]
		case "zones":
			zns := strings.Split(kv[1], " ")
			v.Zones = "\"" + strings.Join(zns, "\",\"") + "\""
		case "service_account_id":
			v.ServiceAccountID = kv[1]
		case "private":
			b, _ := strconv.ParseBool(kv[1])
			v.Private = b
		case "use-istio":
			b, _ := strconv.ParseBool(kv[1])
			v.Private = b
		default:
			return errors.New("Unknownn field:" + kv[0])
		}
	}
	return nil
}

func (Vars) Type() string {
	return "Vars"
}

// TFVars is used to generate a Terraform Vars sections
func (v Vars) TFVars() []byte {
	out := fmt.Sprintf("project_name=\"%s\"\nproject_id=\"%s\"\nservice_account_id=\"%s\"",
		v.ProjectName,
		v.ProjectID,
		v.ServiceAccountID,
	)

	if v.Zones != "" {
		out = out + "\nzones=[" + v.Zones + "]"
	}
	if v.Region != "" {
		out = out + "\nregion=\"" + v.Region + "\""
	}

	return []byte(out)
}

// String implements the Stringer interface
func (v Vars) String() string {
	return fmt.Sprintf(
		`project_name=\"%s\",project_id=\"%s\",region=\"%s\",zones=\"%s\",service_account_id="\%s\",private=\"%t\",use-istio=\"%t\"`,
		v.ProjectName,
		v.ProjectID,
		v.Region,
		v.Zones,
		v.ServiceAccountID,
		v.Private,
		v.UseIstio,
	)
}
