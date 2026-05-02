package photosync

import (
	"io"
	"os"
	"os/exec"
)

func RunPhotofs(path string) os.Cmd, error {
	cmd := exec.Command("photofs", "--path", path)

	return cmd, cmd.Start()
}

func Kill(cmd os.Cmd) error {
	err := cmd.Process.Kill()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func Mount(path string) error {
	return exec.Command("9mount", "tcp!localhost!8000", path).Run()
}

func Umount(path string) error {
	return exec.Command("9umount", path).Run()
}

func Sync(src string, dest string, w io.Writer) error {
	cmd := exec.Command("rsync", "-n", src, dest)
	cmd.Stdout = w
	cmd.Stderr = w

	return cmd.Run()
}
