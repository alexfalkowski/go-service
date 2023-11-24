# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## [v1.122.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.0) - 2023-11-24

- [`8295436`](https://github.com/alexfalkowski/go-service/commit/8295436ccb00ae0fd3a4dc8c5c4383a8557f0278) feat: add new linters (#455)

## [v1.121.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.121.0) - 2023-11-23

- [`c3e6bba`](https://github.com/alexfalkowski/go-service/commit/c3e6bbaf609953c879d8803003e943f2d491673f) feat(transport): do not pass config and use options (#454)
- [`f4418f5`](https://github.com/alexfalkowski/go-service/commit/f4418f5b0b1463536444e702e861c91a2358a4f0) build(deps): add bin (#453)
- [`12ca5ea`](https://github.com/alexfalkowski/go-service/commit/12ca5eab6aef28c9650008313425f61853fc4c45) build(deps): update bin (#452)
- [`98b609c`](https://github.com/alexfalkowski/go-service/commit/98b609c58bb45019df9faf5fbee8a95fdb65a278) docs: add some prerequisites and steps (#451)

## [v1.120.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.120.0) - 2023-11-21

- [`de98542`](https://github.com/alexfalkowski/go-service/commit/de985423eb7a4398f5cfd0ccbf98ebb0c7e94db1) feat(http): add ability to extend http server with middleware (#450)

## [v1.119.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.119.0) - 2023-11-19

- [`c048bcc`](https://github.com/alexfalkowski/go-service/commit/c048bcc56ded6bb8bcd7d12c476d18a98406af25) feat(security): add ability to configure token kind (#449)

## [v1.118.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.118.0) - 2023-11-19

- [`7dc5655`](https://github.com/alexfalkowski/go-service/commit/7dc5655b3e1762363db73b440723db54e9ac3595) feat(security): remove register as it can be done through grpc (#448)

## [v1.117.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.117.1) - 2023-11-17

- [`9284f61`](https://github.com/alexfalkowski/go-service/commit/9284f61c0507cb409dbb3ffd31b51d05bdf4ce90) fix(security): move register to constructor (#447)

## [v1.117.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.117.0) - 2023-11-17

- [`4a023f6`](https://github.com/alexfalkowski/go-service/commit/4a023f6b969ae990645db024e1350867ddf7f8dd) feat(security): add ability to register generators and verifiers (#446)
- [`c46ebb4`](https://github.com/alexfalkowski/go-service/commit/c46ebb4fff64cb587769dfd09eda6ad82365dc4c) docs(security): add enabled (#445)

## [v1.116.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.116.2) - 2023-11-17

- [`fee4f3a`](https://github.com/alexfalkowski/go-service/commit/fee4f3aa7eb3fb412526532329bc6186eebfc519) fix(security): add enabled config (#444)

## [v1.116.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.116.1) - 2023-11-17

- [`bc95770`](https://github.com/alexfalkowski/go-service/commit/bc95770e6876280552a97b11e2f5346cdb80a884) fix(deps): bump the otel group with 4 updates (#443)

## [v1.116.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.116.0) - 2023-11-16

- [`6732898`](https://github.com/alexfalkowski/go-service/commit/6732898584c98a14fdfe36ebd424c616788be54b) feat(deps): update github.com/jackc/pgx/v5 to v5.5.0 (#442)

## [v1.115.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.115.0) - 2023-11-16

- [`3c5ef21`](https://github.com/alexfalkowski/go-service/commit/3c5ef21215c1d7be3f2567f58ae76561ed9841f9) feat(deps): update github.com/jackc/pgx/v5 to v5.3.1 (#441)

## [v1.114.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.114.1) - 2023-11-16

- [`fed6f7a`](https://github.com/alexfalkowski/go-service/commit/fed6f7ae830fa96edfd2275f26695697919f952e) fix(deps): bump github.com/klauspost/compress from 1.17.2 to 1.17.3 (#440)

## [v1.114.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.114.0) - 2023-11-16

- [`258d2d8`](https://github.com/alexfalkowski/go-service/commit/258d2d8fe494e141dc0071797988be571e30bd08) feat: add keepalive (#439)

## [v1.113.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.113.1) - 2023-11-15

- [`8039f9b`](https://github.com/alexfalkowski/go-service/commit/8039f9bb22862a30023977b6fd40569ea5513028) fix(grpc): remove tags (#438)
- [`48ff14b`](https://github.com/alexfalkowski/go-service/commit/48ff14bb5e8b6faf42a3eea1791f489fde72e1c0) build(deps): update bin (#437)

## [v1.113.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.113.0) - 2023-11-15

- [`2df826d`](https://github.com/alexfalkowski/go-service/commit/2df826d3545466a3ad533b6f1b66b1cf43f8f880) feat(security): enable mtls (#436)
- [`94431fa`](https://github.com/alexfalkowski/go-service/commit/94431fa10333389c8ae1ee174a793fdefd7d5b1e) build(deps): update bin (#434)

## [v1.112.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.112.3) - 2023-11-14

- [`5f0df1e`](https://github.com/alexfalkowski/go-service/commit/5f0df1e997cf3e28cafbc6c83d731707148edf6d) fix(metrics): use int64 (#433)

## [v1.112.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.112.2) - 2023-11-14

- [`7b256df`](https://github.com/alexfalkowski/go-service/commit/7b256df3a73f7a11ef35480fe8d11e3f1684cdfa) fix(telemetry): make sure we instrument all observables (#432)

## [v1.112.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.112.1) - 2023-11-13

- [`fe5c8a2`](https://github.com/alexfalkowski/go-service/commit/fe5c8a260cbcbfd396ba09f9b67a2e97d845ea0c) fix(transport): error on shutdown (#431)

## [v1.112.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.112.0) - 2023-11-13

- [`0751baa`](https://github.com/alexfalkowski/go-service/commit/0751baa45057a091bcfa28f28bd9b27894944517) feat(os): add get from env (#430)

## [v1.111.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.111.0) - 2023-11-12

- [`7b46648`](https://github.com/alexfalkowski/go-service/commit/7b46648b8444c240c0d5a7059aaed0e6e56507a9) feat(debug): add gopsutil (#428)
- [`2257bc0`](https://github.com/alexfalkowski/go-service/commit/2257bc0186eac5965f3f311c7f9016f70deafe1a) ci(deps): use postgres 15 (#426)
- [`a0fbfe8`](https://github.com/alexfalkowski/go-service/commit/a0fbfe855b70e0b5b4732b5806ca11b7ca6898e9) docs(debug): add pprof (#425)

## [v1.110.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.110.0) - 2023-11-11

- [`149b102`](https://github.com/alexfalkowski/go-service/commit/149b102db1038dce3ca75762c5c7ed27eb51bf23) feat(debug): add pprof (#424)

## [v1.109.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.109.0) - 2023-11-11

- [`1ddf125`](https://github.com/alexfalkowski/go-service/commit/1ddf125875548b0a4a80c55635d71f1eaca68ba0) feat(debug): add statsviz (#423)

## [v1.108.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.108.6) - 2023-11-10

- [`4e6d48e`](https://github.com/alexfalkowski/go-service/commit/4e6d48e4980e7435c6ed6889785bb0830e904c31) fix(deps): bump the otel group with 4 updates (#420)
- [`8cc94b2`](https://github.com/alexfalkowski/go-service/commit/8cc94b22bf6d270c6b455c821175643c71f75cc4) build(deps): group otel (#419)
- [`d3ef706`](https://github.com/alexfalkowski/go-service/commit/d3ef706411e15a28ba76fc5b6b280dc3a7a0bd27) ci: fix versions (#418)

## [v1.108.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.108.5) - 2023-11-09

- [`221d4a3`](https://github.com/alexfalkowski/go-service/commit/221d4a3a18bed98f4b61e6a5a05a8b2461013366) fix(transport): check metadata (#417)

## [v1.108.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.108.4) - 2023-11-09

- [`bce1b02`](https://github.com/alexfalkowski/go-service/commit/bce1b021043c42b8a19397412351201e6178b769) fix(grpc): set the user-agent (#416)

## [v1.108.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.108.3) - 2023-11-09

- [`d47fa6f`](https://github.com/alexfalkowski/go-service/commit/d47fa6fd6f075d3dc2d80021a2994a0441ca0604) fix(transport): use user-agent (#415)
- [`13f2701`](https://github.com/alexfalkowski/go-service/commit/13f2701217bc34526941f7d6c66ba61e0e3c3c1a) test: for sql (#414)
- [`9bd7fad`](https://github.com/alexfalkowski/go-service/commit/9bd7fadebf5499e84e5cf46eb476e7138f88027f) test: write config file (#413)

## [v1.108.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.108.2) - 2023-11-09

- [`27c2d59`](https://github.com/alexfalkowski/go-service/commit/27c2d594e3ce795f80351498df386c9286f688bf) fix(deps): bump github.com/hashicorp/go-retryablehttp (#412)

## [v1.108.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.108.1) - 2023-11-09

- [`715808f`](https://github.com/alexfalkowski/go-service/commit/715808f64fc38450c1762c90d3401bcd0260fb1c) fix(config): unmarshal in return (#411)

## [v1.108.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.108.0) - 2023-11-08

- [`a74a0c1`](https://github.com/alexfalkowski/go-service/commit/a74a0c1232fb7fd70c59158a9afbe7e6dbd7ab1a) feat(security): remove oauth (#410)

## [v1.107.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.107.0) - 2023-11-08

- [`c2954ce`](https://github.com/alexfalkowski/go-service/commit/c2954ce02e4c5e4a63d93e3401271407760d17d4) feat(security): have a generic way to use auth (#409)
- [`cdbce53`](https://github.com/alexfalkowski/go-service/commit/cdbce537916f0c3efbfe1de277c75ef90f46a264) build(deps): update bin (#408)

## [v1.106.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.106.1) - 2023-11-07

- [`03fb44a`](https://github.com/alexfalkowski/go-service/commit/03fb44a752f106082e05083382c760ba6988f119) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#407)

## [v1.106.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.106.0) - 2023-11-06

- [`c4109a1`](https://github.com/alexfalkowski/go-service/commit/c4109a16b0f463c585dd62f2fe6c9ab71f5e7bd2) feat(nsq): move consumer/producer (#406)

## [v1.105.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.105.0) - 2023-11-06

- [`cd52873`](https://github.com/alexfalkowski/go-service/commit/cd5287312343b89e00c4570182f0af852f293728) feat(transport): allow clients to be created externally (#405)

## [v1.104.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.104.1) - 2023-11-06

- [`2a328f8`](https://github.com/alexfalkowski/go-service/commit/2a328f8a9cdecf870e280a85436f89f78a29c594) fix(grpc): rename packages (#404)

## [v1.104.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.104.0) - 2023-11-06

- [`4d897b3`](https://github.com/alexfalkowski/go-service/commit/4d897b3ce8c14aaaa1180bcd608ef212dc99d774) feat(telemetry): want to sync all the time in developement mode (#403)

## [v1.103.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.103.3) - 2023-11-06

- [`50320b6`](https://github.com/alexfalkowski/go-service/commit/50320b6c2fe3a38a23f41900b35ca09a11ed4c62) fix(deps): bump github.com/spf13/cobra from 1.7.0 to 1.8.0 (#402)
- [`dc6dff9`](https://github.com/alexfalkowski/go-service/commit/dc6dff94bbfc22b82bcbaef9157b3c06b7070e2b) build(deps): update bin (#401)
- [`f5edd4d`](https://github.com/alexfalkowski/go-service/commit/f5edd4d488012c55021eefa1ce9c33366097c790) build(deps): update bin (#400)

## [v1.103.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.103.2) - 2023-11-03

- [`bcd5921`](https://github.com/alexfalkowski/go-service/commit/bcd59213731cc85b56d0d8d54b85c43421129759) fix(transport): logger level (#399)

## [v1.103.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.103.1) - 2023-11-02

- [`71b426c`](https://github.com/alexfalkowski/go-service/commit/71b426c6cdcee2d4bbae3f2373f7a8a3093bd30b) fix(meta): remove remote address (#398)
- [`acb3bb3`](https://github.com/alexfalkowski/go-service/commit/acb3bb3d6042645a64fc2dc88eb0f378d55c9905) test: move way from httpstatus (#397)

## [v1.103.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.103.0) - 2023-11-02

- [`a1a41ec`](https://github.com/alexfalkowski/go-service/commit/a1a41ecf3256597246a05f2ddbef68748cf8628d) feat(security): move auth0 to oauth (#395)

## [v1.102.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.102.0) - 2023-11-02

- [`4ba0951`](https://github.com/alexfalkowski/go-service/commit/4ba095168e878550302008832288babfb63afd93) feat(config): encode config in env variable (#396)

## [v1.101.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.101.0) - 2023-11-01

- [`f844bef`](https://github.com/alexfalkowski/go-service/commit/f844bef6909baf25bd50077f745723e12541b938) feat(http): add marshal options to mux (#394)

## [v1.100.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.100.2) - 2023-10-27

- [`75f53e5`](https://github.com/alexfalkowski/go-service/commit/75f53e5bd8a4f260a9b8b32a349c65e22263221b) fix(deps): bump github.com/google/uuid from 1.3.1 to 1.4.0 (#393)

## [v1.100.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.100.1) - 2023-10-27

- [`09a6123`](https://github.com/alexfalkowski/go-service/commit/09a6123d0e0a9e653aba1d560298d8af0012c06c) fix(deps): bump github.com/vmihailenco/msgpack/v5 from 5.4.0 to 5.4.1 (#392)

## [v1.100.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.100.0) - 2023-10-26

- [`2354314`](https://github.com/alexfalkowski/go-service/commit/2354314b5d29aef6651dd97353df24d5223ec935) feat: add environment (#391)

## [v1.99.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.99.0) - 2023-10-26

- [`43f7a1c`](https://github.com/alexfalkowski/go-service/commit/43f7a1c8ac9ad02266ba551545cf9110cf1ed0c9) feat(logger): use kind (#390)

## [v1.98.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.98.0) - 2023-10-26

- [`38cea71`](https://github.com/alexfalkowski/go-service/commit/38cea719c4b058fa0a78e5b9a8e70202733a0883) feat(grpc): register reflection (#389)

## [v1.97.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.97.1) - 2023-10-24

- [`71c9a4b`](https://github.com/alexfalkowski/go-service/commit/71c9a4ba88891f166261fc357333bb05ce5e8467) fix(deps): bump github.com/klauspost/compress from 1.17.1 to 1.17.2 (#388)

## [v1.97.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.97.0) - 2023-10-19

- [`893057d`](https://github.com/alexfalkowski/go-service/commit/893057dec5bafd65d7c2bfbbe5deb014fb4aa016) feat(logger): add ability to configure (#387)

## [v1.96.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.96.5) - 2023-10-18

- [`6004340`](https://github.com/alexfalkowski/go-service/commit/6004340c3cbcf446756a91c1a60dc5d6657fb42b) fix(deps): bump google.golang.org/grpc from 1.58.3 to 1.59.0 (#385)

## [v1.96.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.96.4) - 2023-10-18

- [`5aca179`](https://github.com/alexfalkowski/go-service/commit/5aca179dc55ba0bbd8d519f7989d3dc4d5a6f675) fix(deps): bump go.uber.org/fx from 1.20.0 to 1.20.1 (#386)

## [v1.96.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.96.3) - 2023-10-16

- [`b3673b3`](https://github.com/alexfalkowski/go-service/commit/b3673b3691537a4cd4822d18127fdb89629eb4c2) fix(config): make sure we unmarshal before we give a config (#384)

## [v1.96.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.96.2) - 2023-10-16

- [`fa53376`](https://github.com/alexfalkowski/go-service/commit/fa53376e541c8d691a307ba9d1507f198b5e453b) fix(deps): bump github.com/klauspost/compress from 1.17.0 to 1.17.1 (#383)

## [v1.96.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.96.1) - 2023-10-14

- [`d4484eb`](https://github.com/alexfalkowski/go-service/commit/d4484ebfa08bc39d01e8a81475bf7d43861698c4) fix: deps mod (#382)

## [v1.96.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.96.0) - 2023-10-14

- [`c5ef5c5`](https://github.com/alexfalkowski/go-service/commit/c5ef5c5392d1560286b0141b4449920484bb4b9d) feat(telemetry): use otel for metrics (#379)

## [v1.95.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.95.2) - 2023-10-11

- [`e563f7d`](https://github.com/alexfalkowski/go-service/commit/e563f7de05bb158ecbaff118dce408c1eaa14a94) fix(deps): bump google.golang.org/grpc from 1.58.2 to 1.58.3 (#381)

## [v1.95.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.95.1) - 2023-10-11

- [`113f432`](https://github.com/alexfalkowski/go-service/commit/113f432324decdc83eca76791836a1da3c53d1a5) fix(deps): bump golang.org/x/net from 0.16.0 to 0.17.0 (#380)
- [`0c7a299`](https://github.com/alexfalkowski/go-service/commit/0c7a299c895a0799cf5ceae0c90f8afd49e9b047) build: remove tools (#378)
- [`6940911`](https://github.com/alexfalkowski/go-service/commit/69409116a7ca6a015a7e117ecbcacbf4e71de39e) build: update bin (#377)
- [`456f42c`](https://github.com/alexfalkowski/go-service/commit/456f42ced3da48103643ba26b963a09642bdfc37) ci: use resource_class large (#376)

## [v1.95.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.95.0) - 2023-10-06

- [`6a470bd`](https://github.com/alexfalkowski/go-service/commit/6a470bd9faa9032744bd798b31efb77904acd078) feat(grpc): add ability to enable or disable (#375)

## [v1.94.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.94.2) - 2023-10-06

- [`691ea14`](https://github.com/alexfalkowski/go-service/commit/691ea14cf367b0bd8719aae782ece66529828439) fix(transport): validate port (#374)

## [v1.94.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.94.1) - 2023-10-06

- [`9025d66`](https://github.com/alexfalkowski/go-service/commit/9025d668c613b0f2fcc8bf80937ca98de604b70f) fix(deps): bump golang.org/x/net from 0.15.0 to 0.16.0 (#373)

## [v1.94.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.94.0) - 2023-10-05

- [`a6b9db1`](https://github.com/alexfalkowski/go-service/commit/a6b9db162ff93c42552dfc2444a8217ba949b4e6) feat: remove some params (#372)
- [`15aa7ae`](https://github.com/alexfalkowski/go-service/commit/15aa7ae4195aba691978cb48516d24602815603d) docs: change ports for transport (#371)

## [v1.93.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.93.0) - 2023-10-05

- [`a9158e8`](https://github.com/alexfalkowski/go-service/commit/a9158e8fb63d9b7931ce01e41824d0a99722480b) feat(transport): separate ports (#370)

## [v1.92.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.92.2) - 2023-10-02

- [`9057ccc`](https://github.com/alexfalkowski/go-service/commit/9057cccd2a21b96bf246c45fb9b9e451a6ec9ce7) fix(deps): bump github.com/rs/cors from 1.10.0 to 1.10.1 (#368)

## [v1.92.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.92.1) - 2023-10-02

- [`16ff5df`](https://github.com/alexfalkowski/go-service/commit/16ff5df9dcf6411c764177a0ecf2a4e0b3afc11f) fix(deps): bump github.com/vmihailenco/msgpack/v5 from 5.3.5 to 5.4.0 (#369)

## [v1.92.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.92.0) - 2023-09-29

- [`2ce4910`](https://github.com/alexfalkowski/go-service/commit/2ce4910e77145eae4abb13173654c210c18c421f) feat(deps): update otel to v1.19.0 (#367)

## [v1.91.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.91.1) - 2023-09-28

- [`78ed3da`](https://github.com/alexfalkowski/go-service/commit/78ed3da10404d1d368a4684a58e01dfae7b08756) fix(deps): bump github.com/prometheus/client_golang (#361)

## [v1.91.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.91.0) - 2023-09-22

- [`58aae85`](https://github.com/alexfalkowski/go-service/commit/58aae85bebade7d31b86f06743cdf92066d54a54) feat(telemetry): move to a single package (#360)

## [v1.90.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.90.4) - 2023-09-22

- [`fec0fcc`](https://github.com/alexfalkowski/go-service/commit/fec0fccde7b379e7f3117feab8043a5dc02fdd86) fix(deps): bump google.golang.org/grpc from 1.58.1 to 1.58.2 (#359)

## [v1.90.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.90.3) - 2023-09-21

- [`0d0cac5`](https://github.com/alexfalkowski/go-service/commit/0d0cac51b90a2662e67631b102d0bd32cec76e7b) fix(errors): use multierr (#358)

## [v1.90.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.90.2) - 2023-09-21

- [`203bc96`](https://github.com/alexfalkowski/go-service/commit/203bc96adfb631eff6c2e665796c4fbfc8d7c5cb) fix(metrics): use consistent naming (#357)

## [v1.90.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.90.1) - 2023-09-21

- [`0bc2866`](https://github.com/alexfalkowski/go-service/commit/0bc28664547c018ecb75915caaad9a18ff44a927) fix(fx): simplify modules (#356)
- [`d5d14fc`](https://github.com/alexfalkowski/go-service/commit/d5d14fc97f66383b94361fc3433417e01e66a728) Revert "feat(metrics): register through DI (#354)" (#355)

## [v1.90.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.90.0) - 2023-09-21

- [`d56e71a`](https://github.com/alexfalkowski/go-service/commit/d56e71a47a40d5dced836ded6b7aaf4c8eefbb18) feat(metrics): register through DI (#354)

## [v1.89.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.89.5) - 2023-09-20

- [`19dd81d`](https://github.com/alexfalkowski/go-service/commit/19dd81db96b61648141609c02d097c23bfa49fa7) fix(deps): bump github.com/rs/cors from 1.9.0 to 1.10.0 (#353)

## [v1.89.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.89.4) - 2023-09-19

- [`f7abdca`](https://github.com/alexfalkowski/go-service/commit/f7abdcad80963162e7bda842847e3a9e0267029c) fix(deps): bump github.com/klauspost/compress from 1.16.7 to 1.17.0 (#348)

## [v1.89.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.89.3) - 2023-09-19

- [`96543de`](https://github.com/alexfalkowski/go-service/commit/96543de269801ef1130ee314e73e83bee90015f9) fix(deps): bump google.golang.org/grpc from 1.58.0 to 1.58.1 (#349)

## [v1.89.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.89.2) - 2023-09-19

- [`5a914c5`](https://github.com/alexfalkowski/go-service/commit/5a914c5309f7d987f6af6f7e3ee5690febbba192) fix(deps): bump go.uber.org/zap from 1.25.0 to 1.26.0 (#350)

## [v1.89.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.89.1) - 2023-09-19

- [`a0f9e5d`](https://github.com/alexfalkowski/go-service/commit/a0f9e5dd0c112244c561b6c1f05591ac0a5461aa) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#352)

## [v1.89.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.89.0) - 2023-09-13

- [`2b292fe`](https://github.com/alexfalkowski/go-service/commit/2b292fe561e3d54f6a7f6992b0294e93e97fd3f9) feat(otel): support native protocol (#345)
- [`dc97a62`](https://github.com/alexfalkowski/go-service/commit/dc97a628e0397a865fd947e644d5b138fff844e4) docs: fix section for runtime (#344)

## [v1.88.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.88.0) - 2023-09-12

- [`1b6d5e7`](https://github.com/alexfalkowski/go-service/commit/1b6d5e76e58fc60f4bd2668be235d4da588650e0) feat: add go.uber.org/automaxprocs (#343)

## [v1.87.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.87.3) - 2023-08-23

- [`c88cea5`](https://github.com/alexfalkowski/go-service/commit/c88cea50c1919ae1f17f90e7c358f3ba8cbc3582) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#342)

## [v1.87.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.87.2) - 2023-08-22

- [`60b30ac`](https://github.com/alexfalkowski/go-service/commit/60b30accaf73d99d05fd6522375143f4c1d5a456) fix(deps): bump github.com/google/uuid from 1.3.0 to 1.3.1 (#341)
- [`8a96df5`](https://github.com/alexfalkowski/go-service/commit/8a96df58395973b83e02c48f57163111e8828de8) test: use net for random port (#340)

## [v1.87.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.87.1) - 2023-08-10

- [`542c0c8`](https://github.com/alexfalkowski/go-service/commit/542c0c8b8f74494b4c46b5bb109e03f2eb6b03e7) fix(deps): bump github.com/alexfalkowski/go-health from 1.12.2 to 1.13.0 (#339)

## [v1.87.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.87.0) - 2023-08-10

- [`92acfb5`](https://github.com/alexfalkowski/go-service/commit/92acfb5e583c950febdf5f219e94058b9609804c) feat: update go to version 1.21 (#338)

## [v1.86.46](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.46) - 2023-08-07

- [`d6adebd`](https://github.com/alexfalkowski/go-service/commit/d6adebde3ec55306bd6f28263b6a2aad34311864) fix(deps): bump golang.org/x/net from 0.13.0 to 0.14.0 (#337)

## [v1.86.45](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.45) - 2023-08-03

- [`97a8ff4`](https://github.com/alexfalkowski/go-service/commit/97a8ff45ef664af053534baaccd368877ddd863c) fix(deps): bump go.uber.org/zap from 1.24.0 to 1.25.0 (#336)

## [v1.86.44](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.44) - 2023-08-02

- [`4f3b0ba`](https://github.com/alexfalkowski/go-service/commit/4f3b0baeda00a1c9c9d9c1353ee4d6e5de7b07c7) fix(deps): bump golang.org/x/net from 0.12.0 to 0.13.0 (#335)

## [v1.86.43](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.43) - 2023-07-31

- [`7a6e9e8`](https://github.com/alexfalkowski/go-service/commit/7a6e9e8cad71bfce21ab5da0edf48f7a2824a6ec) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#334)

## [v1.86.42](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.42) - 2023-07-27

- [`5474bf8`](https://github.com/alexfalkowski/go-service/commit/5474bf8444093dca82be7efe433153a5aebc580b) fix(deps): bump google.golang.org/grpc from 1.56.2 to 1.57.0 (#333)

## [v1.86.41](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.41) - 2023-07-25

- [`824f04c`](https://github.com/alexfalkowski/go-service/commit/824f04ca9af6c83671d6fc0eac822b7ce7a36c83) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#331)
- [`17a14e8`](https://github.com/alexfalkowski/go-service/commit/17a14e8c4816bd0b6d2d384d93da40b1ddac89ae) build(deps): update bin (#332)

## [v1.86.40](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.40) - 2023-07-07

- [`6632993`](https://github.com/alexfalkowski/go-service/commit/66329930443897b8ce522f48967b7b308dfd8e0b) fix(deps): bump google.golang.org/grpc from 1.56.1 to 1.56.2 (#330)

## [v1.86.39](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.39) - 2023-07-06

- [`b9a40a0`](https://github.com/alexfalkowski/go-service/commit/b9a40a03640ef87b51507df8ca665a923409be66) fix(deps): bump golang.org/x/net from 0.11.0 to 0.12.0 (#329)

## [v1.86.38](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.38) - 2023-07-04

- [`d9bc5df`](https://github.com/alexfalkowski/go-service/commit/d9bc5df9a70ee7f00bebdf6a307b76107cc9e424) fix(deps): bump github.com/klauspost/compress from 1.16.6 to 1.16.7 (#328)

## [v1.86.37](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.37) - 2023-06-30

- [`caedb32`](https://github.com/alexfalkowski/go-service/commit/caedb32393de961f3cd54ade3cb0495a0927da89) fix(deps): bump github.com/alexfalkowski/go-health from 1.12.1 to 1.12.2 (#327)

## [v1.86.36](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.36) - 2023-06-29

- [`8da16e6`](https://github.com/alexfalkowski/go-service/commit/8da16e60537c10eb26ae6b5ef9b488169bed7969) fix(deps): bump github.com/smartystreets/goconvey from 1.8.0 to 1.8.1 (#326)

## [v1.86.35](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.35) - 2023-06-27

- [`3dc7e16`](https://github.com/alexfalkowski/go-service/commit/3dc7e160b18cdf51167c38206e199709a02d407b) fix(deps): bump google.golang.org/protobuf from 1.30.0 to 1.31.0 (#325)

## [v1.86.34](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.34) - 2023-06-22

- [`7f083ba`](https://github.com/alexfalkowski/go-service/commit/7f083ba83a4effbba16c21189a061ec975589c32) fix(deps): bump google.golang.org/grpc from 1.56.0 to 1.56.1 (#324)

## [v1.86.33](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.33) - 2023-06-16

- [`191b14d`](https://github.com/alexfalkowski/go-service/commit/191b14d4ff1a3c4574d61aaf33ef49b9c2939b3d) fix(deps): bump github.com/prometheus/client_golang (#323)

## [v1.86.32](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.32) - 2023-06-16

- [`3aa7efa`](https://github.com/alexfalkowski/go-service/commit/3aa7efae2dd600a6e55490e14980b774ff499afc) fix(deps): bump google.golang.org/grpc from 1.55.0 to 1.56.0 (#322)

## [v1.86.31](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.31) - 2023-06-14

- [`61fbb7d`](https://github.com/alexfalkowski/go-service/commit/61fbb7d11f095882a43b81c9aef25fad977757f9) fix(deps): bump golang.org/x/net from 0.10.0 to 0.11.0 (#321)

## [v1.86.30](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.30) - 2023-06-14

- [`b94f2bc`](https://github.com/alexfalkowski/go-service/commit/b94f2bc5208d0f049ff80e9ad8e92b04afb6721b) fix(deps): bump github.com/klauspost/compress from 1.16.5 to 1.16.6 (#320)

## [v1.86.29](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.29) - 2023-06-12

- [`3352a0c`](https://github.com/alexfalkowski/go-service/commit/3352a0c589533c7a8035f41af2232c0cc4be8fa6) fix(deps): bump go.uber.org/fx from 1.19.3 to 1.20.0 (#318)
- [`21e4553`](https://github.com/alexfalkowski/go-service/commit/21e45535abee60835bdc1b16ff0d14e5d0905589) build(deps): update bin (#319)

## [v1.86.28](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.28) - 2023-06-12

- [`93bd88c`](https://github.com/alexfalkowski/go-service/commit/93bd88cdf6f8f23e979a08fa1c32d6d8de520a43) fix(deps): bump github.com/BurntSushi/toml from 1.3.1 to 1.3.2 (#316)

## [v1.86.27](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.27) - 2023-06-12

- [`6e52a6d`](https://github.com/alexfalkowski/go-service/commit/6e52a6d5f3ededdaaec72b071af578c15f306316) fix(deps): bump github.com/golang-migrate/migrate/v4 (#317)

## [v1.86.26](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.26) - 2023-06-07

- [`82b5086`](https://github.com/alexfalkowski/go-service/commit/82b5086d3bc3bf1c9db0aaf6eac75861d3c796d3) fix(deps): bump github.com/hashicorp/go-retryablehttp (#315)

## [v1.86.25](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.25) - 2023-06-07

- [`d8bed4d`](https://github.com/alexfalkowski/go-service/commit/d8bed4dc6fba12904608b34cffa3190aa4ad4208) fix(deps): bump github.com/BurntSushi/toml from 1.3.0 to 1.3.1 (#314)

## [v1.86.24](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.24) - 2023-06-06

- [`1be74ae`](https://github.com/alexfalkowski/go-service/commit/1be74ae59a6601d6ba6d6b734e247c202a51aaf9) fix(deps): bump github.com/golang-migrate/migrate/v4 (#313)

## [v1.86.23](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.23) - 2023-06-01

- [`8bc1dc0`](https://github.com/alexfalkowski/go-service/commit/8bc1dc0c53897c7b306063e53cb5491a15a22d00) fix(deps): bump github.com/BurntSushi/toml from 1.2.1 to 1.3.0 (#311)

## [v1.86.22](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.22) - 2023-06-01

- [`3ad82fd`](https://github.com/alexfalkowski/go-service/commit/3ad82fdb023867f28cbd3750d087b57ff9d79781) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#312)

## [v1.86.21](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.21) - 2023-05-29

- [`1b71042`](https://github.com/alexfalkowski/go-service/commit/1b71042541cea5383f099040b7f675aad147bdd9) fix(deps): bump github.com/golang-migrate/migrate/v4 (#310)

## [v1.86.20](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.20) - 2023-05-25

- [`f8b71df`](https://github.com/alexfalkowski/go-service/commit/f8b71df1bac5bcf473806160191ad72170f12a86) fix(deps): bump github.com/ulule/limiter/v3 from 3.11.1 to 3.11.2 (#309)

## [v1.86.19](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.19) - 2023-05-23

- [`36d7219`](https://github.com/alexfalkowski/go-service/commit/36d7219a54764546d490ee1132c3176b87f45a61) fix(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#306)

## [v1.86.18](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.18) - 2023-05-23

- [`5adbfaa`](https://github.com/alexfalkowski/go-service/commit/5adbfaa5cba73f26ece472ec59eac03f23e210e1) fix(deps): bump go.opentelemetry.io/otel/sdk from 1.15.1 to 1.16.0 (#307)

## [v1.86.17](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.17) - 2023-05-09

- [`2a4462a`](https://github.com/alexfalkowski/go-service/commit/2a4462acacea1f0b3e7323b6c25c868317397397) fix(deps): bump go.uber.org/fx from 1.19.2 to 1.19.3 (#303)

## [v1.86.16](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.16) - 2023-05-09

- [`7d0ca18`](https://github.com/alexfalkowski/go-service/commit/7d0ca1835b45c748d3781373aa2008ddecc9739a) fix(deps): bump golang.org/x/net from 0.9.0 to 0.10.0 (#304)

## [v1.86.15](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.15) - 2023-05-05

- [`35300e8`](https://github.com/alexfalkowski/go-service/commit/35300e88e6d8793ef56ae5804bf9c5d779ad7d5b) fix(deps): bump google.golang.org/grpc from 1.54.0 to 1.55.0 (#302)

## [v1.86.14](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.14) - 2023-05-04

- [`2ded4af`](https://github.com/alexfalkowski/go-service/commit/2ded4afd2223890b93b5231d6decfccfded3b79c) fix(deps): bump github.com/prometheus/client_golang (#301)

## [v1.86.13](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.13) - 2023-05-03

- [`2c2a892`](https://github.com/alexfalkowski/go-service/commit/2c2a892486adcdbab33a536e99fd7a2a672b3227) fix(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#300)

## [v1.86.12](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.12) - 2023-05-03

- [`524bae5`](https://github.com/alexfalkowski/go-service/commit/524bae59a5082040148c058f95d78d5713729c94) fix(deps): bump go.opentelemetry.io/otel/sdk from 1.15.0 to 1.15.1 (#299)

## [v1.86.11](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.11) - 2023-05-01

- [`9e1762c`](https://github.com/alexfalkowski/go-service/commit/9e1762cdaefaab51eab1972ea16bd45bf1d15a60) fix(deps): bump go.opentelemetry.io/otel/exporters/jaeger (#293)

## [v1.86.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.10) - 2023-05-01

- [`1d1fe61`](https://github.com/alexfalkowski/go-service/commit/1d1fe61b00c049966691af381d84671b6197b64e) fix(deps): bump go.opentelemetry.io/otel/sdk from 1.14.0 to 1.15.0 (#295)

## [v1.86.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.9) - 2023-04-18

- [`0034b0e`](https://github.com/alexfalkowski/go-service/commit/0034b0e3ffbab187f4d6a5b6cdcef66d58704e8f) fix(deps): bump github.com/klauspost/compress from 1.16.4 to 1.16.5 (#292)

## [v1.86.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.8) - 2023-04-14

- [`289e26b`](https://github.com/alexfalkowski/go-service/commit/289e26b37e8ceb8fd059fbbc4cbf65e7b3534e9f) fix(deps): bump github.com/prometheus/client_golang (#291)

## [v1.86.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.7) - 2023-04-14

- [`ced9627`](https://github.com/alexfalkowski/go-service/commit/ced9627d966d6606ee00c59dd215d78616ccd26e) fix(deps): bump github.com/rs/cors from 1.8.3 to 1.9.0 (#290)
- [`1c893ef`](https://github.com/alexfalkowski/go-service/commit/1c893ef4e1e82e65385cc1fc637354ee3c5e8801) test: update grpc-gateway to v2.15.0-1 (#289)

## [v1.86.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.6) - 2023-04-12

- [`7244b60`](https://github.com/alexfalkowski/go-service/commit/7244b605a8b0a594c44b01598129574b95dcd68d) fix(deps): bump github.com/alexfalkowski/go-health from 1.12.0 to 1.12.1 (#288)
- [`da5a45c`](https://github.com/alexfalkowski/go-service/commit/da5a45cd3d3d3bbc77480f980421316f931c3379) build(deps): update bin (#287)
- [`3c920ab`](https://github.com/alexfalkowski/go-service/commit/3c920aba8707c9426e148abf72a5b44d72f26a6c) build(deps): update bin (#286)

## [v1.86.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.5) - 2023-04-10

- [`cd52c6e`](https://github.com/alexfalkowski/go-service/commit/cd52c6ef84f5f2f04a3870753c3888fa60f59dbc) fix(deps): bump github.com/smartystreets/goconvey from 1.7.2 to 1.8.0 (#285)

## [v1.86.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.4) - 2023-04-07

- [`c49f55d`](https://github.com/alexfalkowski/go-service/commit/c49f55d21726c9ce366553bbb9ff75009e36d641) fix(deps): bump golang.org/x/net from 0.8.0 to 0.9.0 (#284)

## [v1.86.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.3) - 2023-04-06

- [`1f55fb4`](https://github.com/alexfalkowski/go-service/commit/1f55fb4955a77e1e3dac0503d2482ff765b7aa99) fix(deps): bump github.com/klauspost/compress from 1.16.3 to 1.16.4 (#283)

## [v1.86.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.2) - 2023-04-05

- [`0dc3c6e`](https://github.com/alexfalkowski/go-service/commit/0dc3c6e509a8b1156bae8ec8f0b4bf66af70a1de) fix(deps): bump github.com/spf13/cobra from 1.6.1 to 1.7.0 (#282)

## [v1.86.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.1) - 2023-03-31

- [`5ed4922`](https://github.com/alexfalkowski/go-service/commit/5ed492247c46520c60cb83f72b6d21dbea11d9f6) fix: the deps (#281)

## [v1.86.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.86.0) - 2023-03-31

- [`88baddc`](https://github.com/alexfalkowski/go-service/commit/88baddc9fb9be66dfe704e4ef5758fe12a112ca3) feat: replace opentracing with otel (#279)

## [v1.85.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.85.2) - 2023-03-30

- [`7ec3bb3`](https://github.com/alexfalkowski/go-service/commit/7ec3bb3e55af82cd04f0f00cf1f734152bcda44b) fix(deps): bump go.uber.org/multierr from 1.10.0 to 1.11.0 (#280)

## [v1.85.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.85.1) - 2023-03-22

- [`0a828bd`](https://github.com/alexfalkowski/go-service/commit/0a828bde1878f5659fe9fd9d776ee9377104ece1) fix(deps): bump google.golang.org/grpc from 1.53.0 to 1.54.0 (#277)
- [`bb5e1b1`](https://github.com/alexfalkowski/go-service/commit/bb5e1b12adb848804338351e33b08875d7edebf0) style: ignore unused (#278)

## [v1.85.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.85.0) - 2023-03-18

- [`1b4a160`](https://github.com/alexfalkowski/go-service/commit/1b4a1602cf8956ac68c864bdc336f5c82df62eb9) feat(transport): move strings (#276)

## [v1.84.20](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.20) - 2023-03-16

- [`e55123c`](https://github.com/alexfalkowski/go-service/commit/e55123c17640c91884e2fc1dff023080e8dd6066) fix(deps): bump google.golang.org/protobuf from 1.29.1 to 1.30.0 (#275)

## [v1.84.19](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.19) - 2023-03-16

- [`c54a587`](https://github.com/alexfalkowski/go-service/commit/c54a587c32ece27ba1e89e843dfbc816714a2a9c) fix(deps): bump github.com/grpc-ecosystem/go-grpc-middleware (#274)
- [`9ad9928`](https://github.com/alexfalkowski/go-service/commit/9ad9928d0276ae18714af0bb3877d9a71e416596) build(deps): update bin (#273)
- [`c0c5ab5`](https://github.com/alexfalkowski/go-service/commit/c0c5ab52f384d9daaaca73c1c89171bcdb4803bd) build(deps): update bin (#272)

## [v1.84.18](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.18) - 2023-03-15

- [`6658e99`](https://github.com/alexfalkowski/go-service/commit/6658e99277c47ee2e338a69f0dd98016c6877765) fix(deps): bump google.golang.org/protobuf from 1.29.0 to 1.29.1 (#271)

## [v1.84.17](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.17) - 2023-03-14

- [`4e81b3d`](https://github.com/alexfalkowski/go-service/commit/4e81b3daf25080102c6fe45809611b8a7327a1af) fix(deps): bump github.com/klauspost/compress from 1.16.0 to 1.16.3 (#270)

## [v1.84.16](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.16) - 2023-03-08

- [`5fa4c77`](https://github.com/alexfalkowski/go-service/commit/5fa4c775747b850e46dabfd3f298540a883e0842) fix(deps): bump google.golang.org/protobuf from 1.28.1 to 1.29.0 (#268)

## [v1.84.15](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.15) - 2023-03-08

- [`2e73630`](https://github.com/alexfalkowski/go-service/commit/2e73630fcb772b46589832e89fa2d541c3ad4563) fix(deps): bump go.uber.org/multierr from 1.9.0 to 1.10.0 (#269)

## [v1.84.14](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.14) - 2023-03-08

- [`67e3f46`](https://github.com/alexfalkowski/go-service/commit/67e3f463616aef5969fb5f2feb4c9a2c1d71405b) fix(deps): bump github.com/ulule/limiter/v3 from 3.11.0 to 3.11.1 (#267)

## [v1.84.13](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.13) - 2023-03-06

- [`b3c7273`](https://github.com/alexfalkowski/go-service/commit/b3c72730d65a3edcc69c8b1493243fcb485cee5d) fix(deps): bump golang.org/x/net from 0.7.0 to 0.8.0 (#266)

## [v1.84.12](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.12) - 2023-02-28

- [`d5a0094`](https://github.com/alexfalkowski/go-service/commit/d5a0094cd40cefc50e06b44c2d3ff8dbd45a2dce) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#265)

## [v1.84.11](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.11) - 2023-02-28

- [`e90ba41`](https://github.com/alexfalkowski/go-service/commit/e90ba41fc1b65a3e237f9335e6352b013a35073b) fix(deps): bump github.com/jackc/pgx/v4 from 4.18.0 to 4.18.1 (#263)

## [v1.84.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.10) - 2023-02-28

- [`7adebdf`](https://github.com/alexfalkowski/go-service/commit/7adebdf60d230ee5bbac7b2afa38fa3c8c04ab68) fix(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.47.0 to 1.48.0 (#264)

## [v1.84.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.9) - 2023-02-27

- [`c56dacb`](https://github.com/alexfalkowski/go-service/commit/c56dacb9a2a7c2934ef2d667ee00d80550dc7580) fix(deps): bump github.com/klauspost/compress from 1.15.15 to 1.16.0 (#262)

## [v1.84.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.8) - 2023-02-23

- [`255c44e`](https://github.com/alexfalkowski/go-service/commit/255c44e5d3d7726719522d9ed89cbc6cfd5041a0) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#261)

## [v1.84.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.7) - 2023-02-22

- [`28c8c60`](https://github.com/alexfalkowski/go-service/commit/28c8c60794c752d0e0f6c315c6bc08382cb8d437) fix(deps): bump go.uber.org/fx from 1.19.1 to 1.19.2 (#260)

## [v1.84.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.6) - 2023-02-20

- [`357c8a1`](https://github.com/alexfalkowski/go-service/commit/357c8a13b01b466636bb9324a34b997e6a51b806) fix(deps): bump github.com/golang-jwt/jwt/v4 from 4.4.3 to 4.5.0 (#259)

## [v1.84.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.5) - 2023-02-15

- [`820bafa`](https://github.com/alexfalkowski/go-service/commit/820bafa73bbe442243798bd30ef0cac18002d3e0) fix(deps): bump golang.org/x/net from 0.6.0 to 0.7.0 (#256)
- [`380c873`](https://github.com/alexfalkowski/go-service/commit/380c8730ecc0ada808ca34da9dd91da635a35f38) build(deps): update bin (#258)
- [`0d76a6b`](https://github.com/alexfalkowski/go-service/commit/0d76a6b914d120c46c8105aece2356251a862690) test: ignore sqlclosecheck (#257)

## [v1.84.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.4) - 2023-02-13

- [`d249c32`](https://github.com/alexfalkowski/go-service/commit/d249c3294ae33dd223ed26791f5722a7e9ca1298) fix(deps): bump github.com/jackc/pgx/v4 from 4.17.2 to 4.18.0 (#255)

## [v1.84.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.3) - 2023-02-09

- [`a1de21f`](https://github.com/alexfalkowski/go-service/commit/a1de21f533720b070ea28504ab8d3d127f2b3ad6) fix(deps): bump github.com/alexfalkowski/go-health from 1.11.0 to 1.12.0 (#254)

## [v1.84.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.2) - 2023-02-09

- [`06eabcd`](https://github.com/alexfalkowski/go-service/commit/06eabcd3249c23a1b7cf4890f7f8e561a6a93959) fix(deps): bump golang.org/x/net from 0.5.0 to 0.6.0 (#253)

## [v1.84.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.1) - 2023-02-08

- [`e02c985`](https://github.com/alexfalkowski/go-service/commit/e02c98513b7702ac9bf7572cbe3d2f42753f47bf) fix(deps): bump google.golang.org/grpc from 1.52.3 to 1.53.0 (#252)

## [v1.84.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.84.0) - 2023-02-03

- [`6e4d26b`](https://github.com/alexfalkowski/go-service/commit/6e4d26b4a875bb74916396600ec864671f333f67) feat(go): update to v1.20 (#251)

## [v1.83.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.83.2) - 2023-02-01

- [`23ce4cc`](https://github.com/alexfalkowski/go-service/commit/23ce4cce941f4a10113de9599ef9052dd0bc921e) fix(deps): bump github.com/ulule/limiter/v3 from 3.10.0 to 3.11.0 (#249)

## [v1.83.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.83.1) - 2023-02-01

- [`55ce626`](https://github.com/alexfalkowski/go-service/commit/55ce626e562f7d47bc6668b0b3ef0cb692c941af) fix(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.46.1 to 1.47.0 (#250)

## [v1.83.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.83.0) - 2023-01-30

- [`4db0d8b`](https://github.com/alexfalkowski/go-service/commit/4db0d8b05b48a319cc6f2d9e09789c96ff767875) feat(observability): change endpoints (#248)

## [v1.82.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.82.4) - 2023-01-26

- [`8af7491`](https://github.com/alexfalkowski/go-service/commit/8af74916f6906f8b03717c8ae28ed1499df799c3) fix(deps): bump google.golang.org/grpc from 1.52.1 to 1.52.3 (#247)

## [v1.82.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.82.3) - 2023-01-25

- [`1390f3c`](https://github.com/alexfalkowski/go-service/commit/1390f3c50153adae387c6934526241d97f8f7170) fix(deps): bump google.golang.org/grpc from 1.52.0 to 1.52.1 (#246)

## [v1.82.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.82.2) - 2023-01-23

- [`0d7f52b`](https://github.com/alexfalkowski/go-service/commit/0d7f52b249bf0cc62259f9100e19bbc03e02771b) fix(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.46.0 to 1.46.1 (#244)

## [v1.82.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.82.1) - 2023-01-23

- [`e27fd55`](https://github.com/alexfalkowski/go-service/commit/e27fd55c863b6256fff48a20b7989a12fa875673) fix(deps): bump github.com/klauspost/compress from 1.15.14 to 1.15.15 (#245)

## [v1.82.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.82.0) - 2023-01-17

- [`7e5f549`](https://github.com/alexfalkowski/go-service/commit/7e5f54993a8db25b1c7540e256818b807241781e) feat(cmd): use env variable rather than mem flag (#243)
- [`f0d10ff`](https://github.com/alexfalkowski/go-service/commit/f0d10ff1158a59e53ac9198e0b74dfa59d974fa8) build: add github.com:alexfalkowski/bin (#242)
- [`4ac3330`](https://github.com/alexfalkowski/go-service/commit/4ac33308a31945c4d3ed89889caed6eac99a9477) ci(dependabot): change commit message (#241)
- [`8338d58`](https://github.com/alexfalkowski/go-service/commit/8338d58d12c44274a6323b8d2699ceea474bd40b) ci: use release 3.1 (#240)
- [`f238289`](https://github.com/alexfalkowski/go-service/commit/f23828907e13e72ee29a338c86342ab8bc4e761e) ci: use release 3.0 (#239)

## [v1.81.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.81.3) - 2023-01-11

- [`226b066`](https://github.com/alexfalkowski/go-service/commit/226b06629ea966fcad4b044b6772e576f48d27e3) build(deps): bump go.uber.org/fx from 1.18.2 to 1.19.1 (#238)

## [v1.81.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.81.2) - 2023-01-11

- [`70e63e1`](https://github.com/alexfalkowski/go-service/commit/70e63e1ddfdf768974da08fb73a4286a250ddfd9) build(deps): bump google.golang.org/grpc from 1.51.0 to 1.52.0 (#237)

## [v1.81.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.81.1) - 2023-01-10

- [`9dc31f1`](https://github.com/alexfalkowski/go-service/commit/9dc31f131452cea4012700457135a947f9a52fc7) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.45.1 to 1.46.0 (#236)

## [v1.81.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.81.0) - 2023-01-08

- [`1a9484d`](https://github.com/alexfalkowski/go-service/commit/1a9484d0b55ccde55b7c398b1e1bf8a219eea01b) feat(cmd): allow config to be resused (#235)

## [v1.80.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.80.0) - 2023-01-07

- [`dc29ef2`](https://github.com/alexfalkowski/go-service/commit/dc29ef2ef792e3a96741fc10689928894b2c78e6) feat(config): add toml support (#234)

## [v1.79.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.79.0) - 2023-01-06

- [`97f6284`](https://github.com/alexfalkowski/go-service/commit/97f628418f2c7e0b3c10bec396d134d6b23ae1ea) feat(config): add kind to marshaller (#233)

## [v1.78.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.78.0) - 2023-01-06

- [`5fbe33f`](https://github.com/alexfalkowski/go-service/commit/5fbe33ff06edab996ddb510672dd7fdfcfa93263) feat(config): allow to read from multiple sources (#232)

## [v1.77.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.8) - 2023-01-05

- [`8894ed4`](https://github.com/alexfalkowski/go-service/commit/8894ed41b8db72b52d7814ba48b5d8b37c854181) fix: downgrade go.uber.org/fx to v1.18.2 (#230)

## [v1.77.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.7) - 2023-01-05

- [`f8a9578`](https://github.com/alexfalkowski/go-service/commit/f8a95781752f4218ce93d5304760cd37a0604533) fix(security): verifier must return registered claims (#229)

## [v1.77.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.6) - 2023-01-04

- [`495c94c`](https://github.com/alexfalkowski/go-service/commit/495c94cf171094ccf944f48ec863e21109a10484) refactor(auth0): rename to key (#228)

## [v1.77.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.5) - 2023-01-04

- [`76dd50a`](https://github.com/alexfalkowski/go-service/commit/76dd50a8785b474ffa3f69032d8f88df6141bbae) build(deps): bump go.uber.org/fx from 1.18.2 to 1.19.0 (#226)

## [v1.77.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.4) - 2023-01-04

- [`6fde783`](https://github.com/alexfalkowski/go-service/commit/6fde783027e35d79346568349bb0b9e57fb99b93) build(deps): bump golang.org/x/net from 0.4.0 to 0.5.0 (#225)

## [v1.77.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.3) - 2023-01-04

- [`b439447`](https://github.com/alexfalkowski/go-service/commit/b4394473b08469551a10844471f7617c38d9b407) build(deps): bump github.com/klauspost/compress from 1.15.13 to 1.15.14 (#224)

## [v1.77.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.2) - 2023-01-04

- [`d82b804`](https://github.com/alexfalkowski/go-service/commit/d82b80432838b8be07b7ee5154e797a0ac2f5317) build(deps): bump github.com/hashicorp/go-retryablehttp (#223)

## [v1.77.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.1) - 2022-12-28

- [`474fc31`](https://github.com/alexfalkowski/go-service/commit/474fc31ca95ebf851433d8552e61f161cce48a9b) build(deps): bump github.com/rs/cors from 1.8.2 to 1.8.3 (#222)

## [v1.77.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.77.0) - 2022-12-22

- [`5a0fea9`](https://github.com/alexfalkowski/go-service/commit/5a0fea9e08592c0a67461d43e46c4c64b3480263) feat(security): parse authorization header (#221)

## [v1.76.27](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.27) - 2022-12-22

- [`ad41511`](https://github.com/alexfalkowski/go-service/commit/ad41511ac19013de24adcdf8b1056044277a3a51) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#220)

## [v1.76.26](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.26) - 2022-12-20

- [`a8e5631`](https://github.com/alexfalkowski/go-service/commit/a8e56310cbe04da8c9faca716b75809a3fc75a88) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.44.1 to 1.45.1 (#219)

## [v1.76.25](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.25) - 2022-12-13

- [`241a0be`](https://github.com/alexfalkowski/go-service/commit/241a0be1d37f155225233686464d0a66795969f7) build(deps): bump go.uber.org/multierr from 1.8.0 to 1.9.0 (#218)

## [v1.76.24](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.24) - 2022-12-12

- [`9abe6f5`](https://github.com/alexfalkowski/go-service/commit/9abe6f5d04a4a256670de473b6f0ebd0fefc84f6) build(deps): bump github.com/klauspost/compress from 1.15.12 to 1.15.13 (#217)

## [v1.76.23](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.23) - 2022-12-09

- [`b791f08`](https://github.com/alexfalkowski/go-service/commit/b791f08bc3406b7da5c6989194a38f1811e06396) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.44.0 to 1.44.1 (#216)

## [v1.76.22](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.22) - 2022-12-08

- [`9c44257`](https://github.com/alexfalkowski/go-service/commit/9c44257084aaef6504ab11dc77ab909f869d37c8) build(deps): bump golang.org/x/net from 0.3.0 to 0.4.0 (#215)

## [v1.76.21](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.21) - 2022-12-07

- [`98f750e`](https://github.com/alexfalkowski/go-service/commit/98f750e9fd69f484089ff32c161272a82edcb538) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.43.1 to 1.44.0 (#213)
- [`ec7c409`](https://github.com/alexfalkowski/go-service/commit/ec7c4097662c90dd883e1d63476e87151701b6fc) build(deps): bump golang.org/x/net from 0.2.0 to 0.3.0 (#214)

## [v1.76.20](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.20) - 2022-11-30

- [`f636231`](https://github.com/alexfalkowski/go-service/commit/f636231ab4683e4722d6feb0bb14839f76228b03) build(deps): bump go.uber.org/zap from 1.23.0 to 1.24.0 (#212)

## [v1.76.19](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.19) - 2022-11-30

- [`1f137b8`](https://github.com/alexfalkowski/go-service/commit/1f137b8dcb0ab976c503954c23d69a0b786c2bdc) build(deps): bump github.com/golang-jwt/jwt/v4 from 4.4.2 to 4.4.3 (#211)

## [v1.76.18](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.18) - 2022-11-28

- [`e0975f1`](https://github.com/alexfalkowski/go-service/commit/e0975f14a6f81526073f954424b108b7e404c4f5) ci: use latest (#210)

## [v1.76.17](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.17) - 2022-11-24

- [`13f6199`](https://github.com/alexfalkowski/go-service/commit/13f619925f1215aa864135d1d78006bb00fa96e9) build(deps): bump github.com/go-redis/cache/v8 from 8.4.3 to 8.4.4 (#209)

## [v1.76.16](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.16) - 2022-11-20

- [`7bfbea8`](https://github.com/alexfalkowski/go-service/commit/7bfbea8b9a5b54270dd57a96b8f9f859ebdb21f0) build(deps): bump google.golang.org/grpc from 1.50.1 to 1.51.0 (#208)

## [v1.76.15](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.15) - 2022-11-16

- [`276e929`](https://github.com/alexfalkowski/go-service/commit/276e9298b56dba0d4556d15a1222d19d04428592) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#207)

## [v1.76.14](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.14) - 2022-11-13

- [`5ccd874`](https://github.com/alexfalkowski/go-service/commit/5ccd8745225460ebd9c9ceded348718a690752ff) builld(deps): update (#206)

## [v1.76.13](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.13) - 2022-11-09

- [`ecec53b`](https://github.com/alexfalkowski/go-service/commit/ecec53b2b7f6e88c619b3fa49198fd9590dfa470) build(deps): bump golang.org/x/net from 0.1.0 to 0.2.0 (#205)

## [v1.76.12](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.12) - 2022-11-09

- [`d5e3b8b`](https://github.com/alexfalkowski/go-service/commit/d5e3b8b9917d834fd44fca5865795ad28cde67bc) build(deps): bump github.com/prometheus/client_golang (#204)

## [v1.76.11](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.11) - 2022-11-02

- [`dae8efb`](https://github.com/alexfalkowski/go-service/commit/dae8efb27ca7adaa93c007f6ef61d22fcfc0e865) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#202)

## [v1.76.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.10) - 2022-10-28

- [`85f1589`](https://github.com/alexfalkowski/go-service/commit/85f158935dbaf3edde2b4ae2767b36751c585e65) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.43.0 to 1.43.1 (#201)

## [v1.76.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.9) - 2022-10-27

- [`14342ee`](https://github.com/alexfalkowski/go-service/commit/14342eeaff43bc8ebc58ab938aba74d4c8e1f12f) build(deps): bump github.com/klauspost/compress from 1.15.11 to 1.15.12 (#200)

## [v1.76.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.8) - 2022-10-26

- [`c7616a2`](https://github.com/alexfalkowski/go-service/commit/c7616a276989a73b2d5e0441c435a9eeaaec8fae) build(deps): bump github.com/golang-jwt/jwt/v4 from 4.1.0 to 4.4.2 (#199)

## [v1.76.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.7) - 2022-10-25

- [`1fa562a`](https://github.com/alexfalkowski/go-service/commit/1fa562a1b76ebb2ff86ebdcfa56153f156dfb61c) build(deps): bump github.com/spf13/cobra from 1.6.0 to 1.6.1 (#197)

## [v1.76.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.6) - 2022-10-25

- [`1020640`](https://github.com/alexfalkowski/go-service/commit/1020640d3960eeba233bda22f413d26fe546589a) build(deps): update use github.com/golang-jwt/jwt/v4 (#198)

## [v1.76.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.5) - 2022-10-18

- [`d2b832e`](https://github.com/alexfalkowski/go-service/commit/d2b832eba0d7323bacd93ecfb5fe712151bdc11a) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#196)

## [v1.76.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.4) - 2022-10-17

- [`5b99642`](https://github.com/alexfalkowski/go-service/commit/5b9964272e18e33bc2ce8aaabdd6911af021240f) build(deps): bump google.golang.org/grpc from 1.50.0 to 1.50.1 (#195)

## [v1.76.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.3) - 2022-10-14

- [`bb0a74f`](https://github.com/alexfalkowski/go-service/commit/bb0a74f851fe14145cf66539224d1840dd6f6df7) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.42.1 to 1.43.0 (#194)

## [v1.76.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.2) - 2022-10-12

- [`c3049bb`](https://github.com/alexfalkowski/go-service/commit/c3049bb38d74019710f0c0ed6aaa299623431d0b) build(deps): bump github.com/spf13/cobra from 1.5.0 to 1.6.0 (#193)

## [v1.76.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.1) - 2022-10-12

- [`eec594f`](https://github.com/alexfalkowski/go-service/commit/eec594f89e996c7ddb0d16271351c28c25ed31ff) build(deps): bump github.com/dgraph-io/ristretto from 0.1.0 to 0.1.1 (#192)

## [v1.76.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.76.0) - 2022-10-09

- [`094ee4d`](https://github.com/alexfalkowski/go-service/commit/094ee4d200f935f13af36281d17b93be18d665ee) feat(cache): add addresses to config (#191)

## [v1.75.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.75.0) - 2022-10-08

- [`22b2aa7`](https://github.com/alexfalkowski/go-service/commit/22b2aa7a045267a283b32a57336a7da64eb8026e) feat(cache): add redis incr (#190)

## [v1.74.20](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.20) - 2022-10-07

- [`06b7ca0`](https://github.com/alexfalkowski/go-service/commit/06b7ca08e47f9a59807ef87dc57a8e534d0f1eab) build(deps): bump google.golang.org/grpc from 1.49.0 to 1.50.0 (#189)

## [v1.74.19](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.19) - 2022-09-29

- [`08fa87f`](https://github.com/alexfalkowski/go-service/commit/08fa87fb153f8e73d0112d1098514b8864addf28) build(deps): bump go.uber.org/fx from 1.18.1 to 1.18.2 (#188)

## [v1.74.18](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.18) - 2022-09-27

- [`75e9adf`](https://github.com/alexfalkowski/go-service/commit/75e9adf6edc97eca8fe1fbe12be4cf81bd042d8f) build(deps): bump github.com/klauspost/compress from 1.15.10 to 1.15.11 (#187)

## [v1.74.17](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.17) - 2022-09-27

- [`41f083a`](https://github.com/alexfalkowski/go-service/commit/41f083a33dedfed6c372c17154d8a495922d20e2) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.42.0 to 1.42.1 (#186)

## [v1.74.16](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.16) - 2022-09-19

- [`6368e44`](https://github.com/alexfalkowski/go-service/commit/6368e44f3088661f6a9d70918f5980954bd6a7a4) build(deps): bump github.com/klauspost/compress from 1.15.9 to 1.15.10 (#185)

## [v1.74.15](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.15) - 2022-09-16

- [`6f05667`](https://github.com/alexfalkowski/go-service/commit/6f056676062fc7175e4829e666f8892e7936ce1f) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.41.1 to 1.42.0 (#183)

## [v1.74.14](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.14) - 2022-09-16

- [`7db71dd`](https://github.com/alexfalkowski/go-service/commit/7db71dd88337be8e641b74c8a8c694e7626fd6a6) style: remove contextcheck (#184)

## [v1.74.13](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.13) - 2022-09-13

- [`c5508e3`](https://github.com/alexfalkowski/go-service/commit/c5508e37172a0abe890d345f2ab87f96409fe864) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.41.0 to 1.41.1 (#182)

## [v1.74.12](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.12) - 2022-09-05

- [`61e9b15`](https://github.com/alexfalkowski/go-service/commit/61e9b153e4893991b52858b16525e7b0d196dfe5) build(deps): bump github.com/jackc/pgx/v4 from 4.17.1 to 4.17.2 (#181)

## [v1.74.11](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.11) - 2022-08-30

- [`49a06b4`](https://github.com/alexfalkowski/go-service/commit/49a06b47b4c63741f086e033e230ae199dad4bec) build(deps): update all (#180)

## [v1.74.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.10) - 2022-08-29

- [`ce736fc`](https://github.com/alexfalkowski/go-service/commit/ce736fc1d86b4e13d21b3a23fb9e16705c1bcdeb) build(deps): bump github.com/jackc/pgx/v4 from 4.17.0 to 4.17.1 (#179)

## [v1.74.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.9) - 2022-08-25

- [`7a1d667`](https://github.com/alexfalkowski/go-service/commit/7a1d667669e75d28edc6e33ac2fba085a844217e) build(deps): bump go.uber.org/zap from 1.22.0 to 1.23.0 (#178)

## [v1.74.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.8) - 2022-08-24

- [`2430864`](https://github.com/alexfalkowski/go-service/commit/2430864a5d8d811e017ce3af2358db9ef20fa289) build(deps): bump google.golang.org/grpc from 1.48.0 to 1.49.0 (#176)

## [v1.74.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.7) - 2022-08-24

- [`b64fd60`](https://github.com/alexfalkowski/go-service/commit/b64fd6031c03cc496468ea487266ca833123bc0b) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#177)

## [v1.74.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.6) - 2022-08-23

- [`bbd262a`](https://github.com/alexfalkowski/go-service/commit/bbd262a6d2696fcc23e88cecd2442328895dc15d) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.40.1 to 1.41.0 (#175)

## [v1.74.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.5) - 2022-08-19

- [`466a843`](https://github.com/alexfalkowski/go-service/commit/466a84357ad43508d340d282bf7eff4ce5bf4eb8) build(deps): update all (#174)

## [v1.74.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.4) - 2022-08-10

- [`881eccb`](https://github.com/alexfalkowski/go-service/commit/881eccb2ed841121c232cd6261b942c292764e0c) build(deps): update all (#173)

## [v1.74.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.3) - 2022-08-09

- [`1a20aaf`](https://github.com/alexfalkowski/go-service/commit/1a20aaf43f672af216feca862212ab125c661777) test: run setup-nsq as part of specs (#172)

## [v1.74.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.2) - 2022-08-09

- [`9cc5ed1`](https://github.com/alexfalkowski/go-service/commit/9cc5ed170330eea9fe215ad78565c76cf985dfff) build(deps): update all (#171)

## [v1.74.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.1) - 2022-08-06

- [`14ca14e`](https://github.com/alexfalkowski/go-service/commit/14ca14e672e92c6510a7800354694cb4b013b8e0) docs: use github.com/radovskyb/watcher (#168)

## [v1.74.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.74.0) - 2022-08-06

- [`10f8e61`](https://github.com/alexfalkowski/go-service/commit/10f8e6103af51b64a8ca0e003b8b9227af588f67) feat(go): update go 1.19 (#167)

## [v1.73.13](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.13) - 2022-08-01

- [`052d0ef`](https://github.com/alexfalkowski/go-service/commit/052d0efb9b8f8e0c2e533d0e9fd2541e48f42559) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#166)

## [v1.73.12](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.12) - 2022-07-29

- [`36164ce`](https://github.com/alexfalkowski/go-service/commit/36164cee502390ce0298f344ee85dc3dc31d1e12) build(deps): bump google.golang.org/protobuf from 1.28.0 to 1.28.1 (#165)

## [v1.73.11](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.11) - 2022-07-22

- [`40978bf`](https://github.com/alexfalkowski/go-service/commit/40978bf00b946323d1ec5a2ed895fed1149750e8) build(deps): bump github.com/klauspost/compress from 1.15.8 to 1.15.9 (#164)

## [v1.73.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.10) - 2022-07-20

- [`a886036`](https://github.com/alexfalkowski/go-service/commit/a886036a5e4620dc78368d8254130edbb76a9255) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#163)

## [v1.73.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.9) - 2022-07-20

- [`2b26221`](https://github.com/alexfalkowski/go-service/commit/2b2622154df9d1e7f748eb40580b7517f9d3db5b) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.40.0 to 1.40.1 (#162)

## [v1.73.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.8) - 2022-07-17

- [`5f9bb62`](https://github.com/alexfalkowski/go-service/commit/5f9bb62cdd886874bd5d021a0c0716c8dbec3486) build(deps): update all (#161)

## [v1.73.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.7) - 2022-07-14

- [`b48d68e`](https://github.com/alexfalkowski/go-service/commit/b48d68e62a66203978f198ab3c64f58a76b0c004) build(deps): bump github.com/klauspost/compress from 1.15.7 to 1.15.8 (#160)

## [v1.73.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.6) - 2022-07-13

- [`81cdecb`](https://github.com/alexfalkowski/go-service/commit/81cdecb0010c6549f45e7c1dd05c9e65589a5110) build(deps): bump google.golang.org/grpc from 1.47.0 to 1.48.0 (#159)

## [v1.73.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.5) - 2022-07-12

- [`dc2988b`](https://github.com/alexfalkowski/go-service/commit/dc2988b0382cf1a5f226390d1a9d6379674dbcfd) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.39.0 to 1.39.1 (#158)

## [v1.73.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.4) - 2022-07-05

- [`f032765`](https://github.com/alexfalkowski/go-service/commit/f0327658b16d24994770b6874104410062f2712c) build(deps): bump github.com/linxGnu/mssqlx from 1.1.6 to 1.1.7 (#157)

## [v1.73.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.3) - 2022-07-04

- [`e2bc0d9`](https://github.com/alexfalkowski/go-service/commit/e2bc0d93f60a09cca40425627b086a177221164e) build(deps): bump gopkg.in/DataDog/dd-trace-go.v1 from 1.38.1 to 1.39.0 (#156)

## [v1.73.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.2) - 2022-06-30

- [`36b9bd0`](https://github.com/alexfalkowski/go-service/commit/36b9bd005f8299f6d81e1435141ca13f0c8746b0) build(deps): bump github.com/klauspost/compress from 1.15.6 to 1.15.7 (#155)

## [v1.73.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.1) - 2022-06-22

- [`9fd750f`](https://github.com/alexfalkowski/go-service/commit/9fd750f817030c8c4f3185d951e3cea54b244b8a) build(deps): bump github.com/spf13/cobra from 1.4.0 to 1.5.0 (#154)

## [v1.73.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.73.0) - 2022-06-14

- [`bc79624`](https://github.com/alexfalkowski/go-service/commit/bc796249c5cb4cf11d4ccc8c7d681239b2cdf83f) feat(time): change to 30s (#153)

## [v1.72.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.72.0) - 2022-06-13

- [`2a1ad50`](https://github.com/alexfalkowski/go-service/commit/2a1ad50f7990afdb4ca14bda2ee345defa3efa90) feat: rename to kind (#152)

## [v1.71.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.10) - 2022-06-10

- [`0a0e4d3`](https://github.com/alexfalkowski/go-service/commit/0a0e4d3d6e3b0a56de61470dd3ce2d5e6541a103) build(deps): bump github.com/alexfalkowski/go-health (#151)

## [v1.71.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.9) - 2022-06-06

- [`32e31ce`](https://github.com/alexfalkowski/go-service/commit/32e31ce04ba034eb33fffae05f6ec90789e70e9d) build(deps): bump github.com/klauspost/compress from 1.15.5 to 1.15.6 (#150)

## [v1.71.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.8) - 2022-06-01

- [`35176f6`](https://github.com/alexfalkowski/go-service/commit/35176f619ac477612c2ff8af36bb54a40e0c3186) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#149)

## [v1.71.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.7) - 2022-06-01

- [`658ce7a`](https://github.com/alexfalkowski/go-service/commit/658ce7a8e95e07806d6828b0225db50a46ec53b6) build(deps): bump google.golang.org/grpc from 1.46.2 to 1.47.0 (#148)

## [v1.71.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.6) - 2022-05-30

- [`92dbd58`](https://github.com/alexfalkowski/go-service/commit/92dbd5883da5eb5d5a2ecd59fbd66bab7f586b18) build(deps): bump github.com/klauspost/compress from 1.15.4 to 1.15.5 (#146)

## [v1.71.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.5) - 2022-05-30

- [`6479a48`](https://github.com/alexfalkowski/go-service/commit/6479a48b01b2d644443ac929fe6bd0cf93e8726a) build(deps): bump gopkg.in/yaml.v3 from 3.0.0 to 3.0.1 (#147)

## [v1.71.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.4) - 2022-05-24

- [`025cb16`](https://github.com/alexfalkowski/go-service/commit/025cb16f2d1bfafe455d66f6ff5435d726979809) build(deps): udpate (#145)

## [v1.71.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.3) - 2022-05-24

- [`c1d40b4`](https://github.com/alexfalkowski/go-service/commit/c1d40b4bc37fcdfe69b2e7c029b54c0376e4f0da) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#144)

## [v1.71.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.2) - 2022-05-23

- [`b960038`](https://github.com/alexfalkowski/go-service/commit/b96003828e768d64b65b5bcdbb5a21888b22c17a) build(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 (#143)

## [v1.71.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.1) - 2022-05-23

- [`cf7ec2f`](https://github.com/alexfalkowski/go-service/commit/cf7ec2fc2283884007b787cce0bc1251d75458d7) ci: add dependabot file (#142)

## [v1.71.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.71.0) - 2022-05-19

- [`4140491`](https://github.com/alexfalkowski/go-service/commit/41404914cf965f4e0de6e18a22807f7547a9c18e) feat(config): watch config changes with polling (#141)
- [`817cfc4`](https://github.com/alexfalkowski/go-service/commit/817cfc468abcbe8295be9f520984934ac6ebce5c) docs: add config (#140)

## [v1.70.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.70.1) - 2022-05-18

- [`ecc2f39`](https://github.com/alexfalkowski/go-service/commit/ecc2f39d43561a299e089454463eeb957ecfc954) refactor(transport): use server (#139)

## [v1.70.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.70.0) - 2022-05-18

- [`d0b2478`](https://github.com/alexfalkowski/go-service/commit/d0b2478a916d83f390119c506bc6a110326bd915) feat(transport): add conn mux (#138)

## [v1.69.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.69.0) - 2022-05-17

- [`8e25565`](https://github.com/alexfalkowski/go-service/commit/8e25565e5a9b4dc8c2fb94928027aac5871c6249) feat(grpc): add limiter (#137)

## [v1.68.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.68.1) - 2022-05-17

- [`1b96263`](https://github.com/alexfalkowski/go-service/commit/1b962639b1ea2a57cff0acc49f5433128a68f300) docs(sql): add mssqlx (#136)

## [v1.68.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.68.0) - 2022-05-17

- [`2b3272b`](https://github.com/alexfalkowski/go-service/commit/2b3272b7fbd287c3e96dbd0ee993a9655a748a5a) feat(sql): add master and slaves (#135)

## [v1.67.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.67.1) - 2022-05-16

- [`74f7f67`](https://github.com/alexfalkowski/go-service/commit/74f7f670303b83fa031f49d3ad672602d2fcc3c9) fix(sql): make sure metrics are by driver (#134)

## [v1.67.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.67.0) - 2022-05-16

- [`31e8d42`](https://github.com/alexfalkowski/go-service/commit/31e8d42aab1ecfa35ccb2bea610cf4970c2fbe1b) feat(sql): add conn pool settings (#133)

## [v1.66.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.8) - 2022-05-16

- [`9988d8c`](https://github.com/alexfalkowski/go-service/commit/9988d8c2ae48c78c53ccc2277fcd29fbd4985b11) refactor(sql): support other drivers (#132)

## [v1.66.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.7) - 2022-05-16

- [`af83963`](https://github.com/alexfalkowski/go-service/commit/af839630ddc3b6fca3a4a11b1a799715a0c047e1) ci: update build to large resource class in config.yml (#131)

## [v1.66.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.6) - 2022-05-16

- [`0321253`](https://github.com/alexfalkowski/go-service/commit/0321253f11f4d282a4cec205e225dc8aaac2246a) fix(trace): make sure we add meta after the call (#130)

## [v1.66.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.5) - 2022-05-15

- [`32251f9`](https://github.com/alexfalkowski/go-service/commit/32251f92bf747f7751353239eccc061638d90752) build: turn on race checking (#129)

## [v1.66.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.4) - 2022-05-15

- [`680f4c7`](https://github.com/alexfalkowski/go-service/commit/680f4c7cd84aa1c64ead70b4f302e849ed0f428a) test: add missing (#127)

## [v1.66.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.3) - 2022-05-14

- [`a570f55`](https://github.com/alexfalkowski/go-service/commit/a570f552b3915e18189ef9b154120b95f5096c6e) refactor(http): use funcs (#126)

## [v1.66.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.2) - 2022-05-14

- [`7718f11`](https://github.com/alexfalkowski/go-service/commit/7718f116bb7f69e57df9e9d4fff8348056965df6) build: use race detection (#125)

## [v1.66.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.1) - 2022-05-14

- [`471d45e`](https://github.com/alexfalkowski/go-service/commit/471d45e9f2c278f68752351f4c3261022cf797af) refactor(pg): separate open and registration (#124)

## [v1.66.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.66.0) - 2022-05-14

- [`1a381f3`](https://github.com/alexfalkowski/go-service/commit/1a381f302d9b8c775923ba9dfaaa8e4055e576de) feat(trace): add database and cache (#123)

## [v1.65.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.65.6) - 2022-05-13

- [`884ec4b`](https://github.com/alexfalkowski/go-service/commit/884ec4bbbef483ba518b49c733ca3795b382aa7a) build(deps): update all (#122)

## [v1.65.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.65.5) - 2022-05-12

- [`946be46`](https://github.com/alexfalkowski/go-service/commit/946be46a046c4cc3ecd0bb8bb87f7ec1bcd33282) style: clean up (#121)

## [v1.65.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.65.4) - 2022-05-11

- [`7a53b5c`](https://github.com/alexfalkowski/go-service/commit/7a53b5cd7d2445ef7bb9183bd6c44a71b8ba5964) test: cleanup (#120)
- [`d846174`](https://github.com/alexfalkowski/go-service/commit/d846174fc43ea7e59e788fd2bc755fbec5c10b5b) style: lint (#119)

## [v1.65.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.65.3) - 2022-05-10

- [`157183c`](https://github.com/alexfalkowski/go-service/commit/157183ce96950e33ed15f5440b274ea15c978685) fix(health): readiness observer (#118)

## [v1.65.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.65.2) - 2022-05-10

- [`c64df3b`](https://github.com/alexfalkowski/go-service/commit/c64df3b39d522e4c1a6fc06f8004971887e9e9c1) build(deps): update github.com/alexfalkowski/go-health v1.10.3 (#117)

## [v1.65.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.65.1) - 2022-05-09

- [`94f235f`](https://github.com/alexfalkowski/go-service/commit/94f235f77a10a35301b29ecef63988b818055909) build(deps): update all (#116)

## [v1.65.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.65.0) - 2022-05-09

- [`2aa04fc`](https://github.com/alexfalkowski/go-service/commit/2aa04fce526d70505486f2413c07750d6fa130a0) feat(prometheus): complete metrics (#115)

## [v1.64.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.64.3) - 2022-05-07

- [`26b30b9`](https://github.com/alexfalkowski/go-service/commit/26b30b9871e4bb6cf05f15e3e40a2516e920613c) fix(config): add random wait time to watcher (#114)

## [v1.64.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.64.2) - 2022-05-06

- [`74fdfe5`](https://github.com/alexfalkowski/go-service/commit/74fdfe5eddf6a07f6721119b2b9b84aa72e457c5) fix(trace): use noop opentracing tracer (#113)

## [v1.64.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.64.1) - 2022-05-06

- [`a15fd61`](https://github.com/alexfalkowski/go-service/commit/a15fd61164647526859edb3d9e7f19f61382b3b7) fix(http): ignore trace if not set (#112)

## [v1.64.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.64.0) - 2022-05-06

- [`355a5d2`](https://github.com/alexfalkowski/go-service/commit/355a5d26b4175f8ca515eb805092b5233ff6abad) feat(trace): have one config for opentracing (#111)

## [v1.63.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.63.0) - 2022-04-28

- [`6058276`](https://github.com/alexfalkowski/go-service/commit/6058276d50ed6ce27a8d511c232dfb5b62bd9c3d) feat(opentracing): append service to name (#110)

## [v1.62.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.62.3) - 2022-04-28

- [`a659377`](https://github.com/alexfalkowski/go-service/commit/a65937777ca785fdf045706f4b835319453552ad) fix: params are not options (#109)

## [v1.62.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.62.2) - 2022-04-28

- [`61ab8a7`](https://github.com/alexfalkowski/go-service/commit/61ab8a7a231b6da5359826cbe692ca9c69f9f690) fix: remove getting ip address (#108)

## [v1.62.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.62.1) - 2022-04-28

- [`3ae3df8`](https://github.com/alexfalkowski/go-service/commit/3ae3df88372b6db629ee966cf8993b8f94a2059a) fix(http): remove version header as we have it in grpc gateway (#107)

## [v1.62.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.62.0) - 2022-04-28

- [`0fb1504`](https://github.com/alexfalkowski/go-service/commit/0fb15048ecba77ffcc27ac55b9f89ece4a8c4c6d) feat(config): read file from env (#106)

## [v1.61.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.61.0) - 2022-04-27

- [`773e1c5`](https://github.com/alexfalkowski/go-service/commit/773e1c5edb35beffe01dfc09db32f3cd5ab891b6) feat: add version (#105)

## [v1.60.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.60.1) - 2022-04-27

- [`99ab73c`](https://github.com/alexfalkowski/go-service/commit/99ab73ca500afa2a86a7ccebbe317c8bf3c86cb6) docs: add config watch (#104)

## [v1.60.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.60.0) - 2022-04-27

- [`ce06ac1`](https://github.com/alexfalkowski/go-service/commit/ce06ac1f76e470bbf3dde42163035abe68b1235b) feat(config): watch for file changes (#103)

## [v1.59.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.59.0) - 2022-04-27

- [`05a0411`](https://github.com/alexfalkowski/go-service/commit/05a0411f41b862ae8622b4a2e9c4f3d2841a2130) feat(nsq): add resilience (#102)

## [v1.58.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.58.0) - 2022-04-27

- [`3224fd2`](https://github.com/alexfalkowski/go-service/commit/3224fd28a499a2c7f5be8609cd13e8675fd34680) feat(opentracing): each package has its own tracer (#101)

## [v1.57.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.57.3) - 2022-04-27

- [`86ebab7`](https://github.com/alexfalkowski/go-service/commit/86ebab7d01cd5463c144408cac2fa3faa5ed3048) test: handle edge cases (#100)

## [v1.57.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.57.2) - 2022-04-26

- [`d94cf1b`](https://github.com/alexfalkowski/go-service/commit/d94cf1b05d76c475fee6c9aa02bc77c164ea17a7) refactor: move marshaller and compressor for each package (#99)

## [v1.57.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.57.1) - 2022-04-26

- [`82d3e7d`](https://github.com/alexfalkowski/go-service/commit/82d3e7deafc1068883a18c82362a3d8c1eee2dd1) fix(config): clean the file path (#98)

## [v1.57.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.57.0) - 2022-04-26

- [`d4a766b`](https://github.com/alexfalkowski/go-service/commit/d4a766b3c855a58f4cbdadd4ac64aa8ab08024a2) feat(http): add cors (#97)

## [v1.56.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.56.2) - 2022-04-26

- [`ce93eae`](https://github.com/alexfalkowski/go-service/commit/ce93eaea84ac751b23e540aa3338addee085c0a7) refactor(cmd): return command so it can be customised (#96)

## [v1.56.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.56.1) - 2022-04-26

- [`327e4d0`](https://github.com/alexfalkowski/go-service/commit/327e4d093fdf72d911c2175cfb3b88a0b54a4a3c) refactor: remove redundant funcs (#95)

## [v1.56.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.56.0) - 2022-04-25

- [`ddf1ad0`](https://github.com/alexfalkowski/go-service/commit/ddf1ad0d4fdb7ac42b62df4d7316eb485851757c) feat(transport): add remote address (#94)

## [v1.55.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.55.0) - 2022-04-24

- [`c70f44b`](https://github.com/alexfalkowski/go-service/commit/c70f44b970b7ac4dad2fd8a3fec16d7f0da05805) feat(http): add headers (#93)

## [v1.54.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.54.0) - 2022-04-23

- [`12fb3df`](https://github.com/alexfalkowski/go-service/commit/12fb3df192588d425b192d240ba834a9d32405f9) feat: add marshaller and compressor (#92)

## [v1.53.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.53.0) - 2022-04-19

- [`1c15ba8`](https://github.com/alexfalkowski/go-service/commit/1c15ba8986c4f06a9f3b79834f2d1335143adc89) feat: use one tracer (#91)

## [v1.52.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.52.1) - 2022-04-18

- [`c74d526`](https://github.com/alexfalkowski/go-service/commit/c74d526d5195a548d02fe1e1a8810c1555e0e24f) refactor(opentracing): use same name (#90)

## [v1.52.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.52.0) - 2022-04-18

- [`b4bc23a`](https://github.com/alexfalkowski/go-service/commit/b4bc23a0377874f782db83e9430eb75a8b021c21) feat(trace): create multiple for different traces (#89)

## [v1.51.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.51.3) - 2022-04-14

- [`4ecabcc`](https://github.com/alexfalkowski/go-service/commit/4ecabccc36bb8f74227167e8c692f2b8fd038110) fix(grpc): remove rate limit as it will be handled elsewhere (#88)

## [v1.51.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.51.2) - 2022-04-13

- [`fc59d53`](https://github.com/alexfalkowski/go-service/commit/fc59d5354fd97197494347834afac77d5fd05c68) docs: update readme to communicate how to use the service (#87)

## [v1.51.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.51.1) - 2022-04-08

- [`16e5a4a`](https://github.com/alexfalkowski/go-service/commit/16e5a4ad1575ea74f416f4ca839be0ca727980e3) fix(jaeger): do not log spans (#86)

## [v1.51.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.51.0) - 2022-04-08

- [`a3493f5`](https://github.com/alexfalkowski/go-service/commit/a3493f515d3d646087f91438c8d2e6b66e509fa6) feat(redis): use zap logger (#85)

## [v1.50.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.50.1) - 2022-04-07

- [`0f5bd60`](https://github.com/alexfalkowski/go-service/commit/0f5bd6066ea1f68bacc93e09edc2243b239a64f5) refactor: remove unused code (#84)

## [v1.50.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.50.0) - 2022-04-07

- [`cefc29d`](https://github.com/alexfalkowski/go-service/commit/cefc29d2045038adc5b2bbc0ae8b8c10c1ec5da6) feat: add client options (#83)

## [v1.49.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.49.1) - 2022-04-07

- [`7fac713`](https://github.com/alexfalkowski/go-service/commit/7fac713bf9ebf06506d78258e767a56c6d684d49) fix(opentracing): fix full method (#82)

## [v1.49.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.49.0) - 2022-04-07

- [`8d03d3c`](https://github.com/alexfalkowski/go-service/commit/8d03d3c3911137f7793da2954014c52a4d33b9b7) feat(health): add redis checker (#81)

## [v1.48.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.48.0) - 2022-04-07

- [`32bc3e4`](https://github.com/alexfalkowski/go-service/commit/32bc3e4f358416351f63b55f61588d22813fca07) feat(health): add all the errors to the response (#80)

## [v1.47.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.47.0) - 2022-04-02

- [`eb102bf`](https://github.com/alexfalkowski/go-service/commit/eb102bfe483ed73d620c2c453cc68513e2c69b57) feat(http): add server handlers (#79)

## [v1.46.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.46.0) - 2022-04-02

- [`e4b690d`](https://github.com/alexfalkowski/go-service/commit/e4b690dd1936c0fe8fd68b1e329587935a543ff5) feat(grpc): add local client (#78)

## [v1.45.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.45.1) - 2022-04-01

- [`a57e1bc`](https://github.com/alexfalkowski/go-service/commit/a57e1bc098e22fe05136907b766278ea55ea008b) fix(grpc): remove recovery (#77)

## [v1.45.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.45.0) - 2022-04-01

- [`5887aa4`](https://github.com/alexfalkowski/go-service/commit/5887aa4899252734d055ba2a46debf1e1383c5cb) feat(grpc): add recovery handler to server (#76)

## [v1.44.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.44.0) - 2022-03-31

- [`ab84c75`](https://github.com/alexfalkowski/go-service/commit/ab84c75db81b5ddbdf6fbf5442635d28c3d1f775) feat(config): add ability to marshal to map (#75)

## [v1.43.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.43.0) - 2022-03-30

- [`a1f93ba`](https://github.com/alexfalkowski/go-service/commit/a1f93ba0b0ca17d950d225aca09ddad1d71ce881) feat(config): write to a different location (#74)

## [v1.42.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.42.0) - 2022-03-30

- [`c18d981`](https://github.com/alexfalkowski/go-service/commit/c18d9811928137188491cda309a43bfcf32dee66) feat(trace): add http client for datadog (#73)

## [v1.41.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.41.0) - 2022-03-30

- [`25c0520`](https://github.com/alexfalkowski/go-service/commit/25c052078e29a0dc7edce6a2e5146d3af959d867) feat(config): add write file (#72)

## [v1.40.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.40.0) - 2022-03-29

- [`21c89a6`](https://github.com/alexfalkowski/go-service/commit/21c89a68e562f67dcb1f8fc83976da3ca0b428c2) feat(logger): add nsq and opentracing (#71)

## [v1.39.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.39.2) - 2022-03-29

- [`e7cd1b8`](https://github.com/alexfalkowski/go-service/commit/e7cd1b8280ec27e74a0c61b9245930ad1ca22f9e) refactor(nsq): remove returning ctx (#70)

## [v1.39.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.39.1) - 2022-03-29

- [`14e1654`](https://github.com/alexfalkowski/go-service/commit/14e1654f2109469fe2add553e2fd3e27aa448e47) fix(cmd): add ability to run command (#69)

## [v1.39.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.39.0) - 2022-03-29

- [`59f1297`](https://github.com/alexfalkowski/go-service/commit/59f129702b5606fbc4480df0b8c5605d874071dd) feat(cmd): add run client (#68)

## [v1.38.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.38.0) - 2022-03-27

- [`83b3a10`](https://github.com/alexfalkowski/go-service/commit/83b3a10d2ae1b0f57cfedf4169a14abcb3078f56) feat(ratelimit): add ttl (#67)

## [v1.37.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.37.0) - 2022-03-26

- [`ce18d71`](https://github.com/alexfalkowski/go-service/commit/ce18d71926d3d4b43171d8fa99f5c445c9c127c5) feat(rate): use cache to expire limiters (#66)

## [v1.36.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.36.0) - 2022-03-25

- [`27cfe83`](https://github.com/alexfalkowski/go-service/commit/27cfe8335a03f95fa49a345dbde183d251466458) feat(cmd): change the command to server from serve (#65)

## [v1.35.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.35.1) - 2022-03-20

- [`0cca44b`](https://github.com/alexfalkowski/go-service/commit/0cca44ba3512822d5054732542dd9416ead7b0e4) build(deps): update github.com/alexfalkowski/go-health v1.9.0 (#64)

## [v1.35.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.35.0) - 2022-03-20

- [`2465076`](https://github.com/alexfalkowski/go-service/commit/2465076b738ac919998121c598a595020dfd0a42) feat: update to go 1.18 (#63)

## [v1.34.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.34.0) - 2022-03-19

- [`7e18ab8`](https://github.com/alexfalkowski/go-service/commit/7e18ab8742313a534247dd739fa76bd842829be8) feat: remove pkg folder (#62)

## [v1.33.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.33.0) - 2022-03-19

- [`c78f277`](https://github.com/alexfalkowski/go-service/commit/c78f2771b5f35cb815524d10e8318179825d6983) feat(trace): add tracing for sql and cache (#61)

## [v1.32.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.32.3) - 2022-03-18

- [`76c9d04`](https://github.com/alexfalkowski/go-service/commit/76c9d04bdfaa4b594e92aaa49cfc718b1f229d34) style(lint): add new (#60)

## [v1.32.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.32.2) - 2022-03-04

- [`b4a97a2`](https://github.com/alexfalkowski/go-service/commit/b4a97a2c6c4c6145b4845357cf61f1c506f26a62) build(deps): update all (#59)

## [v1.32.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.32.1) - 2021-10-29

- [`21c1d68`](https://github.com/alexfalkowski/go-service/commit/21c1d68042016d27a6d251963f654b8ef3ad8e32) refactor(grpc): remove health duplication (#58)

## [v1.32.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.32.0) - 2021-10-29

- [`7559397`](https://github.com/alexfalkowski/go-service/commit/7559397e0d7023bc9429b1943a1c552153c1b9a6) feat(grpc): add rate limiting (#57)

## [v1.31.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.31.0) - 2021-10-28

- [`db837a1`](https://github.com/alexfalkowski/go-service/commit/db837a1bccd1ee69929086cb2c4c78eb94f5e7f5) feat(transport): add user agent (#56)

## [v1.30.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.30.0) - 2021-10-27

- [`940bbe0`](https://github.com/alexfalkowski/go-service/commit/940bbe04199f7c544ff5061d6a0b55d7165c8e10) feat(transport): add retry config for http and grpc (#55)

## [v1.29.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.29.2) - 2021-10-26

- [`21266b2`](https://github.com/alexfalkowski/go-service/commit/21266b200e7d2d6dbb1cf9b9dbb5771be12101db) build(deps): update github.com/alexfalkowski/go-health v1.7.3 (#54)

## [v1.29.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.29.1) - 2021-10-26

- [`87660eb`](https://github.com/alexfalkowski/go-service/commit/87660eb87813117b48c76925d940141833b1d953) fix(http): make sure we create a new client (#53)

## [v1.29.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.29.0) - 2021-10-23

- [`e42e3d6`](https://github.com/alexfalkowski/go-service/commit/e42e3d6ca1c24041afb06ed90b6f0cca5cf5986e) feat: add circuit breaker (#52)

## [v1.28.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.28.0) - 2021-10-23

- [`ac9f783`](https://github.com/alexfalkowski/go-service/commit/ac9f783ce49121362bc9f3255ad8f451eed5f749) feat(http): add client retry (#51)

## [v1.27.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.27.2) - 2021-10-20

- [`c13033d`](https://github.com/alexfalkowski/go-service/commit/c13033d003a82f51fd508b9e92f473b0c1011946) feat(health) add liveness and readiness (#50)

## [v1.27.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.27.1) - 2021-10-11

- [`d0d728c`](https://github.com/alexfalkowski/go-service/commit/d0d728c98073fe6f632cfca16f08f25bccdf9b9b) build(deps): clean up go.mod (#49)

## [v1.27.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.27.0) - 2021-10-10

- [`1a95d71`](https://github.com/alexfalkowski/go-service/commit/1a95d719c05c0dc7b4b2f9f2d39328870fec080d) feat: add ability to read config (#48)

## [v1.26.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.26.0) - 2021-09-21

- [`f48958c`](https://github.com/alexfalkowski/go-service/commit/f48958cda73b86147e7d439ae976b393af97e2a5) feat: update go to v.1.17 (#47)

## [v1.25.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.25.0) - 2021-08-16

- [`0ede766`](https://github.com/alexfalkowski/go-service/commit/0ede7662d3a54e108dd02b123148d46630c7724a) feat(deps): update (#46)

## [v1.24.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.24.1) - 2021-07-19

- [`35f2837`](https://github.com/alexfalkowski/go-service/commit/35f2837dbfe07e1c149e1b969e315cc7b022891e) ci: fix release (#45)

## [v1.24.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.24.0) - 2021-06-26

- [`8e253a9`](https://github.com/alexfalkowski/go-service/commit/8e253a9b81e72acff19027cb08d2c4885a7c3b29) feat: add ability to store jwt token (#44)

## [v1.23.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.23.0) - 2021-06-06

- [`e019e0c`](https://github.com/alexfalkowski/go-service/commit/e019e0c772a8149915760e206556b065df40468d) feat(deps): update (#43)

## [v1.22.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.22.0) - 2021-05-17

- [`a3e5311`](https://github.com/alexfalkowski/go-service/commit/a3e531125ef0078232540643e634d3a31ab0eab5) feat(deps): update (#42)
- [`40f31d7`](https://github.com/alexfalkowski/go-service/commit/40f31d7278b7462f1ec46234d1a9ce3cfb68c9fa) feat: update linter to use revive (#41)

## [v1.21.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.21.2) - 2021-05-14

- [`e252d18`](https://github.com/alexfalkowski/go-service/commit/e252d18b21b957fa4dc14ae3c30ad9fdfbf54d0f) fix(grpc): default port to 9090 (#40)

## [v1.21.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.21.1) - 2021-05-14

- [`30f28e5`](https://github.com/alexfalkowski/go-service/commit/30f28e59d814325633da9fc04cf2a81a3280a5a2) fix(cmd): do not pass context (#39)

## [v1.21.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.21.0) - 2021-05-14

- [`7b0a1a5`](https://github.com/alexfalkowski/go-service/commit/7b0a1a5dd64938b19f1b8c05cc1da08a8317dd0b) feat(cmd): add serve and worker commands (#38)

## [v1.20.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.20.3) - 2021-05-13

- [`2eefd7c`](https://github.com/alexfalkowski/go-service/commit/2eefd7c74fed9e2bd451d95e9eec3d3d49dfb82e) fix(health): follow module pattern (#37)

## [v1.20.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.20.2) - 2021-05-13

- [`1650bef`](https://github.com/alexfalkowski/go-service/commit/1650bef374e9f3ca5dad8db49193ed9ba037bd0c) test(security): verify failing tokens (#36)

## [v1.20.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.20.1) - 2021-05-13

- [`a16df8c`](https://github.com/alexfalkowski/go-service/commit/a16df8c087d3ff262c6cfbfc30b1656eb4732346) test(cmd): verify shutdown after 5 secs (#35)

## [v1.20.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.20.0) - 2021-05-13

- [`e86b49b`](https://github.com/alexfalkowski/go-service/commit/e86b49b4033a8989825e2cc901cdce9a5534bf69) feat(http): create server (#34)

## [v1.19.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.19.1) - 2021-05-12

- [`2f0c58b`](https://github.com/alexfalkowski/go-service/commit/2f0c58b82ec870a3551796f7f48e7e9db6fe0846) fix(trace): conform to naming (#33)

## [v1.19.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.19.0) - 2021-05-12

- [`fbdfb8e`](https://github.com/alexfalkowski/go-service/commit/fbdfb8e7cb43e2be87241125cad523fe125f8a5e) feat: move config to appropriate packages (#32)

## [v1.18.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.18.1) - 2021-05-12

- [`7f6b7a7`](https://github.com/alexfalkowski/go-service/commit/7f6b7a752c3770997bfebd684aaba46478f269aa) fix(http): filter sensitive info (#31)

## [v1.18.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.18.0) - 2021-05-12

- [`0faf9fc`](https://github.com/alexfalkowski/go-service/commit/0faf9fc93dddeaf121c1afeb55f0d0c65e37d7cc) feat(security): add auth0 support (#30)

## [v1.17.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.17.0) - 2021-05-11

- [`d1db573`](https://github.com/alexfalkowski/go-service/commit/d1db57329380fdc44c5e17eb63be6f499e743dfa) feat(grpc): allow params to be passed in (#29)

## [v1.16.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.16.2) - 2021-05-11

- [`5725faf`](https://github.com/alexfalkowski/go-service/commit/5725fafd04463f99b70d29652e4317732c96719f) test(metrics): validate collector (#28)

## [v1.16.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.16.1) - 2021-05-11

- [`c47ca87`](https://github.com/alexfalkowski/go-service/commit/c47ca87dfbad2ff1ff1c3291649d98b89ac72762) docs: add prometheus (#27)

## [v1.16.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.16.0) - 2021-05-11

- [`7ad4bf4`](https://github.com/alexfalkowski/go-service/commit/7ad4bf44614b42caabf187c308dc5c399d96d76a) feat(metrics): add redis and ristretto (#26)

## [v1.15.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.15.0) - 2021-05-10

- [`4af1515`](https://github.com/alexfalkowski/go-service/commit/4af1515e26f1403a81e0a5ebf21e901329cc6d7d) feat(sql): add prometheus metrics (#25)

## [v1.14.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.14.0) - 2021-05-08

- [`aeb2811`](https://github.com/alexfalkowski/go-service/commit/aeb2811ea8b8cc355fc886161020f8c3ebb3a0ff) feat(nsq): add producer and consumer (#23)

## [v1.13.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.13.0) - 2021-05-07

- [`cc8f08e`](https://github.com/alexfalkowski/go-service/commit/cc8f08ee90cb24be55097e94c3e902eaa1465b40) feat: error on client with empty token (#22)

## [v1.12.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.12.0) - 2021-05-07

- [`e5b5751`](https://github.com/alexfalkowski/go-service/commit/e5b5751f651760c57ec99cb044fc878b255f7d63) feat: change logger to pass config (#21)

## [v1.11.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.11.0) - 2021-05-06

- [`263ad70`](https://github.com/alexfalkowski/go-service/commit/263ad705ca85153ab7e7c2a7b4efc05bed6393d1) feat: move to transport (#20)

## [v1.10.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.10.1) - 2021-05-06

- [`0646f1e`](https://github.com/alexfalkowski/go-service/commit/0646f1e6cf579294ea8ce4c4edd1c52b2997c9fc) test: wait for cache to accept (#19)
- [`9d048d4`](https://github.com/alexfalkowski/go-service/commit/9d048d40d3b80345fcb59999b8185ff1d8fae8a7) test: rename (#18)

## [v1.10.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.10.0) - 2021-05-06

- [`ce20f14`](https://github.com/alexfalkowski/go-service/commit/ce20f14ca4ae166c7862d0be7c77fb6e4fe2ca0e) feat(datadog): pass more options (#17)

## [v1.9.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.9.1) - 2021-05-06

- [`07945ec`](https://github.com/alexfalkowski/go-service/commit/07945ec4ac798539b14587debcd04cdd73a3415c) docs: fix wording (#16)

## [v1.9.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.9.0) - 2021-05-06

- [`71c7d63`](https://github.com/alexfalkowski/go-service/commit/71c7d6314a94ce0e475ed288d58edd6c60d67095) feat: move the deps (#15)

## [v1.8.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.8.0) - 2021-05-06

- [`f8932ed`](https://github.com/alexfalkowski/go-service/commit/f8932ed1bddee41a7a4eef730593735822622cc2) feat(datadog): add tracer (#14)

## [v1.7.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.7.0) - 2021-05-06

- [`06bd01f`](https://github.com/alexfalkowski/go-service/commit/06bd01fa5251ce6d95bb4baea2d4f63f848d19fa) feat(ristretto): add in memory cache (#13)

## [v1.6.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.6.0) - 2021-05-06

- [`8cea703`](https://github.com/alexfalkowski/go-service/commit/8cea7034bb4cc543b5eda699cb9636868ffdf240) feat(config): change key names (#12)

## [v1.5.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.5.0) - 2021-05-05

- [`7cf4f77`](https://github.com/alexfalkowski/go-service/commit/7cf4f77cf89ee8e683e367184e1c503e1bfbbb9f) feat(grpc): pass insecure to client with no options (#11)

## [v1.4.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.4.2) - 2021-05-05

- [`07591ab`](https://github.com/alexfalkowski/go-service/commit/07591abfb465b45217ec83871784dc47f9923138) test: verify grpc and http client token errors (#10)

## [v1.4.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.4.1) - 2021-05-05

- [`17c7a94`](https://github.com/alexfalkowski/go-service/commit/17c7a94a92c479d315dd2d121a1dd486e7baedfd) test: add more tests for security (#9)

## [v1.4.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.4.0) - 2021-05-05

- [`0da0143`](https://github.com/alexfalkowski/go-service/commit/0da01438437b0a6a8e773c356c6e7030c33d8f28) feat(security): allow the generation and verification of tokens (#8)

## [v1.3.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.3.0) - 2021-05-05

- [`3ec8c3d`](https://github.com/alexfalkowski/go-service/commit/3ec8c3d7e8dd8513068eee6d8e05336f8bd36d57) feat(cache): use snappy to compress (#7)

## [v1.2.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.2.0) - 2021-05-04

- [`168b8e6`](https://github.com/alexfalkowski/go-service/commit/168b8e667bf89194b89be4b400e8f764572cf34b) feat: move packages (#6)

## [v1.1.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.1.4) - 2021-05-03

- [`17edb4e`](https://github.com/alexfalkowski/go-service/commit/17edb4eebea49ce29528e97249389581ce6a9bd2) test: move test code (#5)

## [v1.1.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.1.3) - 2021-05-03

- [`8072525`](https://github.com/alexfalkowski/go-service/commit/80725250b2f7035db74dbd139afbabae5d0ddf59) test(http): make sure the client works with grpc gateway (#4)

## [v1.1.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.1.2) - 2021-05-02

- [`7086bb2`](https://github.com/alexfalkowski/go-service/commit/7086bb2697051b18c4a9e87633fc4b8e57093af5) fix: make sure we get request id correctly (#3)

## [v1.1.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.1.1) - 2021-05-02

- [`e50af18`](https://github.com/alexfalkowski/go-service/commit/e50af1884d8d3f21f855de824a2a7a06a313f651) fix(cache): remove lru (#2)

## [v1.1.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.1.0) - 2021-05-02

- [`1d98e32`](https://github.com/alexfalkowski/go-service/commit/1d98e32e6a54b8c32b6d6a02f6413a884d949f71) feat: move the implementation of the package from go-nonnative-example (#1)
- [`d6a997d`](https://github.com/alexfalkowski/go-service/commit/d6a997d9c87a980d9938b6430f5f56698cacbd9b) Initial commit
