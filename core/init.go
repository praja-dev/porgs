package core

func (p *Plugin) GetInit() error {
	loadData()
	return nil
}
