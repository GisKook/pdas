package db

import (
	"fmt"
	//"github.com/golang/protobuf/proto"
	"github.com/giskook/pdas/base"
	"log"
)

const (
	SQL_INSERT_TABLE string = "INSERT INTO t_indoor_intensity (imei, bt_id, rssi, orientation, loc, finger_id) values (%s, %d,%d,%s, %s, %s);"
)

func (db_socket *DbSocket) ProccessTransaction() {
	for {
		select {
		case transactions := <-db_socket.TransactionChan:
			tx, err := db_socket.Db.Begin()
			if err != nil {
				log.Println("ProccessTransaction")
				log.Println(err)
			}

			for trans := range transactions {
				insert_sql := fmt.Sprintf(SQL_INSERT_TABLE, trans.TagMac, base.GetMac(trans.RingMac), trans.Rssi, trans.Orientation, base.GetLoc(trans.X, trans.Y), trans.FingerID)
				log.Println(insert_sql)

				tx.Exec(insert_sql)
			}
			err = tx.Commit()
			if err != nil {
				log.Println("ProccessTransaction commit error")
				log.Println(err)
			}
			close(transactions)
		}
	}
}
