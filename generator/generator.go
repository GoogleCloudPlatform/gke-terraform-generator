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
	Regions          string
	Zones            string
	ServiceAccountID string
	ServiceAccount   string

	Private  bool
	UseIstio bool
}

// NewVars creates creates a pointer to a new
// instance of Vars with default values set
func NewVars() *Vars {
	return &Vars{
		ServiceAccount: "${var.service_account_id}@${var.project_id}.iam.gserviceaccount.com",
	}
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
		case "regions":
			v.Regions = kv[1]
		case "zones":
			v.Zones = kv[1]
		case "service_account_id":
			v.ServiceAccountID = kv[1]
		case "service_account":
			v.ServiceAccount = kv[1]
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
	if v.ServiceAccount == "" {
		v.ServiceAccount = "${var.service_account_id}@${var.project_id}.iam.gserviceaccount.com"
	}
	return []byte(fmt.Sprintf(`
variable "project_name" {
  description = ""
	default = "%s"
}

variable "project_id" {
  description = ""
	default = "%s"
}

variable "region" {
  description = ""
  default = [%s]
}

variable "zones" {
  description = ""
  default = [%s]
}

variable "service_account_id" {
  description = ""
  default = "%s"
}

variable "service_account" {
  description = ""
  default = "%s"
}`,
		v.ProjectName,
		v.ProjectID,
		v.Regions,
		v.Zones,
		v.ServiceAccountID,
		v.ServiceAccount,
	))
}

// String implements the Stringer interface
func (v Vars) String() string {
	if v.ServiceAccount == "" {
		v.ServiceAccount = "project-service-account@${var.project_id}.iam.gserviceaccount.com"
	}
	return fmt.Sprintf(
		`project_name=%s,project_id=%s,regions=%s,zones=%s,service_account_id=%s,service_account=%s,private=%t,use-istio=%t`,
		v.ProjectName,
		v.ProjectID,
		v.Regions,
		v.Zones,
		v.ServiceAccountID,
		v.ServiceAccount,
		v.Private,
		v.UseIstio,
	)
}
