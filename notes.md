Options pattern: returns a function with the new property added to the client.
type Option func(*Client)
WithHTTPTimeout
WithRetryCount
WithRetrier
WithHTTPClient
We can chain as many methods as we want with this method:
hystrixClient := hystrix.NewClient(
    hystrix.WithHTTPTimeout(timeout),
    hystrix.WithCommandName("MyCommand"),
    hystrix.WithHystrixTimeout(1100*time.Millisecond),
    hystrix.WithMaxConcurrentRequests(100),
    hystrix.WithErrorPercentThreshold(25),
    hystrix.WithSleepWindow(10),
    hystrix.WithRequestVolumeThreshold(10),
    hystrix.WithHTTPClient(&myHTTPClient{
        // replace with custom HTTP client
        client: http.Client{Timeout: 25 * time.Millisecond},
    }),
)

There are 2 common ways of using options pattern:
Chaining methods. You return the type. ⇒ logrus.
Having a method which uses a list of methods. All of them are called line by line. Hence each parameter returns a method ⇒ heimdall.
