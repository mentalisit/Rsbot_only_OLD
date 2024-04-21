package storage

import (
	"Rsbot_only/internal/config"
	"Rsbot_only/internal/models"
	"Rsbot_only/internal/storage/mongo"
	"Rsbot_only/internal/storage/postgres"
	"Rsbot_only/internal/storage/words"
	"fmt"
	"github.com/mentalisit/logger"
	"go.uber.org/zap"
)

type Storage struct {
	log               *zap.Logger
	debug             bool
	ConfigRs          ConfigRs
	TimeDeleteMessage TimeDeleteMessage
	Temp              *mongo.DB
	Words             *words.Words
	Subscribe         Subscribe
	Emoji             Emoji
	Count             Count
	Top               Top
	Update            Update
	Timers            Timers
	DbFunc            DbFunc
	Event             Event
	CorpConfigRS      map[string]models.CorporationConfig
}

func NewStorage(log *logger.Logger, cfg *config.ConfigBot) *Storage {

	//инициализируем и читаем репозиторий из облока конфига конфигурации
	mongoDB := mongo.InitMongoDB(log)

	//подключаю языковой пакет
	w := words.NewWords()

	//инициализируем локальный репозиторий
	local := postgres.NewDb(log, cfg)

	s := &Storage{
		TimeDeleteMessage: mongoDB,
		ConfigRs:          mongoDB,
		Temp:              mongoDB,
		Words:             w,
		Subscribe:         local,
		Emoji:             local,
		Count:             local,
		Top:               local,
		Update:            local,
		Timers:            local,
		DbFunc:            local,
		Event:             local,
		CorpConfigRS:      make(map[string]models.CorporationConfig),
	}

	go s.loadDbArray()
	return s
}
func (s *Storage) loadDbArray() {
	var c = 0
	var rslist string
	rs := s.ConfigRs.ReadConfigRs()
	for _, r := range rs {
		s.CorpConfigRS[r.CorpName] = r
		c++
		rslist = rslist + fmt.Sprintf("%s, ", r.CorpName)
	}
	fmt.Printf("Загружено конфиг RsBot %d : %s\n", c, rslist)
}
func (s *Storage) ReloadDbArray() {
	CorpConfigRS := make(map[string]models.CorporationConfig)

	s.CorpConfigRS = CorpConfigRS

	rs := s.ConfigRs.ReadConfigRs()
	for _, r := range rs {
		s.CorpConfigRS[r.CorpName] = r
	}
}
