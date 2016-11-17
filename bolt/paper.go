package bolt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/bobinette/papernet"
)

var paperBucket = []byte("papers")

// PaperRepository is used to store and retrieve papers from a bolt database.
type PaperRepository struct {
	Driver *Driver
}

// Get retrieves the paper defined by id in the database. If no paper can be found with the
// given id, Get returns nil.
func (r *PaperRepository) Get(ids ...int) ([]*papernet.Paper, error) {
	papers := make([]*papernet.Paper, 0, len(ids))
	err := r.Driver.store.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(paperBucket)

		for _, id := range ids {
			data := bucket.Get(itob(id))
			if data == nil {
				continue
			}

			var paper papernet.Paper
			if err := json.Unmarshal(data, &paper); err != nil {
				return err
			}
			papers = append(papers, &paper)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return papers, nil
}

// Upsert inserts or update a paper in the database, depending on paper.ID.
func (r *PaperRepository) Upsert(paper *papernet.Paper) error {
	return r.Driver.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(paperBucket)

		if paper.ID <= 0 {
			id, err := bucket.NextSequence()
			if err != nil {
				return fmt.Errorf("error incrementing id: %v", err)
			}
			paper.ID = int(id)
		}

		data, err := json.Marshal(paper)
		if err != nil {
			return err
		}

		return bucket.Put(itob(paper.ID), data)
	})
}

func (r *PaperRepository) Delete(id int) error {
	return r.Driver.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(paperBucket)
		return bucket.Delete(itob(id))
	})
}

func (r *PaperRepository) List() ([]*papernet.Paper, error) {
	var papers []*papernet.Paper

	err := r.Driver.store.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(paperBucket)

		c := bucket.Cursor()
		for id, data := c.First(); id != nil; id, data = c.Next() {
			var paper papernet.Paper
			if err := json.Unmarshal(data, &paper); err != nil {
				return err
			}
			papers = append(papers, &paper)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return papers, nil
}

// ------------------------------------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------------------------------------

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
