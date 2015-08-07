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
	"strconv"
	"strings"
)

func main() {

	app := cli.NewApp()
	app.Name = "devicefarm-cli"
	app.Usage = "allows you to interact with AWS devicefarm from the command line"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "projects",
			Usage: "manage the projects",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the projects", // of an account
					Action: func(c *cli.Context) {
						listProjects()
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
					Action: func(c *cli.Context) {
						listArtifacts()
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
					Action: func(c *cli.Context) {
						listDevicePools()
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
						listDevices()
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
					Action: func(c *cli.Context) {
						listJobs()
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
					Action: func(c *cli.Context) {
						listRuns()
					},
				},
				{
					Name:  "schedule",
					Usage: "schedule a run",
					Action: func(c *cli.Context) {
						scheduleRun()
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
						listSuites()
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
						listTests()
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
						listUniqueProblems()
					},
				},
			},
		},
		{
			Name:  "upload",
			Usage: "manages the uploads",
			Subcommands: []cli.Command{
				{
					Name:  "ipa",
					Usage: "uploads an ipa",
					Action: func(c *cli.Context) {
						uploadCreate()
					},
				},
				{
					Name:  "put",
					Usage: "uploads an ipa 2",
					Action: func(c *cli.Context) {
						uploadPut()
					},
				},
				{
					Name:  "list",
					Usage: "lists all uploads", // of a Project
					Action: func(c *cli.Context) {
						listUploads()
					},
				},
				{
					Name:  "info",
					Usage: "info about uploads",
					Action: func(c *cli.Context) {
						uploadInfo()
					},
				},
			},
		},
	}

	app.Run(os.Args)

}

func listProjects() {

	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	resp, err := svc.ListProjects(nil)
	failOnErr(err, "error listing projects")

	fmt.Println(awsutil.Prettify(resp))
}

func listDevicePools() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})

	// CURATED: A device pool that is created and managed by AWS Device Farm.
	// PRIVATE: A device pool that is created and managed by the device pool developer.

	pool := &devicefarm.ListDevicePoolsInput{
		ARN: aws.String(projectArn),
	}
	resp, err := svc.ListDevicePools(pool)

	failOnErr(err, "error listing device pools")
	fmt.Println(awsutil.Prettify(resp))
}

func listDevices() {

	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})

	input := &devicefarm.ListDevicesInput{}
	resp, err := svc.ListDevices(input)
	failOnErr(err, "error listing devices")

	fmt.Println(awsutil.Prettify(resp))
}

func listUploads() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	listReq := &devicefarm.ListUploadsInput{
		ARN: aws.String(projectArn),
	}

	resp, err := svc.ListUploads(listReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func listRuns() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	listReq := &devicefarm.ListRunsInput{
		ARN: aws.String(projectArn),
	}

	resp, err := svc.ListRuns(listReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func listTests() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	// Project -> Runs -> Tests
	// Runs ARN
	arn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	listReq := &devicefarm.ListTestsInput{
		ARN: aws.String(arn),
	}

	resp, err := svc.ListTests(listReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func listUniqueProblems() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	// Project -> Runs -> Tests
	// Runs ARN
	runArn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})

	listReq := &devicefarm.ListUniqueProblemsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListUniqueProblems(listReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func listSuites() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	// Project -> Runs -> Tests
	// Runs ARN
	runArn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	listReq := &devicefarm.ListSuitesInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListSuites(listReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func scheduleRun() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
	appUploadArn := "arn:aws:devicefarm:us-west-2:110440800955:upload:f7952cc6-5833-47f3-afef-c149fb4e7c76/dbbb0b81-5c53-42b1-b1b8-e6239124ca3a"
	devicePoolArn := "arn:aws:devicefarm:us-west-2:110440800955:devicepool:f7952cc6-5833-47f3-afef-c149fb4e7c76/b51ed696-ab6a-440b-8de0-c947ba442d53"
	//testUploadArn :=

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

	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	runReq := &devicefarm.ScheduleRunInput{
		// Documentation is pretty horrible , it says AppArn
		AppARN:        aws.String(appUploadArn),
		DevicePoolARN: aws.String(devicePoolArn),
		Name:          aws.String("test me - w00t"),
		ProjectARN:    aws.String(projectArn),
		Test: &devicefarm.ScheduleRunTest{
			Type: aws.String("BUILTIN_FUZZ"),

			//TestPackageArn: aws.String(testUploadArn)
			//Parameters: // test parameters
			//Filter: // filter to pass to tests
		},
	}

	resp, err := svc.ScheduleRun(runReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func listArtifacts() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	// Runs ARN
	runArn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	listReq := &devicefarm.ListArtifactsInput{
		ARN: aws.String(runArn),
		//Type: aws.String("LOG"),
		//Type: aws.String("FILE"),
		Type: aws.String("SCREENSHOT"),
		// ?? Unclear aws doc
		// https://github.com/aws/aws-sdk-go/blob/master/apis/devicefarm/2015-06-23/api-2.json
		//Name:      aws.String("Calabash JSON Output"),
		//Extension: aws.String("json"),
		//  Extension: ".png", -> wierdos!
		// Name: "i_take_a_screenshot_0",
	}

	resp, err := svc.ListArtifacts(listReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func listJobs() {
	//failOnErr(err, "error in test call with key"+string(aws_key))

	// Runs ARN
	runArn := "arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	listReq := &devicefarm.ListJobsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListJobs(listReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.Prettify(resp))
}

func uploadCreate() {

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

	name := "test-upload"
	uploadType := "IOS_APP"
	projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	uploadReq := &devicefarm.CreateUploadInput{
		Name:       aws.String(name),
		ProjectARN: aws.String(projectArn),
		Type:       aws.String(uploadType),
	}

	resp, err := svc.CreateUpload(uploadReq)

	if err != nil {
		panic(err)
	}

	fmt.Println(awsutil.Prettify(resp))
}

func uploadInfo() {

	uploadArn := "arn:aws:devicefarm:us-west-2:110440800955:upload:f7952cc6-5833-47f3-afef-c149fb4e7c76/1d1a6d6e-554d-48d1-b53f-21f80ef94a14"
	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})
	uploadReq := &devicefarm.GetUploadInput{
		ARN: aws.String(uploadArn),
	}

	resp, err := svc.GetUpload(uploadReq)

	if err != nil {
		panic(err)
	}

	fmt.Println(awsutil.Prettify(resp))

}

func uploadPut() {

	fileToUpload := "a.ipa"
	fileToUploadBasename := "a.ipa"
	uploadType := "IOS_APP"

	projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"

	file, err := os.Open(fileToUploadBasename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()

	buffer := make([]byte, fileSize)

	// read file content to buffer
	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type

	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})

	uploadReq := &devicefarm.CreateUploadInput{
		Name:        aws.String(fileToUpload),
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

	fmt.Println("")
	fmt.Println("")
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

// Not used yet, but who knows we need to get correct results
func amazonEscape(s string) string {
	hexCount := 0

	for i := 0; i < len(s); i++ {
		if amazonShouldEscape(s[i]) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		if c := s[i]; amazonShouldEscape(c) {
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		} else {
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func amazonShouldEscape(c byte) bool {
	return !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') || c == '_' || c == '-' || c == '~' || c == '.' || c == '/' || c == ':')
}
