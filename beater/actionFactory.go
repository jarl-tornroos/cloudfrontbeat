package beater

import (
	"errors"
	"fmt"
)

// actionFactory contains a map[string] of all available actions
// index value is the returned string from the actions Name() function
type actionFactory struct {
	actions map[string]Action
}

// Constructor function for the action factory
func NewActionFactory() *actionFactory {
	factory := &actionFactory{}
	factory.actions = make(map[string]Action)
	return factory
}

// Add Action
func (factory *actionFactory) AddAction(action Action) *actionFactory {
	factory.actions[action.Name()] = action
	return factory
}

// Return the specified action
func (factory *actionFactory) GetAction(name *string) (Action, error) {

	var err error

	action, registered := factory.actions[*name]
	if !registered {
		errorMessage := fmt.Sprintf("Action, %s has not been implemented. Available Actions:", *name)
		for _, val := range factory.actions {
			errorMessage += fmt.Sprintf("\n %s", val.Name())
		}

		err = errors.New(errorMessage)
	}

	return action, err
}
