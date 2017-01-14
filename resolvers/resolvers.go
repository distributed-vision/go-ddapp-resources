package resolvers

type Selector interface {
	Key() interface{}
	Test(entity interface{})
}

type Resolver interface {
	Resolve(selector Selector) (chan interface{}, chan error)
}

/*
func NewSelector(values, select, opts) *Selector {
  if (arguments.length === 1) {
    if (typeof values === 'function') {
      this.selector = values
    } else {
      this.values = values
    }
  } else if (arguments.length === 2) {
    this.values = values
    if (typeof select === 'function')
      this.selector = select
    else {
      this.opts = select
    }
  } else {
    this.values = values
    this.selector = select
    this.opts = opts
  }
}
*/
