# alpha

## todo 26-07-12 tarde
@/internal/config/config.go

15:29 maybe we should let every module uses its own logger?

## adv 26-07-12 night
@global

23:01 fine i introduced gorm and sqlite.

00:42 nxtd; damn i have never used sql or http. wtf are these gorm key syntax?? im a damn noob..

## adv 26-07-13 tarde
@/internal/model/user.go

17:14 okay seems finished user model.

## bigadv 26-07-13 night
@global

/internal/crypto,

/internal/service,

/internal/config,

/internal/db,

/internal/model,

/internal/handler,

/etc

## adv 26-07-13 night
@/test/data.db

00:09 nxtd. temply added something for test. the /config.yaml and its included keys are just for test.

## adv 26-07-14 noon & tarde
@global

17:13 added logic of router, handler, service. did some stuff also.

## adv 26-07-14 night
@global

20:12 added /internal/constants; stuff; adv /internal/db; added /internal/crypto/salt.go

## adv 26-07-23
@global

23:16 不同的 registerway 所需求的参数不一，但不论如何 validate 函数都始终忠实地校验并上报每一个问题，
根据具体 registerway 判断所需参数，这是调用者的任务

