package database

import "log"

const createTableAlertNodeCardano = `
	CREATE TABLE if not exists "alert_node" (
	"alert_id"		INTEGER NOT NULL,
	"normal_from"	REAL NOT NULL,
	"normal_to"		REAL NOT NULL,
	"critical_from"	REAL NOT NULL,
	"critical_to"	REAL NOT NULL,
	"frequncy"		TEXT NOT NULL,
	"node_uuid"		TEXT NOT NULL,
	FOREIGN KEY("alert_id") REFERENCES "alerts"("id"))
`

const createTableAlertsCardano = `
	CREATE TABLE if not exists "alerts" (
	"id"			INTEGER NOT NULL UNIQUE,
	"name"			TEXT NOT NULL,
	"checked_field"	TEXT NOT NULL,
	"type_checker"	TEXT NOT NULL,
	PRIMARY KEY("id"))
`

func CreateTables() error {
	if _, err := dbConn.Exec(createTableAlertNodeCardano); err != nil {
		log.Println("CreateNodeAuthTable", err)
		return err
	}

	if _, err := dbConn.Exec(createTableAlertsCardano); err != nil {
		log.Println("CreateNodeAuthTable", err)
		return err
	}

	return nil
}
