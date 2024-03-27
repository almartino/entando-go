package conference

import (
	"github.com/almartino/entando-go/oidc"
	"gorm.io/gorm"
	"time"
)

func principalFrom(tx *gorm.DB) (*oidc.Principal, bool) {
	ctx := tx.Statement.Context
	return oidc.PrincipalFrom[oidc.Principal](ctx)
}

type Conference struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	CreateBy string    `json:"createBy"`
	CreateTs time.Time `json:"createTs" gorm:"autoCreateTime"`
	UpdateBy string    `json:"updateBy"`
	UpdateTs time.Time `json:"updateTs" gorm:"autoUpdateTime"`
}

// BeforeSave implement the gorm hook
func (c *Conference) BeforeSave(tx *gorm.DB) (err error) {
	return c.setPrincipalID(tx)
}

// BeforeCreate implement the gorm hook
func (c *Conference) BeforeCreate(tx *gorm.DB) (err error) {
	return c.setPrincipalID(tx)
}

// setPrincipalID sets the principal id on CreateBy or UpdateBy fields.
func (c *Conference) setPrincipalID(tx *gorm.DB) (err error) {
	p, ok := principalFrom(tx)
	if !ok {
		// skip
		return nil
	}
	if c.ID > 0 {
		c.UpdateBy = p.ID
		return
	}
	c.CreateBy = p.ID
	return
}
