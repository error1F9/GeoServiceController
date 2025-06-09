package modules

import (
	userRepository "GeoService/internal/modules/user/repository"
	"github.com/uptrace/bun"
)

type Repositories struct {
	UserRepository userRepository.Userer
}

func NewRepositories(db *bun.DB) *Repositories {
	return &Repositories{
		UserRepository: userRepository.NewRepository(db),
	}
}
