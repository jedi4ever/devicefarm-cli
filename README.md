# Description

> A Command-Line interface to AWS Devicefarm , because humans don't use ARNs. Written in golang for maximum deployability

** If you find this app useful, consider sponsoring the AWS costs for doing test runs, devicefarm ain't cheapo (open a ticket and let me know) **

Current State: 
- usable but no good err checking and output

Things I want to cover:
- have correct exit codes
- support a config file & environment settings & cli options
- have formatters like json, junit etc..
- poll for test results
- provided binaries & packages & docker instance
- make it slack friendly for reports

Cleanup will happen soon!

# Usage

## Configuring access to AWS
- You'll need to sign up with AWS and create an account.
- This app uses the standard AWS tool configuration. It either reads from ~/.aws/config.

Set:
```
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
```

## Listing projects
Note: currently projects need to be create via the console

```
$ ./devicefarm-cli list projects
+----------+-------------------------------+----------------------------------------------------------------------------------------+
|   NAME   |            CREATED            |                                          ARN                                           |
+----------+-------------------------------+----------------------------------------------------------------------------------------+
| samplejr | 2015-07-29 15:38:10 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76 |
+----------+-------------------------------+----------------------------------------------------------------------------------------+
```

## Listing devices
```
$ ./devicefarm-cli list devices
+---------------------------------------+-------+----------+--------+-----------------------------------------------------------------------+
|                 NAME                  |  OS   | PLATFORM |  FORM  |                                  ARN                                  |
+---------------------------------------+-------+----------+--------+-----------------------------------------------------------------------+
| Apple iPad Mini 2                     | 7.1.2 | IOS      | TABLET | arn:aws:devicefarm:us-west-2::device:3B33A0062E6D47B4A50437C48F9141F6 |
| Apple iPod Touch 5th Gen              | 8.1.2 | IOS      | TABLET | arn:aws:devicefarm:us-west-2::device:8BF08BA1178D4DF7B46AB882A62FFD68 |
| Samsung Galaxy Note II (AT&T)         | 4.4.2 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:A06AFA2C23C24230A74E69E14E790659 |
| Motorola DROID Ultra (Verizon)        | 4.4.4 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:B6B524CF9BF84CA891FFEF1C88E9A279 |
| Samsung Galaxy S4 Active (AT&T)       | 4.4.2 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:577DC08D6B964346B86610CFF090CD59 |
| Apple iPhone 4S                       | 8.1.2 | IOS      | PHONE  | arn:aws:devicefarm:us-west-2::device:F7508F91AE044EE0843985B190679B62 |
| Samsung Galaxy S5 Active (AT&T)       | 4.4.2 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:9710D509338C4639ADEFC5D6E99F45E6 |
| Apple iPhone 5c                       | 8.0.2 | IOS      | PHONE  | arn:aws:devicefarm:us-west-2::device:6A9072D1240B4AE09213819618C58BB4 |
| Samsung Galaxy Note (AT&T)            | 4.1.2 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:CFBBF54B2A724C6791D5F267B580C05B |
| Apple iPad Mini 1st Gen               | 7.1.2 | IOS      | TABLET | arn:aws:devicefarm:us-west-2::device:306ABA42C96044ED9AC3EE8684B56C54 |
| ASUS Nexus 7 - 2nd Gen (WiFi)         | 5.0.2 | ANDROID  | TABLET | arn:aws:devicefarm:us-west-2::device:0ACA52D055B34FB18F481AFE2F1A4661 |
| Apple iPhone 6 Plus                   | 8.3   | IOS      | PHONE  | arn:aws:devicefarm:us-west-2::device:5FD56B0CDB324C2E84FAAF04126CAC79 |
| Apple iPhone 6                        | 8.3   | IOS      | PHONE  | arn:aws:devicefarm:us-west-2::device:02FC89B471F4439CA76CADA591DE9867 |
| Samsung Galaxy S4 (AT&T)              | 4.2.2 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:9A8C3412A5E64CC0900331A317B62B93 |
| Sony Xperia Z3 Compact (GSM)          | 4.4.4 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:133050132D324AA8B264180506232FC1 |
| HTC One M8 (Sprint)                   | 4.4.4 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:19692012F8434F3BB06F701276FAF6F9 |
| LG Optimus Fuel (TracFone)            | 4.4   | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:2C1D29BF2F704280A1D4CCB76DE90D8F |
| Apple iPhone 6 Plus                   | 8.1   | IOS      | PHONE  | arn:aws:devicefarm:us-west-2::device:3993D9928C604961B0A0364A35AFD4AE |
| ASUS Nexus 7 - 2nd Gen (WiFi)         | 4.3.1 | ANDROID  | TABLET | arn:aws:devicefarm:us-west-2::device:208FE7EE973042EA97DEC2EEF31CD10A |
| Samsung Galaxy S5 (Verizon)           | 4.4.4 | ANDROID  | PHONE  | arn:aws:devicefarm:us-west-2::device:C30737D1E582482C9D06BC4878E7F795 |
......
```

