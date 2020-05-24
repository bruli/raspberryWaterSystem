package weather

type Getter struct {
	repo Repository
}

func NewGetter(repo Repository) *Getter {
	return &Getter{repo: repo}
}

func (g *Getter) Get() (temp, hum float32, err error) {
	temp, hum, err = reader(g.repo)
	if err != nil {
		return 0, 0, err
	}

	return temp, hum, nil
}
