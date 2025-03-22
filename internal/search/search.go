package search

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	baseURL       = "http://www.oregonliquorsearch.com/"
	searchURL     = "http://www.oregonliquorsearch.com/servlet/FrontController"
	ageBtnFormURL = "http://www.oregonliquorsearch.com/servlet/WelcomeController"
)

// User agent strings to cycle through
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
}

// LiquorItem represents a found liquor item
// with only the information we care about
type LiquorItem struct {
	Name  string
	Code  string
	Store string
	Date  time.Time
	Price string
}

// ProductInfo represents all the possible information about a liquor item
// including the information we don't really care about
type ProductInfo struct {
	ItemCode    string
	Name        string
	BottlePrice string
	CasePrice   string
	Size        string
	Proof       string
	Category    string
}

// Searcher provides functionality to search for liquor items
type Searcher struct {
	client     *http.Client
	userAgent  string
	cycleAgent bool
}

// NewSearcher creates a new searcher with cookie support
func NewSearcher(userAgent string) *Searcher {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	bigLenUserAgents := new(big.Int)
	bigLenUserAgents.SetInt64(int64(len(userAgents))) // Convert int to int64 first
	randUserAgent, _ := rand.Int(rand.Reader, bigLenUserAgents)
	cycleAgent := userAgent == ""
	if cycleAgent {
		userAgent = userAgents[randUserAgent.Int64()]
	}

	return &Searcher{
		client:     client,
		userAgent:  userAgent,
		cycleAgent: cycleAgent,
	}
}

// updateUserAgent sets a new random user agent if cycling is enabled
func (s *Searcher) updateUserAgent() {
	if s.cycleAgent {
		bigLenUserAgents := new(big.Int)
		bigLenUserAgents.SetInt64(int64(len(userAgents))) // Convert int to int64 first
		randUserAgent, _ := rand.Int(rand.Reader, bigLenUserAgents)
		s.userAgent = userAgents[randUserAgent.Int64()]
		log.Debugf("Using user agent: %s", s.userAgent)
	}
}

// AgeVerification performs the age verification
func (s *Searcher) AgeVerification() error {
	// First get the page to get session cookies
	// nosemgrep: problem-based-packs.insecure-transport.go-stdlib.http-customized-request.http-customized-request
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}
	defer resp.Body.Close()

	// Parse the form for the age verification
	_, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse page: %w", err)
	}

	// Prepare the form submission for age verification
	formData := url.Values{}
	formData.Set("ageCheck", "true")
	formData.Set("action", "search")

	// Submit the form
	if viper.GetBool("debug") {
		log.Debugf("AgeVerification() POSTing %v\n", formData)
	}
	// nosemgrep: problem-based-packs.insecure-transport.go-stdlib.http-customized-request.http-customized-request
	req, err = http.NewRequest("POST", ageBtnFormURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create form submission request: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ageBtnFormURL)

	resp, err = s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit age verification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("age verification failed with status: %s", resp.Status)
	}

	return nil
}

// SearchItem searches for a specific liquor item by name or code
func (s *Searcher) SearchItem(ctx context.Context, item string, zipcode string, distance int) ([]LiquorItem, error) {
	s.updateUserAgent()

	// Perform age verification before search
	if err := s.AgeVerification(); err != nil {
		return nil, fmt.Errorf("age verification failed: %w", err)
	}

	// Prepare search form data
	formData := url.Values{}
	formData.Set("view", "global")
	formData.Set("action", "search")
	formData.Set("radiusSearchParam", fmt.Sprintf("%d", distance))
	formData.Set("productSearchParam", item)
	formData.Set("locationSearchParam", zipcode)
	formData.Set("btnSearch", "Search")

	// Submit search form
	if viper.GetBool("debug") {
		log.Debugf("SearchItem() POSTing formData %v\n", formData)
	}
	// nosemgrep: problem-based-packs.insecure-transport.go-stdlib.http-customized-request.http-customized-request
	req, err := http.NewRequest("POST", searchURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", searchURL)

	// Perform search request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %s", resp.Status)
	}

	// Generate goquery document from response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to generate goquery document from search query response: %w", err)
	}

	// Extract product information
	product := extractProductInfo(doc)

	// Extract results from the table and generate list of found LiquorItem
	results := extractResults(doc, product)

	return results, nil
}

// extractResults extracts found products from the table and creates a list of found liquor item results
func extractResults(doc *goquery.Document, product ProductInfo) []LiquorItem {
	var results []LiquorItem
	currentStore := ""

	doc.Find("tr.row, tr.alt-row").Each(func(i int, s *goquery.Selection) {
		// Check if the store has stock
		qtyText := strings.TrimSpace(s.Find("td.qty").Text())
		if qtyText == "0" {
			return // Skip stores with no stock
		}

		currentStore = strings.TrimSpace(s.Find("td").Eq(2).Text())

		if currentStore != "" {
			results = append(results, LiquorItem{
				Name:  product.Name,
				Code:  product.ItemCode,
				Store: currentStore,
				Date:  time.Now(),
				Price: product.BottlePrice,
			})
		}
	})

	return results
}

// extractProductInfo extracts product details from the product-details table
func extractProductInfo(doc *goquery.Document) ProductInfo {
	product := ProductInfo{}

	// Extract product name and item code from the product description
	productDesc := strings.TrimSpace(doc.Find("#product-desc h2").Text())
	if productDesc != "" {
		// Parse "Item 99900733075(7330B): MICHTER'S STRAIGHT RYE"
		parts := strings.SplitN(productDesc, ":", 2)
		if len(parts) == 2 {
			// Extract the item code
			itemParts := strings.Split(parts[0], " ")
			if len(itemParts) >= 2 {
				fullCode := itemParts[1]
				// Extract the code in parentheses if it exists
				codeInParens := ""
				if i := strings.Index(fullCode, "("); i != -1 {
					if j := strings.Index(fullCode, ")"); j != -1 && j > i {
						codeInParens = fullCode[i+1 : j]
					}
				}

				if codeInParens != "" {
					product.ItemCode = codeInParens
				} else {
					product.ItemCode = fullCode
				}
			}

			// Extract the product name
			product.Name = strings.TrimSpace(parts[1])
		}
	}

	// Extract bottle price
	doc.Find("#product-details tr").Each(func(i int, s *goquery.Selection) {
		s.Find("th").Each(func(j int, th *goquery.Selection) {
			label := strings.TrimSpace(th.Text())
			switch label {
			case "Bottle Price:":
				product.BottlePrice = strings.TrimSpace(th.Next().Text())
			case "Case Price:":
				product.CasePrice = strings.TrimSpace(th.Next().Text())
			case "Size:":
				product.Size = strings.TrimSpace(th.Next().Text())
			case "Proof:":
				product.Proof = strings.TrimSpace(th.Next().Text())
			case "Category:":
				product.Category = strings.TrimSpace(th.Next().Text())
			}
		})
	})

	return product
}
