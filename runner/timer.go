package runner

import (
	"github.com/rs/zerolog/log"
)

func timerRoutine(r Runner) {
	log.Debug().Msg("Timer routine started")
	for {
		select {
		case <-r.runnerBase().timer.C:
			log.Debug().Msg("triggering an upload via timer")
			_ = TriggerUpload(r)

		case <-r.runnerBase().resetTimerChan:
			log.Debug().Msg("Resetting timer")
			if !r.runnerBase().timer.Stop() && len(r.runnerBase().timer.C) > 0 {
				<-r.runnerBase().timer.C
			}
			r.runnerBase().timer.Reset(r.runnerBase().sampleRate)
		case <-r.runnerBase().ctx.Done():
			r.runnerBase().timer.Stop()
			log.Debug().Msg("Stopping timer...")
			return
		}
	}
}
