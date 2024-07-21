package aerospike

import (
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	as "github.com/aerospike/aerospike-client-go"

	"github.com/geniusrabbit/adcorelib/models"
)

var (
	writePolicy = as.NewWritePolicy(0, 0)
	scanPolicy  = as.NewScanPolicy()
	readPolicy  = as.NewPolicy()
)

var (
	ErrUnderSynchronisation = errors.New(`[aerospike] under synchronisation`)
)

type synchronizer struct {
	// isSincing state
	isSincing int32

	// Connect aerospike client
	client *as.Client

	// current namespace name
	namespace string

	// Time of last sinc
	lastSyncTime time.Time
}

// NewSynchronizer returns new sync manager of the ads
func NewSynchronizer(ac *as.Client, namespace string) *synchronizer {
	return &synchronizer{client: ac, namespace: namespace}
}

// func (sync *synchronizer) Sync(reader storage.Reader) error {
// 	if sync.startSync() {
// 		return ErrUnderSynchronisation
// 	}
// 	defer sync.setSyncState(false)

// 	list, err := reader.CampaignList(nil)
// 	if err != nil {
// 		return err
// 	}

// 	namespace := sync.currentNamespace()
// 	for _, camp := range list {
// 		var (
// 			formats      []uint
// 			countriesArr gosql.NullableOrderedNumberArray[uint]
// 			languagesArr gosql.NullableOrderedNumberArray[uint]
// 			hours, err   = types.HoursByString(camp.Hours.String)
// 		)

// 		if err != nil {
// 			return err
// 		}

// 		// Format list
// 		for _, ad := range camp.Ads {
// 			formats = append(formats, uint(ad.ID))
// 		}

// 		// Countries filter
// 		if camp.Geos.Len() > 0 {
// 			for _, cc := range camp.Geos {
// 				countriesArr = append(countriesArr, uint(gogeo.CountryByCode2(cc).ID))
// 			}
// 			countriesArr.Sort()
// 		}

// 		// Languages filter
// 		if len(camp.Languages) > 0 {
// 			for _, lg := range camp.Languages {
// 				languagesArr = append(languagesArr, languages.GetLanguageIdByCodeString(lg))
// 			}
// 			languagesArr.Sort()
// 		}

// 		private := camp.Private.IsPrivate()
// 		key, _ := as.NewKey(namespace, setName(private), campaignKey(camp))
// 		err = sync.client.PutBins(writePolicy, key,
// 			as.NewBin("id", camp.ID),
// 			as.NewBin("account_id", camp.AccountID),

// 			as.NewBin("weight", 0),

// 			as.NewBin("daily_test_budget", billing.MoneyFloat(camp.DailyTestBudget).Int64()),
// 			as.NewBin("test_budget", billing.MoneyFloat(camp.TestBudget).Int64()),
// 			as.NewBin("daily_budget", billing.MoneyFloat(camp.DailyBudget).Int64()),
// 			as.NewBin("budget", billing.MoneyFloat(camp.Budget).Int64()),

// 			as.NewBin("context", camp.Context.String()),

// 			as.NewBin("formats", formats),
// 			// as.NewBin("keywords", camp.Keywords),
// 			as.NewBin("zones", camp.Zones.Ordered()),
// 			as.NewBin("domains", camp.Domains),
// 			as.NewBin("categories", camp.Categories.Ordered()),
// 			as.NewBin("countries", countriesArr),
// 			as.NewBin("languages", languagesArr),
// 			as.NewBin("browsers", camp.Browsers.Ordered()),
// 			as.NewBin("os", camp.Os.Ordered()),
// 			as.NewBin("device_types", camp.DeviceTypes.Ordered()),
// 			as.NewBin("devices", camp.Devices.Ordered()),
// 			as.NewBin("hours", hours.String()),
// 			as.NewBin("sex", camp.Sex.Ordered()),
// 			as.NewBin("age", camp.Age.Ordered()),
// 			as.NewBin("trace", camp.Trace),
// 			as.NewBin("trace_percent", camp.TracePercent),
// 		)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// deleteSet and all records inside
func (sync *synchronizer) deleteSet(namespace, name string) (int, error) {
	var (
		record       *as.Record
		count        int
		records, err = sync.client.ScanAll(scanPolicy, namespace, name)
	)
	if err != nil {
		return 0, err
	}
	defer records.Close()
	for {
		record, err = records.Read()
		if err == nil && record != nil {
			_, err = sync.client.Delete(writePolicy, record.Key)
		}
		if err != nil || record == nil {
			break
		}
		count++
	}
	if err == io.EOF {
		err = nil
	}
	return count, err
}

func (sync *synchronizer) namespaces() []string {
	return sync.client.GetNodeNames()
}

func (sync *synchronizer) setNamespace(_ string) error {
	return nil
}

func (sync *synchronizer) currentNamespace() string {
	key, err := as.NewKey("adserver", "options", "namespace")
	if err != nil {
		return ""
	}
	record, err := sync.client.Get(readPolicy, key)
	if err != nil {
		return ""
	}
	value, _ := record.Bins["value"]
	fmt.Println("Value", value)
	return ""
}

func (sync *synchronizer) isSyncState() bool {
	return atomic.LoadInt32(&sync.isSincing) == 1
}

func (sync *synchronizer) startSync() bool {
	return atomic.CompareAndSwapInt32(&sync.isSincing, 0, 1)
}

func (sync *synchronizer) setSyncState(sincing bool) {
	if sincing {
		atomic.StoreInt32(&sync.isSincing, 1)
	} else {
		atomic.StoreInt32(&sync.isSincing, 0)
	}
}

func campaignKey(camp *models.Campaign) any {
	return camp.ID
}

func setName(private bool) string {
	if private {
		return "private"
	}
	return "public"
}
