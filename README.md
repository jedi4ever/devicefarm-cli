# Description
This is the start (poc code currently) of an AWS Device Farm cli utility.

It will provide a more human & continous friendly CLI tool.

Have a look in the code.

Things I want to cover:
- upload of all elements
- have correct exit codes
- support a config file & environment settings & cli options
- have formatters like json, junit etc..
- poll for test results
- provided binaries & packages & docker instance
- make it slack friendly for reports

cleanup will happen soon!

# Gotchas so far:
- The upload urls aws provides are pre-signed s3 urls. when you directly pass this to golang newhttprequest is converts the path & query string internal. Therefore the request to get it gets a bad signature error. You need to use URL.raw_query & URL.opaque when creating the request
- The user on aws it runs under is called 'rwx'
- Files are kept for 30days around
- You can not delete/update devicepools
- The docs are confusing at best: many calls specify that you need an ARN . for example List Artificats you need to specify an Artifact ARN, this is wrong and needs to be the Run ARN
- The doc use of ..Arn.. in params is not consistent with the golang where the param in the struct is often ARN (uppercase)
- the extention of test results is inconsistent: sometimes it's 'xml' 'json' but '.png' (with the dot in the extension)
- listArtifacts API doc is incorrect: for the type you can only specify "LOG", "FILE", "SCREENSHOT"
- listArtificats doesn't take the options Name & Extension to filter on specific one
- Pi found in the api of certain devices `
API Result #aws #devicefarm :

       CPU: {
         Architecture: "foo",
         Clock: 3.14159,
         Frequency: "foo"
       },
       `
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
