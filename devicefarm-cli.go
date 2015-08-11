package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {

	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})

	app := cli.NewApp()
	app.Name = "devicefarm-cli"
	app.Usage = "allows you to interact with AWS devicefarm from the command line"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{Name: "Patrick Debois",
			Email: "Patrick.Debois@jedi.be",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "list",
			Usage: "list various elements on devicefarm",
			Subcommands: []cli.Command{
				{
					Name:  "projects",
					Usage: "list the projects", // of an account
					Action: func(c *cli.Context) {
						listProjects(svc)
					},
				},
				{
					Name:  "devices",
					Usage: "list the devices", // globally
					Action: func(c *cli.Context) {
						listDevices(svc)
					},
				},
				{
					Name:  "samples",
					Usage: "list the samples",
					Action: func(c *cli.Context) {
						// Not yet implemented
						// listSamples()
					},
				},
				{
					Name:  "jobs",
					Usage: "list the jobs", // of a test
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")

						listJobs(svc, runArn)
					},
				},
				{
					Name:  "uploads",
					Usage: "lists all uploads", // of a Project
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
					},
					Action: func(c *cli.Context) {
						projectArn := c.String("project")
						listUploads(svc, projectArn)
					},
				},
				{
					Name:  "artifacts",
					Usage: "list the artifacts", // of a test
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
						cli.StringFlag{
							Name:   "job",
							EnvVar: "DF_JOB",
							Usage:  "job arn or run description",
						},
						cli.StringFlag{
							Name:   "type",
							EnvVar: "DF_ARTIFACT_TYPE",
							Usage:  "type of the artifact [LOG,FILE,SCREENSHOT]",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						jobArn := c.String("job")

						filterArn := ""
						if runArn != "" {
							filterArn = runArn
						} else {
							filterArn = jobArn
						}

						artifactType := c.String("type")
						listArtifacts(svc, filterArn, artifactType)
					},
				},
				{
					Name:  "suites",
					Usage: "list the suites",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
						cli.StringFlag{
							Name:   "job",
							EnvVar: "DF_JOB",
							Usage:  "job arn or run description",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						jobArn := c.String("job")
						filterArn := ""
						if runArn != "" {
							filterArn = runArn
						} else {
							filterArn = jobArn
						}
						listSuites(svc, filterArn)
					},
				},
				{
					Name:  "devicepools",
					Usage: "list the devicepools", //globally
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
					},
					Action: func(c *cli.Context) {

						projectArn := c.String("project")
						listDevicePools(svc, projectArn)
					},
				},
				{
					Name: "problems",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
					},
					Usage: "list the problems", // of Test
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						listUniqueProblems(svc, runArn)
					},
				},
				{
					Name:  "tests",
					Usage: "list the tests", // of a Run
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
						cli.StringFlag{
							Name:   "job",
							EnvVar: "DF_JOB",
							Usage:  "job arn or run description",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						jobArn := c.String("job")
						filterArn := ""
						if runArn != "" {
							filterArn = runArn
						} else {
							filterArn = jobArn
						}
						listTests(svc, filterArn)
					},
				},
				{
					Name:  "runs",
					Usage: "list the runs",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
					},
					Action: func(c *cli.Context) {
						projectArn := c.String("project")
						listRuns(svc, projectArn)
					},
				},
			},
		},
		{
			Name:  "download",
			Usage: "download various devicefarm elements",
			Subcommands: []cli.Command{
				{
					Name:  "artifacts",
					Usage: "download the artifacts",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
						cli.StringFlag{
							Name:   "job",
							EnvVar: "DF_JOB",
							Usage:  "job arn or run description",
						},
						cli.StringFlag{
							Name:   "type",
							EnvVar: "DF_ARTIFACT_TYPE",
							Usage:  "type of the artifact [LOG,FILE,SCREENSHOT]",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						jobArn := c.String("job")

						filterArn := ""
						if runArn != "" {
							filterArn = runArn
						} else {
							filterArn = jobArn
						}

						artifactType := c.String("type")
						downloadArtifacts(svc, filterArn, artifactType)
					},
				},
			},
		},
		{
			Name:  "status",
			Usage: "get the status of a run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "run",
					EnvVar: "DF_RUN",
					Usage:  "run arn or run description",
				},
			},
			Action: func(c *cli.Context) {
				runArn := c.String("run")
				runStatus(svc, runArn)
			},
		},
		{
			Name:  "report",
			Usage: "get report about a run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "run",
					EnvVar: "DF_RUN",
					Usage:  "run arn or run description",
				},
			},
			Action: func(c *cli.Context) {
				runArn := c.String("run")
				runReport(svc, runArn)
			},
		},
		{
			Name:  "schedule",
			Usage: "schedule a run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "project",
					EnvVar: "DF_PROJECT",
					Usage:  "project arn or project description",
				},
				cli.StringFlag{
					Name:   "device-pool",
					EnvVar: "DF_DEVICE_POOL",
					Usage:  "devicepool arn or devicepool name",
				},
				cli.StringFlag{
					Name:   "device",
					EnvVar: "DF_DEVICE",
					Usage:  "device arn or devicepool name to run the test on",
				},
				cli.StringFlag{
					Name:   "name",
					EnvVar: "DF_RUN_NAME",
					Usage:  "name to give to the run that is scheduled",
				},
				cli.StringFlag{
					Name:   "app-file",
					EnvVar: "DF_APP_FILE",
					Usage:  "path of the app file to be executed",
				},
				cli.StringFlag{
					Name:   "app-type",
					EnvVar: "DF_APP_TYPE",
					Usage:  "type of app [ANDROID_APP,IOS_APP]",
				},
				cli.StringFlag{
					Name:   "test-file",
					EnvVar: "DF_TEST_FILE",
					Usage:  "path of the test file to be executed",
				},
				cli.StringFlag{
					Name:   "test-type",
					EnvVar: "DF_TEST_TYPE",
					//Usage:  "type of test [APPIUM_JAVA_JUNIT_TEST_PACKAGE, INSTRUMENTATION_TEST_PACKAGE, UIAUTOMATION_TEST_PACKAGE, APPIUM_JAVA_TESTNG_TEST_PACKAGE, IOS_APP, CALABASH_TEST_PACKAGE, ANDROID_APP, UIAUTOMATOR_TEST_PACKAGE, XCTEST_TEST_PACKAGE, EXTERNAL_DATA]",
					Usage: "type of test [UIAUTOMATOR, CALABASH, APPIUM_JAVA_TESTNG, UIAUTOMATION, BUILTIN_FUZZ, INSTRUMENTATION, APPIUM_JAVA_JUNIT, BUILTIN_EXPLORER, XCTEST]",
				},
				cli.StringFlag{
					Name:   "test",
					Usage:  "arn or name of the test upload to schedule",
					EnvVar: "DF_TEST",
				},
				cli.StringFlag{
					Name:   "app",
					Usage:  "arn or name of the app upload to schedule",
					EnvVar: "DF_APP",
				},
			},
			Action: func(c *cli.Context) {
				projectArn := c.String("project")
				runName := c.String("name")
				deviceArn := c.String("device")
				devicePoolArn := c.String("device-pool")
				appArn := c.String("app")
				appFile := c.String("app-file")
				appType := c.String("app-type")
				testPackageArn := c.String("test-package")
				testPackageType := c.String("test-type")
				testPackageFile := c.String("test-file")
				scheduleRun(svc, projectArn, runName, deviceArn, devicePoolArn, appArn, appFile, appType, testPackageArn, testPackageFile, testPackageType)
			},
		},
		{
			Name:  "create",
			Usage: "creates various devicefarm elements",
			Subcommands: []cli.Command{
				{
					Name:  "upload",
					Usage: "creates an upload",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "name of the upload",
						},
						cli.StringFlag{
							Name:  "type",
							Usage: "type of upload [ANDROID_APP,IOS_APP,EXTERNAL_DATA,APPIUM_JAVA_JUNIT_TEST_PACKAGE,APPIUM_JAVA_TESTNG_TEST_PACKAGE,CALABASH_TEST_PACKAGE,INSTRUMENTATION_TEST_PACKAGE,UIAUTOMATOR_TEST_PACKAGE,XCTEST_TEST_PACKAGE",
						},
					},
					Action: func(c *cli.Context) {
						uploadName := c.String("name")
						uploadType := c.String("type")
						projectArn := c.String("project")
						uploadCreate(svc, uploadName, uploadType, projectArn)
					},
				},
			},
		},
		{
			Name:  "info",
			Usage: "get detailed info about various devicefarm elements",
			Subcommands: []cli.Command{
				{
					Name:  "run",
					Usage: "get info about a run",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						runInfo(svc, runArn)
					},
				},
				{
					Name:  "upload",
					Usage: "info about uploads",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "upload",
							EnvVar: "DF_UPLOAD",
							Usage:  "upload arn or description",
						},
					},
					Action: func(c *cli.Context) {
						uploadArn := c.String("upload")
						uploadInfo(svc, uploadArn)
					},
				},
			},
		},
		{
			Name:  "upload",
			Usage: "uploads an app, test and data",
			Subcommands: []cli.Command{
				{
					Name:  "file",
					Usage: "uploads an file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "name of the upload",
						},
						cli.StringFlag{
							Name:  "file",
							Usage: "path to the file to upload",
						},
						cli.StringFlag{
							Name:  "type",
							Usage: "type of upload [ANDROID_APP,IOS_APP,EXTERNAL_DATA,APPIUM_JAVA_JUNIT_TEST_PACKAGE,APPIUM_JAVA_TESTNG_TEST_PACKAGE,CALABASH_TEST_PACKAGE,INSTRUMENTATION_TEST_PACKAGE,UIAUTOMATOR_TEST_PACKAGE,XCTEST_TEST_PACKAGE",
						},
					},
					Action: func(c *cli.Context) {
						uploadType := c.String("type")
						projectArn := c.String("project")
						uploadFilePath := c.String("file")
						uploadName := c.String("name")
						uploadPut(svc, uploadFilePath, uploadType, projectArn, uploadName)
					},
				},
			},
		},
	}

	app.Run(os.Args)

}

