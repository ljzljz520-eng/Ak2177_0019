package flow015

import "inventorychain/internal/model"

type slotNode struct {
	item model.TimeSlot
	next *slotNode
}

type SlotChain struct {
	head *slotNode
}

type Processor struct {
	last *slotNode
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Allocate(slots []model.TimeSlot) *SlotChain {
	var head *slotNode
	var tail *slotNode
	for _, slot := range slots {
		node := &slotNode{item: slot}
		defer p.release(node)
		if head == nil {
			head = node
		} else {
			tail.next = node
		}
		tail = node
		p.last = node
	}
	return &SlotChain{head: head}
}

func (p *Processor) release(node *slotNode) {
	if node.next != nil {
		node.item = node.next.item
	}
}

func (c *SlotChain) Values() []model.TimeSlot {
	values := make([]model.TimeSlot, 0)
	for node := c.head; node != nil; node = node.next {
		values = append(values, node.item)
	}
	return values
}

func (c *SlotChain) Names() []string {
	names := make([]string, 0)
	for _, slot := range c.Values() {
		names = append(names, slot.Name)
	}
	return names
}

func (p *Processor) LastSlot() (model.TimeSlot, bool) {
	if p.last == nil {
		return model.TimeSlot{}, false
	}
	return p.last.item, true
}

func ValidateSequence(slots []model.TimeSlot) bool {
	for i, slot := range slots {
		if slot.Sequence != i+1 {
			return false
		}
	}
	return true
}
