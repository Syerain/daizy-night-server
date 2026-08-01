## Registercode string
every register request must use a Registercode.

## handlers
handler 's only job is to handle a request and pass it to a service.
it shouldnt contain any business logic

dbware is as the same.

## about Providers
他们的字段应该关于他们“是什么”，而不是“从哪里来”
据此应该尽可能避免不合理的引用，而是直接取字段值