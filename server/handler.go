package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"sync"

	"github.com/muraenateam/necrobrowser/action/navigation"
	"github.com/muraenateam/necrobrowser/core"
	"github.com/muraenateam/necrobrowser/log"
	"github.com/muraenateam/necrobrowser/zombie"
)

var (
	// starting necrobrowser from the home directory of a classic Linux GUI distro,
	// the directory where files downloaded by Zombie are placed is /home/userX/Downloads
	jobs       sync.Map
	jobsStatus []InstrumentationJob
)

var (
	ErrUnsupportedProvider = errors.New("provider not supported")
)

// this goes in the Sync.Map jobs
type InstrumentationJob struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	Context  context.Context
}

type KnownBrowserRequest struct {
	JobID    string `json:"jobId" binding:"required"`
	URL      string `json:"url"`
	Selector string `json:"selector"`
}

type InstrumentationRequest struct {
	Provider      string `json:"provider" binding:"required"`
	DebuggingPort int    `json:"debugPort" binding:"required"`

	// classic credentials including 2fa token if any
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`

	// authenticated session cookie retrieved from Hive which allows for direct
	// session riding without the need to authenticate via credentials
	SessionCookies []zombie.SessionCookie `json:"sessionCookies"`

	// keywords to search through emails or in general search bars
	// for example: password, credentials, access, https://, vpn, certificate, credit card, etc..
	Keywords []string `json:"keywords"`
}

func StatusHandler(c *gin.Context) {
	if !Authenticate(c) {
		return
	}

	jobs.Range(func(k, v interface{}) bool {
		_, job := k.(string), v.(InstrumentationJob)
		jobsStatus = append(jobsStatus, job)
		return true
	})

	c.JSON(http.StatusOK, gin.H{"jobs": jobsStatus})

}

func KnownBrowserHandler(c *gin.Context) {
	if !Authenticate(c) {
		return
	}

	var toInstrument KnownBrowserRequest
	err := c.BindJSON(&toInstrument)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err})
		return
	}

	err = instrumentKnownBrowser(toInstrument)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"success": false, "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}

func logError(err error) {
	log.Error("ERROR: %+v", err)
}

func NewBrowserHandler(c *gin.Context) {
	if !Authenticate(c) {
		return
	}

	// Retrieve config options
	options := GetOptions(c)

	var toInstrument InstrumentationRequest
	err := c.BindJSON(&toInstrument)
	if err != nil {
		logError(err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if !core.IsValidZombie(toInstrument.Provider) {
		logError(ErrUnsupportedProvider)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": ErrUnsupportedProvider})
		return
	}

	err = toInstrument.instrumentNewBrowser(options)
	if err != nil {
		logError(err)
		c.JSON(http.StatusExpectationFailed, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "target": toInstrument.Username, "extruded": "TODO"})
	return
}

// instrumentNewBrowser uses Chrome DevTools Protocol (CDP) to instrument Chrome

func (i *InstrumentationRequest) instrumentNewBrowser(options core.Options) (err error) {
	log.Info("Instructing zombie for %s", i.Provider)

	// The zombie
	t := zombie.Target{
		Context:  context.Background(),
		Cookies:  i.SessionCookies,
		Username: i.Username,
		Password: i.Password,
	}
	z, err := core.GetZombie(i.Provider, t, options)
	if err != nil {
		return err
	}

	// Fetch the loot path, it's used in docker to mount the volume
	loot := core.GetZombieLootPath(options.LootPath, zombie.GetTag(z.Name()))
	if _, err = core.CheckLoot(loot); err != nil {
		return
	}

	// check if headless or gui mode
	var allocCtx context.Context
	var cancelCtx context.CancelFunc

	if *options.Headless == true {
		log.Info("Going HEADLESS mode")

		// Create a new Docker container
		name := fmt.Sprintf("%s_%s", i.Provider, Random(10))
		log.Info("Creating a new container %s", name)
		c, err := core.NewContainer(name, *options.DockerImage, strconv.Itoa(i.DebuggingPort), loot)
		defer c.Cancel()
		if err != nil {
			return err
		}

		log.Info("A new container is alive at %s", c.Target.WebSocketDebuggerUrl)

		// create new remote allocator pointing to docker debug port
		allocCtx, cancelCtx = chromedp.NewRemoteAllocator(context.Background(), c.Target.WebSocketDebuggerUrl)
		defer cancelCtx()
	} else {
		log.Info("Going GUI mode")
		opts := []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			// chromedp.UserDataDir(loot),
			chromedp.WindowSize(1920, 1080),
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"),
		}

		// create new local exec allocator in gui mode
		allocCtx, cancelCtx = chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancelCtx()
	}

	id := uuid.NewV4().String()
	log.Info("Adding new Instrumentation job with id: %s", id)
	jobs.Store(id, InstrumentationJob{id, i.Provider, allocCtx})

	// Instrument the zombie
	zombieCtx, _ := chromedp.NewContext(allocCtx)

	// force a context timeout
	//timeoutCtx, _ := context.WithTimeout(zombieCtx, 30 * time.Second)

	z.SetContext(zombieCtx)
	_, err = z.Instrument()
	if err != nil {
		return err
	}

	//err = chromeZombie.Wait()
	//if err != nil {
	//	return err
	//}

	//err = chromeZombie.Shutdown(dockerContext)
	//if err != nil {
	//	return err
	//}

	return nil

}

func instrumentKnownBrowser(req KnownBrowserRequest) (err error) {

	j, found := jobs.Load(req.JobID)
	if !found {
		err = fmt.Errorf("job id %s NOT FOUND", req.JobID)
		logError(err)
		return err
	}
	job := j.(InstrumentationJob)
	log.Info("InstrumentationJob: %+v", job)

	//
	// Navigator
	//
	n := &navigation.Navigation{
		URL: req.URL,
	}
	n.Target = zombie.Target{Context: job.Context}
	if err = n.Navigate(); err != nil {
		logError(err)
		return err
	}

	return nil
}
