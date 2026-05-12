# Changelog

## 1.0.0 (2026-05-12)


### Features

* add stdlib-based zoom cli wrapper with auth and service commands ([3833532](https://github.com/Hawkeye-Claims/go-zoom/commit/3833532cb49511aff444fd6edce1aa942240a291))
* add Zoom workflow skills ([4ecf88d](https://github.com/Hawkeye-Claims/go-zoom/commit/4ecf88dea648bb6f8912ce812147e215316cdcf0))
* **cli:** add create subcommands with json and json-file input ([3fbb66b](https://github.com/Hawkeye-Claims/go-zoom/commit/3fbb66bd470c545adf454c90db96848af8f6c082))
* **cli:** add phone recording download commands ([0274a5c](https://github.com/Hawkeye-Claims/go-zoom/commit/0274a5cd569f06d371adf624fd3d5d40c6452203))
* **cli:** support query JSON for get commands ([965cc3d](https://github.com/Hawkeye-Claims/go-zoom/commit/965cc3d9d5e842fc56384ec164c07868cc7abdb1))
* included event for call recording webhooks ([1f0ac32](https://github.com/Hawkeye-Claims/go-zoom/commit/1f0ac32892a076ccac6a43b2113afe9173a75ed4))
* included missing methods from existing sdk methods ([53571b5](https://github.com/Hawkeye-Claims/go-zoom/commit/53571b57690ffc89f2168ffbbb8670366eb99d12))
* support CLI OAuth authorization links ([5e340a0](https://github.com/Hawkeye-Claims/go-zoom/commit/5e340a0c46e7363ef3712c9dab598e7f67ff4097))


### Bug Fixes

* added call history id to recording model ([b9737be](https://github.com/Hawkeye-Claims/go-zoom/commit/b9737bef60d31d5d57f176c234ac8cbbe3a81d2a))
* ApprovalType should not omit when zero'd ([22423a3](https://github.com/Hawkeye-Claims/go-zoom/commit/22423a31f6e507149eecf024902261e78059cfce))
* check oauth token unlock error ([914aad9](https://github.com/Hawkeye-Claims/go-zoom/commit/914aad99c7a0c58c5f96c2e52fea3b0471bc07d2))
* escape meetingId and userId with url.PathEscape in meeting.go URLs ([7dc0ca7](https://github.com/Hawkeye-Claims/go-zoom/commit/7dc0ca7a34955290186f2b3fe4a7abcc22997534))
* escape userId with url.PathEscape in URL path segments ([885d33d](https://github.com/Hawkeye-Claims/go-zoom/commit/885d33d45e203aae3cebc4a2e00583e229148f7e))
* initialize models.User pointer in UsersService.Update to unmarshal response ([7e6907a](https://github.com/Hawkeye-Claims/go-zoom/commit/7e6907adaaa8c83946f4fb0b6b1c01df43c0f05e))
* make Feature a pointer and add omitempty to all UserUpdateAttributes fields ([0198910](https://github.com/Hawkeye-Claims/go-zoom/commit/01989108cf3df5f21f6c6fbe25b9d8e4c8ce8543))
* **models:** fixed ext types that weren't using the available enums ([b3008f2](https://github.com/Hawkeye-Claims/go-zoom/commit/b3008f2f603988855554c2d5a3292b0857c62388))
* omit empty meeting settings fields in payloads ([07611df](https://github.com/Hawkeye-Claims/go-zoom/commit/07611df67aaf53ae1ef3d664340cb09ea31c8e9b))
* omit empty nested fields in meeting and user payloads ([33dfe0a](https://github.com/Hawkeye-Claims/go-zoom/commit/33dfe0af2359e4e8c26fe8ef94f85723178dce41))
* omit empty nested fields in meeting create payload ([5c99b2d](https://github.com/Hawkeye-Claims/go-zoom/commit/5c99b2d5a9b6ca7897376fd550182df964c9153a))
* omit empty tracking fields in meeting create payload ([8234830](https://github.com/Hawkeye-Claims/go-zoom/commit/82348303da225e580166e99a0a312aba243bb452))
* **phone/models:** added missing fileds in CallPath and CallElement structs ([5c6b52e](https://github.com/Hawkeye-Claims/go-zoom/commit/5c6b52e4b2a2ec923e8f161e007aa050500f4660))
* **phone:** added missing fields to Call History and Call Path objects ([8998c7c](https://github.com/Hawkeye-Claims/go-zoom/commit/8998c7c284b5845043826bb56a10f946e758d1a3))
* **phone:** corrected type of channel marks ([d391664](https://github.com/Hawkeye-Claims/go-zoom/commit/d391664c9efa80e2115bfafe67de2fe275a797dc))
* **phone:** fixed version type ([3cc73de](https://github.com/Hawkeye-Claims/go-zoom/commit/3cc73de893e2af1f138cfd801ff84e87312cd2bd))
* **phone:** formatting of TS fields to string ([4c47748](https://github.com/Hawkeye-Claims/go-zoom/commit/4c477484976e1b516c97515a4b30ff3463cf8467))
* **phone:** improve OAuth token error handling and call history id mapping ([2f297a0](https://github.com/Hawkeye-Claims/go-zoom/commit/2f297a0ba34b826f9fcde458e4aa4b1368c34531))
* rename CallRecorddingGetOptions to CallRecordingGetOptions ([f932b4a](https://github.com/Hawkeye-Claims/go-zoom/commit/f932b4a221908f809483c0999867a5fc41719e40))
* **server:** correct Zoom webhook validation parsing ([1e0ed15](https://github.com/Hawkeye-Claims/go-zoom/commit/1e0ed15247e062e57cda09cfc18402f54365f1a3))
* track cli entrypoint and validate build ([ec79861](https://github.com/Hawkeye-Claims/go-zoom/commit/ec798614ff9a4ffcf0cbde01a7f1ba53193a340b))
* validate oauth listen address includes port ([ec90cc7](https://github.com/Hawkeye-Claims/go-zoom/commit/ec90cc7630ead99452ebedacaadd5b676f34ab5d))
