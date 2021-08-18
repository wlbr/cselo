package postgresql

import (
	"context"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

func (s *Postgres) GetFlashStatsForPlayer(p *elo.Player, iv *elo.Intervall) *elo.Player {
	var eflash, tflash, sflash int
	fp := p
	// if p.SteamID == "" {
	// 	fp = s.GetPlayerByID(p.ID)
	// 	p.Name = fp.Name
	// 	p.SteamID = fp.SteamID
	// }
	if fp != nil {
		row := s.Db.QueryRow(context.Background(), "select count(case when victimtype='enemy' then 1 end) as enemyflashes, "+
			"count(case when victimtype='teammate' then 1 end) as teammateflash, "+
			"count(case when victimtype='self' then 1 end) as selfflash "+
			"from blindings "+
			//"WHERE timestmp > current_date - make_interval(days => $1) AND actor = $2 ;", 700, fp.ID)
			"WHERE timestmp > $1 and timestmp < $2 AND actor = $3 ;", iv.Start, iv.End, fp.ID)
		err := row.Scan(&eflash, &tflash, &sflash)
		if err != nil {
			log.Error("Cannot get player flash from PostgresQL database: %v", err)
		}
		log.Warn("Flashs for player %v: enemyflash: %v  -  teammateflash: %v  -  selfflash: %v", fp, eflash, tflash, sflash)

	} else {
		log.Error("ha? %v", fp)
	}
	return p
}
