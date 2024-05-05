package server

import (
	"bufio"
	"log"
	"os/exec"
	"strings"

	"github.com/orlandokj/just/config"
)

type JavaConfig struct {
    JavaHome string `json:"javaHome"`
    JavaOpts string `json:"javaOpts"`
    JarFile string `json:"jarFile"`
}

type JavaServer struct {
    config JavaConfig
}


func (js JavaServer) Build() error {
    // build the java server
    cmd := exec.Command("mvn", "package", "-DskipTests")
    cmd.Env = append(cmd.Env, "JAVA_HOME=" + js.config.JavaHome)
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

func (js JavaServer) Run() error {
    // run the java server
    cmd := exec.Command("java")
    cmd.Path = js.config.JavaHome + "/bin/java"
    arguments := strings.Split(js.config.JavaOpts, " ")
    cmd.Args = append(cmd.Args, arguments...)
    cmd.Args = append(cmd.Args, "-jar", js.config.JarFile)
    log.Printf("Running java server with command %s", cmd.String())
    cmd.Env = append(cmd.Env, "JAVA_HOME=" + js.config.JavaHome)
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

func CreateJavaServer(config config.Config) (Server, error) {
    javaConfig := JavaConfig{}
    err := config.ToConfigType(&javaConfig)
    if err != nil {
        return nil, err
    }

    return JavaServer{
        config: javaConfig,
    }, nil
}