// --- internal API starts here

/* List all Projects */
func listProjects(svc *devicefarm.DeviceFarm) {

	resp, err := svc.ListProjects(nil)
	failOnErr(err, "error listing projects")

	//fmt.Println(awsutil.Prettify(resp))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Created", "ARN"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, m := range resp.Projects {
		line := []string{*m.Name, time.Time.String(*m.Created), *m.ARN}
		table.Append(line)
	}
	table.Render() // Send output
}

/* List all DevicePools */
func listDevicePools(svc *devicefarm.DeviceFarm, projectArn string) {
	// CURATED: A device pool that is created and managed by AWS Device Farm.
	// PRIVATE: A device pool that is created and managed by the device pool developer.

	pool := &devicefarm.ListDevicePoolsInput{
		ARN: aws.String(projectArn),
	}
	resp, err := svc.ListDevicePools(pool)

	failOnErr(err, "error listing device pools")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all Devices */
func listDevices(svc *devicefarm.DeviceFarm) {

	input := &devicefarm.ListDevicesInput{}
	resp, err := svc.ListDevices(input)

	failOnErr(err, "error listing devices")
	//fmt.Println(awsutil.Prettify(resp))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Os", "Platform", "Form", "ARN"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, m := range resp.Devices {
		line := []string{*m.Name, *m.Os, *m.Platform, *m.FormFactor, *m.ARN}
		table.Append(line)
	}
	table.Render() // Send output

	/*
	   	    {
	         ARN: "arn:aws:devicefarm:us-west-2::device:A0E6E6E1059E45918208DF75B2B7EF6C",
	         CPU: {
	           Architecture: "ARMv7",
	           Clock: 2265,
	           Frequency: "MHz"
	         },
	         FormFactor: "PHONE",
	         HeapSize: 0,
	         Image: "NA",
	         Manufacturer: "LG",
	         Memory: 17179869184,
	         Model: "G2",
	         Name: "LG G2 (Sprint)",
	         Os: "4.2.2",
	         Platform: "ANDROID",
	         Resolution: {
	           Height: 1920,
	           Width: 1080
	         }
	       }
	*/

}

