package metadata

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	nnsmetadata "github.com/nnsprotocol/nns-metadata"
	"github.com/nnsprotocol/nns-metadata/logger"
	"github.com/nnsprotocol/nns-metadata/thegraph"
)

type Attribute struct {
	TraitType   string `json:"trait_type"`
	DisplayType string `json:"display_type"`
	Value       any    `json:"value"`
}

type Metadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Attributes  []Attribute `json:"attributes"`
	URL         string      `json:"url"`
	Version     int         `json:"version"`
	ImageURL    string      `json:"image_url"`
	Image       string      `json:"image"`
}

type Handler struct {
	WebURL      string
	APIURL      string
	GraphClient thegraph.Client
}

const (
	registrarAddress = "0x2e12d051c3Be2aAA932EB09CE4Cba8F58fa1860e"
)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log := logger.New("metadata")
	var parts []string
	for _, p := range strings.Split(r.URL.Path, "/") {
		if p == "" {
			continue
		}
		parts = append(parts, p)
	}

	if len(parts) < 3 {
		log.WithField("path", r.URL.Path).
			Error("unexpected parts")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.
		WithField("parts", parts).
		Debug("valid path")

	network := parts[0]
	contract := parts[1]
	token := parts[2]
	if token == "" {
		log.Error("empty token")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	net, ok := nnsmetadata.ParseNetwork(network)
	if !ok {
		log.WithField("network", network).
			Error("unexpected network")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// }
	// if contract != registrarAddress {
	// 	log.WithField("contract", contract).
	// 		Error("unexpected contract")
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	thex, ok := formatHashToHex(token)
	if !ok {
		log.
			WithField("token", token).
			Error("invalid token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d, found, err := h.GraphClient.Domain(net, thex)
	if err != nil {
		log.WithError(err).
			WithField("token", token).
			Error("error getting domain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		log.WithField("token", token).
			WithField("hex", thex).
			Error("domain not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.respond(w, parts, d, network, contract, token)
}

func (h Handler) respond(w http.ResponseWriter, parts []string, d thegraph.Domain, network, contract, token string) {
	log := logger.New("metadata")
	if len(parts) == 4 && parts[3] == "image" {
		image, err := SVGImage(d.Name)
		if err != nil {
			log.WithError(err).Error("error creating image")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		w.Write(image)
		log.Info("returned image")
		return
	}

	imgURL := fmt.Sprintf("%s/%s/%s/%s/image", h.APIURL, network, contract, token)
	md := Metadata{
		Name:        d.Name,
		Description: fmt.Sprintf("%s.⌐◨-◨ an NNS name", d.Name),
		Attributes: []Attribute{
			{TraitType: "Length", DisplayType: "number", Value: len(d.Name)},
		},
		URL:      fmt.Sprintf("%s/name/%s.⌐◨-◨", h.WebURL, d.Name),
		Version:  0,
		ImageURL: imgURL,
		Image:    imgURL,
	}
	b, err := json.Marshal(md)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	log.Info("returned metadata")
}

func formatHashToHex(v string) (string, bool) {
	n := new(big.Int)
	n, ok := n.SetString(v, 0)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("0x%064x", n), true
}
