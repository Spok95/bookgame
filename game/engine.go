package game

func (p *Player) AddItem(item string) bool {
	if len(p.Inventory) >= 5 {
		return false
	}
	for _, i := range p.Inventory {
		if i == item {
			return false
		}
	}
	p.Inventory = append(p.Inventory, item)
	return true
}
