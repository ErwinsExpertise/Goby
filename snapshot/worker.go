package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type body struct {
	Type string `json:"type, omitempty"` // Default snapshot
	Name string `json:"name, omitempty"` // Snapshot name
}

type Snapshot struct {
	Snapshots []Snap `json:"snapshots, omitempty"`
}

type Snap struct {
	ID         int    `json:"id, omitempty"`
	Name       string `json:"name, omitempty"`
	Created_At string `json:"created_at, omitempty"`
	Type       string `json:"type, omitempty"`
}

const (
	POST   = "POST"
	GET    = "GET"
	DELETE = "DELETE"
)

// SnapshotDroplet performs snapshot for the Droplet. Droplet is name as: DropletID - Current time in RFC822
func SnapshotDroplet(token, dropletID string) error {
	endpoint := "https://api.digitalocean.com/v2/droplets/" + dropletID + "/actions"

	newBody := body{
		Type: "snapshot",
		Name: dropletID + "-" + time.Now().Format(time.RFC822),
	}
	_, err := newBody.makeRequest(token, endpoint, POST)
	if err != nil {
		return err
	}
	log.Printf("Created snapshot %s\n", newBody.Name)

	return nil
}

// ListSnapshots returns struct containing all of the droplets snapshots
func ListSnapshots(token, dropletID string) Snapshot {
	endpoint := "https://api.digitalocean.com/v2/droplets/" + dropletID + "/snapshots"

	var newBody body

	resp, err := newBody.makeRequest(token, endpoint, GET)
	if err != nil {
		fmt.Println(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var snapshots Snapshot

	json.Unmarshal(respBody, &snapshots)
	return snapshots
}

// Clean will delete any droplet snapshots that are older than the given keep time
func (s *Snapshot) Clean(token string, keepTime int64) {
	currentTime := time.Now()

	for _, snap := range s.Snapshots {
		dropletTime, err := time.Parse(time.RFC3339, snap.Created_At)
		if err != nil {
			fmt.Println(err)
			continue
		}
		isOld := compareTimes(currentTime, dropletTime, keepTime)

		if isOld == true {
			snap.delete(token)
		}
	}

}

func (s *Snap) delete(token string) {
	endpoint := "https://api.digitalocean.com/v2/snapshots/" + strconv.Itoa(s.ID)
	var newBody body

	_, err := newBody.makeRequest(token, endpoint, DELETE)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Deleted snapshot: %s\n", s.Name)
}

func compareTimes(t1, t2 time.Time, keepTime int64) bool {
	newTime := t2.AddDate(0, 0, int(keepTime))

	if newTime.Before(t1) {
		return true
	}

	return false
}

func (b *body) makeRequest(token, endpoint, method string) (*http.Response, error) {
	payload, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, errors.New("could not contact DigitalOcean API response recieved: " + resp.Status)
	}

	return resp, err

}
