package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	nnsmetadata "github.com/nnsprotocol/nns-metadata"
	"github.com/nnsprotocol/nns-metadata/metadata"
	"github.com/nnsprotocol/nns-metadata/thegraph"
)

func main() {
	tgc := thegraph.Client{
		GraphNames: map[nnsmetadata.Network]string{
			nnsmetadata.Goerli:  "apbigcod/nns-goerli",
			nnsmetadata.Mainnet: "apbigcod/nns",
		},
	}
	h := metadata.Handler{
		GraphClient: tgc,
		WebURL:      "https://nns.xyz",
	}
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// on Lambda
		h.APIURL = "https://metadata.nns.xyz"
		lambda.Start(httpadapter.NewV2(h).ProxyWithContext)
	} else {
		// Local
		h.APIURL = "http://localhost:8081"
		http.ListenAndServe(":8081", h)
	}
}
