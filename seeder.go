package gormseeder

import (
	"errors"
	"slices"
	"sort"
	"time"

	"gorm.io/gorm"
)

var (
	ErrSeedersTableNotFound = errors.New("seeders table not found")
)

type seederRecord struct {
	Key       string `gorm:"primaryKey"`
	Name      *string
	CreatedAt time.Time
}

func (seederRecord) TableName() string {
	return "seeders"
}

type SeederFunc func(tx *gorm.DB) (err error)

type Seed struct {
	Key      string
	Name     string
	Seed     SeederFunc
	Rollback SeederFunc
}

type Seeder struct {
	db      *gorm.DB
	seeders []*Seed
}

func New(db *gorm.DB, seeders []*Seed) *Seeder {
	sort.Slice(seeders, func(i, j int) bool {
		return seeders[i].Key < seeders[j].Key
	})

	s := &Seeder{
		db:      db,
		seeders: seeders,
	}

	return s
}

func (s *Seeder) Add(seeders *Seed) {
	s.seeders = append(s.seeders, seeders)
}

func (s *Seeder) Seed() error {
	if err := s.createTable(); err != nil {
		return err
	}

	for _, seeder := range s.seeders {
		if err := s.seed(seeder); err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) Rollback() error {
	if !s.db.Migrator().HasTable(seederRecord{}) {
		return ErrSeedersTableNotFound
	}

	// reverse seeders order, so we can start from the last seeders
	slices.Reverse(s.seeders)

	for _, seeder := range s.seeders {
		err := s.rollback(seeder)
		if err == gorm.ErrRecordNotFound {
			// the seeder may have been deleted, so we continue for the next record
			continue
		}
		if err != nil {
			return err
		}

		// return on success, so we only rollback one seeder.
		return nil
	}

	// need to reverse reverse back to original order
	slices.Reverse(s.seeders)

	return nil
}

func (s *Seeder) createTable() error {
	if s.db.Migrator().HasTable(seederRecord{}) {
		return nil
	}

	return s.db.AutoMigrate(&seederRecord{})
}

func (s *Seeder) seed(seeders *Seed) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	result := tx.First(&seederRecord{}, seeders.Key)
	if result.Error == nil {
		return nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	var seederName *string

	if seeders.Name != "" {
		seederName = &seeders.Name
	}

	result = tx.Create(&seederRecord{Key: seeders.Key, Name: seederName})
	if result.Error != nil {
		return result.Error
	}

	if err := seeders.Seed(tx); err != nil {
		return err
	}

	return tx.Commit().Error
}

func (s *Seeder) rollback(seeders *Seed) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	result := tx.First(&seederRecord{}, seeders.Key)
	if result.Error != nil {
		return result.Error
	}

	err := seeders.Rollback(tx)
	if err != nil {
		return err
	}

	result = tx.Delete(&seederRecord{}, seeders.Key)
	if result.Error != nil {
		return result.Error
	}

	return tx.Commit().Error
}
