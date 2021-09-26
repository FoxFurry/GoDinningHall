package supervisor

import (
	"github.com/foxfurry/go_dining_hall/internal/service/table"
	"github.com/foxfurry/go_dining_hall/internal/service/waiter"
	"sync"
)

type ISupervisor interface {
	GetTables() []table.Table
	GetTablesPointer() []*table.Table
	GenerateTables(num int, menu int)
	InitializeTables()

	GenerateWaiter(num int)
	StartWaiters()
	FreeTable(idx int)
}

type DiningSupervisor struct {
	tablesMutex sync.Mutex
	tables      []table.Table

	waitersMutex sync.Mutex
	waiters      []waiter.Waiter
}

func (s *DiningSupervisor) StartWaiters() {
	for idx, _ := range s.waiters {
		go s.waiters[idx].WatchTables()
	}
}

func (s *DiningSupervisor) GenerateWaiter(num int) {
	tables := s.GetTablesPointer()
	s.waiters = nil
	for idx := 0; idx < num; idx++ {
		s.waiters = append(s.waiters, waiter.NewWaiter(idx, tables))
	}
}

func (s *DiningSupervisor) GenerateTables(num int, menu int) {
	s.tablesMutex.Lock()
	s.tables = nil
	for idx := 0; idx < num; idx++ {
		s.tables = append(s.tables, table.NewTable(menu, idx))
	}
	s.tablesMutex.Unlock()
}

func (s *DiningSupervisor) GetTables() []table.Table {
	var tmp []table.Table
	s.tablesMutex.Lock()
	tmp = s.tables
	s.tablesMutex.Unlock()
	return tmp
}

func (s *DiningSupervisor) InitializeTables() {
	s.tablesMutex.Lock()
	for idx, _ := range s.tables {
		go s.tables[idx].Simulate()
	}
	s.tablesMutex.Unlock()
}

func (s *DiningSupervisor) GetTablesPointer() []*table.Table {
	var tmp []*table.Table
	s.tablesMutex.Lock()
	for idx, _ := range s.tables {
		tmp = append(tmp, &s.tables[idx])
	}
	s.tablesMutex.Unlock()
	return tmp
}

func (s *DiningSupervisor) FreeTable(idx int) {
	s.tablesMutex.Lock()
	s.tables[idx].SetFree()
	s.tablesMutex.Unlock()
}
