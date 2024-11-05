package log

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
)

var (
	debug    *log.Logger
	info     *log.Logger
	notice   *log.Logger
	warning  *log.Logger
	err      *log.Logger
	crit     *log.Logger
	priority syslog.Priority
	writer   *syslog.Writer
)

func ParseLogLevel(s string) (syslog.Priority, error) {
	var prio syslog.Priority
	switch strings.ToLower(s) {
	case "debug":
		prio = syslog.LOG_DEBUG
	case "info":
		prio = syslog.LOG_INFO
	case "notice":
		prio = syslog.LOG_NOTICE
	case "warning":
		prio = syslog.LOG_WARNING
	case "err":
		prio = syslog.LOG_ERR
	case "crit":
		prio = syslog.LOG_CRIT
	default:
		return 0, fmt.Errorf("invalid log level %q", s)
	}
	return prio, nil
}

func InitLogging(prio syslog.Priority) error {
	priority = prio
	if os.Getenv("INVOCATION_ID") != "" {
		// executed by systemd
		w, err := syslog.New(syslog.LOG_DAEMON, "")
		if err != nil {
			return err
		}
		writer = w
	} else {
		flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
		debug = log.New(os.Stderr, "DEBUG ", flags)
		info = log.New(os.Stderr, "INFO  ", flags)
		notice = log.New(os.Stderr, "NOTICE  ", flags)
		warning = log.New(os.Stderr, "WARNING  ", flags)
		err = log.New(os.Stderr, "ERR ", flags)
		crit = log.New(os.Stderr, "CRIT ", flags)
	}
	return nil
}

func format(message string, args ...any) string {
	return fmt.Sprintf(strings.TrimSpace(message)+"\n", args...)
}

func Debugf(message string, args ...any) {
	if priority < syslog.LOG_DEBUG {
		return
	}
	msg := format(message, args...)
	if writer != nil {
		writer.Debug(msg)
	} else {
		debug.Output(2, msg)
	}
}

func Infof(message string, args ...any) {
	if priority < syslog.LOG_INFO {
		return
	}
	msg := format(message, args...)
	if writer != nil {
		writer.Info(msg)
	} else {
		info.Output(2, msg)
	}
}

func Noticef(message string, args ...any) {
	if priority < syslog.LOG_NOTICE {
		return
	}
	msg := format(message, args...)
	if writer != nil {
		writer.Notice(msg)
	} else {
		notice.Output(2, msg)
	}
}

func Warnf(message string, args ...any) {
	if priority < syslog.LOG_WARNING {
		return
	}
	msg := format(message, args...)
	if writer != nil {
		writer.Warning(msg)
	} else {
		warning.Output(2, msg)
	}
}

func Errf(message string, args ...any) {
	if priority < syslog.LOG_ERR {
		return
	}
	msg := format(message, args...)
	if writer != nil {
		writer.Err(msg)
	} else {
		err.Output(2, msg)
	}
}

func Critf(message string, args ...any) {
	if priority < syslog.LOG_CRIT {
		return
	}
	msg := format(message, args...)
	if writer != nil {
		writer.Crit(msg)
	} else {
		crit.Output(2, msg)
	}
}
