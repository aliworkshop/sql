package sql

import (
	"github.com/aliworkshop/configer"
	"github.com/aliworkshop/dbcore"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type Category struct {
	Id           int       `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EnglishTitle string    `json:"english_title"`
	IsPublished  bool      `json:"is_published"`
	Level        int       `json:"level"`
	ParentId     int       `json:"parent_id"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
}

func (Category) TableName() string {
	return "accounting"
}

func TestRepo_Count(t *testing.T) {
	registry := configer.New()
	registry.SetConfigType("yaml")
	f, err := os.Open("./config.sample.yaml")
	if err != nil {
		panic("cannot read config: " + err.Error())
	}
	err = registry.ReadConfig(f)
	if err != nil {
		panic("cannot read config" + err.Error())
	}

	db := NewRepository(registry, nil)
	err = db.Initialize()
	assert.Nil(t, err)

	q := dbcore.NewQuery().WithModelFunc(func() dbcore.Modeler {
		return new(Category)
	}).WithSelect("SUM(credit_amount) - SUM(dept_amount) as balance").WithFilter(dbcore.NewFilter().WithOrMatch(&dbcore.Match{
		Key:   "user_id",
		Value: 12,
	})).WithSorts(dbcore.SortItem{Field: "play_count", ReplaceWith: "playCount(id)"})
	count, err := db.Get(q)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}
