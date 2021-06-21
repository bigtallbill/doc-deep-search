package services

import (
	"gitlab.com/wuup-lab/deckabot-go/pkg/database"
	"gitlab.com/bigtallbill/expenses/pkg/docs"
	"gopkg.in/mgo.v2/bson"
)

type Shop struct {
	conn *database.Connection
}

func NewShop(conn *database.Connection) Shop {
	return Shop{conn: conn}
}

func (s *Shop) Save(shopName string) {
	s.conn.GetWriteDatabase().C("shops").Insert(docs.Shop{Id: bson.NewObjectId(), Name: shopName})
}

func (s *Shop) List() (shops []docs.Shop) {
	s.conn.GetWriteDatabase().C("shops").Find(bson.M{}).All(&shops)
	return
}