## Listing runs
```
$ ./devicefarm-cli list runs
+----------------------------------+-------------+--------------+---------+-----------+-------------------------------+-------------------------------------------------------------------------------------------------------------------------+
|               NAME               |  PLATFORM   |     TYPE     | RESULT  |  STATUS   |             DATE              |                                                           ARN                                                           |
+----------------------------------+-------------+--------------+---------+-----------+-------------------------------+-------------------------------------------------------------------------------------------------------------------------+
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-10 19:47:07 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/06d07b26-b547-443e-a55d-a5a9f370e427 |
| Android rulez                    | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-10 18:34:48 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/3226f872-d7f8-49ab-98b1-3a1755bea233 |
| Android rulez                    | ANDROID_APP | CALABASH     | PASSED  | COMPLETED | 2015-08-10 16:00:50 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/d8d7b4e1-eb92-4cf1-91a0-751e21ee3034 |
| test me - w00t                   | ANDROID_APP | BUILTIN_FUZZ | ERRORED | COMPLETED | 2015-08-07 13:01:10 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/d6db8610-9d94-4f1b-819d-dec4cc663857 |
| samplejr-appstore-1.0.1-325.ipa  | IOS_APP     | BUILTIN_FUZZ | ERRORED | COMPLETED | 2015-08-04 19:50:11 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/da290ade-6df3-40be-8de1-4324fa9f89f6 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | ERRORED | COMPLETED | 2015-08-01 18:23:34 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/32313b78-95c5-461c-adc7-9d95d2f49755 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | PASSED  | COMPLETED | 2015-08-01 17:59:38 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/6959f092-8fda-452b-83df-c12ba222cf02 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | PASSED  | COMPLETED | 2015-08-01 17:51:27 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/291dc007-4226-4e15-bf07-8d463c11287d |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | PASSED  | COMPLETED | 2015-08-01 17:45:28 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/c43f4a7a-049b-4e08-8966-e39531ef2eb3 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 15:39:57 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/2c5ecc4b-622a-4a19-9f2f-677adab32a3b |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 15:30:36 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/62eb307a-4ed2-4acc-b6da-29244cbb6599 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 15:25:04 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/9f641918-b68a-49dc-bc35-856797695744 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | ERRORED | COMPLETED | 2015-08-01 15:23:05 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/bee2e870-8b7d-4481-9941-18d19b8aae66 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 15:12:00 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/29fe49fc-5c30-4f83-b569-d341e7966eb8 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 14:56:04 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/809fd10a-30d3-42a6-8ea8-b5ccec0fbbb3 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 14:42:23 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/5781eae8-a980-42f7-8ce1-a57d9ff07435 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 14:30:57 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/3480c752-9c60-4728-83d3-640afb48b4ab |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | ERRORED | COMPLETED | 2015-08-01 14:25:02 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/e870a56b-8167-4f05-b7bc-cc368e6d14f5 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 13:55:44 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/3da475d1-f463-431a-8ce2-14a07ed17151 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 13:25:45 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/85b479f6-7770-4286-90d1-46b9c884d6aa |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 13:00:22 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/9f881042-0b05-413a-8e0e-7390206cb88d |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 12:54:49 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/16769413-8e95-4aae-97e5-44bce75dd188 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 12:34:54 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/8635bdad-6e2a-4964-8626-b2c8f80d7f72 |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-08-01 12:19:07 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/8b69c0d5-cfbe-4977-a483-95c490ea8bdc |
| app-samplejr-staging-release.apk | ANDROID_APP | CALABASH     | FAILED  | COMPLETED | 2015-07-30 09:37:05 +0000 UTC | arn:aws:devicefarm:us-west-2:110440800955:run:f7952cc6-5833-47f3-afef-c149fb4e7c76/d1c72f3c-294e-45de-9e91-e6f840b8c06d |
+----------------------------------+-------------+--------------+---------+-----------+-------------------------------+-------------------------------------------------------------------------------------------------------------------------+
```

