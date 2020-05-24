package execution

type Programs []*Program

func (p *Programs) Add(pgrm *Program) {
	*p = append(*p, pgrm)
}

func (p *Programs) GetPrograms() map[string]*Programs {
	currentPrograms := make(map[string]*Programs)
	for _, prgm := range *p {
		hour := prgm.getHour()
		currentPrgms := currentPrograms[hour]
		if currentPrgms == nil {
			var new Programs
			new.Add(prgm)
			currentPrograms[hour] = &new
		} else {
			current := *currentPrgms
			for i, pr := range current {
				if pr.hasSameZones(prgm) {
					best := pr.getBestProgram(prgm)
					if best != pr {
						current[i] = best
					}
				} else {
					currentPrgms.Add(prgm)
				}
			}
			currentPrograms[hour] = currentPrgms
		}
	}
	return currentPrograms
}
