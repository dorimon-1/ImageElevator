package runner

import (
	"github.com/rs/zerolog/log"
)

func (r *Runner) timerRoutine() {
	log.Debug().Msg("Timer routine started")
	for {
		select {
		case <-r.timer.C:
			log.Debug().Msg("triggering an upload via timer")
			_ = r.TriggerUpload()

		case <-r.resetTimerChan:
			log.Debug().Msg("Resetting timer")
			if !r.timer.Stop() && len(r.timer.C) > 0 {
				<-r.timer.C
			}
			r.timer.Reset(r.sampleRate)
		case <-r.ctx.Done():
			r.timer.Stop()
			log.Debug().Msg("Stopping timer...")
			return
		}
	}
}
