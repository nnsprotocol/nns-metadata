package main

import (
	"net/http"
	"os"

	nnsmetadata "github.com/apbigcod/nns-metadata"
	"github.com/apbigcod/nns-metadata/metadata"
	"github.com/apbigcod/nns-metadata/thegraph"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
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
