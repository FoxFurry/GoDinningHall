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

	panicOnTimeout bool = true
)

type Table struct {
	id int

	menuMutex sync.Mutex
	menu      int

	currentStateMutex sync.Mutex
	currentState      State

	currentOrderMutex sync.Mutex
	currentOrder      *dto.Order
	currentOrderWait int
}

func NewTable(newMenu int, newID int) Table {
	return Table{
		id:           newID,
		menu:         newMenu,
		currentState: NotReady,
		currentOrder: &dto.Order{},
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
		food := rand.Intn(menu)
		t.pushFood(food)
	}

	t.currentOrderMutex.Lock()
	t.currentOrder.TableID = t.id
	t.currentOrder.OrderID = rand.Intn(1000)
	t.currentOrder.MaxWait = 100					// TODO: Add full menu and max wait calculation
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
			logger.LogTable(t.id, "Waiting to be picked up")

		case Waiting:
			t.currentOrderWait++

			if t.currentOrderWait >= t.GetCurrentOrder().MaxWait {
				logger.LogTable(t.id, "Time is out!")

				if panicOnTimeout {
					logger.LogPanic("Wait time is out!")
				}

			}
		}
		time.Sleep(time.Second)
	}
}

func (t *Table) PickUp() dto.Order {
	t.setState(Waiting)
	t.currentOrderWait = 0
	return *t.GetCurrentOrder()
}

func (t *Table) SetFree() {
	t.setState(NotReady)
}
