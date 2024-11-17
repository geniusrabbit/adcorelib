//
// @project GeniusRabbit corelib 2016 – 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2024
//

package adresponse

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	openrtb "github.com/bsm/openrtb"
	natresp "github.com/bsm/openrtb/native/response"
	"golang.org/x/net/html/charset"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// BidResponse RTB record
type BidResponse struct {
	Src         adtype.Source
	Req         *adtype.BidRequest
	Application *admodels.Application
	Target      admodels.Target
	BidResponse openrtb.BidResponse
	context     context.Context
	optimalBids []*openrtb.Bid
	ads         []adtype.ResponserItemCommon
}

// AuctionID response
func (r *BidResponse) AuctionID() string {
	return r.BidResponse.ID
}

// AuctionType of request
func (r *BidResponse) AuctionType() types.AuctionType {
	return r.Req.AuctionType
}

// Source of response
func (r *BidResponse) Source() adtype.Source {
	return r.Src
}

// Prepare bid response
func (r *BidResponse) Prepare() {
	// Prepare URLs and markup for response
	for i, seat := range r.BidResponse.SeatBid {
		for i, bid := range seat.Bid {
			if imp := r.Req.ImpressionByIDvariation(bid.ImpID); imp != nil {
				// Prepare date for bid W/H
				if bid.W == 0 && bid.H == 0 {
					bid.W, bid.H = imp.W, imp.H
				}

				if imp.IsDirect() {
					// Custom direct detect
					if bid.AdMarkup == "" {
						bid.AdMarkup, _ = customDirectURL(bid.Ext)
					}
					if strings.HasPrefix(bid.AdMarkup, `<?xml`) {
						bid.AdMarkup, _ = decodePopMarkup([]byte(bid.AdMarkup))
					}
				}
			}

			replacer := r.newBidReplacer(&bid)
			bid.AdMarkup = replacer.Replace(bid.AdMarkup)
			bid.NURL = prepareURL(bid.NURL, replacer)
			bid.BURL = prepareURL(bid.BURL, replacer)

			seat.Bid[i] = bid
		}

		r.BidResponse.SeatBid[i] = seat
	} // end for

	for _, bid := range r.OptimalBids() {
		imp := r.Req.ImpressionByIDvariation(bid.ImpID)
		if imp == nil {
			continue
		}

		if imp.IsDirect() {
			format := imp.FormatByType(types.FormatDirectType)
			if format == nil {
				continue
			}
			r.ads = append(r.ads, &ResponseBidItem{
				ItemID:     imp.ID,
				Src:        r.Src,
				Req:        r.Req,
				Imp:        imp,
				FormatType: types.FormatDirectType,
				RespFormat: format,
				Bid:        bid,
				ActionLink: bid.AdMarkup,
			})
			continue
		}

		for _, format := range imp.Formats() {
			if bid.ImpID != imp.IDByFormat(format) {
				continue
			}
			switch {
			case format.IsNative():
				native, err := decodeNativeMarkup([]byte(bid.AdMarkup))
				// TODO parse native request
				if err == nil {
					r.ads = append(r.ads, &ResponseBidItem{
						ItemID:     imp.ID,
						Src:        r.Src,
						Req:        r.Req,
						Imp:        imp,
						FormatType: types.FormatNativeType,
						RespFormat: format,
						Bid:        bid,
						Native:     native,
						ActionLink: native.Link.URL,
					})
				}
			case format.IsBanner() || format.IsProxy():
				r.ads = append(r.ads, &ResponseBidItem{
					ItemID:     imp.ID,
					Src:        r.Src,
					Req:        r.Req,
					Imp:        imp,
					FormatType: bannerFormatType(bid.AdMarkup),
					RespFormat: format,
					Bid:        bid,
				})
			}
			break
		}
	}
}

// Request information
func (r *BidResponse) Request() *adtype.BidRequest {
	return r.Req
}

// Ads list
func (r *BidResponse) Ads() []adtype.ResponserItemCommon {
	return r.ads
}

// Item by impression code
func (r *BidResponse) Item(impid string) adtype.ResponserItemCommon {
	for _, it := range r.Ads() {
		if it.ImpressionID() == impid {
			return it
		}
	}
	return nil
}

// Price for response
func (r *BidResponse) Price() billing.Money {
	var price billing.Money
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			price += billing.MoneyFloat(bid.Price)
		}
	}
	return price
}

// Count bids
func (r *BidResponse) Count() int {
	return len(r.Bids())
}

// Validate response
func (r *BidResponse) Validate() error {
	if r == nil {
		return adtype.ErrResponseEmpty
	}
	err := r.BidResponse.Validate()
	if err == nil {
		for _, seat := range r.BidResponse.SeatBid {
			if seat.Group == 1 {
				return adtype.ErrResponseInvalidGroup
			}
		}
	}
	return err
}

