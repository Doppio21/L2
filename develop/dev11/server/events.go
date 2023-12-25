package server

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/huandu/skiplist"
)

type nodeKey struct {
	day    int
	month  int
	year   int
	userID int
}

type nodeValue struct {
	events map[string]string
}

type Events struct {
	events *skiplist.SkipList
	mu     sync.RWMutex
}

type Event struct {
	Date        string
	Name        string
	Description string
}

func NewEvents() *Events {
	list := skiplist.New(skiplist.GreaterThanFunc(func(k1, k2 interface{}) int {
		e1 := k1.(*nodeKey)
		e2 := k2.(*nodeKey)

		switch {
		case e1.userID != e2.userID:
			return e1.userID - e2.userID
		case e1.year != e2.year:
			return e1.year - e2.year
		case e1.month != e2.month:
			return e1.month - e2.month
		case e1.day != e2.day:
			return e1.day - e2.day
		}

		return 0
	}))
	return &Events{
		events: list,
	}
}

func splitDate(date string) (int, int, int, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, 0, 0, err
	}

	return t.Year(), int(t.Month()), t.Day(), nil
}

func (e *Events) getByKey(key *nodeKey) (*nodeValue, error) {
	nv, ok := e.events.GetValue(key)
	if !ok {
		return nil, errors.New("not exists")
	}

	val, ok := nv.(*nodeValue)
	if !ok {
		panic("invalid node value")
	}

	return val, nil
}

func (e *Events) Create(userID int, date string, name, description string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	year, month, day, err := splitDate(date)
	if err != nil {
		return err
	}

	key := &nodeKey{
		day:    day,
		month:  month,
		year:   year,
		userID: userID,
	}

	nv, ok := e.events.GetValue(key)
	if !ok {
		e.events.Set(key, &nodeValue{events: map[string]string{
			name: description,
		}})
		return nil
	}

	val, ok := nv.(*nodeValue)
	if !ok {
		panic("invalid node value")
	}

	if _, ok := val.events[name]; ok {
		return errors.New("already exists")
	}

	val.events[name] = description
	return nil
}

func (e *Events) Update(userid int, date string, name, description string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	year, month, day, err := splitDate(date)
	if err != nil {
		return err
	}

	key := &nodeKey{
		day:    day,
		month:  month,
		year:   year,
		userID: userid,
	}

	val, err := e.getByKey(key)
	if err != nil {
		return err
	}

	if _, ok := val.events[name]; !ok {
		return errors.New("not exists")
	}

	val.events[name] = description
	return nil
}

func (e *Events) Delete(userid int, date string, name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	year, month, day, err := splitDate(date)
	if err != nil {
		return err
	}

	key := &nodeKey{
		day:    day,
		month:  month,
		year:   year,
		userID: userid,
	}

	val, err := e.getByKey(key)
	if err != nil {
		return err
	}

	if _, ok := val.events[name]; !ok {
		return errors.New("not exists")
	}

	delete(val.events, name)
	return nil
}

func (e *Events) Get(userid int, date string, t string) ([]Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var ret []Event

	curDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	key := &nodeKey{
		day:    curDate.Day(),
		month:  int(curDate.Month()),
		year:   curDate.Year(),
		userID: userid,
	}

	thresholdDate := curDate

	switch t {
	case "day":
		thresholdDate = thresholdDate.AddDate(0, 0, 0)
	case "week":
		thresholdDate = thresholdDate.AddDate(0, 0, 7)
	case "month":
		thresholdDate = thresholdDate.AddDate(0, 1, 0)
	default:
		return nil, errors.New("invalid type")
	}

	nv := e.events.Get(key)
	if nv == nil {
		return nil, nil
	}

	for curDate.Before(thresholdDate) || curDate.Equal(thresholdDate) {
		date = strDate(curDate)
		val, ok := nv.Value.(*nodeValue)
		if !ok {
			panic("invalid node value")
		}

		for name, desc := range val.events {
			ret = append(ret, Event{
				Date:        date,
				Name:        name,
				Description: desc,
			})
		}

		if nv = nv.Next(); nv == nil {
			break
		}

		key, ok := nv.Key().(*nodeKey)
		if !ok {
			panic("invalid node key")
		}

		if key.userID != userid {
			break
		}

		curDate = time.Date(key.year, time.Month(key.month), key.day,
			0, 0, 0, 0, time.Local)
	}

	return ret, nil
}

func strDate(date time.Time) string {
	d := date.String()
	ret := strings.Split(d, " ")
	return ret[0]
}