/* List all uploads */
func listUploads(svc *devicefarm.DeviceFarm, projectArn string) {

	listReq := &devicefarm.ListUploadsInput{
		ARN: aws.String(projectArn),
	}

	resp, err := svc.ListUploads(listReq)

	failOnErr(err, "error listing uploads")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all runs */
func listRuns(svc *devicefarm.DeviceFarm, projectArn string) {

	listReq := &devicefarm.ListRunsInput{
		ARN: aws.String(projectArn),
	}

	resp, err := svc.ListRuns(listReq)

	failOnErr(err, "error listing runs")
	//fmt.Println(awsutil.Prettify(resp))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Platform", "Type", "Result", "Status", "Date", "ARN"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, m := range resp.Runs {
		line := []string{*m.Name, *m.Platform, *m.Type, *m.Result, *m.Status, time.Time.String(*m.Created), *m.ARN}
		table.Append(line)
	}
	table.Render() // Send output

}

/* List all tests */
func listTests(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListTestsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListTests(listReq)

	failOnErr(err, "error listing tests")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all unique problems */
func listUniqueProblems(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListUniqueProblemsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListUniqueProblems(listReq)

	failOnErr(err, "error listing problems")
	fmt.Println(awsutil.Prettify(resp))
}

/* List suites */
func listSuites(svc *devicefarm.DeviceFarm, filterArn string) {

	listReq := &devicefarm.ListSuitesInput{
		ARN: aws.String(filterArn),
	}

	resp, err := svc.ListSuites(listReq)

	failOnErr(err, "error listing suites")
	//fmt.Println(awsutil.Prettify(resp))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Result", "Message"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(120)

	for _, m := range resp.Suites {
		line := []string{*m.Name, *m.Status, *m.Result, *m.Message}
		table.Append(line)
	}
	table.Render() // Send output

}

func guessAppType(fileName string) (appType string, err error) {

	lowerCaseFileName := strings.ToLower(fileName)

	extension := filepath.Ext(lowerCaseFileName)

	if extension == ".apk" {
		return "ANDROID_APP", nil
	}

	if extension == ".ipa" {
		return "IOS_APP", nil
	}

	return "", errors.New("Can't guess App Type")

}

func lookupTestPackageType(testType string) (testPackageType string, err error) {

	if testType == "APPIUM_JAVA_JUNIT" {
		return "APPIUM_JAVA_JUNIT_TEST_PACKAGE", nil
	}

	if testType == "INSTRUMENTATION" {
		return "INSTRUMENTATION_TEST_PACKAGE", nil
	}

	if testType == "UIAUTOMATION" {
		return "UIAUTOMATION_TEST_PACKAGE", nil
	}

	if testType == "APPIUM_JAVA_TESTNG" {
		return "APPIUM_JAVA_TESTNG_TEST_PACKAGE", nil
	}

	if testType == "CALABASH" {
		return "CALABASH_TEST_PACKAGE", nil
	}

	if testType == "UIAUTOMATER" {
		return "UIAUTOMATER_TEST_PACKAGE", nil
	}

	if testType == "XCTEST" {
		return "XCTEST_TEST_PACKAGE", nil
	}

	// BUILTIN_EXPLORER: For Android, an app explorer that will traverse an Android app, interacting with it and capturing screenshots at the same time.
	// BUILTIN_FUZZ: The built-in fuzz type.
	return "", errors.New("Could not guess test type, you can use the BUILTIN_FUZZ or BUILTIN_EXPLORER")

}

/* Schedule Run */
func scheduleRun(svc *devicefarm.DeviceFarm, projectArn string, runName string, deviceArn string, devicePoolArn string, appArn string, appFile string, appType string, testPackageArn string, testPackageFile string, testType string) {

	// Upload the app file if there is one
	if appFile != "" {
		uploadApp, err := uploadPut(svc, appFile, appType, projectArn, "")

		// Try to guess the upload type based on the filename
		if appType == "" {
			appType, err = guessAppType(appFile)
		}

		failOnErr(err, "error scheduling run")
		appArn = *uploadApp.ARN
	}

	// Try to guess the upload type based on the filename
	/*
		if testType == "" {
			testType, err := guessTestType(testPackageFile)
		}
	*/

	testPackageType, err := lookupTestPackageType(testType)

	// Upload the testPackage file if there is one
	if testPackageFile != "" {

		uploadTestPackage, err := uploadPut(svc, testPackageFile, testPackageType, projectArn, "")
		failOnErr(err, "error scheduling run")
		testPackageArn = *uploadTestPackage.ARN
	}

	runTest := &devicefarm.ScheduleRunTest{
		Type:           aws.String(testType),
		TestPackageARN: aws.String(testPackageArn),
		//Parameters: // test parameters
		//Filter: // filter to pass to tests
	}

	fmt.Println(appArn)
	fmt.Println(devicePoolArn)
	fmt.Println(runName)
	fmt.Println(testPackageArn)
	fmt.Println(testPackageType)
	fmt.Println(projectArn)
	/*
		if testPackageArn != "" {
			runTest.TestPackageARN = aws.String(testPackageArn)
		}
	*/

	// Wait for the uploads to be processed
	// Better is to check the upload status PROCESSING
	time.Sleep(4 * time.Second)

	runReq := &devicefarm.ScheduleRunInput{
		AppARN:        aws.String(appArn),
		DevicePoolARN: aws.String(devicePoolArn),
		Name:          aws.String(runName),
		ProjectARN:    aws.String(projectArn),
		Test:          runTest,
	}

	fmt.Println(awsutil.Prettify(runReq))

	resp, err := svc.ScheduleRun(runReq)

	failOnErr(err, "error scheduling run")
	fmt.Println(awsutil.Prettify(resp))

	// Now we wait for the run status to go COMPLETED
	runArn := *resp.Run.ARN

	status := ""
	for status != "COMPLETED" {
		time.Sleep(4 * time.Second)
		infoReq := &devicefarm.GetRunInput{
			ARN: aws.String(runArn),
		}

		fmt.Print(".")
		resp, err := svc.GetRun(infoReq)
		failOnErr(err, "error scheduling run")
		status = *resp.Run.Status
	}

	// Generate report
	runReport(svc, runArn)

}

/* List Artifacts */

func listArtifacts(svc *devicefarm.DeviceFarm, filterArn string, artifactType string) {

	fmt.Println(filterArn)

	listReq := &devicefarm.ListArtifactsInput{
		ARN: aws.String(filterArn),
	}

	listReq.Type = aws.String("LOG")
	resp, err := svc.ListArtifacts(listReq)
	failOnErr(err, "error listing artifacts")
	fmt.Println(awsutil.Prettify(resp))

	listReq.Type = aws.String("SCREENSHOT")
	resp, err = svc.ListArtifacts(listReq)
	failOnErr(err, "error listing artifacts")

	fmt.Println(awsutil.Prettify(resp))

	listReq.Type = aws.String("FILE")
	resp, err = svc.ListArtifacts(listReq)
	failOnErr(err, "error listing artifacts")

	fmt.Println(awsutil.Prettify(resp))
}

/* Download Artifacts */
func downloadArtifacts(svc *devicefarm.DeviceFarm, filterArn string, artifactType string) {

	fmt.Println(filterArn)

	listReq := &devicefarm.ListArtifactsInput{
		ARN: aws.String(filterArn),
	}

	types := []string{"LOG", "SCREENSHOT", "FILE"}

	for _, each := range types {
		listReq.Type = aws.String(each)

		resp, err := svc.ListArtifacts(listReq)
		failOnErr(err, "error listing artifacts")

		for index, artifact := range resp.Artifacts {
			fileName := fmt.Sprintf("report/%d-%s.%s", index, *artifact.Name, *artifact.Extension)
			downloadArtifact(fileName, artifact)
		}
	}

}

func downloadArtifact(fileName string, artifact *devicefarm.Artifact) {

	url := *artifact.URL

	dirName := path.Dir(fileName)
	err := os.MkdirAll(dirName, 0777)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("Downloading [%s] -> [%s]\n", url, fileName)

	downloadURL(url, fileName)
}

func downloadURL(url string, fileName string) {

	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(url) // add a filter to check redirect

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s with %v bytes downloaded", fileName, size)

}

/* List Jobs */
func listJobs(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListJobsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListJobs(listReq)

	failOnErr(err, "error listing jobs")
	fmt.Println(awsutil.Prettify(resp))
}

/* Create an upload */
func uploadCreate(svc *devicefarm.DeviceFarm, uploadName string, uploadType string, projectArn string) {

	uploadReq := &devicefarm.CreateUploadInput{
		Name:       aws.String(uploadName),
		ProjectARN: aws.String(projectArn),
		Type:       aws.String(uploadType),
	}

	resp, err := svc.CreateUpload(uploadReq)

	failOnErr(err, "error creating upload")
	fmt.Println(awsutil.Prettify(resp))
}

/* Get Run Info */
func runInfo(svc *devicefarm.DeviceFarm, runArn string) {

	infoReq := &devicefarm.GetRunInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.GetRun(infoReq)

	failOnErr(err, "error getting run info")
	fmt.Println(awsutil.Prettify(resp))
}

/* Get Run Report */
func runReport(svc *devicefarm.DeviceFarm, runArn string) {

	infoReq := &devicefarm.GetRunInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.GetRun(infoReq)

	failOnErr(err, "error getting run info")
	fmt.Printf("%s\n", *resp.Run.Name)
	//fmt.Println(awsutil.Prettify(resp))

	jobReq := &devicefarm.ListJobsInput{
		ARN: aws.String(runArn),
	}

	// Find all artifacts
	artifactReq := &devicefarm.ListArtifactsInput{
		ARN: aws.String(runArn),
	}

	types := []string{"LOG", "SCREENSHOT", "FILE"}
	artifacts := map[string][]devicefarm.ListArtifactsOutput{}

	for _, artifactType := range types {

		artifactReq.Type = aws.String(artifactType)

		artifactResp, err := svc.ListArtifacts(artifactReq)
		failOnErr(err, "error getting run info")

		// Store type artifacts
		artifacts[artifactType] = append(artifacts[artifactType], *artifactResp)
		/*
			for _, artifact := range artifactResp.Artifacts {
				fmt.Println(awsutil.Prettify(artifact))
			}
		*/
	}

	respJob, err := svc.ListJobs(jobReq)
	failOnErr(err, "error getting jobs")

	for _, job := range respJob.Jobs {

		fmt.Println("==========================================")
		time.Sleep(2 * time.Second)

		jobFriendlyName := fmt.Sprintf("%s - %s - %s", *job.Name, *job.Device.Model, *job.Device.Os)

		//fmt.Println(awsutil.Prettify(job))

		suiteReq := &devicefarm.ListSuitesInput{
			ARN: aws.String(*job.ARN),
		}
		suiteResp, err := svc.ListSuites(suiteReq)
		failOnErr(err, "error getting run info")

		for _, suite := range suiteResp.Suites {
			message := ""
			if suite.Message != nil {
				message = *suite.Message
			}

			fmt.Printf("%s -> %s : %s \n----> %s\n", jobFriendlyName, *suite.Name, message, *suite.ARN)
			dirPrefix := fmt.Sprintf("report/%s/%s/", jobFriendlyName, *suite.Name)
			downloadArtifactsForSuite(dirPrefix, artifacts, *suite)
		}

		//fmt.Println(awsutil.Prettify(suiteResp))
	}

}

func downloadArtifactsForSuite(dirPrefix string, allArtifacts map[string][]devicefarm.ListArtifactsOutput, suite devicefarm.Suite) {
	suiteArn := *suite.ARN
	artifactTypes := []string{"LOG", "SCREENSHOT", "FILE"}

	r := strings.NewReplacer(":suite:", ":artifact:")
	artifactPrefix := r.Replace(suiteArn)

	for _, artifactType := range artifactTypes {
		typedArtifacts := allArtifacts[artifactType]
		for _, artifactList := range typedArtifacts {
			count := 0
			for _, artifact := range artifactList.Artifacts {
				if strings.HasPrefix(*artifact.ARN, artifactPrefix) {
					fmt.Printf("[%s] %s.%s\n", artifactType, *artifact.Name, *artifact.Extension)
					//pathFull := strings.Split(suiteArn, ":")[6]
					//pathSuffix := strings.Split(pathFull, "/")
					//runId := pathSuffix[0]
					//jobId := pathSuffix[1]
					//suiteId := pathSuffix[2]
					//artifactId := pathSuffix[3]
					fileName := fmt.Sprintf("%s/%d_%s.%s", dirPrefix, count, *artifact.Name, *artifact.Extension)
					//fileName := fmt.Sprintf("%s/%s/%s/%s.%s", dirPrefix, suiteId, artifactId, *artifact.Name, *artifact.Extension)
					fmt.Printf("[%s] %s\n%s\n", artifactType, fileName, *artifact.URL)
					downloadArtifact(fileName, artifact)
					count++
				}
			}
		}
	}

}

/* Get Run Status */
func runStatus(svc *devicefarm.DeviceFarm, runArn string) {

	infoReq := &devicefarm.GetRunInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.GetRun(infoReq)

	failOnErr(err, "error getting run info")
	fmt.Println(*resp.Run.Status)
}

/* Get Job Info */
func jobInfo(svc *devicefarm.DeviceFarm, jobArn string) {

	infoReq := &devicefarm.GetJobInput{
		ARN: aws.String(jobArn),
	}

	resp, err := svc.GetJob(infoReq)

	failOnErr(err, "error getting job info")
	fmt.Println(awsutil.Prettify(resp))
}

/* Get Suite Info */
func suiteInfo(svc *devicefarm.DeviceFarm, suiteArn string) {

	infoReq := &devicefarm.GetJobInput{
		ARN: aws.String(suiteArn),
	}

	resp, err := svc.GetJob(infoReq)

	failOnErr(err, "error getting suite info")
	fmt.Println(awsutil.Prettify(resp))
}

/* Get Upload Info */
func uploadInfo(svc *devicefarm.DeviceFarm, uploadArn string) {

	uploadReq := &devicefarm.GetUploadInput{
		ARN: aws.String(uploadArn),
	}

	resp, err := svc.GetUpload(uploadReq)

	failOnErr(err, "error getting upload info")
	fmt.Println(awsutil.Prettify(resp))
}

/* Upload a file */
func uploadPut(svc *devicefarm.DeviceFarm, uploadFilePath string, uploadType string, projectArn string, uploadName string) (upload *devicefarm.Upload, err error) {

	// Read File
	file, err := os.Open(uploadFilePath)

	if err != nil {
		return nil, err
		fmt.Println(err)
	}

	defer file.Close()

	// Get file size
	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()

	// read file content to buffer
	buffer := make([]byte, fileSize)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type

	// Prepare upload
	if uploadName == "" {
		uploadName = path.Base(uploadFilePath)
	}

	uploadReq := &devicefarm.CreateUploadInput{
		Name:        aws.String(uploadName),
		ProjectARN:  aws.String(projectArn),
		Type:        aws.String(uploadType),
		ContentType: aws.String("application/octet-stream"),
	}

	uploadResp, err := svc.CreateUpload(uploadReq)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println(awsutil.Prettify(uploadResp))

	uploadInfo := uploadResp.Upload

	upload_url := *uploadInfo.URL

	fmt.Println(upload_url)

	req, err := http.NewRequest("PUT", upload_url, fileBytes)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Remove Host and split to get [0] = path & [1] = querystring
	strippedUrl := strings.Split(strings.Replace(upload_url, "https://prod-us-west-2-uploads.s3-us-west-2.amazonaws.com/", "/", -1), "?")
	req.URL.Opaque = strippedUrl[0]
	req.URL.RawQuery = strippedUrl[1]

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.FormatInt(fileSize, 10))

	// Debug Request to AWS
	debug(httputil.DumpRequestOut(req, false))

	client := &http.Client{}

	res, err := client.Do(req)

	dump, _ := httputil.DumpResponse(res, true)
	log.Printf("} -> %s\n", dump)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()

	return uploadResp.Upload, nil
}

/*
 * Helper page to exit on error with a nice message
 */
func failOnErr(err error, reason string) {
	if err != nil {
		log.Fatal("Failed: %s, %s\n\n", reason, err)
		os.Exit(-1)
	}

	return
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