## Schedule run
- To schedule a run use the following syntax (soon it will even be simpler).
- You can also set params through environment variables with `DF_` prefix

```
$ ./devicefarm-cli schedule --project <project-arn> --device-pool <device-pool-arn> --app-file app-samplejr-staging-release.apk --test-file calabash_tests.zip --name "A new test, a new hope"
- Uploading app-file app-samplejr-staging-release.apk of type ANDROID_APP .
- Uploading test-file calabash_tests.zip of type CALABASH_TEST_PACKAGE.
- Initiating test run
- Waiting until the tests complete ..................................................................................
- Generating report
Reporting on run app-samplejr-staging-release.apk
- [FILE] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Setup Suite/0_Logcat.logcat
- [SCREENSHOT] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/0_i_take_a_screenshot_0..png
- [SCREENSHOT] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/1_i_take_a_screenshot_0..png
- [SCREENSHOT] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/2_i_take_a_screenshot_0..png
- [SCREENSHOT] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/3_i_take_a_screenshot_0..png
- [SCREENSHOT] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/4_i_take_a_screenshot_0..png
- [FILE] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/0_Calabash Pretty Output.txt
- [FILE] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/1_Calabash Standard Output.txt
- [FILE] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/2_Calabash JUnit XML Output.xml
- [FILE] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/3_Calabash JSON Output.json
- [FILE] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Homepage feature/4_Logcat.logcat
- [FILE] report/Samsung Nexus 10 (WiFi) - Nexus 10 - 4.2.2/Teardown Suite/0_Logcat.logcat
```

```
NAME:
   schedule - schedule a run

USAGE:
   command schedule [command options] [arguments...]

OPTIONS:
   --project            project arn or project description [$DF_PROJECT]
   --device-pool        devicepool arn or devicepool name [$DF_DEVICE_POOL]
   --device             device arn or devicepool name to run the test on [$DF_DEVICE]
   --name               name to give to the run that is scheduled [$DF_RUN_NAME]
   --app-file           path of the app file to be executed [$DF_APP_FILE]
   --app-type           type of app [ANDROID_APP,IOS_APP] [$DF_APP_TYPE]
   --test-file          path of the test file to be executed [$DF_TEST_FILE]
   --test-type          type of test [UIAUTOMATOR, CALABASH, APPIUM_JAVA_TESTNG, UIAUTOMATION, BUILTIN_FUZZ, INSTRUMENTATION, APPIUM_JAVA_JUNIT, BUILTIN_EXPLORER, XCTEST] [$DF_TEST_TYPE]
   --test               arn or name of the test upload to schedule [$DF_TEST]
   --app                arn or name of the app upload to schedule [$DF_APP]
```

## Report
Report will download the results (artifacts) in a standard structure (html format will be improved soon)

