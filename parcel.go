package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// добавление строки в таблицу parcel
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// возвращаем идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// чтение строки по заданному number
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// чтение строк из таблицы parcel по заданному client
	row, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer row.Close()
	// заполняем срез Parcel данными из таблицы
	var res []Parcel

	for row.Next() {
		var parcels Parcel

		err := row.Scan(&parcels.Number, &parcels.Client, &parcels.Status, &parcels.Address, &parcels.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, parcels)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}
