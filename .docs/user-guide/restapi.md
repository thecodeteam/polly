# REST API

The Polly REST API documentation at Apiary!

## Administrative
All of the REST interface documentation is available [here](http://docs.polly1.apiary.io/)
for the administrative interface of Polly.

- Volume Create/Remove
- Offer Create/Remove
- Label Create/Remove

## Scheduler and Offers
This interface is responsible for integration directly with schedulers. This
is not a defined interface by Polly.

- Offer/Acceptance
- Claim/Unclaim

## libStorage
All storage operations and integration from storage orchestrators such as
[REX-Ray] take place through this API. It is documented
[here](https://docs.libstorage.apiary.io). Polly is considered a libStorage
server so any compatible libStorage client may be a directly integrated with
Polly services.
