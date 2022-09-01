package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/apbigcod/nns-metadata/logger"
	"github.com/apbigcod/nns-metadata/thegraph"
)

type Attribute struct {
	TraitType   string `json:"trait_type"`
	DisplayType string `json:"display_type"`
	Value       any
}

type Metadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Attributes  []Attribute `json:"attributes"`
	URL         string      `json:"url"`
	Version     int         `json:"version"`
	ImageURL    string      `json:"image_url"`
}

type Handler struct {
	WebURL      string
	APIURL      string
	GraphClient thegraph.Client
}

const (
	mainnet          = "mainnet"
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
	if network != mainnet {
		log.WithField("network", network).
			Error("unexpected network")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if contract != registrarAddress {
		log.WithField("contract", contract).
			Error("unexpected contract")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if token == "" {
		log.Error("empty token")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.
		WithField("token", token).
		Debug("processing request")

	d, found, err := h.GraphClient.Domain(token)
	if err != nil {
		log.WithError(err).
			Error("error getting domain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		log.Error("domain not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	image, err := SVGImage(d.Name)
	if err != nil {
		log.WithError(err).
			Error("error generating image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(parts) == 4 && parts[3] == "image" {
		w.Header().Add("content-type", "image/svg+xml")
		w.WriteHeader(http.StatusOK)
		w.Write(image)
		log.Info("returned image")
		return
	}

	md := Metadata{
		Name:        fmt.Sprintf("%s.⌐◨-◨", d.Name),
		Description: fmt.Sprintf("%s.⌐◨-◨ an NNS name", d.Name),
		Attributes: []Attribute{
			{TraitType: "Registered Date", DisplayType: "date", Value: d.CreatedDate.UnixMilli()},
			{TraitType: "Length", DisplayType: "number", Value: len(d.Name)},
		},
		URL:      fmt.Sprintf("%s/name/%s.⌐◨-◨", h.WebURL, d.Name),
		Version:  0,
		ImageURL: fmt.Sprintf("%s/%s/%s/%s/image", h.APIURL, mainnet, registrarAddress, token),
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
