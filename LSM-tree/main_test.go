package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemTable(t *testing.T) {
	memTable := NewMemTable()

	// Test Set and Get
	memTable.Set("key1", "value1")
	value, ok := memTable.Get("key1")
	assert.True(t, ok, "Expected key1 to exist")
	assert.Equal(t, "value1", value, "Expected value1")

	// Test Get non-existent key
	_, ok = memTable.Get("key2")
	assert.False(t, ok, "Expected key2 to not exist")
}

func TestMemTableFlush(t *testing.T) {
	memTable := NewMemTable()
	memTable.Set("key1", "value1")
	memTable.Set("key2", "value2")

	filename := "test_sstable.txt"
	defer os.Remove(filename)

	err := memTable.Flush(filename)
	assert.NoError(t, err, "Failed to flush MemTable")

	// Verify file contents
	file, err := os.Open(filename)
	assert.NoError(t, err, "Failed to open SSTable file")
	defer file.Close()

	stat, err := file.Stat()
	assert.NoError(t, err, "Failed to stat SSTable file")
	assert.NotEqual(t, 0, stat.Size(), "Expected non-empty SSTable file")

	// Verify file contents
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	assert.NoError(t, scanner.Err(), "Failed to read SSTable file")
	assert.Equal(t, 2, len(lines), "Expected 2 lines in SSTable file")
	assert.Equal(t, "key1:value1", lines[0], "Expected key1:value1")
	assert.Equal(t, "key2:value2", lines[1], "Expected key2:value2")
}

func TestSSTable(t *testing.T) {
	filename := "test_sstable.txt"
	defer os.Remove(filename)

	// Create a test SSTable file
	file, err := os.Create(filename)
	assert.NoError(t, err, "Failed to create SSTable file")
	file.WriteString("key1:value1\nkey2:value2\n")
	file.Close()

	sstable, err := NewSSTable(filename)
	assert.NoError(t, err, "Failed to create SSTable")

	// Test Get
	value, ok := sstable.Get("key1")
	assert.True(t, ok, "Expected key1 to exist")
	assert.Equal(t, "value1", value, "Expected value1")

	// Test Get non-existent key
	_, ok = sstable.Get("key3")
	assert.False(t, ok, "Expected key3 to not exist")
}

func TestLSMTree(t *testing.T) {
	lsmTree := NewLSMTree()

	// Test Set and Get
	lsmTree.Set("key1", "value1")
	value, ok := lsmTree.Get("key1")
	assert.True(t, ok, "Expected key1 to exist")
	assert.Equal(t, "value1", value, "Expected value1")

	// Test Flush
	err := lsmTree.Flush()
	assert.NoError(t, err, "Failed to flush LSMTree")

	// Test Get after flush
	value, ok = lsmTree.Get("key1")
	assert.True(t, ok, "Expected key1 to exist after flush")
	assert.Equal(t, "value1", value, "Expected value1 after flush")

	// Test Get non-existent key
	_, ok = lsmTree.Get("key2")
	assert.False(t, ok, "Expected key2 to not exist")
}
