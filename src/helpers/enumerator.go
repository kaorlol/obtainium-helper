package helpers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"obtainium-helper/src/utils"
)

func getSplitIdentifier(identifier utils.Identifier) ([]int, error) {
	var identifiers []int

	for _, v := range strings.Split(identifier.Latest, ".") {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		identifiers = append(identifiers, num)
	}

	return identifiers, nil
}

func identifierToString(identifiers []int) string {
	var str strings.Builder
	for i, v := range identifiers {
		if i > 0 {
			str.WriteString(".")
		}
		str.WriteString(fmt.Sprintf("%d", v))
	}
	return str.String()
}

func incrementIdentifierBy(identifier []int, count, limit int) {
	for i := len(identifier) - 1; i >= 0; i-- {
		count += identifier[i]
		identifier[i] = count % (limit + 1)
		count /= (limit + 1)
	}
}

func Enumerate(ctx context.Context, app utils.Download) (string, error) {
	initialIdentifier, err := getSplitIdentifier(app.Identifier)
	if err != nil {
		return "", err
	}

	numWorkers := app.Identifier.EnumLimit / 10
	results := make(chan string, numWorkers)
	var wg sync.WaitGroup
	var identifierCounter int
	identifierMutex := &sync.Mutex{}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					identifierMutex.Lock()
					if identifierCounter >= app.Identifier.EnumLimit {
						identifierMutex.Unlock()
						cancel()
						return
					}

					currentIdentifier := make([]int, len(initialIdentifier))
					copy(currentIdentifier, initialIdentifier)
					incrementIdentifierBy(currentIdentifier, identifierCounter, app.Identifier.IncrementLimit)
					identifierCounter++
					identifierMutex.Unlock()

					newIdentifier := identifierToString(currentIdentifier)
					resp, err := utils.Request(fmt.Sprintf("%s%s.apk", app.URL, newIdentifier), app.Agent)
					if err != nil {
						continue
					}

					if resp.StatusCode == http.StatusOK {
						results <- newIdentifier
						cancel()
						return
					}
				}
			}
		}(i)
	}

	var foundIdentifier string
	select {
	case foundIdentifier = <-results:
		cancel()
	case <-ctx.Done():
	}

	wg.Wait()
	close(results)

	if foundIdentifier == "" {
		return "", nil
	}

	url := fmt.Sprintf("%s%s.apk", app.URL, foundIdentifier)
	return url, nil
}
