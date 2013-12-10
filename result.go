package hue

// Result encapsulates the standard response message that the
// bridge returns
type Result struct {
	Success map[string]interface{} `json:"success"`
}
