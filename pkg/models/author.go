package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)
type DateOnly struct {
    time.Time
}

const layout = "2006-01-02"

func (d *DateOnly) UnmarshalJSON(b []byte) error {
    s := strings.Trim(string(b), `"`)
    t, err := time.Parse(layout, s)
    if err != nil {
        return err
    }
    d.Time = t
    return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
    return []byte(`"` + d.Time.Format(layout) + `"`), nil
}

func (d *DateOnly) Scan(value interface{}) error {
    if val, ok := value.(time.Time); ok {
        d.Time = val
        return nil
    }
    return fmt.Errorf("cannot scan value %v into DateOnly", value)
}

func (d DateOnly) Value() (driver.Value, error) {
    return d.Time, nil
}
type Author struct {
    UserID    uint      `json:"-" gorm:"primaryKey;not null;uniqueIndex;constraint:OnDelete:CASCADE;"`
    User      *UserDB   `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Firstname string    `json:"Firstname" binding:"required"`
	Lastname  string    `json:"Lastname"  binding:"required"`
	Birthday  DateOnly  `json:"Birthday"  binding:"required" gorm:"type:date"`
	Books     []Book    `gorm:"foreignKey:AuthorID"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type AuthorResp struct {
    UserID    uint      `json:"user_id"`
    Firstname string    `json:"firstname"`
    Lastname  string    `json:"lastname"`
    Birthday  DateOnly  `json:"birthday"`
}