package elevator

import (
	"github.com/rs/zerolog/log"
)

func timerRoutine(r Elevator) {
	log.Debug().Msg("Timer routine started")
	for {
		select {
		case <-r.baseElevator().timer.C:
			log.Debug().Msg("triggering an upload via timer")
			_ = TriggerUpload(r)

		case <-r.baseElevator().resetTimerChan:
			log.Debug().Msg("Resetting timer")
			if !r.baseElevator().timer.Stop() && len(r.baseElevator().timer.C) > 0 {
				<-r.baseElevator().timer.C
			}
			r.baseElevator().timer.Reset(r.baseElevator().sampleRate)
		case <-r.baseElevator().ctx.Done():
			r.baseElevator().timer.Stop()
			log.Debug().Msg("Stopping timer...")
			return
		}
	}
}
