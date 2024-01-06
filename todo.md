## TODO LIST

mobility-events

- "evuuid": "1c68267f-0182-53e5-a3bd-3940b1f0c47e"
- "evstart": "2022-05-10 00:00:00.000+0000"
- "evend": "2022-05-11 00:00:00.000+0000",

tourism-events

- "Id": "BFEB2DDB0FD54AC9BC040053A5514A92_REDUCED"
- "DateBegin": "2022-06-01T00:00:00"
- "DateEnd": "2022-06-01T00:00:00"

event:
    - id: evuuid / Id
    - start_date: evstart / DateBegin
    - end_date: evend / DateEnd

Thought process:
Mobility

- focus on events
- using specific path /v2/flat%2Cevent
- important details - id, start_date, end_date

TODOs:

- [ ] allow complex json response
- [ ] move config.json to root of the project, consider adding support for yaml as it may be more readable. Cool tool to convert json to yaml: <https://github.com/TomWright/dasel>
- [ ] would be nice to make the whole thing a mini-app that ties together all the poings
  - then we could test end to end (request to the router -> router -> service -> mapper -> response)
git