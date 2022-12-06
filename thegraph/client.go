package thegraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	nnsmetadata "github.com/nnsprotocol/nns-metadata"
)

// Client integrates with thegraph.com api.
type Client struct {
	GraphNames map[nnsmetadata.Network]string
}

// Domain holds information about a domain.
type Domain struct {
	Name        string
	CreatedDate time.Time
}

// DomainName returns the domain name from a label hash.
func (c Client) Domain(network nnsmetadata.Network, labelHash string) (Domain, bool, error) {
	gn, ok := c.GraphNames[network]
	if !ok {
		return Domain{}, false, fmt.Errorf("invalid network %q", network)
	}

	body := map[string]interface{}{
		"query": fmt.Sprintf(`{
			domains(first:1, where:{labelhash:%q}){
				labelName
				createdAt
			}
		}
		`, labelHash),
	}
	b, err := json.Marshal(body)
	if err != nil {
		return Domain{}, false, fmt.Errorf("error marshalling body: %v", err)
	}
	url := fmt.Sprintf("https://api.thegraph.com/subgraphs/name/%s", gn)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return Domain{}, false, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Domain{}, false, fmt.Errorf("error executing request: %v", err)
	}
	defer res.Body.Close()

	var resB struct {
		Data struct {
			Domains []struct {
				LabelName string `json:"labelName"`
				CreatedAt string `json:"createdAt"`
			} `json:"domains"`
		} `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&resB)
	if err != nil {
		return Domain{}, false, fmt.Errorf("error parsing response: %v", err)
	}
	if len(resB.Data.Domains) == 0 {
		return Domain{}, false, nil
	}
	d := resB.Data.Domains[0]
	if d.LabelName == "" {
		return Domain{}, false, nil
	}

	ct, err := strconv.ParseInt(d.CreatedAt, 10, 64)
	if err != nil {
		return Domain{}, false, fmt.Errorf("creation date is not a number: %v", err)
	}
	return Domain{
		Name:        d.LabelName,
		CreatedDate: time.Unix(ct, 0),
	}, true, nil
}
