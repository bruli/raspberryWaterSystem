package zone

import (
	"fmt"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
)

type Remover struct {
	repo Repository
	log  logger.Logger
}

func NewRemover(repo Repository, log logger.Logger) *Remover {
	return &Remover{repo: repo, log: log}
}

func (r *Remover) Remove(zoneID string) error {
	zon := r.repo.Find(zoneID)
	if zon == nil {
		err := NewNotFound(zoneID)
		r.log.Fatal(err.err)
		return err
	}
	zons := r.repo.GetZones()
	zons.remove(zoneID)
	err := r.repo.Save(*zons)
	if err != nil {
		return fmt.Errorf("failed to remove zoneID %s: %s", zoneID, err)
	}
	return nil
}

type NotFound struct {
	err string
}

func NewNotFound(zoneID string) *NotFound {
	return &NotFound{fmt.Sprintf("zone ID '%s' does not exits", zoneID)}
}

func (n NotFound) Error() string {
	return n.err
}
