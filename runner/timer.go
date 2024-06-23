package runner

import (
	"github.com/rs/zerolog/log"
)

func timerRoutine(r Runner) {
	log.Debug().Msg("Timer routine started")
	for {
		select {
		case <-r.basicRunner().timer.C:
			log.Debug().Msg("triggering an upload via timer")
			_ = TriggerUpload(r)

		case <-r.basicRunner().resetTimerChan:
			log.Debug().Msg("Resetting timer")
			if !r.basicRunner().timer.Stop() && len(r.basicRunner().timer.C) > 0 {
				<-r.basicRunner().timer.C
			}
			r.basicRunner().timer.Reset(r.basicRunner().sampleRate)
		case <-r.basicRunner().ctx.Done():
			r.basicRunner().timer.Stop()
			log.Debug().Msg("Stopping timer...")
			return
		}
	}
}
