package cli

type Output struct {
	Success      bool     `json:"success"`
	ErrorMessage string   `json:"error-message,omitempty"`
	Nodes        []string `json:"nodes,omitempty"`
}
