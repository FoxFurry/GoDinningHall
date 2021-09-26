package waiter

import (
	"bytes"
	"encoding/json"
	"github.com/foxfurry/go_dining_hall/internal/service/table"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

type Waiter struct {
	id int
	tables []*table.Table
}

func NewWaiter(newID int, newTables[]*table.Table) Waiter {
	return Waiter{
		id: newID,
		tables: newTables,
	}
}

func (w *Waiter) WatchTables() {
	for {
		for idx, _ := range w.tables {
			if w.tables[idx].GetState() == table.Ready {
				order := w.tables[idx].PickUp()

				jsonBody, err := json.Marshal(order)
				if err != nil {
					log.Panic(err)
				}
				contentType := "application/json"

				log.Printf("Waiter %v picked up an order from table %v", w.id, order.TableID)

				http.Post(viper.GetString("kitchen_host") + "/order", contentType, bytes.NewReader(jsonBody))

				time.Sleep(time.Second)
			}
		}
		time.Sleep(time.Second)
	}
}