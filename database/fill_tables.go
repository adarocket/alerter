package database

import (
	"log"
)

// FIXME: где интерфейс?

const fillAlertsTable = `
	INSERT OR IGNORE INTO alerts 
	(id, name, checked_field, type_checker) 
	VALUES (0, 'sizeCache check', 'SizeCache', 'change_up'),
	       (1, 'cpuState check', 'CpuState', 'interval'),
	       (2, 'blocks check', 'Blocks', 'more')
`

const fillAlertsNodeTable = `
	INSERT OR IGNORE INTO alert_node
	(alert_id, normal_from, normal_to, critical_from, critical_to, frequncy)
	VALUES (0, 0.0, 10.0, 10.0, 20.0, 'normal'),
	       (1, 5.0, 16.0, 16.0, 20.0, 'normal'),
	       (2, 0.0, 10.0, 10.0, 20.0, 'normal')
`

func (p sqlite) FillTables() error {
	_, err := p.dbConn.Exec(fillAlertsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = p.dbConn.Exec(fillAlertsNodeTable)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
