package docker

import "fmt"

// Docker holds list of data for all dockers.
type Docker []struct {
	Docker Model `yaml:"docker"`
}

// LookupRepo lookup for a docker in a bunch of dockers that match given
// id name. then return it if founded otherwise error would be non nil.
func (d *Docker) LookupRepo(id string) (Model, error) {
	for _, s := range *d {
		if s.Docker.Id == id {
			return s.Docker, nil
		}
	}

	return Model{}, fmt.Errorf("docker id not found")
}

// Sanitization loop through all dockers and run their each sanitization.
func (d *Docker) Sanitization() error {
	for i, dd := range *d {
		if err := dd.Docker.Sanitization(); err != nil {
			return fmt.Errorf("failed sanitizing #%d docker in config: %s", i+1, err)
		}
	}

	return nil
}
