package ege

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type Notificator struct {
	grabber       *Grabber
	isStarted     bool
	ticker        *time.Ticker
	stopChannel   chan bool
	msgFunc       func(s string)
	lastError     error
	errorsCounter int
}

func NewNotificator(delay int, grabber *Grabber, msgFunc func(s string)) *Notificator {
	n := &Notificator{
		ticker:  time.NewTicker(time.Duration(delay) * time.Second),
		msgFunc: msgFunc,
		grabber: grabber,
	}
	n.stopChannel = make(chan bool)

	return n
}

func (n *Notificator) Start() error {
	if n.isStarted {
		return NewError(NotificatorIsAlreadyOn, 0, nil)
	}
	n.msgFunc("Запуск менеджера уведомлений")
	if err := n.notifyIfDifferent(); err != nil {
		return err
	}
	n.errorsCounter = 0
	n.isStarted = true
	n.msgFunc("Менеджер уведомлений запущен")
	go func() {
		for {
			select {
			case <-n.stopChannel:
				n.isStarted = false
				return
			case <-n.ticker.C:
				if err := n.notifyIfDifferent(); err != nil {
					log.Print(err)
					n.lastError = err
					n.errorsCounter++
					if n.errorsCounter >= 10 {
						n.msgFunc("Отключение из-за лимита ошибок")
						n.ticker.Stop()
						n.isStarted = false
						return
					}
					continue
				}
				if n.errorsCounter != 0 {
					n.errorsCounter = 0
				}
			}
		}
	}()
	return nil
}

func (n *Notificator) notifyIfDifferent() error {
	result, err := n.grabber.GetIfDifferent()
	if err != nil {
		return err
	}
	if result != nil {
		n.msgFunc("@all, получены новые данные: \n" + result.String())
	}
	return nil
}

func (n *Notificator) Stop() error {
	if !n.isStarted {
		return NewError(NotificatorIsAlreadyOff, 0, nil)
	}
	n.msgFunc("Менеджер уведомлений остановлен")
	n.ticker.Stop()
	n.stopChannel <- true
	n.isStarted = false
	return nil
}

type Stats struct {
	isStarted     bool
	lastError     error
	errorsCounter int
}

func (s *Stats) String() string {
	errStr := "нет"
	if s.lastError != nil {
		errStr = s.lastError.Error()
	}

	return fmt.Sprintf("Запущен? %t\nПоследняя ошибка: %s\nСчетчик ошибок: %s\n",
		s.isStarted,
		errStr,
		strconv.Itoa(s.errorsCounter))
}

func (n *Notificator) GetStats() *Stats {
	return &Stats{
		isStarted:     n.isStarted,
		lastError:     n.lastError,
		errorsCounter: n.errorsCounter,
	}
}
