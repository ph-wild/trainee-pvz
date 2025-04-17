package models

import "time"

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type PVZ struct {
	ID               string    `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}
type Reception struct {
	ID       string    `db:"id"`
	DateTime time.Time `db:"datetime"`
	PVZID    string    `db:"pvz_id"`
	Status   string    `db:"status"`
}

type Product struct {
	ID          string    `db:"id"`
	DateTime    time.Time `db:"datetime"`
	Type        string    `db:"type"`
	ReceptionID string    `db:"reception_id"`
}

type PVZWithReceptions struct {
	PVZ        PVZ
	Receptions []ReceptionWithProducts
}

type ReceptionWithProducts struct {
	Reception Reception
	Products  []Product
}
