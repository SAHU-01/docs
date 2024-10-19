package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
)

type PricingData struct {
	Documentation string `json:"documentation"`
	EntireRow     struct {
		Feature          string `json:"Feature"`
		SubscriptionTier string `json:"Subscription Tier"`
	} `json:"entire_row"`
}

func featureSupport(from io.Reader, to io.Writer) error {
	// Read the pricing data
	jsonFile, err := ioutil.ReadFile("data/pricing_data.json")
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	var pricingData []PricingData
	err = json.Unmarshal(jsonFile, &pricingData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Read the content
	content, err := io.ReadAll(from)
	if err != nil {
		return err
	}

	// Parse the content
	page, err := pageparser.Parse(strings.NewReader(string(content)), pageparser.Config{})
	if err != nil {
		return err
	}

	// Extract the front matter from the parsed page
	meta := page.Meta() // Meta contains the front matter and other metadata
	if meta == nil {
		return nil // No front matter found
	}

	// Check if the front matter contains a 'features' key
	features, ok := meta["features"].([]interface{})
	if !ok {
		return nil // No features to process
	}

	// Prepare the new content
	var newContent strings.Builder
	newContent.Write(content)

	// Add shortcodes for each feature
	for _, feature := range features {
		featureStr, ok := feature.(string)
		if !ok {
			continue
		}

		// Find the pricing data for this feature
		var tierInfo string
		for _, pd := range pricingData {
			if pd.EntireRow.Feature == featureStr {
				tierInfo = pd.EntireRow.SubscriptionTier
				break
			}
		}

		if tierInfo == "" {
			continue
		}

		shortcode := fmt.Sprintf("\n{{< feature-support feature_name=\"%s\" >}}\n", featureStr)
		newContent.WriteString(shortcode)
	}

	// Write the new content
	_, err = to.Write([]byte(newContent.String()))
	return err
}

func main() {
	// This is a placeholder for any logic you may need to initialize
	// or run the featureSupport function as required by your app.
}
