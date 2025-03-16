package httpsrv

import (
	"goapp/internal/pkg/watcher"
)

func (s *Server) addWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	s.watchers[w.GetWatcherId()] = w
}

func (s *Server) removeWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	// Print satistics before removing watcher.
	for i := range s.sessionStats {
		if s.sessionStats[i].id == w.GetWatcherId() {
			s.sessionStats[i].print()

			// Remove the entry from the slice by replacing it with the last one
			// and truncating the slice
			lastIdx := len(s.sessionStats) - 1
			s.sessionStats[i] = s.sessionStats[lastIdx]
			s.sessionStats = s.sessionStats[:lastIdx]

			break
		}
	}
	// Remove watcher.
	delete(s.watchers, w.GetWatcherId())
}

func (s *Server) notifyWatchers(str string) {
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()

	// Send message to all watchers and increment stats.
	for id := range s.watchers {
		s.watchers[id].Send(str)
		s.incStats(id)
	}
}
