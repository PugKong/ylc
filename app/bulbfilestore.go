package app

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
)

type BulbFileStore struct {
	bulbs map[string]Bulb
	dir   string
}

func NewBulbFileStore(dir string) *BulbFileStore {
	return &BulbFileStore{
		bulbs: make(map[string]Bulb),
		dir:   dir,
	}
}

func (b *BulbFileStore) Init() error {
	filePath := b.bulbsPath()
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("read data from %q: %w", filePath, err)
	}

	err = json.Unmarshal(bytes, &b.bulbs)
	if err != nil {
		return fmt.Errorf("decode data from %q: %w", filePath, err)
	}

	return nil
}

func (b *BulbFileStore) All() []Bulb {
	bulbs := make([]Bulb, 0, len(b.bulbs))
	for _, bulb := range b.bulbs {
		bulbs = append(bulbs, bulb)
	}

	slices.SortFunc(bulbs, func(x, y Bulb) int { return cmp.Compare(x.Name, y.Name) })

	return bulbs
}

func (b *BulbFileStore) AllNames() []string {
	bulbs := b.All()
	names := make([]string, 0, len(bulbs))
	for _, bulb := range bulbs {
		names = append(names, bulb.Name)
	}

	return names
}

var ErrBulbNotFound = errors.New("not found")

func (b *BulbFileStore) FindByID(id string) (Bulb, error) {
	bulb, ok := b.bulbs[id]
	if ok {
		return bulb, nil
	}

	return Bulb{}, ErrBulbNotFound
}

func (b *BulbFileStore) FindByName(name string) (Bulb, error) {
	for _, bulb := range b.bulbs {
		if bulb.Name == name {
			return bulb, nil
		}
	}

	return Bulb{}, ErrBulbNotFound
}

func (b *BulbFileStore) Save(bulb Bulb) {
	b.bulbs[bulb.ID] = bulb
}

func (b *BulbFileStore) Delete(bulb Bulb) {
	delete(b.bulbs, bulb.ID)
}

func (b *BulbFileStore) Flush() error {
	bytes, err := json.Marshal(b.bulbs)
	if err != nil {
		return fmt.Errorf("encode data: %w", err)
	}

	filePath := b.bulbsPath()
	if err := os.WriteFile(filePath, bytes, 0o600); err != nil {
		return fmt.Errorf("save data to %q: %w", filePath, err)
	}

	return nil
}

func (b *BulbFileStore) bulbsPath() string {
	return path.Join(b.dir, "bulbs.json")
}
