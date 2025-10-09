# AdSource

The AdSource package is a crucial component of the AdEngine, designed to manage and interact with various ad sources. It provides a unified interface for all ad sources and includes a collection of standard ad sources such as in-memory, database, OpenRTB (Real-Time Bidding), and more.

## Key Features

- **Unified Interface**: Offers a consistent interface for different types of ad sources, simplifying integration and management.
- **Parallel Bid Requests**: Supports parallel processing of bid requests using a worker pool, improving efficiency and performance.
- **Tracing and Logging**: Integrates with tracing (using OpenTracing) and logging (using zap) to provide detailed insights into the bidding process and performance.
- **Metrics Collection**: Collects and reports metrics for monitoring the performance and health of ad sources.

## Main Components

### MultisourceWrapper

The `MultisourceWrapper` is the core abstraction in this package. It manages multiple ad sources and controls where to send requests and how to handle responses.

#### Features

- **Main Source**: A primary source that is called every time a bid request is made.
- **Source List**: A list of external platforms to send requests to.
- **Execution Pool**: A pool for executing bid requests in parallel.
- **Request Timeout**: Manages the duration for request timeouts.
- **Max Parallel Requests**: Limits the maximum number of parallel requests.
- **Metrics Accessor**: Accesses and updates metrics related to bidding.

#### Example Usage

```go
wrapper, err := adsource.NewMultisourceWrapper(options...)
if err != nil {
    log.Fatal(err)
}

response := wrapper.Bid(request)
if response.Error() != nil {
    log.Println("Bid request failed:", response.Error())
} else {
    log.Println("Bid request succeeded:", response.Ads())
}
```

## OpenRTB Driver

The `openrtb` package provides an implementation of an OpenRTB (Real-Time Bidding) driver. This driver allows the AdEngine to interact with OpenRTB-compliant ad exchanges.

### Features (oepnrtb.Driver)

- **Supports OpenRTB Versions 2.5 and 3.0**: Handles bid requests and responses in both OpenRTB 2.5 and 3.0 formats.
- **Latency Metrics**: Measures and reports the latency of bid requests.
- **Error Handling**: Manages errors and retries for bid requests.
- **RPS (Requests Per Second) Limiting**: Controls the rate of requests to comply with source limitations.

#### Example Usage (oepnrtb.Driver)

```go
driver, err := openrtb.newDriver(context.Background(), source, netClient)
if err != nil {
    log.Fatal(err)
}

request := &bidrequest.BidRequest{/*...*/}
response := driver.Bid(request)
if response.Error() != nil {
    log.Println("Bid request failed:", response.Error())
} else {
    log.Println("Bid request succeeded:", response.Ads())
}
```

### Key Methods

- **ID**: Returns the ID of the source.
- **Protocol**: Returns the protocol of the source.
- **Test**: Validates the request before processing.
- **Bid**: Handles a bid request and processes it through the OpenRTB exchange.
- **ProcessResponseItem**: Processes individual response items.
- **Metrics**: Returns metrics information for the platform.

### Dependencies

The package relies on several external libraries to provide its functionality:

- `github.com/bsm/openrtb`: For handling OpenRTB bid requests and responses.
- `github.com/demdxx/gocast/v2`: For type casting.
- `github.com/geniusrabbit/adcorelib/*`: Various modules from the AdCoreLib for context, event tracking, fast time, and more.
- `go.uber.org/zap`: For logging.

## Error Handling

The package defines standard error messages for common error scenarios, ensuring consistency and clarity in error reporting.

```go
var (
    ErrSourcesCantBeNil = errors.New("[SSP] sources can't be nil")
)
```

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests.
