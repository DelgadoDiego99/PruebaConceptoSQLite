package models

import (
	"database/sql"
	"strconv"

	_ "modernc.org/sqlite"
)

type Estudiante struct {
	Codigo   int    `json:"Codigo"`
	Nombre   string `json:"Nombre"`
	Apellido string `json:"Apellido"`
	Ingreso  string `json:"Fecha_Ingreso"`
}

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite", "./universidad.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func GetEstudiantes(count int) ([]Estudiante, error) {

	rows, err := DB.Query("SELECT Codigo, Nombre, Apellido, Ingreso from Estudiante LIMIT " + strconv.Itoa(count))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	people := make([]Estudiante, 0)
	for rows.Next() {
		singleEstudiante := Estudiante{}
		err = rows.Scan(&singleEstudiante.Codigo, &singleEstudiante.Nombre, &singleEstudiante.Apellido, &singleEstudiante.Ingreso)

		if err != nil {
			return nil, err
		}

		people = append(people, singleEstudiante)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return people, err
}

func GetEstudianteById(id string) (Estudiante, error) {

	stmt, err := DB.Prepare("SELECT Codigo, Nombre, Apellido, Ingreso from Estudiante WHERE Codigo = ?")
	if err != nil {
		return Estudiante{}, err
	}

	person := Estudiante{}
	sqlErr := stmt.QueryRow(id).Scan(&person.Codigo, &person.Nombre, &person.Apellido, &person.Ingreso)
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Estudiante{}, nil
		}
		return Estudiante{}, sqlErr
	}
	return person, nil
}

func AddEstudiante(newEstudiante Estudiante) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO Estudiante (Codigo, Nombre, Apellido) VALUES (?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newEstudiante.Codigo, newEstudiante.Nombre, newEstudiante.Apellido)
	if err != nil {
		return false, err
	}
	tx.Commit()

	return true, nil
}

func UpdateEstudiante(ourEstudiante Estudiante, Codigo int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE Estudiante SET Nombre = ?, Apellido = ? WHERE Codigo = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ourEstudiante.Nombre, ourEstudiante.Apellido, Codigo)

	if err != nil {
		return false, err
	}
	tx.Commit()

	return true, nil
}
