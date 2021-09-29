package table

import (
	"github.com/foxfurry/go_dining_hall/internal/dto"
	"github.com/foxfurry/go_dining_hall/internal/infrastructure/logger"
	"github.com/foxfurry/go_dining_hall/internal/infrastructure/table_helper"
	"math/rand"
	"sync"
	"time"
)

type State int

const (
	NotReady State = iota
	Ready
	Waiting
)

const (
	orderProbability     = 0.15
	maxFoodCount     int = 6
)

type Table struct {
	id int

	menuMutex sync.Mutex
	menu      int

	currentStateMutex sync.Mutex
	currentState      State

	currentOrderMutex sync.Mutex
	currentOrder      *dto.Order

	pickupTimeMutex sync.Mutex
	pickupTime      *time.Time
}

func NewTable(newMenu int, newID int) Table {
	return Table{
		id:           newID,
		menu:         newMenu,
		currentState: NotReady,
		currentOrder: &dto.Order{},
		pickupTime:   nil,
	}
}

func (t *Table) pushFood(food int) {
	t.currentOrderMutex.Lock()
	t.currentOrder.Items = append(t.currentOrder.Items, food)
	t.currentOrderMutex.Unlock()
}

func (t *Table) SetMenu(newMenu int) {
	t.menuMutex.Lock()
	t.menu = newMenu
	t.menuMutex.Unlock()
}

func (t *Table) setState(newState State) {
	t.currentStateMutex.Lock()
	t.currentState = newState
	t.currentStateMutex.Unlock()
}

func (t *Table) getMenu() int {
	var tmp int
	t.menuMutex.Lock()
	tmp = t.menu
	t.menuMutex.Unlock()
	return tmp
}

func (t *Table) GetState() State {
	var tmp State
	t.currentStateMutex.Lock()
	tmp = t.currentState
	t.currentStateMutex.Unlock()
	return tmp
}

func (t *Table) getPickupTime() *time.Time {
	var tmp *time.Time
	t.pickupTimeMutex.Lock()
	tmp = t.pickupTime
	t.pickupTimeMutex.Unlock()
	return tmp
}

func (t *Table) setPickupTime(newTime time.Time) {
	t.pickupTimeMutex.Lock()
	t.pickupTime = &newTime
	t.pickupTimeMutex.Unlock()
}

func (t *Table) GetCurrentOrder() *dto.Order {
	var tmp *dto.Order
	t.currentOrderMutex.Lock()
	tmp = t.currentOrder
	t.currentOrderMutex.Unlock()
	return tmp
}

func (t *Table) GenerateOrder() {
	var menu = t.getMenu()
	var count = rand.Intn(maxFoodCount)+1

	for idx := 0; idx < count; idx++ {
		t.pushFood(rand.Intn(menu))
	}

	t.currentOrderMutex.Lock()
	t.currentOrder.TableID = t.id
	t.currentOrder.OrderID = rand.Intn(1000)
	t.currentOrderMutex.Unlock()

	t.setState(Ready)
}

func (t *Table) Simulate() {
	for {
		switch t.GetState() {
		case NotReady:
			if table_helper.CoinFlip(orderProbability) {
				t.GenerateOrder()
				logger.LogTableF(t.id, "Order generated: %v", t.GetCurrentOrder().Items)
			} else {
				//log.Printf("Table %v: Order not generated!", t.id)
			}
		case Ready:
			logger.LogTableF(t.id, "Waiting to be picked up")
		}

		time.Sleep(time.Second)
	}
}

func (t *Table) PickUp() dto.Order {
	t.setState(Waiting)
	t.setPickupTime(time.Now())

	return *t.GetCurrentOrder()
}

func (t *Table) SetFree() {
	t.setState(NotReady)
}
