package callback

import (
	"git.multiverse.io/eventkit/kit/log"
	"sync"
)

// HookHandleFunc is a controller that used to hook handle function registration.
type HookHandleFunc struct {
	Name   string
	Func   func() error
	IsSync bool
}

var (
	// initHookHandleFuncMap defines callback function type variables, sed-sdk callback after receiving the message
	initHookHandleFuncMap = make(map[string]*HookHandleFunc)
	locker                sync.RWMutex
)

// RegisterInitHookFunc provides an application registration start hook function
func RegisterInitHookFunc(name string, handleFunc func() error, isSync bool) {
	locker.Lock()
	defer locker.Unlock()
	log.Infosf("Register init hook function, function name[%s], handleFunc[%++v], isSync[%v]", name, handleFunc, isSync)
	initHookHandleFuncMap[name] = &HookHandleFunc{
		Name:   name,
		Func:   handleFunc,
		IsSync: isSync,
	}
}

// RunHookHandleFunc is start all the hooked functions when service start.
func RunHookHandleFunc() error {
	locker.RLock()
	defer locker.RUnlock()
	if len(initHookHandleFuncMap) > 0 {
		var wg sync.WaitGroup

		log.Infosf("Start run hook ...")
		for k, v := range initHookHandleFuncMap {
			log.Infosf("Start running hook[%s], hook info:[%++v]", k, v)
			if v.IsSync {
				// sync call
				if err := v.Func(); nil != err {
					log.Errorsf("Running sync hook function[%s] failed, error=%++v", k, err)
					return err
				}
				log.Infosf("Run hook[%s] end!", k)
			} else {
				wg.Add(1)
				// async call
				go func(hookName string, hookHandleFunc *HookHandleFunc) {
					defer func() {
						log.Infosf("Service config, run hook[%s] end!", hookName)
						wg.Done()
					}()
					if err := hookHandleFunc.Func(); nil != err {
						log.Errorsf("Running async hook function[%s] failed, error=%++v", hookName, err)
					}
				}(k, v)
			}
		}
		wg.Wait()
		log.Infosf("Running hook function successfully!")
	} else {
		log.Infosf("FSM hasn't registered any function, skip run hook...")
	}

	return nil
}
