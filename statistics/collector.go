package statistics

// StatsCollector interface defines the methods for a collector.
// A collector is a temporary storage of incoming statistics until they can be
// dumped in a more permanent storage, usually a MySQL or a time-series database.
// Statistics are stored on per application basis. So all the methods takes the id of
// the application as a parameter.
type StatsCollector interface {
	// HandleWebsocketMessage collects all the outgoing websocket messages.
	HandleWebsocketMessage(appId string)

	// HandleApiMessage collects all the incoming api message.
	HandleApiMessage(appId string)

	// HandleConnection collects all the new connections made to the server.
	HandleConnection(appId string)

	// HandleDisconnection collects when the connection drops from the server.
	HandleDisconnection(appId string)

	// Flush will empty all the stored statistics. Usually should be called when
	// stats are stored permanently in the database.
	Flush()

	// GetAllStatistics will return all the statistics across all the apps.
	GetAllStatistics() []Statistic

	// GetAppStatistics will return the current stored statistics for the appId.
	GetAppStatistics(appId string) Statistic
}