```
$ ./devicefarm-cli report --run <run-arn>
report/
├── ASUS\ Nexus\ 7\ -\ 1st\ Gen\ (WiFi)\ -\ Nexus\ 7\ -\ 1st\ Gen\ -\ 4.2.1
│   ├── Homepage\ feature
│   │   ├── 0_Calabash\ Pretty\ Output.txt
│   │   ├── 0_i_take_a_screenshot_0..png
│   │   ├── 1_Calabash\ Standard\ Output.txt
│   │   ├── 1_i_take_a_screenshot_0..png
│   │   ├── 2_Calabash\ JUnit\ XML\ Output.xml
│   │   ├── 2_i_take_a_screenshot_0..png
│   │   ├── 3_Calabash\ JSON\ Output.json
│   │   ├── 3_i_take_a_screenshot_0..png
│   │   ├── 4_Logcat.logcat
│   │   └── 4_i_take_a_screenshot_0..png
│   ├── Setup\ Suite
│   │   └── 0_Logcat.logcat
│   └── Teardown\ Suite
│       └── 0_Logcat.logcat
├── ASUS\ Nexus\ 7\ -\ 1st\ Gen\ (WiFi)\ -\ Nexus\ 7\ -\ 1st\ Gen\ -\ 4.4.2
│   ├── Homepage\ feature
│   │   ├── 0_Calabash\ Pretty\ Output.txt
│   │   ├── 0_i_take_a_screenshot_0..png
│   │   ├── 1_Calabash\ Standard\ Output.txt
│   │   ├── 1_i_take_a_screenshot_0..png
│   │   ├── 2_Calabash\ JUnit\ XML\ Output.xml
│   │   ├── 2_i_take_a_screenshot_0..png
│   │   ├── 3_Logcat.logcat
│   │   ├── 3_i_take_a_screenshot_0..png
│   │   ├── 4_Calabash\ JSON\ Output.json
│   │   └── 4_i_take_a_screenshot_0..png
│   ├── Setup\ Suite
│   │   └── 0_Logcat.logcat
│   └── Teardown\ Suite
│       └── 0_Logcat.logcat
```

# CLI

```
NAME:
   devicefarm-cli - allows you to interact with AWS devicefarm from the command line

USAGE:
   devicefarm-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR(S):
   Patrick Debois <Patrick.Debois@jedi.be>

COMMANDS:
   list         list various elements on devicefarm
   download     download various devicefarm elements
   status       get the status of a run
   report       get report about a run
   schedule     schedule a run
   create       creates various devicefarm elements
   info         get detailed info about various devicefarm elements
   upload       uploads an app, test and data
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```


# Gotchas so far:
- The user on aws it runs under is called 'rwx'
- Files are kept for 30days around
- You can not delete/update devicepools

- The docs are confusing at best: many calls specify that you need an ARN . for example List Artificats you need to specify an Artifact ARN, this is wrong and needs to be the Run ARN
- The doc use of ..Arn.. in params is not consistent with the golang where the param in the struct is often ARN (uppercase)
- the extention of test results is inconsistent: sometimes it's 'xml' 'json' but '.png' (with the dot in the extension)
- listArtifacts doesn't take the options Name & Extension to filter on specific one

- The test devices sometimes suffer from a DNS server not responding
- Be sure to crank up the timeouts in your test as devices can be slow
- I'd love to have the tests tagable like servers so I can do the billing per customer
- You can run arbitrary code from your test script - if all fails you can create a reverse remoteshell for debugging , using calabash ruby foo
- The machines run debian
- The aws console shows you a certain calabash version (0.7) but in reality the servers have a different one (0.8)
- Sometimes devices are not available during your tests


# Helpful links for fixing/finding issues
- <https://github.com/mitchellh/goamz/blob/caaaea8b30ee15616494ee68abd5d8ebbbef05cf/s3/s3.go>
- <https://github.com/aws/aws-sdk-go/blob/67dc9f948602be9f85cb640f89b0adec994ccbda/internal/signer/v4/v4.go#L263>
