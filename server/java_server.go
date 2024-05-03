package server

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
)

type JavaConfig struct {
    JavaHome string `json:"javaHome"`
    JavaOpts string `json:"javaOpts"`
    JarFile string `json:"jarFile"`
}

func JavaServerBuild(config JavaConfig) error {
    // build the java server
    cmd := exec.Command("mvn", "package", "-DskipTests")
    cmd.Env = append(cmd.Env, "JAVA_HOME=" + config.JavaHome)
    stdout, err := cmd.StdoutPipe()
    cmd.Stderr = cmd.Stdout
    if err != nil {
        return err
    }

    reader := bufio.NewReader(stdout)
    err = cmd.Start()
    if err != nil {
        return err
    }

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        print(line)
    }
    return cmd.Wait()
}

func JavaServerRun(config JavaConfig) error {
    // run the java server
    cmd := exec.Command("java")
    cmd.Path = config.JavaHome + "/bin/java"
    arguments := strings.Split(config.JavaOpts, " ")
    cmd.Args = append(cmd.Args, arguments...)
    cmd.Args = append(cmd.Args, "-jar", config.JarFile)
    log.Printf("Running java server with command %s", cmd.String())
    cmd.Env = append(cmd.Env, "JAVA_HOME=" + config.JavaHome)
    stdout, err := cmd.StdoutPipe()
    cmd.Stderr = cmd.Stdout
    if err != nil {
        return err
    }

    reader := bufio.NewReader(stdout)
    err = cmd.Start()
    if err != nil {
        return err
    }

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        print(line)
    }
    return cmd.Wait()
}
