package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofurry/reqrisk"
)

func main() {
	assessor := reqrisk.New()

	result := assessor.Evaluate(reqrisk.Features{
		Entropy: &reqrisk.EntropyFeature{
			Value:      7.9,
			SampleSize: 1024,
			Source:     "body",
		},
		Complexity: &reqrisk.ComplexityFeature{
			Depth:          7,
			FieldCount:     42,
			MaxArrayLength: 18,
			Source:         "json",
		},
		Charset: &reqrisk.CharsetFeature{
			NonASCIIRatio:  0.82,
			ZeroWidthCount: 3,
			MixedScripts:   true,
			Source:         "body",
		},
		Fingerprint: &reqrisk.FingerprintFeature{
			Rarity:     0.91,
			Volatility: 0.63,
			Source:     "request-shape",
		},
		Meta: reqrisk.Meta{
			Target: "body",
			Route:  "/api/comment",
			Method: "POST",
		},
	})

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "encode result: %v\n", err)
		os.Exit(1)
	}
}
