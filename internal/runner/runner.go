package runner

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/toozej/go-find-liquor/internal/notification"
	"github.com/toozej/go-find-liquor/internal/search"
	"github.com/toozej/go-find-liquor/pkg/config"
)

// Runner executes periodic searches
type Runner struct {
	config    config.Config
	searcher  *search.Searcher
	notifier  *notification.NotificationManager
	stopChan  chan struct{}
	runningCh chan struct{}
}

// NewRunner creates a new runner with the given configuration
func NewRunner(cfg config.Config) (*Runner, error) {
	// Initialize the searcher
	searcher := search.NewSearcher(cfg.UserAgent)

	// Initialize notification manager
	notifier, err := notification.NewNotificationManager(cfg.Notifications)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification manager: %w", err)
	}

	return &Runner{
		config:    cfg,
		searcher:  searcher,
		notifier:  notifier,
		stopChan:  make(chan struct{}),
		runningCh: make(chan struct{}, 1),
	}, nil
}

// Start begins periodic searches
func (r *Runner) Start(ctx context.Context) error {
	// Initial search
	go func() {
		r.runningCh <- struct{}{}
		defer func() {
			<-r.runningCh
		}()

		if err := r.runSearch(ctx); err != nil {
			log.Errorf("Search failed: %v", err)
		}
	}()

	// Setup ticker for recurring searches
	ticker := time.NewTicker(r.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if we're already running
			select {
			case r.runningCh <- struct{}{}:
				// We got the semaphore, run the search
				go func() {
					defer func() {
						<-r.runningCh
					}()

					if err := r.runSearch(ctx); err != nil {
						log.Errorf("Search failed: %v", err)
					}
				}()
			default:
				// A search is already running, skip this tick
				log.Warnf("Previous search still running, skipping")
			}
		case <-r.stopChan:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// runSearch performs a single search for all items
func (r *Runner) runSearch(ctx context.Context) error {
	if len(r.config.Items) == 0 {
		return fmt.Errorf("no items to search for")
	}

	if r.config.Zipcode == "" {
		return fmt.Errorf("zipcode is required")
	}

	log.Infof("Starting search for %d items within %d miles of %s",
		len(r.config.Items), r.config.Distance, r.config.Zipcode)

	for _, item := range r.config.Items {
		// Create a context with timeout for this item
		itemCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		log.Infof("Searching for item: %s", item)

		// Search for the item
		results, err := r.searcher.SearchItem(itemCtx, item, r.config.Zipcode, r.config.Distance)
		if err != nil {
			log.Errorf("Failed to search for %s: %v", item, err)
			continue
		}

		log.Infof("Found %d results for %s", len(results), item)

		// Process and notify about results
		for _, result := range results {
			if err := r.notifier.NotifyFound(itemCtx, result); err != nil {
				log.Warnf("Failed to send notification: %v", err)
			}
		}

		// Random wait between searches
		if len(r.config.Items) > 1 && item != r.config.Items[len(r.config.Items)-1] {
			randTimeBig := new(big.Int)
			randTimeBig.SetInt64(int64(30))
			randTime, _ := rand.Int(rand.Reader, randTimeBig)
			waitTime := time.Duration(randTime.Int64()) * time.Second
			log.Debugf("Waiting %s before next search", waitTime)

			select {
			case <-time.After(waitTime):
				// Continue to next item
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	log.Infof("Search completed, next search in %s", r.config.Interval)
	return nil
}

// Stop halts the runner
func (r *Runner) Stop() {
	close(r.stopChan)
}

// RunOnce performs a single search and returns
func (r *Runner) RunOnce(ctx context.Context) error {
	return r.runSearch(ctx)
}
