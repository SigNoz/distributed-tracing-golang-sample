package datastore

type InsertParams struct {
	Query string
	Vars  []interface{}
}

type SelectParams struct {
	Query   string
	Filters []interface{}
	Result  []interface{}
}

type UpdateParams struct {
	Query string
	Vars  []interface{}
}

type DB interface {
	InsertOne(InsertParams) (int64, error)
	SelectOne(SelectParams) error
	UpdateOne(UpdateParams) error
	Close()
}
