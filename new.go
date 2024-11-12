package selenium

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func NewChromeDriver() (d *WebDriver, err error) {
	return doNewChromeDriver(nil, nil)
}

func NewChromeDriverCap(cap Capabilities, chromeCap ChromeCapabilities) (d *WebDriver, err error) {
	return doNewChromeDriver(&cap, &chromeCap)
}

func doNewChromeDriver(cap *Capabilities, chromeCap *ChromeCapabilities) (d *WebDriver, err error) {
	var (
		service    *selenium.Service
		port       int
		driverPath string
	)
	if port, err = FindFreePort(); err != nil {
		return nil, err
	}
	driverPath, err = exec.LookPath("chromedriver")
	if err != nil {
		return
	}

	if service, err = selenium.NewChromeDriverService(driverPath, port); nil != err {
		return nil, fmt.Errorf("start a chromedriver service falid %v", err)
	}
	log.Printf("start a chromedriver service on 127.0.0.1:%v", port)

	if cap == nil {
		cap = NewCapabilities()
	}

	if driver, err := selenium.NewRemote(newCaps(chromeCap), fmt.Sprintf("http://localhost:%d/wd/hub", port)); err != nil {
		service.Stop()
		return nil, fmt.Errorf("failed to open session: %v", err)
	} else {
		return &WebDriver{Data: driver, Ser: service}, nil
	}
}

func newCaps(c *ChromeCapabilities) (caps selenium.Capabilities) {
	caps = selenium.Capabilities{"browserName": "chrome"}
	if c != nil {
		caps.AddChrome(chrome.Capabilities(*c))
	}
	return
}
