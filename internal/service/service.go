package service

import iDB "cashier/internal/repository/database"

type service struct {
	db iDB.IDatabase
}

func New(db iDB.IDatabase) IService {
	return &service{
		db: db,
	}
}
