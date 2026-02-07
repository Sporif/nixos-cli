package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/nix-community/nixos-cli/internal/logger"
)

type System interface {
	CommandRunner
	IsNixOS() bool
	IsRemote() bool
	FS() Filesystem
}

// Invoke the `nix-copy-closure` command to copy between two types of
// systems.
func CopyClosures(src System, dest System, paths []string, extraArgs ...string) error {
	log := src.Logger()

	if len(paths) == 0 {
		log.Debugf("no store paths to copy")
		return nil
	}

	argv := []string{"nix-copy-closure"}
	sshopts := []string{os.Getenv("NIX_SSHOPTS")}

	srcIsRemote := src.IsRemote()
	destIsRemote := dest.IsRemote()

	var commandRunner CommandRunner

	// All type asserts must work here, otherwise the IsRemote() method is
	// implemented incorrectly for a given platform or the conditions are
	// put together incorrectly.
	//
	// There are/will be no other system types implemented, so casting
	// directly is fine here.
	if srcIsRemote && destIsRemote {
		// remote -> remote, so treat the source as a store and use the local
		// machine as the command runner.
		//
		// This should either be running as a trusted user or as root, so
		// remote store access should exist.
		commandRunner = NewLocalSystem(log)
		srcAddr := fmt.Sprintf("ssh://%s", dest.(*SSHSystem).Address())
		destAddr := dest.(*SSHSystem).Address()
		argv = append(argv, "--store", srcAddr, "--to", destAddr)
		sshopts = append(sshopts, src.(*SSHSystem).Sshopts()...)
		sshopts = append(sshopts, dest.(*SSHSystem).Sshopts()...)
	} else if srcIsRemote && !destIsRemote {
		// remote -> local, so use --from and run on the local host (dest), since there
		// is no reliable way to run this on the remote while determining how
		// the local address appears to it.
		commandRunner = dest
		srcAddr := src.(*SSHSystem).Address()
		argv = append(argv, "--from", srcAddr)
		sshopts = append(sshopts, src.(*SSHSystem).Sshopts()...)
	} else if !srcIsRemote && destIsRemote {
		// local -> remote, so run this command on the local host.
		commandRunner = src
		destAddr := dest.(*SSHSystem).Address()
		argv = append(argv, "--to", destAddr)
		sshopts = append(sshopts, dest.(*SSHSystem).Sshopts()...)
	} else {
		// local -> local, no-op
		log.Debugf("both systems are local, skipping copy")
		return nil
	}

	argv = append(argv, extraArgs...)
	if log.GetLogLevel() == logger.LogLevelDebug {
		argv = append(argv, "-v")
	}

	argv = append(argv, paths...)

	log.CmdArray(argv)

	cmd := NewCommand(argv[0], argv[1:]...)
	sshoptsEnv := strings.TrimSpace(strings.Join(sshopts, " "))
	log.Debugf("before NIX_SSHOPTS='%v'", sshoptsEnv)
	if sshoptsEnv != "" {
		cmd.SetEnv("NIX_SSHOPTS", sshoptsEnv)
		log.Debugf("using NIX_SSHOPTS='%v'", sshoptsEnv)
	}
	_, err := commandRunner.Run(cmd)
	return err
}
