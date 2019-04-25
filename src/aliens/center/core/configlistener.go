/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/4/13
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

func NewConfigListener() *ConfigListener {
	return &ConfigListener{make(map[string][]func(data []byte))}
}

type ConfigListener struct {
	handlers map[string][]func(data []byte)
}

func (this *ConfigListener) Contains(configType string) bool {
	_, ok := this.handlers[configType]
	return ok
}

func (this *ConfigListener) ConfigChange(configType string, content []byte) {
	handlers := this.handlers[configType]
	if handlers != nil {
		for _, handler := range handlers {
			handler(content)
		}
	}
}

func (this *ConfigListener) AddHandler(configType string, handler func(data []byte)) {
	handlers := this.handlers[configType]
	if handlers == nil {
		handlers = []func(data []byte){handler}
	} else {
		handlers = append(handlers, handler)
	}
	this.handlers[configType] = handlers
}
