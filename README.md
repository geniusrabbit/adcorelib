# Corelib advertisement project

Contains common lib pockages

## Prices math

ComissionShareFactor - represents internal commision of the system.
RevenueShareFactor - represents revenue of the publisher or advertiser

RevenueShareFactor = 1. - ComissionShareFactor

For example if
ComissionShareFactor = 10% = 0.1 then
RevenueShareFactor = 90% = 0.9
RevenueShareReduce = 10%

* Source
  - `ComissionShareFactor` and `RevenueShareReduce`
    - newPrice = (price * (1 - `RevenueShareReduce`)) - price * `ComissionShareFactor` 1 * (1 - 0.1) * (1 - 0.1)  = 0.81
* Target (Zone + Site, AccessPoint{DSP})
  - `ComissionShareFactor` and `RevenueShareReduce`
    - publisherPrice = price - price * `ComissionShareFactor` - price * `RevenueShareReduce`
  - If target have fixed view price then can be used that value instead
  - If target is AccessPoint then `RevenueShareReduce` will reduce descrepancy

So we have two commissions:
  1 - from source to reduce descrepancy between buyer(DSP) and saller
  2 - from target
