package ecipuri

// well known constants from 24-ECIPURI
const (
	SpecsFile      = ".well-known/ecips/specs.json"
	RegistriesFile = ".well-known/ecips/known.json"
)

// Spec is the meta data of Specification from 24-ECIPURI
type Spec struct {
	URI           string `json:uri`           // The URI for the specification
	DocumentURL   string `json:documentUrl`   // The location of the authoritative source for the specification.
	DiscussionURL string `json:discussionUrl` // The location of discussion for the specification.
	Status        string `json:status`        // A status description for the specification.
	Author        string `json:author`        // Contact of the author for the specification.
	CreatedAt     string `json:createdAt`     // Date when the specification was created.
}