// Error of the response
func (r *BidResponse) Error() error {
	return r.Validate()
}

// Bids list
func (r *BidResponse) Bids() []*openrtb.Bid {
	result := make([]*openrtb.Bid, 0, len(r.BidResponse.SeatBid))
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			result = append(result, &bid)
		}
	}
	return result
}

// OptimalBids list (the most expensive)
func (r *BidResponse) OptimalBids() []*openrtb.Bid {
	if len(r.optimalBids) > 0 {
		return r.optimalBids
	}

	bids := make(map[string]*openrtb.Bid, len(r.BidResponse.SeatBid))
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			if obid, ok := bids[bid.ImpID]; !ok || obid.Price < bid.Price {
				bids[bid.ImpID] = &bid
			}
		}
	}

	r.optimalBids = make([]*openrtb.Bid, 0, len(bids))
	for _, b := range bids {
		r.optimalBids = append(r.optimalBids, b)
	}
	return r.optimalBids
}

// BidPosition returns index from OpenRTB bid
func (r *BidResponse) BidPosition(b *openrtb.Bid) int {
	idx := 0
	for _, seat := range r.BidResponse.SeatBid {
		for _, bid := range seat.Bid {
			if bid.ImpID == b.ImpID {
				return idx
			}
			idx++
		}
	}
	return -1
}

// UpdateBid object
func (r *BidResponse) UpdateBid(b *openrtb.Bid) {
	for _, seat := range r.BidResponse.SeatBid {
		for j, bid := range seat.Bid {
			if bid.ImpID == b.ImpID {
				seat.Bid[j] = *b
			}
		}
	}
}

// Context of response
func (r *BidResponse) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		r.context = ctx[0]
	}
	if r.context == nil {
		return r.Req.Ctx
	}
	return r.context
}

// Get context value
func (r *BidResponse) Get(key string) any {
	if r.context != nil {
		return r.context.Value(key)
	}
	return nil
}

func (r *BidResponse) newBidReplacer(bid *openrtb.Bid) *strings.Replacer {
	return strings.NewReplacer(
		"${AUCTION_AD_ID}", bid.AdID,
		"${AUCTION_ID}", r.BidResponse.ID,
		"${AUCTION_BID_ID}", r.BidResponse.BidID,
		"${AUCTION_IMP_ID}", bid.ImpID,
		"${AUCTION_PRICE}", fmt.Sprintf("%.6f", bid.Price),
		"${AUCTION_CURRENCY}", "USD",
	)
}

// Release response and all linked objects
func (r *BidResponse) Release() {
	if r == nil {
		return
	}
	if r.Req != nil {
		r.Req.Release()
		r.Req = nil
	}
	r.ads = r.ads[:0]
	r.optimalBids = r.optimalBids[:0]
	r.Application = nil
	r.Target = nil
	r.BidResponse.SeatBid = r.BidResponse.SeatBid[:0]
	r.BidResponse.Ext = r.BidResponse.Ext[:0]
}

func decodePopMarkup(data []byte) (val string, err error) {
	var item struct {
		URL string `xml:"popunderAd>url"`
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&item); err == nil {
		val = item.URL
	}
	return val, err
}

func customDirectURL(data []byte) (val string, err error) {
	var item struct {
		URL         string `json:"url"`
		LandingPage string `json:"landingpage"`
		Link        string `json:"link"`
	}
	if err = json.Unmarshal(data, &item); err == nil {
		val = max(item.URL, item.LandingPage, item.Link)
	}
	return val, err
}

func decodeNativeMarkup(data []byte) (*natresp.Response, error) {
	var (
		native struct {
			Native natresp.Response `json:"native"`
		}
		err error
	)
	if bytes.Contains(data, []byte(`"native"`)) {
		err = json.Unmarshal(data, &native)
	} else {
		err = json.Unmarshal(data, &native.Native)
	}
	if err != nil {
		err = json.Unmarshal(data, &native.Native)
	}
	if err != nil {
		return nil, err
	}
	return &native.Native, nil
}

func bannerFormatType(markup string) types.FormatType {
	if strings.HasPrefix(markup, "http://") ||
		strings.HasPrefix(markup, "https://") ||
		(strings.HasPrefix(markup, "//") && !strings.ContainsAny(markup, "\n\t")) ||
		strings.Contains(markup, "<iframe") {
		return types.FormatProxyType
	}
	return types.FormatBannerType
}

func prepareURL(surl string, replacer *strings.Replacer) string {
	if surl == "" {
		return surl
	}
	if u, err := url.QueryUnescape(surl); err == nil {
		surl = u
	}
	return replacer.Replace(surl)
}

var (
	_ adtype.Responser = &BidResponse{}
)
