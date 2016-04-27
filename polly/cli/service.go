package cli

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gotil"

	"github.com/emccode/polly/daemon"
	"github.com/emccode/polly/util"
)

var (
	useSystemDForSCMCmds = gotil.FileExists(util.UnitFilePath) &&
		getInitSystemType() == SystemD
)

func (c *CLI) start() {
	if !c.fg && useSystemDForSCMCmds {
		startViaSystemD()
		return
	}

	checkOpPerms("started")

	log.WithField("os.Args", os.Args).Debug("invoking service start")

	pidFile := util.PidFilePath()
	_ = util.LogDirPath()

	if gotil.FileExists(pidFile) {
		pid, pidErr := util.ReadPidFile()
		if pidErr != nil {
			fmt.Printf("Error reading Polly PID file at %s\n", pidFile)
			panic(1)
		}

		rrproc, err := findProcess(pid)
		if err != nil {
			fmt.Printf("Error finding process for PID %d", pid)
			panic(1)
		}

		if rrproc != nil {
			fmt.Printf("Polly already running at PID %d\n", pid)
			panic(1)
		}

		if err := os.RemoveAll(pidFile); err != nil {
			fmt.Println("Error removing Polly PID file")
			panic(1)
		}
	}

	if c.fg || c.client != "" {
		c.startDaemon()
	} else {
		c.tryToStartDaemon()
	}
}

func failOnError(err error) {
	if err != nil {
		fmt.Printf("FAILED!\n  %v\n", err)
		panic(err)
	}
}

func startViaSystemD() {
	execSystemDCmd("start")
	statusViaSystemD()
}

func stopViaSystemD() {
	execSystemDCmd("stop")
	statusViaSystemD()
}

func statusViaSystemD() {
	execSystemDCmd("status")
}

func execSystemDCmd(cmdType string) {
	cmd := exec.Command("systemctl", cmdType, "-l", "polly")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				panic(status.ExitStatus())
			}
		}
	}
}

func (c *CLI) startDaemon() {

	var out io.Writer = os.Stdout
	if !log.IsTerminal() {
		logFile, logFileErr := util.LogFile("polly.log")
		failOnError(logFileErr)
		out = io.MultiWriter(os.Stdout, logFile)
	}
	log.SetOutput(out)

	fmt.Fprintf(out, "%s\n", pollyLogoASCII)
	util.PrintVersion(out)
	fmt.Fprintln(out)

	var success []byte
	var failure []byte
	var conn net.Conn

	if !c.fg {

		success = []byte{0}
		failure = []byte{1}

		var dialErr error

		log.Printf("dialing %s", c.client)
		conn, dialErr = net.Dial("unix", c.client)
		if dialErr != nil {
			panic(dialErr)
		}
	}

	writePidErr := util.WritePidFile(-1)
	if writePidErr != nil {
		if conn != nil {
			conn.Write(failure)
		}
		panic(writePidErr)
	}

	defer func() {
		r := recover()
		os.Remove(util.PidFilePath())
		if r != nil {
			panic(r)
		}
	}()

	log.Printf("created pid file, pid=%d", os.Getpid())

	init := make(chan error)
	sigc := make(chan os.Signal, 1)
	stop := make(chan os.Signal)

	signal.Notify(sigc,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		daemon.Start(c.p.Config, init, stop)
	}()

	var initErrors []error

	for initErr := range init {
		initErrors = append(initErrors, initErr)
		log.Println(initErr)
	}

	if conn != nil {
		if len(initErrors) == 0 {
			conn.Write(success)
		} else {
			conn.Write(failure)
		}

		conn.Close()
	}

	if len(initErrors) > 0 {
		return
	}

	sigv := <-sigc
	log.Printf("received shutdown signal %v", sigv)
	stop <- sigv
}

func (c *CLI) tryToStartDaemon() {
	_, _, thisAbsPath := gotil.GetThisPathParts()

	fmt.Print("Starting Polly...")

	signal := make(chan byte)
	client := fmt.Sprintf("%s/%s.sock", os.TempDir(), gotil.RandomString(32))
	log.WithField("client", client).Debug("trying to start service")

	l, lErr := net.Listen("unix", client)
	failOnError(lErr)

	go func() {
		conn, acceptErr := l.Accept()
		if acceptErr != nil {
			fmt.Printf("FAILED!\n  %v\n", acceptErr)
			panic(acceptErr)
		}
		defer conn.Close()
		defer os.Remove(client)

		log.Debug("accepted connection")

		buff := make([]byte, 1)
		conn.Read(buff)

		log.Debug("received data")

		signal <- buff[0]
	}()

	cmdArgs := []string{
		"start",
		fmt.Sprintf("--client=%s", client),
		fmt.Sprintf("--logLevel=%v", c.logLevel())}

	cmd := exec.Command(thisAbsPath, cmdArgs...)

	cmdErr := cmd.Start()
	failOnError(cmdErr)

	sigVal := <-signal
	if sigVal != 0 {
		fmt.Println("FAILED!")
		panic(1)
	}

	pid, _ := util.ReadPidFile()
	fmt.Printf("SUCCESS!\n\n")
	fmt.Printf("  The Polly daemon is now running at PID %d. To\n", pid)
	fmt.Printf("  shutdown the daemon execute the following command:\n\n")
	fmt.Printf("    sudo %s stop\n\n", thisAbsPath)
}

func stop() {
	if useSystemDForSCMCmds {
		stopViaSystemD()
		return
	}

	checkOpPerms("stopped")

	if !gotil.FileExists(util.PidFilePath()) {
		fmt.Println("Polly is already stopped")
		panic(1)
	}

	fmt.Print("Shutting down Polly...")

	pid, pidErr := util.ReadPidFile()
	failOnError(pidErr)

	proc, procErr := os.FindProcess(pid)
	failOnError(procErr)

	killErr := proc.Signal(syscall.SIGHUP)
	failOnError(killErr)

	if err := os.RemoveAll(util.PidFilePath()); err != nil {
		fmt.Println("Error removing Polly PID file")
		panic(1)
	}

	fmt.Println("SUCCESS!")
}

func (c *CLI) status() {
	if useSystemDForSCMCmds {
		statusViaSystemD()
		return
	}

	pidFile := util.PidFilePath()

	if !gotil.FileExists(pidFile) {
		fmt.Println("Polly is stopped")
		return
	}

	pid, pidErr := util.ReadPidFile()
	if pidErr != nil {
		fmt.Printf("Error reading Polly PID file at %s\n", pidFile)
		panic(1)
	}

	rrproc, err := findProcess(pid)

	if err != nil || rrproc == nil {
		if err := os.RemoveAll(pidFile); err != nil {
			fmt.Println("Error removing stale Polly PID file")
			panic(1)
		}
		fmt.Println("Polly is stopped")
		return
	}

	fmt.Printf("Polly is running at PID %d\n", pid)
	return
}

func (c *CLI) restart() {
	checkOpPerms("restarted")

	if gotil.FileExists(util.PidFilePath()) {
		stop()
	}

	c.start()
}

func checkOpPerms(op string) error {
	// if os.Geteuid() != 0 {
	// 	return goof.Newf("Polly can only be %s by root", op)
	// }

	return nil
}

const pollyLogoASCII = `<Polymorphic Storage Scheduling>`
