package priorityselect

import "reflect"

type Selector struct {
	cases   []reflect.SelectCase
	buffers []*box
}

type box struct {
	v interface{}
}

type ChanClosedErr struct {
	Chan interface{}
}

func (self ChanClosedErr) Error() string {
	return "chan closed"
}

func New(chans ...interface{}) *Selector {
	selector := new(Selector)
	selector.cases = append(selector.cases, reflect.SelectCase{
		Dir: reflect.SelectDefault,
	})
	for _, c := range chans {
		selector.cases = append(selector.cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		})
	}
	selector.buffers = make([]*box, len(chans))
	return selector
}

func (self *Selector) Select() (interface{}, error) {
	nCases := len(self.cases)
sel:
	i := 1
	for i < nCases && self.buffers[i-1] == nil {
		i++
	}
	if i == nCases { // no buffered value
		n, v, ok := reflect.Select(self.cases[1:i]) // no default case
		if !ok {
			return nil, ChanClosedErr{Chan: self.cases[n+1].Chan}
		}
		self.buffers[n] = &box{v.Interface()}
		goto sel
	} else { // has buffered value
		n, v, ok := reflect.Select(self.cases[:i]) // default case at index 0
		if !ok && n > 0 {
			return nil, ChanClosedErr{Chan: self.cases[n].Chan}
		}
		if n > 0 { // higher priority chan received
			self.buffers[n-1] = &box{v.Interface()}
			goto sel
		}
		// default
		for i, buf := range self.buffers {
			if buf != nil {
				self.buffers[i] = nil
				return buf.v, nil
			}
		}
	}
	panic("impossible")
	return nil, nil
}
