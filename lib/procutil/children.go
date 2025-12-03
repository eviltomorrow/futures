package procutil

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/eviltomorrow/futures/lib/netutil"
)

func RunAppWithChildren(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("panic: invalid os.Args format")
	}

	var (
		name    = args[0]
		newArgs = func() []string {
			data := make([]string, 0, len(args)-2)
			// for _, arg := range args[1:] {
			// 	switch arg {
			// 	case "-d", "--daemon":

			// 	default:
			// 		data = append(data, arg)
			// 	}
			// }
			data = append(data, "--disable-stdlog")
			return data
		}()
	)

	cmd := exec.Command(name, newArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	port, err := netutil.GetAvailablePort()
	if err != nil {
		return fmt.Errorf("generate available port failure, nest error: %v", err)
	}

	address := net.JoinHostPort(pingbackHost, fmt.Sprintf("%d", port))
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen %s failure, nest error: %v", address, err)
	}
	defer ln.Close()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return err
	}

	success, errCh := make(chan struct{}), make(chan error)
	go func() {
		defer stdin.Close()

		bi := &BootInfo{
			ChallengeKey: key,
			ListenPort:   port,
		}
		buf, err := bi.Marshal()
		if err != nil {
			errCh <- fmt.Errorf("marshal bootinfo failure, nest error: %v", err)
			return
		}
		_, err = stdin.Write(buf)
		if err != nil {
			errCh <- fmt.Errorf("stdin write failure, nest error: %v", err)
			return
		}
	}()

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				errCh <- fmt.Errorf("accept tcp/ip failure, nest error: %v", err)
				break
			}

			if err = handlePingbackConn(conn, key); err != nil {
				errCh <- fmt.Errorf("handle pingback data failure, nest error: %v", err)
			} else {
				close(success)
			}

		}
	}()

	go func() {
		err := cmd.Wait()
		errCh <- err
	}()

	select {
	case <-success:
		printRunning(cmd.Process.Pid)
	case err := <-errCh:
		return fmt.Errorf("startup children process failure, nest error: %v", err)
	}
	return nil
}
