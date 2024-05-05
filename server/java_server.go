package server

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/orlandokj/just/config"
)

type JavaConfig struct {
    JavaHome string `json:"javaHome"`
    JavaOpts string `json:"javaOpts"`
    JarFile string `json:"jarFile"`
    workDir string
}


type JavaServer struct {
    config JavaConfig
    logIngestion LogFunc
}

func (js JavaServer) Build() (ServerProcess, error) {
    cmd := exec.Command("mvn", "package", "-DskipTests")
    cmd.Path = js.config.workDir
    cmd.Env = append(cmd.Env, "JAVA_HOME=" + js.config.JavaHome)
    return RunCmd(cmd, js.logIngestion)
}

func (js JavaServer) Run() (ServerProcess, error) {
    cmd := exec.Command("java")
    cmd.Path = js.config.JavaHome + "/bin/java"
    cmd.Dir = js.config.workDir
    arguments := strings.Split(js.config.JavaOpts, " ")
    cmd.Args = append(cmd.Args, arguments...)
    cmd.Args = append(cmd.Args, "-jar", js.config.JarFile)
    js.logIngestion(fmt.Sprintf("Running java server with command %s", cmd.String()))
    return RunCmd(cmd, js.logIngestion)
}

func CreateJavaServer(config config.Config, logIngestion LogFunc) (Server, error) {
    javaConfig := JavaConfig{
        workDir: config.WorkDir,
    }
    err := config.ToConfigType(&javaConfig)
    if err != nil {
        return nil, err
    }

    return JavaServer{
        config: javaConfig,
        logIngestion: logIngestion,
    }, nil
}

