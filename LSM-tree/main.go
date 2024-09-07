package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// MemTable represents an in-memory key-value store.
type MemTable struct {
	data map[string]string
}

// NewMemTable initializes a new MemTable.
func NewMemTable() *MemTable {
	return &MemTable{data: make(map[string]string)}
}

// Set adds a key-value pair to the MemTable.
func (m *MemTable) Set(key, value string) {
	m.data[key] = value
}

// Get retrieves a value by key from the MemTable.
func (m *MemTable) Get(key string) (string, bool) {
	value, ok := m.data[key]
	return value, ok
}

// Flush writes the MemTable contents to a new SSTable.
func (m *MemTable) Flush(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Get the keys and sort them
	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	sort.Strings(keys) // Sort keys

	// Write sorted key-value pairs to the SSTable
	for _, key := range keys {
		value := m.data[key]
		_, err := writer.WriteString(fmt.Sprintf("%s:%s\n", key, value))
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

// SSTable represents an immutable file-based key-value store.
type SSTable struct {
	filename string
	index    []string
	offsets  []int64
}

// NewSSTable initializes a new SSTable and builds an in-memory index of keys.
func NewSSTable(filename string) (*SSTable, error) {
	sstable := &SSTable{filename: filename}
	err := sstable.buildIndex()
	if err != nil {
		return nil, err
	}
	return sstable, nil
}

// buildIndex builds an in-memory index of keys and their file offsets for binary search.
func (s *SSTable) buildIndex() error {
	file, err := os.Open(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var offset int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			s.index = append(s.index, parts[0])
			s.offsets = append(s.offsets, offset)
		}
		offset += int64(len(line)) + 1 // +1 for the newline character
	}
	return scanner.Err()
}

// Get retrieves a value by key from the SSTable using binary search and direct file access.
func (s *SSTable) Get(key string) (string, bool) {
	// Perform binary search on the in-memory index
	i := sort.SearchStrings(s.index, key)
	if i < len(s.index) && s.index[i] == key {
		// If key is found in the index, seek to the corresponding offset in the file
		file, err := os.Open(s.filename)
		if err != nil {
			return "", false
		}
		defer file.Close()

		// Seek sets the offset for the next Read or Write on file to offset
		_, err = file.Seek(s.offsets[i], 0)
		if err != nil {
			return "", false
		}

		// Read the line from the file
		reader := bufio.NewReader(file)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", false
		}

		parts := strings.Split(line, ":")
		if parts[0] == key {
			// Trim the newline character from the value
			return strings.TrimSpace(parts[1]), true
		}
	}
	return "", false
}

// LSMTree represents the high-level LSM Tree structure.
type LSMTree struct {
	memTable *MemTable
	sstables []*SSTable
}

// NewLSMTree initializes a new LSMTree.
func NewLSMTree() *LSMTree {
	return &LSMTree{
		memTable: NewMemTable(),
		sstables: []*SSTable{},
	}
}

// Set adds a key-value pair to the LSM Tree.
func (l *LSMTree) Set(key, value string) {
	l.memTable.Set(key, value)
}

// Get retrieves a value by key from the LSM Tree.
func (l *LSMTree) Get(key string) (string, bool) {
	if value, ok := l.memTable.Get(key); ok {
		return value, true
	}

	// Search in SSTables (in reverse order for the most recent data)
	for i := len(l.sstables) - 1; i >= 0; i-- {
		if value, ok := l.sstables[i].Get(key); ok {
			return value, true
		}
	}
	return "", false
}

// Flush writes the MemTable to a new SSTable and clears the MemTable.
func (l *LSMTree) Flush() error {
	filename := "sstable_" + strconv.Itoa(len(l.sstables)) + ".txt"
	err := l.memTable.Flush(filename)
	if err != nil {
		return err
	}

	sstable, err := NewSSTable(filename)
	if err != nil {
		return err
	}

	l.sstables = append(l.sstables, sstable)
	l.memTable = NewMemTable() // Reset the MemTable after flushing
	return nil
}

func main() {
	// Create a new LSM Tree
	lsmTree := NewLSMTree()

	// Insert key-value pairs
	lsmTree.Set("key1", "value1")
	lsmTree.Set("key2", "value2")
	lsmTree.Set("key3", "value3")

	// Flush MemTable to SSTable
	err := lsmTree.Flush()
	if err != nil {
		fmt.Println("Error flushing MemTable:", err)
		return
	}

	// Insert more key-value pairs
	lsmTree.Set("key4", "value4")
	lsmTree.Set("key5", "value5")

	// Retrieve a value
	if value, ok := lsmTree.Get("key3"); ok {
		fmt.Println("Found key3:", value)
	} else {
		fmt.Println("key3 not found")
	}

	// Flush again
	err = lsmTree.Flush()
	if err != nil {
		fmt.Println("Error flushing MemTable:", err)
		return
	}

	// Retrieve values from flushed SSTables
	if value, ok := lsmTree.Get("key2"); ok {
		fmt.Println("Found key2:", value)
	} else {
		fmt.Println("key2 not found")
	}

	if value, ok := lsmTree.Get("key4"); ok {
		fmt.Println("Found key4:", value)
	} else {
		fmt.Println("key4 not found")
	}
}
