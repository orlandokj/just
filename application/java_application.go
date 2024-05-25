package application

import (
	"fmt"
	"os/exec"
	"strings"
)

type JavaConfig struct {
    JavaHome string `json:"javaHome"`
    JavaOpts string `json:"javaOpts"`
    JarFile string `json:"jarFile"`
    workDir string
}


type JavaApplication struct {
    config JavaConfig
    logIngestion LogFunc
}

func (js JavaApplication) Build() (RunningProcess, error) {
    cmd := exec.Command("mvn", "package", "-DskipTests")
    cmd.Path = js.config.workDir
    cmd.Env = append(cmd.Env, "JAVA_HOME=" + js.config.JavaHome)
    return RunCmd(cmd, js.logIngestion)
}

func (js JavaApplication) Run() (RunningProcess, error) {
    cmd := exec.Command("java")
    cmd.Path = js.config.JavaHome + "/bin/java"
    cmd.Dir = js.config.workDir
    arguments := strings.Split(js.config.JavaOpts, " ")
    cmd.Args = append(cmd.Args, arguments...)
    cmd.Args = append(cmd.Args, "-jar", js.config.JarFile)
    js.logIngestion(fmt.Sprintf("Running java application with command %s", cmd.String()))
    return RunCmd(cmd, js.logIngestion)
}

func CreateJavaApplication(config ApplicationConfig, logIngestion LogFunc) (ApplicationHandler, error) {
    javaConfig := JavaConfig{
        workDir: config.WorkDir,
    }
    err := config.ToConfigType(&javaConfig)
    if err != nil {
        return nil, err
    }

    return JavaApplication{
        config: javaConfig,
        logIngestion: logIngestion,
    }, nil
}

