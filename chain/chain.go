package chain

// based on the article here: https://go101.org/article/channel-closing.html
// only sender should close the channel, never the receiver

// chain action channel context
// because the action is executed one by one,
// so mutex lock is not necessary
type context struct {
	// output channel, it will be closed automatically,
	// when execution returned, so don't closed it in action
	outc chan<- []byte
}

func exec(ctx *context) error {
	return nil
}

type action interface {
	run(ctx *context)
}

// Exec the chain actions with an input channel
// return an error channel for execution error handling
// and an []byte channel for some action result data send out
func Exec() (<-chan error, <-chan []byte) {
	errc := make(chan error)
	outc := make(chan []byte)

	ctx := &context{
		outc: outc,
	}

	go func() {
		defer close(errc)
		defer close(outc)

		// execute the first action,
		// and send the last error to error channel
		err := exec(ctx)
		errc <- err
	}()

	return errc, outc
}
