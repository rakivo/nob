package nob

var (
	cs children
)

func (c *Child) Wait() error {
	return c.Cmd.Wait()
}

func (c *Child) MustWait() {
	if err := c.Wait(); err != nil {
		panic(err)
	}
}

func WaitAll() error {
	for c := cs.pop(); c != nil; c = cs.pop() {
		if err := c.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func MustWaitAll() {
	if err := WaitAll(); err != nil {
		panic(err)
	}
}
