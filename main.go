package main

import (
	"github.com/foxfurry/go_dining_hall/internal/domain/table"
	"math/rand"
	"time"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

func main(){
	newTable := table.NewTable(5)

	newTable.StartGenerator()
}
