package libfossil

import "testing"

func TestNopSyncObserverDoesNotPanic(t *testing.T) {
	o := NopSyncObserver()
	o.Started(SessionStart{ProjectCode: "test", Push: true, Pull: true})
	o.RoundStarted(1)
	o.RoundCompleted(1, RoundStats{FilesSent: 1})
	o.Completed(SessionEnd{Rounds: 1})
	o.Error(nil)
	o.HandleStarted(HandleStart{RemoteAddr: "127.0.0.1"})
	o.HandleCompleted(HandleEnd{FilesSent: 1, FilesRecvd: 2})
	o.TableSyncStarted(TableSyncStart{Table: "config"})
	o.TableSyncCompleted(TableSyncEnd{Table: "config", RowsSent: 5})
}

func TestNopCheckoutObserverDoesNotPanic(t *testing.T) {
	o := NopCheckoutObserver()
	o.ExtractStarted(ExtractStart{RID: 1, Dir: "/tmp"})
	o.ExtractFileCompleted("file.txt", ChangeAdded)
	o.ExtractCompleted(ExtractEnd{FilesWritten: 1})
	o.ScanStarted("/tmp")
	o.ScanCompleted(ScanEnd{FilesScanned: 5})
	o.CommitStarted(CommitStart{Comment: "test", User: "user", Files: 1})
	o.CommitCompleted(CommitEnd{UUID: "abc123", RID: 1})
	o.Error(nil)
}

func TestStdoutSyncObserverDoesNotPanic(t *testing.T) {
	o := StdoutSyncObserver()
	o.Started(SessionStart{})
	o.RoundStarted(1)
	o.RoundCompleted(1, RoundStats{})
	o.Completed(SessionEnd{})
	o.Error(nil)
	o.HandleStarted(HandleStart{})
	o.HandleCompleted(HandleEnd{})
	o.TableSyncStarted(TableSyncStart{})
	o.TableSyncCompleted(TableSyncEnd{})
}

func TestStdoutCheckoutObserverDoesNotPanic(t *testing.T) {
	o := StdoutCheckoutObserver()
	o.ExtractStarted(ExtractStart{})
	o.ExtractFileCompleted("file.txt", ChangeAdded)
	o.ExtractCompleted(ExtractEnd{})
	o.ScanStarted("/tmp")
	o.ScanCompleted(ScanEnd{})
	o.CommitStarted(CommitStart{})
	// Test with empty UUID (guard against panic on UUID[:10])
	o.CommitCompleted(CommitEnd{UUID: "", RID: 1})
	o.CommitCompleted(CommitEnd{UUID: "abcdef1234567890", RID: 2})
	o.Error(nil)
}
