package validator

import (
	"regexp"
	"strconv"
	"sync"
)

// regexCache caches compiled regex patterns
var regexCache = &regexpCache{
	cache: make(map[string]*regexp.Regexp),
	mu:    sync.RWMutex{},
}

type regexpCache struct {
	cache map[string]*regexp.Regexp
	mu    sync.RWMutex
}

// getRegex returns a compiled regex, compiling and caching if necessary
func getRegex(pattern string) (*regexp.Regexp, error) {
	regexCache.mu.RLock()
	if re, ok := regexCache.cache[pattern]; ok {
		regexCache.mu.RUnlock()
		return re, nil
	}
	regexCache.mu.RUnlock()

	regexCache.mu.Lock()
	defer regexCache.mu.Unlock()

	// Double-check after acquiring write lock
	if re, ok := regexCache.cache[pattern]; ok {
		return re, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	regexCache.cache[pattern] = re
	return re, nil
}

// PrecompileRegexes precompiles common regex patterns
func PrecompileRegexes() {
	patterns := []string{
		mailFormatRegex,
		alphanumericFormatRegex,
		alphabetFormatRegex,
		numericFormatRegex,
	}

	for _, pattern := range patterns {
		_, _ = getRegex(pattern)
	}
}

// parseIntCache caches parsed integer values
var parseIntCache = &intCache{
	cache: make(map[string]int),
	mu:    sync.RWMutex{},
}

type intCache struct {
	cache map[string]int
	mu    sync.RWMutex
}

// getInt parses and caches integer values
func getInt(s string) (int, error) {
	parseIntCache.mu.RLock()
	if val, ok := parseIntCache.cache[s]; ok {
		parseIntCache.mu.RUnlock()
		return val, nil
	}
	parseIntCache.mu.RUnlock()

	parseIntCache.mu.Lock()
	defer parseIntCache.mu.Unlock()

	// Double-check after acquiring write lock
	if val, ok := parseIntCache.cache[s]; ok {
		return val, nil
	}

	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	parseIntCache.cache[s] = val
	return val, nil
}
