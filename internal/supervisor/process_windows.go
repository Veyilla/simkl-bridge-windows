package supervisor

import (
	"context"
	"errors"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Die Job-Object-Lebensdauer folgt Resurrector 1.2.0 (MIT).
// Angepasst: ausschließlich die eigene EXE, keine frei wählbaren Programme.
type ownedWorker struct {
	job, process windows.Handle
	stop         sync.Once
	stopErr      error
	pid          int
}

func NewOwnedWorker(ctx context.Context) (Process, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, err
	}
	limits := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err = windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		windows.CloseHandle(job)
		return nil, err
	}
	command, err := windows.UTF16PtrFromString(syscall.EscapeArg(exe) + " --bridge-worker")
	if err != nil {
		windows.CloseHandle(job)
		return nil, err
	}
	application, err := windows.UTF16PtrFromString(exe)
	if err != nil {
		windows.CloseHandle(job)
		return nil, err
	}
	si := windows.StartupInfo{Flags: windows.STARTF_USESHOWWINDOW, ShowWindow: windows.SW_HIDE}
	si.Cb = uint32(unsafe.Sizeof(si))
	var pi windows.ProcessInformation
	err = windows.CreateProcess(application, command, nil, nil, false, windows.CREATE_SUSPENDED|windows.CREATE_NEW_PROCESS_GROUP|windows.CREATE_NO_WINDOW, nil, nil, &si, &pi)
	if err != nil {
		windows.CloseHandle(job)
		return nil, err
	}
	defer windows.CloseHandle(pi.Thread)
	if err = windows.AssignProcessToJobObject(job, pi.Process); err == nil {
		_, err = windows.ResumeThread(pi.Thread)
	}
	if err != nil {
		windows.TerminateProcess(pi.Process, 1)
		windows.CloseHandle(pi.Process)
		windows.CloseHandle(job)
		return nil, err
	}
	return &ownedWorker{job: job, process: pi.Process, pid: int(pi.ProcessId)}, nil
}
func (p *ownedWorker) PID() int { return p.pid }
func (p *ownedWorker) Stop() error {
	p.stop.Do(func() { p.stopErr = windows.CloseHandle(p.job) })
	return p.stopErr
}
func (p *ownedWorker) Wait() error {
	defer windows.CloseHandle(p.process)
	defer p.Stop()
	if _, err := windows.WaitForSingleObject(p.process, windows.INFINITE); err != nil {
		return err
	}
	var code uint32
	if err := windows.GetExitCodeProcess(p.process, &code); err != nil {
		return err
	}
	switch code {
	case 0:
		return nil
	case 41:
		return ErrAuthentication
	case 42:
		return ErrConfiguration
	default:
		return errors.New("Bridge-Prozess unerwartet beendet")
	}
}
