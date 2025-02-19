package main

import (
	"testing"
	"time"

	"github.com/YoavIsaacs/pokadexcli/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "poopoo poop                 poo",
			expected: []string{"poopoo", "poop", "poo"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Error: expected %d words, got %d", len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Error: expected %d words, got %d", len(c.expected), len(actual))
			}
		}
	}
}

func TestCacheAddAndGet(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
	testKey := "hello"
	testData := []byte("test-data")

	cache.Add(testKey, testData)

	data, err := cache.Get(testKey)
	if err != nil {
		t.Errorf("Failed to get item from cache: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("Cache returned wrong data. Expected %s, got %s", string(testData), string(data))
	}

	_, err = cache.Get("Non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent key, got nil")
	}
}
