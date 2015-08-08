package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strconv"
	"strings"
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
			Name:  "projects",
			Usage: "manage the projects",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the projects", // of an account
					Action: func(c *cli.Context) {
						fmt.Println(c.Args())
						listProjects(svc)
					},
				},
			},
		},
		{
			Name:  "artifacts",
			Usage: "manage the artifacts",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the artifacts", // of a test
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
							Value:  "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")

						listArtifacts(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "devicepools",
			Usage: "manage the device pools",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the devicepools", //globally
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
							Value:  "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76",
						},
					},
					Action: func(c *cli.Context) {

						projectArn := c.String("project")
						listDevicePools(svc, projectArn)
					},
				},
			},
		},
		{
			Name:  "devices",
			Usage: "manage the devices",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the devices", // globally
					Action: func(c *cli.Context) {
						listDevices(svc)
					},
				},
			},
		},
		{
			Name:  "jobs",
			Usage: "manage the jobs",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the jobs", // of a test
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
							Value:  "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")

						listJobs(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "runs",
			Usage: "manage the runs",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the runs",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
							Value:  "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76",
						},
					},
					Action: func(c *cli.Context) {
						projectArn := c.String("project")
						listRuns(svc, projectArn)
					},
				},
				{
					Name:  "schedule",
					Usage: "schedule a run",
					Action: func(c *cli.Context) {
						projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
						appUploadArn := "arn:aws:devicefarm:us-west-2:110440800955:upload:f7952cc6-5833-47f3-afef-c149fb4e7c76/dbbb0b81-5c53-42b1-b1b8-e6239124ca3a"
						devicePoolArn := "arn:aws:devicefarm:us-west-2:110440800955:devicepool:f7952cc6-5833-47f3-afef-c149fb4e7c76/b51ed696-ab6a-440b-8de0-c947ba442d53"
						testUploadArn := ""
						testType := "BUILTIN_FUZZ"

						/*
							BUILTIN_FUZZ: The built-in fuzz type.
							BUILTIN_EXPLORER: For Android, an app explorer that will traverse an Android app, interacting with it and capturing screenshots at the same time.
							APPIUM_JAVA_JUNIT: The Appium Java JUnit type.
							APPIUM_JAVA_TESTNG: The Appium Java TestNG type.
							CALABASH: The Calabash type.
							INSTRUMENTATION: The Instrumentation type.
							UIAUTOMATION: The uiautomation type.
							UIAUTOMATOR: The uiautomator type.
							XCTEST: The XCode test type.
						*/
						scheduleRun(svc, projectArn, appUploadArn, devicePoolArn, testUploadArn, testType)
					},
				},
			},
		},
		{
			Name:  "samples",
			Usage: "manage the samples",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the samples",
					Action: func(c *cli.Context) {
						//listSamples()
					},
				},
			},
		},
		{
			Name:  "suites",
			Usage: "manage the suites",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the suites",
					Action: func(c *cli.Context) {
						runArn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
						listSuites(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "tests",
			Usage: "manage the tests",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the tests", // of a Run
					Action: func(c *cli.Context) {
						runArn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
						listTests(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "problems",
			Usage: "manage the problems",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the problems", // of Test
					Action: func(c *cli.Context) {
						runArn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
						listUniqueProblems(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "upload",
			Usage: "manages the uploads",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "creates an upload",
					Action: func(c *cli.Context) {
						/*
							ANDROID_APP: An Android upload.
							IOS_APP: An iOS upload.
							EXTERNAL_DATA: An external data upload.
							APPIUM_JAVA_JUNIT_TEST_PACKAGE: An Appium Java JUnit test package upload.
							APPIUM_JAVA_TESTNG_TEST_PACKAGE: An Appium Java TestNG test package upload.
							CALABASH_TEST_PACKAGE: A Calabash test package upload.
							INSTRUMENTATION_TEST_PACKAGE: An instrumentation upload.
							UIAUTOMATOR_TEST_PACKAGE: A uiautomator test package upload.
							XCTEST_TEST_PACKAGE: An XCode test package upload.
						*/

						uploadName := "test-upload"
						uploadType := "IOS_APP"
						projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
						uploadCreate(svc, uploadName, uploadType, projectArn)
					},
				},
				{
					Name:  "file",
					Usage: "uploads an file",
					Action: func(c *cli.Context) {
						uploadType := "IOS_APP"
						projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
						uploadFilePath := "a.ipa"
						uploadPut(svc, uploadFilePath, uploadType, projectArn)
					},
				},
				{
					Name:  "list",
					Usage: "lists all uploads", // of a Project
					Action: func(c *cli.Context) {
						projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
						listUploads(svc, projectArn)
					},
				},
				{
					Name:  "info",
					Usage: "info about uploads",
					Action: func(c *cli.Context) {
						uploadArn := "arn:aws:devicefarm:us-west-2:110440800955:upload:f7952cc6-5833-47f3-afef-c149fb4e7c76/1d1a6d6e-554d-48d1-b53f-21f80ef94a14"
						uploadInfo(svc, uploadArn)
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

	fmt.Println(awsutil.Prettify(resp))
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
	fmt.Println(awsutil.Prettify(resp))
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
	fmt.Println(awsutil.Prettify(resp))
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
func listSuites(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListSuitesInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListSuites(listReq)

	failOnErr(err, "error listing suites")
	fmt.Println(awsutil.Prettify(resp))
}

/* Schedule Run */
func scheduleRun(svc *devicefarm.DeviceFarm, projectArn string, appUploadArn string, devicePoolArn string, testUploadArn string, testType string) {

	runReq := &devicefarm.ScheduleRunInput{
		// Documentation is pretty horrible , it says AppArn
		AppARN:        aws.String(appUploadArn),
		DevicePoolARN: aws.String(devicePoolArn),
		Name:          aws.String("test me - w00t"),
		ProjectARN:    aws.String(projectArn),
		Test: &devicefarm.ScheduleRunTest{
			Type: aws.String(testType),

			//TestPackageArn: aws.String(testUploadArn)
			//Parameters: // test parameters
			//Filter: // filter to pass to tests
		},
	}

	resp, err := svc.ScheduleRun(runReq)

	failOnErr(err, "error scheduling run")
	fmt.Println(awsutil.Prettify(resp))
}

/* List Artifacts */

// ?? Unclear aws doc
// https://github.com/aws/aws-sdk-go/blob/master/apis/devicefarm/2015-06-23/api-2.json
//Name:      aws.String("Calabash JSON Output"),
//Extension: aws.String("json"),
//  Extension: ".png", -> wierdos!
// Name: "i_take_a_screenshot_0",
func listArtifacts(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListArtifactsInput{
		ARN: aws.String(runArn),
		//Type: aws.String("LOG"),
		//Type: aws.String("FILE"),
		Type: aws.String("SCREENSHOT"),
	}

	resp, err := svc.ListArtifacts(listReq)
	failOnErr(err, "error listing artifacts")

	fmt.Println(awsutil.Prettify(resp))
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
func uploadPut(svc *devicefarm.DeviceFarm, uploadFilePath string, uploadType string, projectArn string) {

	// Read File
	file, err := os.Open(uploadFilePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
	uploadFileBasename := path.Base(uploadFilePath)
	uploadReq := &devicefarm.CreateUploadInput{
		Name:        aws.String(uploadFileBasename),
		ProjectARN:  aws.String(projectArn),
		Type:        aws.String(uploadType),
		ContentType: aws.String("application/octet-stream"),
	}

	resp, err := svc.CreateUpload(uploadReq)
	fmt.Println(awsutil.Prettify(resp))

	uploadInfo := resp.Upload

	upload_url := *uploadInfo.URL

	fmt.Println(upload_url)

	req, err := http.NewRequest("PUT", upload_url, fileBytes)

	if err != nil {
		log.Fatal(err)
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
	}
	defer res.Body.Close()

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
