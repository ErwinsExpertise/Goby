package snapshot

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type body struct {
	Type string `json:"type, omitempty"` // Default snapshot
	Name string `json:"name, omitempty"` // Snapshot name
}

// SnapshotDroplet performs snapshot for the Droplet
func SnapshotDroplet(token, dropletID string) error {
	endpoint := "https://api.digitalocean.com/v2/droplets/" + dropletID + "/actions"

	newBody := body{
		Type: "snapshot",
		Name: dropletID + "-" + time.Now().String(),
	}

	payload, err := json.Marshal(newBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	if resp, err := http.DefaultClient.Do(req); resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("could not contact DigitalOcean API response recieved: " + resp.Status)
	} else {
		log.Printf("Created snapshot %s\n", newBody.Name)
		return err
	}

}
