package utils

type Operator struct {
	RefID     string `json:"ref_id"`
	ShortCode string `json:"short_code"`
	Name      string `json:"name"`
}

// FindOperatorRefID searches for an operator by short code and returns its RefID.
func FindOperatorRefID(operators []Operator, shortCode string) string {
	for _, operator := range operators {
		if operator.ShortCode == shortCode {
			return operator.RefID
		}
	}
	return ""
}
