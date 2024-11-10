# adresponse Package

The `adresponse` package is a part of the GeniusRabbit core library, designed to handle advertisement responses within an advertising system. It provides structures and methods for managing ad items, impressions, campaigns, pricing models, and more.

## Overview

The primary structure in this package is `ResponseAdItem`, which represents an advertisement item selected from storage. It encapsulates information about the ad, campaign, impression, pricing, and various identifiers essential for processing and tracking advertisements.

## Installation

To use the `adresponse` package, import it into your Go project:

```go
import "github.com/geniusrabbit/adcorelib/adresponse"
```

Ensure that you have all necessary dependencies installed, such as `admodels`, `adtype`, `billing`, and other related packages.

## Structures

### ResponseAdItem

The `ResponseAdItem` struct represents an advertisement item and contains the following fields:

- **Ctx** (`context.Context`): The context of the request.
- **ItemID** (`string`): Unique identifier of the response item.
- **Src** (`adtype.Source`): Source of the response.
- **Req** (`*adtype.BidRequest`): Bid request information.
- **Imp** (`*adtype.Impression`): Impression associated with the response item.
- **Campaign** (`*admodels.Campaign`): Campaign associated with the response item.
- **Ad** (`*admodels.Ad`): Advertisement data.
- **AdBid** (`*admodels.AdBid`): Bid details for the ad.
- **AdLink** (`admodels.AdLink`): Link for the ad.
- **BidECPM** (`billing.Money`): Bid's effective cost per mille (CPM).
- **BidPrice** (`billing.Money`): Maximum RTB bid price (CPM only).
- **AdPrice** (`billing.Money`): New price for advertisement target actions (click, lead, impression).
- **AdLeadPrice** (`billing.Money`): Price for a lead action.
- **CPMBidPrice** (`billing.Money`): CPM bid price, updated only by the price predictor.
- **SecondAd** (`adtype.SecondAd`): Secondary advertisement information.

## Methods

### Identification and Context

- **ID() string**: Returns the unique identifier of the current response item.
- **AuctionID() string**: Retrieves the auction identifier from the request.
- **Context(ctx ...context.Context) context.Context**: Gets or sets the context value.
- **Get(key string) any**: Retrieves a value associated with a key from the context.

### Impression and Request

- **Impression() *adtype.Impression**: Returns the impression associated with the response item.
- **ImpressionID() string**: Retrieves the unique identifier of the impression.
- **ExtImpressionID() string**: Returns the external impression identifier.
- **ExtTargetID() string**: Retrieves the external target identifier.
- **Request() *adtype.BidRequest**: Provides access to the bid request information.

### Advertisement and Campaign

- **AdDirectLink() string**: Returns the direct link of the advertisement.
- **ContentItemString(name string) string**: Retrieves a content item as a string from the advertisement.
- **ContentItem(name string) any**: Retrieves the ad response data for a specified content item.
- **ContentFields() map[string]any**: Returns the content fields from the advertisement object.
- **MainAsset() *admodels.AdAsset**: Returns the main asset from the advertisement.
- **Asset(name string) *admodels.AdAsset**: Retrieves an asset by name.
- **Assets() admodels.AdAssets**: Returns a list of assets from the advertisement.
- **AdID() uint64**: Retrieves the advertisement identifier.
- **AdIDString() string**: Returns the advertisement identifier as a string.
- **CampaignID() uint64**: Retrieves the campaign identifier.
- **CampaignIDString() string**: Returns the campaign identifier as a string.
- **ProjectID() uint64**: Retrieves the project identifier.
- **AccountID() uint64**: Retrieves the account identifier.
- **CreativeIDString() string**: Returns the creative identifier as a string for reporting content issues or defects.
- **Second() *adtype.SecondAd**: Returns secondary advertisement information.

### Source and Format

- **Source() adtype.Source**: Returns the source of the response.
- **NetworkName() string**: Retrieves the network name associated with the source.
- **Format() *types.Format**: Returns the format object model.
- **PriorityFormatType() types.FormatType**: Determines the primary format type from the current ad.
- **IsDirect() bool**: Indicates whether the response item is direct.
- **ActionURL() string**: Returns the action URL for direct banners.

