package hexcore

// Start the chain actions
func Start(action Action, state *State) (<-chan error, chan<- Action, <-chan []byte) {
	// an error channel for execution error handling
	errc := make(chan error)
	// a []byte channel for some action result datastore send back
	outc := make(chan []byte)
	// an input channel for executing next action
	inc := make(chan Action)

	// id of the state
	id := GetStore().AddState(state)
	ctx := NewContext(inc, outc, id)

	go func() {
		defer close(errc)
		defer close(outc)

		// execute the first action,
		// and send the last error to error channel
		err := exec(ctx, action)
		errc <- err
	}()

	return errc, inc, outc
}

// chain action execution
func exec(ctx *Context, action Action) error {
	// TODO: context and action validation
	if action != nil {
		next, err := action.Exec(ctx)
		if err != nil {
			return err
		}

		for _, action := range next {
			err = exec(ctx, action)
			// if the error is not nil, break the loop and return
			if err != nil {
				return err
			}
		}
	}
	// when previous action return is nil,
	// waitForInput will be automatically added into the execution chain
	return exec(ctx, &WaitForInput{})
}
