package services

import (
	"gitlab.com/wuup-lab/deckabot-go/pkg/database"
	"gitlab.com/bigtallbill/expenses/pkg/docs"
	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	conn *database.Connection
}

func NewProduct(conn *database.Connection) Product {
	return Product{conn: conn}
}

func (p *Product) Save(name string, department string) {
	p.conn.GetWriteDatabase().C("products").Insert(
		docs.Product{
			Id:            bson.NewObjectId(),
			Name:          name,
			Department:    department,
			ShopCodeNames: make(map[bson.ObjectId]string),
		},
	)
}

func (p *Product) List() (products []docs.Product) {
	p.conn.GetReadDatabase().C("products").Find(bson.M{}).All(&products)
	return
}

func (p *Product) FindOne(id bson.ObjectId) (product docs.Product) {
	p.conn.GetReadDatabase().C("products").Find(bson.M{"_id": id}).One(&product)
	return
}