### Dimensions

- **Width() int**: Retrieves the width of the ad.
- **Height() int**: Retrieves the height of the ad.
- **TargetID() uint64**: Retrieves the target identifier.
- **TargetIDString() string**: Returns the target identifier as a string.

### Pricing and Bidding

- **PricingModel() types.PricingModel**: Returns the pricing model of the advertisement.
- **ECPM() billing.Money**: Retrieves the effective cost per mille of the item.
- **Price(action admodels.Action, removeFactors ...adtype.PriceFactor) billing.Money**: Calculates the total price for a specific action (click, lead, impression), considering any price factors.
- **SetCPMPrice(price billing.Money, includeFactors ...adtype.PriceFactor)**: Updates the DSP auction value with the specified price.
- **CPMPrice(removeFactors ...adtype.PriceFactor) billing.Money**: Calculates the price value for DSP auction, removing any specified price factors.
- **AuctionCPMBid() billing.Money**: Returns the bid price without any commission.
- **PurchasePrice(action admodels.Action, removeFactors ...adtype.PriceFactor) billing.Money**: Determines the price of a view from an external resource for a given action.
- **PotentialPrice(action admodels.Action) billing.Money**: Calculates the potential price that could be received from a source but was marked as discrepancy.

### Tracking and Validation

- **ViewTrackerLinks() []string**: Returns tracking links for view actions.
- **ClickTrackerLinks() []string**: Returns third-party tracker URLs to be fired on click of the URL.
- **Validate() error**: Validates the item, returning an error if validation fails.

### Revenue and Commission

- **RevenueShareFactor() float64**: Returns the revenue share percentage.
- **ComissionShareFactor() float64**: Returns the commission share percentage that the system gets from the publisher.

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/geniusrabbit/adcorelib/admodels"
    "github.com/geniusrabbit/adcorelib/adresponse"
    "github.com/geniusrabbit/adcorelib/adtype"
    "github.com/geniusrabbit/adcorelib/billing"
)

func main() {
    // Create a new ResponseAdItem
    adItem := &adresponse.ResponseAdItem{
        ItemID: "unique-item-id",
        Req:    &adtype.BidRequest{ID: "auction-123"},
        Ad: &admodels.Ad{
            ID:            456,
            Price:         1000,
            PricingModel:  admodels.PricingModelCPM,
            Format:        &admodels.Format{Types: admodels.FormatTypes{admodels.FormatBanner}},
            LeadPrice:     500,
            Content:       map[string]any{"title": "Ad Title"},
        },
        Imp: &adtype.Impression{ID: "impression-789", W: 300, H: 250},
        AdPrice:  billing.Money(1500),
        BidPrice: billing.Money(1200),
    }

    // Access basic information
    fmt.Println("Ad Item ID:", adItem.ID())
    fmt.Println("Auction ID:", adItem.AuctionID())
    fmt.Println("Ad ID:", adItem.AdID())
    fmt.Println("Ad Direct Link:", adItem.AdDirectLink())

    // Set CPM Price
    adItem.SetCPMPrice(2000)
    fmt.Println("CPM Bid Price:", adItem.CPMBidPrice)

    // Calculate ECPM
    fmt.Println("ECPM:", adItem.ECPM())

    // Get the price for an impression action
    price := adItem.Price(admodels.ActionImpression)
    fmt.Println("Price for Impression:", price)

    // Validate the ad item
    if err := adItem.Validate(); err != nil {
        fmt.Println("Validation Error:", err)
    } else {
        fmt.Println("Ad Item is valid.")
    }

    // Access content fields
    content := adItem.ContentFields()
    fmt.Println("Ad Content Title:", content["title"])

    // Additional usage as per application logic
}
```

## Notes

- **Dependencies**: The package relies on other components like `admodels`, `adtype`, `billing`, etc. Ensure these are properly imported and initialized in your project.
- **Placeholders**: Some methods may return placeholders or require additional implementation, especially where external data or complex logic is involved.
- **Internal Processing**: Methods like `processParameters` are used internally to replace placeholders in strings with actual values based on the ad item's context and properties.
