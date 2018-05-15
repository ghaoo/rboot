package Rboot

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultRobotName = `Rboot`
)

type Robot struct {
	name     string
	alias    string
	provider map[string]Provider
	store    Store

	handlers   []handler
	signalChan chan os.Signal
}

func New() (*Robot, error) {
	robot := &Robot{
		name:       DefaultRobotName,
		signalChan: make(chan os.Signal, 1),
	}

	err := robot.registerProvider()

	if err != nil {
		return nil, err
	}

	default_store, err := NewStore(robot)

	if err != nil {
		return nil, err
	}
	robot.store = default_store

	return robot, nil
}

func (robot *Robot) SetAlias(alias string) {
	robot.alias = alias
}

func (robot *Robot) SetName(name string) {
	robot.name = name
}

func (robot *Robot) GetName() string {
	return robot.name
}

func (robot *Robot) SetStore(store Store) {
	robot.store = store
}

func (robot *Robot) Handlers() []handler {
	return robot.handlers
}

func (robot *Robot) Receive(msg *Message) error {

	log.Printf(`%v`,robot)

	for _, handler := range robot.handlers {
		response := NewResponse(robot, msg)

		if err := handler.Handle(response); err != nil {
			return err
		}
	}
	return nil
}

func (robot *Robot) Handle(handlers ...handler) {
	for _, h := range handlers {

		robot.handlers = append(robot.handlers, h)
	}
}

func (robot *Robot) Run() {
	log.Printf("starting robot")

	log.Printf("opening %s store connection", robot.store.Name())
	go func() {
		robot.store.Open()
	}()

	for name, prov := range robot.provider {
		log.Printf("starting %s provider", name)
		go prov.Run()
	}

	signal.Notify(robot.signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	stop := false
	for !stop {
		select {
		case sig := <-robot.signalChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				stop = true
			}
		}
	}

	signal.Stop(robot.signalChan)

	robot.Stop()
}

func (robot *Robot) Stop() error {

	for name, prov := range robot.provider {
		log.Printf("stopping %s provider", name)

		if err := prov.Close(); err != nil {
			return err
		}
	}

	log.Printf("closing %s store connection", robot.store.Name())
	if err := robot.store.Close(); err != nil {
		return err
	}

	log.Printf("stopping robot")
	return nil
}
