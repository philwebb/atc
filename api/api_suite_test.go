package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api"
	"github.com/concourse/atc/api/buildserver"
	"github.com/concourse/atc/api/buildserver/buildserverfakes"
	"github.com/concourse/atc/api/containerserver/containerserverfakes"
	"github.com/concourse/atc/api/jobserver/jobserverfakes"
	"github.com/concourse/atc/api/pipes/pipesfakes"
	"github.com/concourse/atc/api/resourceserver/resourceserverfakes"
	"github.com/concourse/atc/api/teamserver/teamserverfakes"
	"github.com/concourse/atc/api/volumeserver/volumeserverfakes"
	"github.com/concourse/atc/api/workerserver/workerserverfakes"
	"github.com/concourse/atc/auth/authfakes"
	"github.com/concourse/atc/config"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/engine/enginefakes"
	"github.com/concourse/atc/worker/workerfakes"
	"github.com/concourse/atc/wrappa"
)

var (
	sink *lager.ReconfigurableSink

	externalURL  = "https://example.com"
	oAuthBaseURL = "https://oauth.example.com"

	authValidator                 *authfakes.FakeValidator
	userContextReader             *authfakes.FakeUserContextReader
	fakeTokenGenerator            *authfakes.FakeTokenGenerator
	providerFactory               *authfakes.FakeProviderFactory
	fakeEngine                    *enginefakes.FakeEngine
	fakeWorkerClient              *workerfakes.FakeClient
	authDB                        *authfakes.FakeAuthDB
	buildsDB                      *buildserverfakes.FakeBuildsDB
	volumesDB                     *volumeserverfakes.FakeVolumesDB
	configDB                      *dbfakes.FakeConfigDB
	workerDB                      *workerserverfakes.FakeWorkerDB
	containerDB                   *containerserverfakes.FakeContainerDB
	pipeDB                        *pipesfakes.FakePipeDB
	pipelineDBFactory             *dbfakes.FakePipelineDBFactory
	pipelinesDB                   *dbfakes.FakePipelinesDB
	teamDB                        *teamserverfakes.FakeTeamDB
	fakeSchedulerFactory          *jobserverfakes.FakeSchedulerFactory
	fakeScannerFactory            *resourceserverfakes.FakeScannerFactory
	configValidationErrorMessages []string
	configValidationWarnings      []config.Warning
	peerAddr                      string
	drain                         chan struct{}
	cliDownloadsDir               string
	logger                        *lagertest.TestLogger

	constructedEventHandler *fakeEventHandlerFactory

	server *httptest.Server
	client *http.Client
)

type fakeEventHandlerFactory struct {
	db      buildserver.BuildsDB
	buildID int

	lock sync.Mutex
}

func (f *fakeEventHandlerFactory) Construct(
	logger lager.Logger,
	db buildserver.BuildsDB,
	buildID int,
) http.Handler {
	f.lock.Lock()
	f.db = db
	f.buildID = buildID
	f.lock.Unlock()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("fake event handler factory was here"))
	})
}

var _ = BeforeEach(func() {
	authDB = new(authfakes.FakeAuthDB)
	buildsDB = new(buildserverfakes.FakeBuildsDB)
	configDB = new(dbfakes.FakeConfigDB)
	pipelineDBFactory = new(dbfakes.FakePipelineDBFactory)
	workerDB = new(workerserverfakes.FakeWorkerDB)
	containerDB = new(containerserverfakes.FakeContainerDB)
	volumesDB = new(volumeserverfakes.FakeVolumesDB)
	pipeDB = new(pipesfakes.FakePipeDB)
	pipelinesDB = new(dbfakes.FakePipelinesDB)
	teamDB = new(teamserverfakes.FakeTeamDB)

	authValidator = new(authfakes.FakeValidator)
	userContextReader = new(authfakes.FakeUserContextReader)
	fakeTokenGenerator = new(authfakes.FakeTokenGenerator)
	providerFactory = new(authfakes.FakeProviderFactory)

	configValidationErrorMessages = []string{}
	configValidationWarnings = []config.Warning{}
	peerAddr = "127.0.0.1:1234"
	drain = make(chan struct{})

	fakeEngine = new(enginefakes.FakeEngine)
	fakeWorkerClient = new(workerfakes.FakeClient)

	fakeSchedulerFactory = new(jobserverfakes.FakeSchedulerFactory)
	fakeScannerFactory = new(resourceserverfakes.FakeScannerFactory)

	var err error

	cliDownloadsDir, err = ioutil.TempDir("", "cli-downloads")
	Expect(err).NotTo(HaveOccurred())

	constructedEventHandler = &fakeEventHandlerFactory{}

	logger = lagertest.NewTestLogger("callbacks")

	sink = lager.NewReconfigurableSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG), lager.DEBUG)
	logger.RegisterSink(sink)

	handlers, err := api.NewHandler(
		logger,

		externalURL,

		[]wrappa.Wrappa{wrappa.NewAPIAuthWrappa(true, authValidator, userContextReader)},

		fakeTokenGenerator,
		providerFactory,
		oAuthBaseURL,

		pipelineDBFactory,
		configDB,

		authDB,
		buildsDB,
		workerDB,
		containerDB,
		volumesDB,
		pipeDB,
		pipelinesDB,
		teamDB,

		func(atc.Config) ([]config.Warning, []string) {
			return configValidationWarnings, configValidationErrorMessages
		},
		peerAddr,
		constructedEventHandler.Construct,
		drain,

		fakeEngine,
		fakeWorkerClient,

		fakeSchedulerFactory,
		fakeScannerFactory,

		sink,

		cliDownloadsDir,
		"1.2.3",
	)
	Expect(err).NotTo(HaveOccurred())
	Expect(handlers).To(HaveLen(1))
	handler := handlers[0]

	server = httptest.NewServer(handler)

	client = &http.Client{
		Transport: &http.Transport{},
	}
})

var _ = AfterEach(func() {
	server.Close()
})

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}
