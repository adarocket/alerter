package sqllite

import (
	"log"
)

const fillAlertsTable = `
	INSERT OR IGNORE INTO alerts 
	(id, name, checked_field, type_checker) 
	VALUES (0, 'epoch check', 'statistic.epoch.epoch_number', 'Interval'),
	       (1, 'cpuState check', 'CpuState', 'interval'),
	       (2, 'blocks check', 'Blocks', 'more')
`

const fillAlertsNodeTable = `
	INSERT OR IGNORE INTO alert_node
	(alert_id, normal_from, normal_to, critical_from, critical_to, frequncy, node_uuid)
	VALUES (0, 0.0, 10.0, 10.0, 20.0, 'Max', '08e792fd-2a19-466f-9a2a-d9fd40bdf9d1'),
	       (1, 5.0, 16.0, 16.0, 20.0, 'Max', '08e792fd-2a19-466f-9a2a-d9fd40bdf9d1'),
	       (2, 0.0, 10.0, 10.0, 20.0, 'Max', '08e792fd-2a19-466f-9a2a-d9fd40bdf9d1')
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
