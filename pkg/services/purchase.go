package services

import (
	"gitlab.com/wuup-lab/deckabot-go/pkg/database"
	"gitlab.com/bigtallbill/expenses/pkg/docs"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Purchase struct {
	conn *database.Connection
}

func NewPurchase(conn *database.Connection) Purchase {
	return Purchase{conn: conn}
}

func (p *Purchase) Save(shopId string, items []docs.Item, subtotal float64, tax float64, total float64) {
	p.conn.GetWriteDatabase().C("purchases").Insert(
		docs.Purchase{
			Id:       bson.NewObjectId(),
			Date:     time.Now(),
			Shop:     bson.ObjectIdHex(shopId),
			Items:    items,
			SubTotal: subtotal,
			Tax:      tax,
			Total:    total,
		},
	)
}

func (p *Purchase) List() (purchases []docs.Purchase) {
	p.conn.GetReadDatabase().C("purchases").Find(bson.M{}).Sort("date-").All(&purchases)
	return
}
