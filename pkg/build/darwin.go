package build

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/google/syzkaller/pkg/osutil"
)

type darwin struct{}

func (ctx darwin) build(targetArch, vmType, kernelDir, outputDir, compiler, userspaceDir,
	cmdlineFile, sysctlFile string, config []byte) error {
	const kernelName = "SYZKALLER"
	confDir := fmt.Sprintf("%v/sys/arch/%v/conf", kernelDir, targetArch)
	compileDir := fmt.Sprintf("%v/sys/arch/%v/compile/%v", kernelDir, targetArch, kernelName)

	if err := ctx.configure(confDir, compileDir, kernelName); err != nil {
		return err
	}

	if err := ctx.make(compileDir, "all"); err != nil {
		return err
	}

	for src, dst := range map[string]string{
		filepath.Join(compileDir, "obj/darwin"):     "kernel",
		filepath.Join(compileDir, "obj/darwin.gdb"): "obj/darwin.gdb",
		filepath.Join(userspaceDir, "image"):        "image",
		filepath.Join(userspaceDir, "key"):          "key",
	} {
		fullDst := filepath.Join(outputDir, dst)
		if err := osutil.CopyFile(src, fullDst); err != nil {
			return fmt.Errorf("failed to copy %v -> %v: %v", src, fullDst, err)
		}
	}
	return nil
}

func (ctx darwin) clean(kernelDir string) error {
	return ctx.make(kernelDir, "", "clean")
}

func (ctx darwin) configure(confDir, compileDir, kernelName string) error {
	conf := []byte(`
include "arch/amd64/conf/GENERIC"
pseudo-device kcov 1
`)
	if err := osutil.WriteFile(filepath.Join(confDir, kernelName), conf); err != nil {
		return err
	}

	if err := osutil.MkdirAll(compileDir); err != nil {
		return err
	}
	makefile := []byte(".include \"../Makefile.inc\"\n")
	if err := osutil.WriteFile(filepath.Join(compileDir, "Makefile"), makefile); err != nil {
		return err
	}
	if err := ctx.make(compileDir, "obj"); err != nil {
		return err
	}
	if err := ctx.make(compileDir, "config"); err != nil {
		return err
	}

	return nil
}

func (ctx darwin) make(kernelDir string, args ...string) error {
	args = append([]string{"-j", strconv.Itoa(runtime.NumCPU())}, args...)
	_, err := osutil.RunCmd(10*time.Minute, kernelDir, "make", args...)
	return err
}
