// Package prices provides embeddable price scopes for advertisement campaigns
// and separate actions performed on them.
//
// # Three prices
//
// Every action of an advertisement resolves to three different money values:
//
//   - PotentialPrice – the maximum price the advertiser could have paid for the
//     action. It is based on the declared maximum bid and is not reduced by any
//     factor. Used to track the discrepancy between the expected and the real price.
//   - AdvertiserPrice – the price which will be charged from the advertiser. It is
//     the current (post auction) bid as is, with no reduction by any correction
//     factor: the advertiser covers the full cost of running traffic across every
//     connected network, discrepancies included, so none of it is ever deducted from
//     their charge.
//   - PublisherPrice – the price which the system has to pay to the publisher. It is
//     the same bid reduced by the source and the target discrepancy correction
//     factors and by the system commission share, or a fixed price of the target if
//     it defines one. Whatever these factors remove stays with the network.
//
// The following invariant always holds for non-fixed publisher prices and factors
// in the 0..1 range:
//
//	PublisherPrice <= AdvertiserPrice <= PotentialPrice
//
// # Network profit
//
// NetworkProfit is the difference between what the advertiser was charged and what
// the publisher was paid:
//
//	NetworkProfit = AdvertiserPrice - PublisherPrice
//
// Since none of the correction factors ever reduce AdvertiserPrice, NetworkProfit
// captures both the system commission share and whatever the source/target
// discrepancy corrections deducted from the publisher payout.
//
// # Units
//
// Impression and view bids are stored in CPM units, i.e. the price of 1000 actions,
// which is the unit used by the external API and by RTB protocols. Click and lead
// bids are stored as the price of a single action. The action dispatch methods of
// the [PriceScope] always operate with the price of one single action, so consumers
// never convert the units themselves. Use [PriceFromCPM] and [CPMFromPrice] if the
// conversion is required outside of the scope.
//
// # No auction rules here
//
// The package performs honest unconditional calculations for every action which has
// a bid defined. It intentionally knows nothing about pricing models, auction types
// or which of the simultaneously active models should actually be charged – that is
// the responsibility of the auction and RTB integration modules.
//
// # Embedding
//
// Each pricing model is a small standalone scope with unique field and method names,
// so any number of them can be embedded into the same structure without selector
// conflicts:
//
//	type Campaign struct {
//		prices.CPMScope
//		prices.CPCScope
//	}
//
// Use the composite [PriceScope] to get all models at once together with the generic
// action dispatch and the three price calculations.
package prices
