# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## [v1.371.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.371.0) - 2025-01-16

- [`5cfa016`](https://github.com/alexfalkowski/go-service/commit/5cfa016ff5e5eb7048b23758720637cf52180171) feat(http): handle errors in a testable way (#1179)

## [v1.370.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.370.2) - 2025-01-16

- [`60cd87e`](https://github.com/alexfalkowski/go-service/commit/60cd87e5a2c8eb93ccc6ba58590ffc71e0d4cc02) fix(transport): geolocation can only be read from headers (#1178)
- [`4a603b4`](https://github.com/alexfalkowski/go-service/commit/4a603b4722353c91344e0dfff02e4ad0a7c83a67) test(events): add missing (#1176)

## [v1.370.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.370.1) - 2025-01-15

- [`1f8d4b1`](https://github.com/alexfalkowski/go-service/commit/1f8d4b15c4c15cb367e7153cf153c5e18a505006) fix(deps): upgraded google.golang.org/protobuf v1.36.2 => v1.36.3 (#1175)
- [`1631e83`](https://github.com/alexfalkowski/go-service/commit/1631e833ecad17a1ffb7b6c457a6ce2b4ab637e1) test(ed25519): add missing (#1174)

## [v1.370.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.370.0) - 2025-01-15

- [`849b23d`](https://github.com/alexfalkowski/go-service/commit/849b23d2215829ceb643df592166c7286ec13d06) feat(hmac): writting to a hash does not produce an error (#1173)
- [`dbbef5a`](https://github.com/alexfalkowski/go-service/commit/dbbef5a2122dfaab618660ec6809df697b64f060) docs(limiter): update kinds to point to transport/meta/key.go (#1172)
- [`e898a31`](https://github.com/alexfalkowski/go-service/commit/e898a31e45c746dfc622e410ba7de6d9928b5787) test(aes): add missing (#1170)

## [v1.369.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.369.1) - 2025-01-15

- [`d3acf6d`](https://github.com/alexfalkowski/go-service/commit/d3acf6dfbca9235825963fb4e0ee24eb1fde8d86) fix(mvc): write header on error (#1169)
- [`6dbee6e`](https://github.com/alexfalkowski/go-service/commit/6dbee6e2bc7bcf30d75e920d04f9f9312e6b4e68) docs(diagrams): update with make diagrams (#1168)
- [`80a6bad`](https://github.com/alexfalkowski/go-service/commit/80a6bad95871ce1bc50bff8f714b6e38f5a3e0d4) test(time): simplify (#1167)
- [`688eec6`](https://github.com/alexfalkowski/go-service/commit/688eec667bb5807b7504b497501bcd623bc51dc8) test(time): add missing (#1166)

## [v1.369.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.369.0) - 2025-01-14

- [`3b86153`](https://github.com/alexfalkowski/go-service/commit/3b861537edb6be3762bf3474d268290a870d4765) feat(time): use recover for nts (#1165)

## [v1.368.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.368.2) - 2025-01-14

- [`b5b40a3`](https://github.com/alexfalkowski/go-service/commit/b5b40a31446da624686c3cd4545015803b0eb08c) fix(feature): use must for metrics (#1164)

## [v1.368.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.368.1) - 2025-01-14

- [`785814a`](https://github.com/alexfalkowski/go-service/commit/785814adebf8e87f1d583e4d75dccea13d2ec58d) fix(rand): use reader for rand.Int (#1163)

## [v1.368.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.368.0) - 2025-01-14

- [`110bce2`](https://github.com/alexfalkowski/go-service/commit/110bce28b2001a3b8154e7f8839beb938895af14) feat(rpc): allow errors to occur if request is nil (#1161)
- [`a7313c5`](https://github.com/alexfalkowski/go-service/commit/a7313c5f56687daaa7878ec738f3e2a2e0ece5d1) test(http): add missing (#1160)
- [`6c78712`](https://github.com/alexfalkowski/go-service/commit/6c7871296f70322913f6cdffde31fad9738109aa) build(ci): update go:1.5 (#1159)

## [v1.367.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.367.0) - 2025-01-14

- [`9f9415a`](https://github.com/alexfalkowski/go-service/commit/9f9415a884d0d42cb721eba74534773c1a736fc7) feat(crypto): add generators (#1158)

## [v1.366.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.366.0) - 2025-01-13

- [`d6e7221`](https://github.com/alexfalkowski/go-service/commit/d6e722120c648df53f15ceed8c0f1a3f30a9e834) feat(token): no need to check header for jwt (#1157)

## [v1.365.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.365.0) - 2025-01-13

- [`34a7832`](https://github.com/alexfalkowski/go-service/commit/34a783294f228b5937817fd51fd11f6cc61fa2c8) feat(crypto): use recover for ed25519 (#1155)
- [`264d062`](https://github.com/alexfalkowski/go-service/commit/264d06244ac21967d813300338cd0dd5fbd2397f) test(server): add missing (#1153)
- [`6288dda`](https://github.com/alexfalkowski/go-service/commit/6288dda798a0c78b6cbaf7721b156af69ed12b5a) test(encoding): add missing (#1152)
- [`89f18ab`](https://github.com/alexfalkowski/go-service/commit/89f18abfc1cf4678a791cdfc405b5b35bb4ff5c2) test(crypto): add missing (#1151)
- [`4e944f1`](https://github.com/alexfalkowski/go-service/commit/4e944f193b06ca9bd0dfa6b37ae714498fa3b0e8) test(sql): add missing (#1149)
- [`46055ae`](https://github.com/alexfalkowski/go-service/commit/46055ae76099a9a1999adcbb3a4e12d0e1bf3544) test(crypto): add missing (#1148)

## [v1.364.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.364.1) - 2025-01-13

- [`a6d765b`](https://github.com/alexfalkowski/go-service/commit/a6d765b093a5f25bb42bd9384dc252c108414656) fix(deps): upgraded google.golang.org/grpc v1.69.2 => v1.69.4 (#1147)

## [v1.364.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.364.0) - 2025-01-12

- [`5adb399`](https://github.com/alexfalkowski/go-service/commit/5adb399a0c0ea2a32f87c1d105ecb84deefa6537) feat(mvc): ignore adding if invalid (#1146)

## [v1.363.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.363.0) - 2025-01-12

- [`c0ddd25`](https://github.com/alexfalkowski/go-service/commit/c0ddd259554e5819b8b9712157335d65e7ffe4fe) feat(limiter): remove interface (#1145)
- [`0b227b2`](https://github.com/alexfalkowski/go-service/commit/0b227b24b767b95925c48fbf557569b911fbc29e) test(limiter): reword criteria (#1144)

## [v1.362.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.362.0) - 2025-01-12

- [`631947e`](https://github.com/alexfalkowski/go-service/commit/631947e231429cb7f8e92990d06f0f2e3ba1bb7d) feat(limiter): create an interface to have better abstraction (#1143)
- [`943fc7f`](https://github.com/alexfalkowski/go-service/commit/943fc7fe00ab28bb9508159f3e62b6a8b798453c) test(redis): add missing (#1142)

## [v1.361.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.361.2) - 2025-01-12

- [`d1e1f3b`](https://github.com/alexfalkowski/go-service/commit/d1e1f3bdd527c2c4e27758f4370987267e7ea89f) fix(env): remove unused func (#1141)

## [v1.361.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.361.1) - 2025-01-12

- [`eb3ffd2`](https://github.com/alexfalkowski/go-service/commit/eb3ffd27ff4a4094454813e889dc21e025377ad2) fix(net): remove unused func (#1140)

## [v1.361.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.361.0) - 2025-01-11

- [`bb38b67`](https://github.com/alexfalkowski/go-service/commit/bb38b67205c623c51cc9f10c441f9cdd1ba5b771) feat(retry): migrate to github.com/sethvargo/go-retry (#1139)

## [v1.360.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.360.0) - 2025-01-10

- [`f79ca59`](https://github.com/alexfalkowski/go-service/commit/f79ca591f49c954e0325f5a7ab4db7815c038472) feat(debug): ignore errors (#1138)
- [`c697caf`](https://github.com/alexfalkowski/go-service/commit/c697cafdcb637ae812fb16e239a96a9e9a4a2f6e) test(crypto): add missing (#1136)
- [`c2a44a9`](https://github.com/alexfalkowski/go-service/commit/c2a44a95eddeb4061d06ec6f6c72f6f23a91ac9b) test(hooks): add missing (#1135)
- [`57ece9b`](https://github.com/alexfalkowski/go-service/commit/57ece9b0c08f5298ff234ad43b51d236e67d95b1) docs(diagrams): update with make diagrams (#1134)
- [`b763992`](https://github.com/alexfalkowski/go-service/commit/b76399280a5c9a6aa10e9662f4bce9e81bb5e0d3) test(all): add missing (#1133)
- [`a86a8de`](https://github.com/alexfalkowski/go-service/commit/a86a8de41201061e2e0a29ccfbdf30179de2db07) test(varnamelen): add more meaningful variable names (#1132)
- [`0cbc569`](https://github.com/alexfalkowski/go-service/commit/0cbc569286bc11f554dc29eb530d7fcc3ad0a35f) test(gochecknoinits): remove exclude (#1131)
- [`91814be`](https://github.com/alexfalkowski/go-service/commit/91814beca4a847980fd2565c45f743df49ad1231) build(fatcontext): remove exclude (#1130)
- [`d337cd3`](https://github.com/alexfalkowski/go-service/commit/d337cd3ce4cd41a155ac43714255cb28bd37c577) test(varnamelen): remove exclude (#1129)

## [v1.359.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.359.0) - 2025-01-09

- [`70460f0`](https://github.com/alexfalkowski/go-service/commit/70460f0c859e0ec014ab38a451bbc18a6b1c641c) feat(varnamelen): add more meaningful variable names (#1128)

## [v1.358.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.358.0) - 2025-01-09

- [`42ba6b1`](https://github.com/alexfalkowski/go-service/commit/42ba6b1e06a24ce03f2c28abb6e97abbdd102fde) feat(linter): add more meaningful variables from varnamelen (#1127)

## [v1.357.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.357.1) - 2025-01-09

- [`3e74f68`](https://github.com/alexfalkowski/go-service/commit/3e74f6824c38818fd7f935a60454288c10815af9) fix(deps): upgraded github.com/go-resty/resty/v2 v2.16.2 => v2.16.3 (#1126)

## [v1.357.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.357.0) - 2025-01-08

- [`4466d39`](https://github.com/alexfalkowski/go-service/commit/4466d39671e1cd45d1c27d5b49f718452d05f240) feat(http): remove handler type (#1125)
- [`b6aae88`](https://github.com/alexfalkowski/go-service/commit/b6aae88b35e9d7cfd5893ad7b6bac0fdd48022e3) build(deps): bump bin from `a00abbe` to `cb313fe` (#1123)
- [`9694d51`](https://github.com/alexfalkowski/go-service/commit/9694d510e82713fd95ddc372ef68c9e862a280c7) build(deps): bump bin from `a278340` to `a00abbe` (#1122)

## [v1.356.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.356.1) - 2025-01-07

- [`b081e7d`](https://github.com/alexfalkowski/go-service/commit/b081e7dd3f8284c5931e9b666dab562416c15de1) fix(deps): upgraded google.golang.org/protobuf v1.36.1 => v1.36.2 (#1121)

## [v1.356.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.356.0) - 2025-01-07

- [`94e7327`](https://github.com/alexfalkowski/go-service/commit/94e73276a00be71b70c4dc0a88ddc7de303fbd3d) feat(deps): bump github.com/KimMachineGun/automemlimit from 0.6.1 to 0.7.0 (#1120)
- [`ab95dca`](https://github.com/alexfalkowski/go-service/commit/ab95dca3e50cbc138e917e2270a3ea0150fedb6a) test(ctx): have a default timeout (#1119)

## [v1.355.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.355.0) - 2025-01-06

- [`f7ef194`](https://github.com/alexfalkowski/go-service/commit/f7ef194990807c817de08c582ba0b81027886801) feat(deps): bump golang.org/x/net from 0.33.0 to 0.34.0 (#1118)

## [v1.354.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.354.0) - 2025-01-06

- [`c1979c9`](https://github.com/alexfalkowski/go-service/commit/c1979c9f27338571543482eaa8f00ba90286c748) feat(deps): bump golang.org/x/crypto from 0.31.0 to 0.32.0 (#1117)
- [`075e7ef`](https://github.com/alexfalkowski/go-service/commit/075e7ef61a70ba8eb14ed2cab43a89f2e9f844a7) build(deps): bump bin from `08600e7` to `a278340` (#1116)

## [v1.353.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.353.1) - 2025-01-04

- [`f863424`](https://github.com/alexfalkowski/go-service/commit/f86342433690f6049300f3a642a3fb85c6b919b2) fix(deps): upgraded github.com/standard-webhooks/standard-webhooks/libraries v0.0.0-20241231132107-9775b90ad9a4 => v0.0.0-20250103171228-b75b9ab8ea1e (#1115)
- [`cfcb0b5`](https://github.com/alexfalkowski/go-service/commit/cfcb0b5d93e50bef5d412c5e0a79e848f813e18b) build(deps): bump bin from `bde1e44` to `08600e7` (#1114)
- [`967700c`](https://github.com/alexfalkowski/go-service/commit/967700c3429fcea3da82a7affc07b70aa46f4239) docs(diagrams): update with make diagrams (#1113)

## [v1.353.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.353.0) - 2025-01-03

- [`fa802e6`](https://github.com/alexfalkowski/go-service/commit/fa802e6afb62f7b493a824f080ee460107f88bbf) feat(fmt): remove Sprintf (#1112)
- [`7bf6c53`](https://github.com/alexfalkowski/go-service/commit/7bf6c53c3c48160686eb1498ad2ce25dd15e0136) build(ci): for loop can be changed to use an integer range (Go 1.22+) (#1111)

## [v1.352.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.352.2) - 2025-01-02

- [`7086fc4`](https://github.com/alexfalkowski/go-service/commit/7086fc45652ac5937f50b8904352cd584fd941e8) fix(deps): github.com/shirou/gopsutil/v4 v4.24.11 => v4.24.12 (#1110)

## [v1.352.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.352.1) - 2024-12-29

- [`48c6272`](https://github.com/alexfalkowski/go-service/commit/48c6272221a46d3db72c02469e70424d72a36023) fix(http): use meta map for health (#1108)

## [v1.352.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.352.0) - 2024-12-28

- [`5d2132f`](https://github.com/alexfalkowski/go-service/commit/5d2132fc1b633e1e9dfb21bbc5a5b5221fa2a2b0) feat(health): return the error as a combined error from health (#1106)

## [v1.351.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.351.1) - 2024-12-24

- [`4b96f19`](https://github.com/alexfalkowski/go-service/commit/4b96f19392b3d1184937c6015e757a8b255455ea) fix(deps): upgraded google.golang.org/protobuf v1.36.0 => v1.36.1 (#1105)
- [`4fbb99c`](https://github.com/alexfalkowski/go-service/commit/4fbb99ce1b285edaf5192daf4d2a4ed35e38c2b9) fix(deps): upgraded github.com/matthewhartstonge/argon2 v1.1.0 => v1.1.1 (#1104)

## [v1.351.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.351.0) - 2024-12-23

- [`cf0a2bc`](https://github.com/alexfalkowski/go-service/commit/cf0a2bc6591795083f2523a726f3f9557a755072) feat(maps): remove any and use types (#1101)
- [`f7348d4`](https://github.com/alexfalkowski/go-service/commit/f7348d4e5bba1c596819bea54d23d30f3689b575) build(deps): bump bin from `a433391` to `bde1e44` (#1100)

## [v1.350.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.350.0) - 2024-12-23

- [`564bf95`](https://github.com/alexfalkowski/go-service/commit/564bf956ecebce4ef65deab31ab597a6f92d1e2a) feat(deps): bump github.com/matthewhartstonge/argon2 from 1.0.3 to 1.1.0 (#1099)

## [v1.349.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.349.0) - 2024-12-23

- [`c743d20`](https://github.com/alexfalkowski/go-service/commit/c743d20c6557726a33480c41886775d77f07ec6d) feat(deps): bump github.com/jackc/pgx/v5 from 5.7.1 to 5.7.2 (#1098)
- [`edca015`](https://github.com/alexfalkowski/go-service/commit/edca0159314fe21da767b715958277cbf23075b8) build(deps): bump bin from `9fa29e2` to `a433391` (#1097)

## [v1.348.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.348.1) - 2024-12-22

- [`ff3432f`](https://github.com/alexfalkowski/go-service/commit/ff3432f2f8ee3055c028815453450a761f883dae) fix(grpc): add removed code (#1096)

## [v1.348.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.348.0) - 2024-12-22

- [`3c8a1ab`](https://github.com/alexfalkowski/go-service/commit/3c8a1abe0d1f4bdaa063924380c9cb0a89e718c1) feat(mem): improve allocations (#1095)
- [`d1ba91f`](https://github.com/alexfalkowski/go-service/commit/d1ba91fe7e099dc44cc17db817403d08520b2111) build(deps): bump bin from `10049eb` to `9fa29e2` (#1093)

## [v1.347.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.347.0) - 2024-12-22

- [`801f805`](https://github.com/alexfalkowski/go-service/commit/801f80592b714263b4c662500fc53fc6ccdcff83) feat(http): use content handlers (#1090)

## [v1.346.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.346.0) - 2024-12-21

- [`5a9ead6`](https://github.com/alexfalkowski/go-service/commit/5a9ead60daf7cc9910aba2074b12b3495e298447) feat(http): generic rest methods (#1089)

## [v1.345.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.345.0) - 2024-12-19

- [`3dcfffc`](https://github.com/alexfalkowski/go-service/commit/3dcfffcae509c445c40d698eccd18cd67e9de669) feat(crypto): add ability to generate random letters (#1088)

## [v1.344.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.344.0) - 2024-12-18

- [`f259975`](https://github.com/alexfalkowski/go-service/commit/f2599758f8ad44cd7bbd878838cba1e03585fc56) feat(deps): upgraded golang.org/x/net v0.32.0 => v0.33.0 (#1087)

## [v1.343.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.343.1) - 2024-12-18

- [`c6c12ce`](https://github.com/alexfalkowski/go-service/commit/c6c12ce92ead3b1de9d0d5463f6d1454264d41cb) fix(deps): upgraded google.golang.org/grpc v1.69.0 => v1.69.2 (#1086)

## [v1.343.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.343.0) - 2024-12-18

- [`7a616c7`](https://github.com/alexfalkowski/go-service/commit/7a616c786556e445e550cc61115b9a6753aa7963) feat(structs): add zero check (#1085)

## [v1.342.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.342.0) - 2024-12-17

- [`63395ab`](https://github.com/alexfalkowski/go-service/commit/63395ab1ed6c0010c573b6ceada39a4029aa4d2d) feat(health): add db checker (#1084)

## [v1.341.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.341.0) - 2024-12-16

- [`7af1b6b`](https://github.com/alexfalkowski/go-service/commit/7af1b6b113a1f950fa307accf8be3558a506a412) feat(deps): bump google.golang.org/protobuf from 1.35.2 to 1.36.0 (#1083)
- [`ab0bbf4`](https://github.com/alexfalkowski/go-service/commit/ab0bbf4f89a22a45aa99171a7e131e91fe7578f6) docs(diagrams): update diagrams (#1082)

## [v1.340.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.340.0) - 2024-12-15

- [`6a87282`](https://github.com/alexfalkowski/go-service/commit/6a872827ca88481ce97e49a5b7d4dce504dd36e9) feat(http): rename to be media type for content (#1081)

## [v1.339.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.339.0) - 2024-12-14

- [`3b69263`](https://github.com/alexfalkowski/go-service/commit/3b6926340de993d4eb8b64d2a491350a9f57e770) feat(http): remove h2c (#1080)

## [v1.338.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.338.0) - 2024-12-13

- [`b1b29a7`](https://github.com/alexfalkowski/go-service/commit/b1b29a7cb7b52736bc78299bd6476b9d126c273a) feat(cmd): remove server/client commands (#1079)

## [v1.337.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.337.0) - 2024-12-12

- [`125131b`](https://github.com/alexfalkowski/go-service/commit/125131bbaa6e162c6add10f9375774bc6fbcd5c1) feat(deps): upgraded go.opentelemetry.io/otel v1.32.0 => v1.33.0 (#1078)

## [v1.336.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.336.0) - 2024-12-12

- [`d71603f`](https://github.com/alexfalkowski/go-service/commit/d71603f5c9936d827c9962a0df40bd7044446724) feat(deps): upgraded google.golang.org/grpc v1.68.1 => v1.69.0 (#1077)

## [v1.335.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.335.1) - 2024-12-12

- [`4cc86dc`](https://github.com/alexfalkowski/go-service/commit/4cc86dc55eaa552d66df0c80b0d412c127139f5e) fix(deps): upgraded github.com/goccy/go-json v0.10.3 => v0.10.4 (#1076)

## [v1.335.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.335.0) - 2024-12-12

- [`3dbb73f`](https://github.com/alexfalkowski/go-service/commit/3dbb73f737b5e18de2062dc99fe872bd9fcad3fc) feat(deps): bump github.com/grpc-ecosystem/go-grpc-middleware/v2 from 2.1.0 to 2.2.0 (#1074)

## [v1.334.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.334.0) - 2024-12-12

- [`133ad61`](https://github.com/alexfalkowski/go-service/commit/133ad616a85d713e94302518c25c8526c3547228) feat(deps): bump github.com/matthewhartstonge/argon2 from 1.0.2 to 1.0.3 (#1073)

## [v1.333.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.333.0) - 2024-12-12

- [`2191a0c`](https://github.com/alexfalkowski/go-service/commit/2191a0c108b9a88e131d0101e6a2ec42b853808a) feat(deps): bump aidanwoods.dev/go-paseto from 1.5.2 to 1.5.3 (#1075)

## [v1.332.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.332.2) - 2024-12-11

- [`b924cf8`](https://github.com/alexfalkowski/go-service/commit/b924cf887ec685e839551d403290b3d7e368bfc2) fix(config): too many prefixes (#1071)

## [v1.332.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.332.1) - 2024-12-11

- [`00f4cbd`](https://github.com/alexfalkowski/go-service/commit/00f4cbdcd5ee227addd5308ec925c64f094d8caf) fix(config): too many prefixes (#1070)

## [v1.332.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.332.0) - 2024-12-11

- [`2169fa3`](https://github.com/alexfalkowski/go-service/commit/2169fa38ff4cda72ea38c8d0aaf2a0ad46c10a72) feat(errors): add prefix (#1069)

## [v1.331.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.331.1) - 2024-12-11

- [`970a34d`](https://github.com/alexfalkowski/go-service/commit/970a34dddb31c36ff314382d74b8dd592c40d944) fix(config): add prefix (#1068)

## [v1.331.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.331.0) - 2024-12-11

- [`a3eb840`](https://github.com/alexfalkowski/go-service/commit/a3eb840840d6c40a091344512dad93d2ace952ee) feat(errors): shorten prefixes (#1067)

## [v1.330.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.330.0) - 2024-12-11

- [`b6861c2`](https://github.com/alexfalkowski/go-service/commit/b6861c233ffe27dd7846e91c22dfcf78c6e0dbbf) feat(deps): bump github.com/open-feature/go-sdk from 1.13.1 to 1.14.0 (#1065)
- [`57dbb8a`](https://github.com/alexfalkowski/go-service/commit/57dbb8a376f735ed2f370115b5f023b20c56c703) build(deps): bump bin from `cf4a7d3` to `10049eb` (#1066)

## [v1.329.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.329.0) - 2024-12-10

- [`6fce124`](https://github.com/alexfalkowski/go-service/commit/6fce12436b5a1a5aed1a5b01825b7b3f925b1cde) feat(config): generic new config (#1064)

## [v1.328.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.328.1) - 2024-12-10

- [`728310d`](https://github.com/alexfalkowski/go-service/commit/728310d23dd499190524d8f75d9de0fb0d5d2b2a) fix(cmd): error on decode with no encoder (#1062)

## [v1.328.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.328.0) - 2024-12-09

- [`b8ffdd0`](https://github.com/alexfalkowski/go-service/commit/b8ffdd0d757bed2a9a0e8ae2d46c08905fb078b3) feat(rest): add error from response (#1061)

## [v1.327.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.327.0) - 2024-12-09

- [`bdd34f3`](https://github.com/alexfalkowski/go-service/commit/bdd34f38b1df2ad02c3b40f8ad38c88ef80bc951) feat(cmd): check string flag (#1060)

## [v1.326.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.326.0) - 2024-12-07

- [`4818000`](https://github.com/alexfalkowski/go-service/commit/4818000ffe435cb14d99d08c8897676922f4e8ae) feat(http): content handler to return error (#1059)

## [v1.325.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.325.0) - 2024-12-06

- [`ed43862`](https://github.com/alexfalkowski/go-service/commit/ed438620667f56372599cff77501d15bda1e3c11) feat(token): add key (#1058)

## [v1.324.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.324.0) - 2024-12-06

- [`738c443`](https://github.com/alexfalkowski/go-service/commit/738c4438862986cb0f5b6ab2b19fd13ec198ad4b) feat(token): add jwt/paseto (#1057)

## [v1.323.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.323.0) - 2024-12-05

- [`3ad2524`](https://github.com/alexfalkowski/go-service/commit/3ad25243c794c5b9c3892ad32cdf4cf87336818d) feat(deps): bump golang.org/x/net from 0.31.0 to 0.32.0 (#1056)

## [v1.322.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.322.0) - 2024-12-04

- [`e7a7790`](https://github.com/alexfalkowski/go-service/commit/e7a779054db2f86f6874a845a5b11859dddd3cee) feat(deps): bump golang.org/x/crypto from 0.29.0 to 0.30.0 (#1055)

## [v1.321.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.321.0) - 2024-12-04

- [`035d30e`](https://github.com/alexfalkowski/go-service/commit/035d30e9c5c1092ec35ae19512abbae8d72e6c5b) feat(deps): bump google.golang.org/grpc from 1.68.0 to 1.68.1 (#1054)
- [`28e1a41`](https://github.com/alexfalkowski/go-service/commit/28e1a41675affd1ffdd9e6e7db20a0e39bb630f1) build(lint): remove disabled linters (#1053)

## [v1.320.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.320.1) - 2024-12-01

- [`d58869c`](https://github.com/alexfalkowski/go-service/commit/d58869cb5c7f9828355169fffaea1dc3dbaaae34) fix(deps): upgraded github.com/shirou/gopsutil/v4 v4.24.10 => v4.24.11 (#1052)

## [v1.320.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.320.0) - 2024-12-01

- [`e999237`](https://github.com/alexfalkowski/go-service/commit/e9992379794ad6bca8cfce7de43c9c85bf575ea5) feat(lint): enable errcheck (#1051)

## [v1.319.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.319.0) - 2024-12-01

- [`9714d16`](https://github.com/alexfalkowski/go-service/commit/9714d162c732205365f54dea31b4b45483077620) feat(lint): enable err113 (#1050)

## [v1.318.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.318.0) - 2024-11-30

- [`c9a3b75`](https://github.com/alexfalkowski/go-service/commit/c9a3b75eec0f89b2060db3a492c4655d7faa6457) feat(lint): enable linters (#1049)
- [`f8873f4`](https://github.com/alexfalkowski/go-service/commit/f8873f4c529bbacbe62f0fa9f3b989d6a053504f) build(lint): enable all and disbale others (#1048)
- [`bf999d7`](https://github.com/alexfalkowski/go-service/commit/bf999d7360e9b1303d6d96035f9d8287e22f05db) test(rest): remove duplication (#1047)

## [v1.317.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.317.0) - 2024-11-29

- [`df70ce5`](https://github.com/alexfalkowski/go-service/commit/df70ce5817cdb30474ba058245cfc514c90b2406) feat(rest): add missing methods (#1046)

## [v1.316.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.316.1) - 2024-11-28

- [`4979dfe`](https://github.com/alexfalkowski/go-service/commit/4979dfe37c2b3da8364a77c57132283088749da1) fix(rest): remove content type (#1045)

## [v1.316.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.316.0) - 2024-11-28

- [`c678459`](https://github.com/alexfalkowski/go-service/commit/c6784593d78f906a890ac03a0da345b60bfb30f3) feat(rest): use resty (#1044)
- [`bc7f1f4`](https://github.com/alexfalkowski/go-service/commit/bc7f1f4c68a2e240f7f99c48d3b34a4b3980cfad) test(rpc): no content (#1043)

## [v1.315.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.315.0) - 2024-11-20

- [`78c087d`](https://github.com/alexfalkowski/go-service/commit/78c087df324a7b6d869db0ffdd838e4f388931cb) feat(rest): add basic rest client (#1042)
- [`7730544`](https://github.com/alexfalkowski/go-service/commit/7730544f7e7a0c38d0c6ca8ce6e3847e8f978dd5) docs(readme): update diagrams (#1041)

## [v1.314.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.314.0) - 2024-11-19

- [`9c85680`](https://github.com/alexfalkowski/go-service/commit/9c8568071fa172ed5cd0c8a907db13efbb5bf7b5) feat(http): add rest (#1040)

## [v1.313.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.313.0) - 2024-11-14

- [`9bd47b0`](https://github.com/alexfalkowski/go-service/commit/9bd47b04336d8ad8c60a1dbd72f32bb7be9fd33f) feat(mvc): add ability to serve static files (#1039)

## [v1.312.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.312.1) - 2024-11-14

- [`10e8069`](https://github.com/alexfalkowski/go-service/commit/10e80698cf6c798c89b5689343e733d388a0cea4) fix(deps): upgraded google.golang.org/protobuf v1.35.1 => v1.35.2 (#1038)

## [v1.312.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.312.0) - 2024-11-12

- [`8be9380`](https://github.com/alexfalkowski/go-service/commit/8be9380e2763b13bac20b8187810d24d4bd61e76) feat(cmp): use or (#1037)

## [v1.311.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.311.0) - 2024-11-11

- [`8e83cf2`](https://github.com/alexfalkowski/go-service/commit/8e83cf29d44dec1a6f9898684e955d1ddbe89cc2) feat(deps): bump github.com/matthewhartstonge/argon2 from 1.0.1 to 1.0.2 (#1036)

## [v1.310.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.310.0) - 2024-11-10

- [`da6898b`](https://github.com/alexfalkowski/go-service/commit/da6898b18815f33952f205ed9b9d40d5edb0efbb) feat(deps): upgraded go.opentelemetry.io/otel v1.31.0 => v1.32.0 (#1035)

## [v1.309.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.309.0) - 2024-11-08

- [`ff1371b`](https://github.com/alexfalkowski/go-service/commit/ff1371b1efa4d811f529470600c593d16dcae1b8) feat(deps): bump golang.org/x/crypto from 0.28.0 to 0.29.0 (#1034)

## [v1.308.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.308.0) - 2024-11-07

- [`6e6ff51`](https://github.com/alexfalkowski/go-service/commit/6e6ff51fbdd3a090e03674860f33b162d0ee9b0e) feat(deps): upgraded google.golang.org/grpc v1.67.1 => v1.68.0 (#1033)

## [v1.307.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.307.2) - 2024-11-01

- [`9d0e980`](https://github.com/alexfalkowski/go-service/commit/9d0e9809082c6fe81f3d9937e5cc2a56ed272313) fix(deps): upgraded github.com/shirou/gopsutil/v4 v4.24.9 => v4.24.10 (#1032)
- [`1e57ea3`](https://github.com/alexfalkowski/go-service/commit/1e57ea3f11d37bd7c5edcd481100e78ed66459cb) build(deps): bump bin from `e8f9d73` to `cf4a7d3` (#1031)

## [v1.307.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.307.1) - 2024-10-18

- [`9daabb3`](https://github.com/alexfalkowski/go-service/commit/9daabb30e66fdf180d0c06b8d87929daeeb27fa2) fix(deps): upgraded github.com/open-feature/go-sdk v1.13.0 => v1.13.1 (#1030)

## [v1.307.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.307.0) - 2024-10-17

- [`5d7f7d3`](https://github.com/alexfalkowski/go-service/commit/5d7f7d3e8ab53216d432ba765f62b74652de3db9) feat(deps): bump github.com/redis/go-redis/v9 from 9.6.2 to 9.7.0 (#1029)

## [v1.306.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.306.1) - 2024-10-16

- [`bfbdbb2`](https://github.com/alexfalkowski/go-service/commit/bfbdbb2dc2e7df8fa25195e5d6a48e8ac7fe73d3) fix(deps): upgraded github.com/prometheus/client_golang v1.20.4 => v1.20.5 (#1028)

## [v1.306.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.306.0) - 2024-10-14

- [`31c4654`](https://github.com/alexfalkowski/go-service/commit/31c46541f894153a8dd01572729462f3fce268c5) feat(deps): upgraded go.opentelemetry.io/contrib/instrumentation v0.55.0 => v0.56.0 (#1026)

## [v1.305.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.305.1) - 2024-10-14

- [`3461aca`](https://github.com/alexfalkowski/go-service/commit/3461acac137155ee27a59218c8402a0831bb2480) fix(deps): upgraded github.com/redis/go-redis/v9 v9.6.1 => v9.6.2 (#1025)

## [v1.305.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.305.0) - 2024-10-11

- [`aeea234`](https://github.com/alexfalkowski/go-service/commit/aeea234dfafec5a2153b1b03ef38bd9cff1b19b0) feat(deps): upgraded go.opentelemetry.io/otel v1.30.0 => v1.31.0 (#1024)

## [v1.304.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.304.0) - 2024-10-11

- [`7d77c8a`](https://github.com/alexfalkowski/go-service/commit/7d77c8ab3c69e2aac04e28654e8201ece08c5107) feat(deps): upgraded go.uber.org/fx v1.22.2 => v1.23.0 (#1023)

## [v1.303.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.303.1) - 2024-10-11

- [`e3856c1`](https://github.com/alexfalkowski/go-service/commit/e3856c1cd0cd620add6961f36effb36a590d4c6f) fix(deps): upgraded github.com/klauspost/compress v1.17.10 => v1.17.11 (#1022)

## [v1.303.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.303.0) - 2024-10-08

- [`dd2d205`](https://github.com/alexfalkowski/go-service/commit/dd2d205cd88496dbd9cdda29c513bb7e327d2cb0) feat(deps): bump google.golang.org/protobuf from 1.34.2 to 1.35.1 (#1021)

## [v1.302.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.302.0) - 2024-10-07

- [`7b84a43`](https://github.com/alexfalkowski/go-service/commit/7b84a43d817f27dec2e67951c58fb57b70c6ba3d) feat(deps): bump github.com/beevik/nts from 0.1.1 to 0.2.0 (#1020)

## [v1.301.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.301.0) - 2024-10-07

- [`32b88e9`](https://github.com/alexfalkowski/go-service/commit/32b88e98cd44791f44ac3bda2aefbc15cd2d4739) feat(deps): bump golang.org/x/net from 0.29.0 to 0.30.0 (#1018)
- [`30a7e90`](https://github.com/alexfalkowski/go-service/commit/30a7e90aad29a558bb86cbbdd6b6f36b0362777f) build(deps): bump bin from `8b87d26` to `e8f9d73` (#1017)

## [v1.300.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.300.2) - 2024-10-01

- [`98a05c3`](https://github.com/alexfalkowski/go-service/commit/98a05c33c5ef8238f33a59d63f681f8c422f05ec) fix(deps): update google.golang.org/grpc to v1.67.1 (#1016)
- [`a97ddd0`](https://github.com/alexfalkowski/go-service/commit/a97ddd05344c93602da36b1913da1d8e069e4462) build(make): seperate targets for benchmarks (#1015)

## [v1.300.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.300.1) - 2024-09-25

- [`85debb9`](https://github.com/alexfalkowski/go-service/commit/85debb9f2591656822a7076dda50f05971be4350) fix(grpc): verify nil in server (#1014)

## [v1.300.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.300.0) - 2024-09-24

- [`4e05cae`](https://github.com/alexfalkowski/go-service/commit/4e05cae8594505a3704cf1692c9a0e5e20c88bf4) feat(http): remove cors (#1012)

## [v1.299.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.299.0) - 2024-09-24

- [`3be9967`](https://github.com/alexfalkowski/go-service/commit/3be9967515898e6e7893e015f6cf60d24a8c1bd8) feat(deps): bump go.uber.org/automaxprocs from 1.5.3 to 1.6.0 (#1011)

## [v1.298.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.298.0) - 2024-09-24

- [`82d2e5e`](https://github.com/alexfalkowski/go-service/commit/82d2e5e3f5d0677be02caa02050cb27e661a1613) feat(deps): bump github.com/klauspost/compress from 1.17.9 to 1.17.10 (#1010)

## [v1.297.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.297.0) - 2024-09-20

- [`9badf58`](https://github.com/alexfalkowski/go-service/commit/9badf5806e70119bfbae8e744bb2d948df0fc49b) feat(deps): update grpc to v1.67.0 (#1009)
- [`0234b78`](https://github.com/alexfalkowski/go-service/commit/0234b78f9d67fb116f1cfa04d05ac7014d4c605d) test(http): ignore logger if not present (#1008)
- [`781f914`](https://github.com/alexfalkowski/go-service/commit/781f9140caa58591925344fac9a661acc7984e42) test(http): benchmark different handlers (#1007)

## [v1.296.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.296.0) - 2024-09-17

- [`eddc0c6`](https://github.com/alexfalkowski/go-service/commit/eddc0c632fff47bed06f262b9702058c81b9def2) feat(cache): remove ristreto (#1006)

## [v1.295.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.295.1) - 2024-09-17

- [`8f43726`](https://github.com/alexfalkowski/go-service/commit/8f4372676dcf2eb0a1f49fae259a11b0a2670198) fix(deps): update github.com/prometheus/client_golang to v1.20.4 (#1004)

## [v1.295.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.295.0) - 2024-09-17

- [`f8d743a`](https://github.com/alexfalkowski/go-service/commit/f8d743a55b4ef2c5c551e0fff6e495b9b84d3128) feat(deps): bump github.com/go-sprout/sprout from 0.5.1 to 0.6.0 (#1003)
- [`8f5360a`](https://github.com/alexfalkowski/go-service/commit/8f5360aa95da049877b0fe67eb4a4cbb7edcd3e0) test(http): add std benchmark (#1002)
- [`0f8a4b0`](https://github.com/alexfalkowski/go-service/commit/0f8a4b069936c54addb52c923441f51edcab37c9) build(ci): add benchmarks (#1000)
- [`d9fd255`](https://github.com/alexfalkowski/go-service/commit/d9fd255495327c6331f154a1c738ea27fd3a65e8) build(deps): bump bin from `b975eab` to `8b87d26` (#1001)

## [v1.294.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.294.0) - 2024-09-13

- [`33dab17`](https://github.com/alexfalkowski/go-service/commit/33dab172102d577dd27b455fae5afe8e58916a2b) feat(deps): bump github.com/matthewhartstonge/argon2 from 1.0.0 to 1.0.1 (#999)

## [v1.293.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.293.0) - 2024-09-12

- [`a8d8912`](https://github.com/alexfalkowski/go-service/commit/a8d89129d6725e8f3570d8f6c1ba2600bd45b299) feat(deps): update go.opentelemetry.io/contrib/instrumentation to  v0.55.0 (#998)

## [v1.292.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.292.1) - 2024-09-11

- [`d1ad370`](https://github.com/alexfalkowski/go-service/commit/d1ad370dbd5eaa974237d7cc26e40c1fe1fe3086) fix(deps): update google.golang.org/grpc to v1.66.2 (#995)

## [v1.292.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.292.0) - 2024-09-11

- [`a8886f4`](https://github.com/alexfalkowski/go-service/commit/a8886f4de5712d1ed7686e01744d437cc5fb9a42) feat(deps): update go.opentelemetry.io/otel to v1.30.0 (#994)
- [`931ea0c`](https://github.com/alexfalkowski/go-service/commit/931ea0cae2ec22747859faa52b35e9b811d81634) build(deps): bump bin from `1eca781` to `b975eab` (#988)
- [`c985687`](https://github.com/alexfalkowski/go-service/commit/c985687ecf9447a933204aeebdfdae5ef60ba907) docs(readme): limiter broken link (#987)
- [`fa1dbd1`](https://github.com/alexfalkowski/go-service/commit/fa1dbd1a2fcec01590ae83fc12cb76569e6b5863) docs(readme): broken link (#986)

## [v1.291.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.291.1) - 2024-09-10

- [`99f0ffc`](https://github.com/alexfalkowski/go-service/commit/99f0ffc165c70dc6fceada95ecb37670e58b40a0) fix(deps): update github.com/jackc/pgx/v5 to v5.7.1 (#985)

## [v1.291.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.291.0) - 2024-09-10

- [`35d3a96`](https://github.com/alexfalkowski/go-service/commit/35d3a96e5a6bd712e9f18bcf942a1a3d9d6fcc60) feat(deps): bump github.com/open-feature/go-sdk from 1.12.0 to 1.13.0 (#983)

## [v1.290.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.290.1) - 2024-09-10

- [`d0c9d43`](https://github.com/alexfalkowski/go-service/commit/d0c9d43c1374c5ef21ffad98f3cf63f8cdb12193) fix(deps): update google.golang.org/grpc to v1.66.1 (#984)

## [v1.290.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.290.0) - 2024-09-09

- [`d77c6a4`](https://github.com/alexfalkowski/go-service/commit/d77c6a4610cbb61d04003b740b6aad417e845177) feat(deps): bump github.com/jackc/pgx/v5 from 5.6.0 to 5.7.0 (#981)

## [v1.289.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.289.0) - 2024-09-06

- [`2489cdf`](https://github.com/alexfalkowski/go-service/commit/2489cdf0e1f91ee961bd4fb237dd2fa5331d3fce) feat(deps): bump golang.org/x/net from 0.28.0 to 0.29.0 (#979)

## [v1.288.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.288.2) - 2024-09-06

- [`9a007d4`](https://github.com/alexfalkowski/go-service/commit/9a007d4cbb867ed4d05948729004ed3ab09bc874) fix(dep): updated github.com/prometheus/client_golang to v1.20.3 (#980)

## [v1.288.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.288.1) - 2024-09-02

- [`74bbc42`](https://github.com/alexfalkowski/go-service/commit/74bbc422cc4b5765f05b3689e0b93a759da00f5d) fix(deps): update github.com/shirou/gopsutil/v4 to v4.24.8 (#976)

## [v1.288.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.288.0) - 2024-09-01

- [`c6e42d8`](https://github.com/alexfalkowski/go-service/commit/c6e42d8179e3b43ac4734467359b9bd384ef013e) feat(maps): add a map with string key and any value (#974)

## [v1.287.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.287.2) - 2024-08-31

- [`7295a29`](https://github.com/alexfalkowski/go-service/commit/7295a29377b54dfd512571533d4d6c9c491983d1) fix(deps): update github.com/felixge/fgprof to v0.9.5 (#973)

## [v1.287.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.287.1) - 2024-08-30

- [`49c0896`](https://github.com/alexfalkowski/go-service/commit/49c08962f80006312f3f3692fd654cbf30ee82c8) fix(deps): update github.com/rs/cors to v1.11.1 (#972)
- [`a7808df`](https://github.com/alexfalkowski/go-service/commit/a7808df5d4f0e7ba760d80319bf21e6b25fb080d) docs(diagrams): update with make diagrams (#970)

## [v1.287.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.287.0) - 2024-08-30

- [`162f34d`](https://github.com/alexfalkowski/go-service/commit/162f34d970cb917874bcbfb3b74a338b81b3aad3) feat(http): move to a content struct (#969)

## [v1.286.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.286.2) - 2024-08-29

- [`62926cb`](https://github.com/alexfalkowski/go-service/commit/62926cbf5fe1d29c5a1e280ea0fd31409faaa55d) fix(http): make sure rpc uses content handler (#968)

## [v1.286.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.286.1) - 2024-08-29

- [`1cf54bd`](https://github.com/alexfalkowski/go-service/commit/1cf54bd63c1f0ab221f80402e6240448c5b93be5) fix(http): make sure we set the media to json if encoder is not found (#967)

## [v1.286.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.286.0) - 2024-08-29

- [`469ffa0`](https://github.com/alexfalkowski/go-service/commit/469ffa0c11e9240c9c80f25605dbfc6f95679ac4) feat(http): reuse content negotiating handler (#966)

## [v1.285.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.285.0) - 2024-08-29

- [`eab95aa`](https://github.com/alexfalkowski/go-service/commit/eab95aa3483e8fd0c4166df04825a6b54c6939ec) feat(http): default to json if we can not find encoder (#965)

## [v1.284.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.284.0) - 2024-08-29

- [`0ae14e6`](https://github.com/alexfalkowski/go-service/commit/0ae14e6aaeead451c5c13dcf13cc7bae95c476c7) feat(health): make sure we use content negotiation (#964)

## [v1.283.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.283.0) - 2024-08-29

- [`812655c`](https://github.com/alexfalkowski/go-service/commit/812655c6a82a6ed23bcc0fecefe6037ecff19458) feat(encoding): remove unmarshal (#963)

## [v1.282.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.282.0) - 2024-08-28

- [`292055a`](https://github.com/alexfalkowski/go-service/commit/292055a22a1aed68829efbce231cc8d4652a61ce) feat(sync): use a buffer pool (#959)
- [`4ebc662`](https://github.com/alexfalkowski/go-service/commit/4ebc662a7556bb62eb8ac6337abb537cefdf24ef) build(deps): bump bin from `4d2cb2a` to `1eca781` (#961)

## [v1.281.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.281.0) - 2024-08-28

- [`e128672`](https://github.com/alexfalkowski/go-service/commit/e1286726a75f9f0c4e5b5d7ac27b078d7e633db7) feat(deps): bump google.golang.org/grpc from 1.65.0 to 1.66.0 (#960)
- [`05ec181`](https://github.com/alexfalkowski/go-service/commit/05ec18119bebca1a85c1530509bf8874df2e4c65) build(deps): bump bin from `41b7c8b` to `4d2cb2a` (#958)
- [`1dd56d6`](https://github.com/alexfalkowski/go-service/commit/1dd56d6590127ea03a6cf730c6153234cf72347f) build(sec): disable gosec lint (#957)

## [v1.280.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.280.1) - 2024-08-26

- [`17f3f0b`](https://github.com/alexfalkowski/go-service/commit/17f3f0bb0ef3c49d0a3faa113d56dfbd5a8d29d4) fix(deps): update github.com/prometheus/client_golang to v1.20.2 (#955)

## [v1.280.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.280.0) - 2024-08-26

- [`9bd0ff4`](https://github.com/alexfalkowski/go-service/commit/9bd0ff470fa973c2d3fd26dcef917a0efba864a0) feat(deps): update opentelemetry (#952)

## [v1.279.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.279.0) - 2024-08-21

- [`6f39414`](https://github.com/alexfalkowski/go-service/commit/6f3941419fddddf1d9d696bbd9d120797b264ab0) feat(deps): bump github.com/prometheus/client_golang from 1.20.0 to 1.20.1 (#946)

## [v1.278.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.278.0) - 2024-08-19

- [`dd4ef2a`](https://github.com/alexfalkowski/go-service/commit/dd4ef2a2d3da9422f4c9d0bf00e80ab90055512e) feat(deps): bump github.com/go-sprout/sprout from 0.5.0 to 0.5.1 (#945)

## [v1.277.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.277.0) - 2024-08-16

- [`a03d8da`](https://github.com/alexfalkowski/go-service/commit/a03d8dace49d3b4e37a5259517f2bfab88bdb3fb) feat(mvc): update github.com/go-sprout/sprout to v0.5.0 (#944)

## [v1.276.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.276.2) - 2024-08-15

- [`5f0ba95`](https://github.com/alexfalkowski/go-service/commit/5f0ba95292ad99b2b4949d96de28f9397a27f592) fix(deps): bump github.com/prometheus/client_golang from 1.19.1 to 1.20.0 (#940)

## [v1.276.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.276.1) - 2024-08-15

- [`7c93f2f`](https://github.com/alexfalkowski/go-service/commit/7c93f2f48a4737377bcdd8b790cb87652260eef5) fix(deps): bump github.com/open-feature/go-sdk-contrib/hooks/open-telemetry from 0.3.3 to 0.3.4 (#941)
- [`2b5a833`](https://github.com/alexfalkowski/go-service/commit/2b5a8336b9f4183d8cb4c01e474ac3f2b3ff2e63) build(deps): bump bin from `ac397e7` to `41b7c8b` (#942)

## [v1.276.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.276.0) - 2024-08-14

- [`f9a5087`](https://github.com/alexfalkowski/go-service/commit/f9a50870671076a0aca37e28cff4ead26ff53684) feat(go): upload to v1.23 (#939)
- [`7e19b29`](https://github.com/alexfalkowski/go-service/commit/7e19b29b89692bb50579c20fdc70fd3c44736b76) build(deps): bump bin from `808558f` to `ac397e7` (#937)
- [`92d03a4`](https://github.com/alexfalkowski/go-service/commit/92d03a43dff90a218c68b94c7d7eb34f869f41a9) docs(diagrams): update (#935)

## [v1.275.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.275.0) - 2024-08-11

- [`a307639`](https://github.com/alexfalkowski/go-service/commit/a3076396003982575a31ede5fe1642abf683a020) docs(telemetry): add headers (#934)
- [`d1026a6`](https://github.com/alexfalkowski/go-service/commit/d1026a6cb8da268da9c643499993132f0ebc6c68) feat(telemetry): move to use headers (#933)

## [v1.274.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.274.0) - 2024-08-10

- [`af87847`](https://github.com/alexfalkowski/go-service/commit/af878478dac2524fb6fb62b288497f677abc20a3) feat(token): remove security package (#932)
- [`0e72feb`](https://github.com/alexfalkowski/go-service/commit/0e72feb4ad97b3b9e2f4fe456569a9928ad126f2) feat(transport): remove security package (#931)

## [v1.273.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.273.1) - 2024-08-09

- [`b0a8381`](https://github.com/alexfalkowski/go-service/commit/b0a83812355944af51198271401ca786c81cd5e7) fix(security): remove argon package (#930)

## [v1.273.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.273.0) - 2024-08-09

- [`9fe7725`](https://github.com/alexfalkowski/go-service/commit/9fe772568d097531c7c77e6b6bee7c88151afdce) feat(crypto): remove argon token (#929)
- [`303cf60`](https://github.com/alexfalkowski/go-service/commit/303cf6005b2bde0e1eb762a2d5f42156c9c54bb1) build(deps): bump bin from `2d4d510` to `808558f` (#928)

## [v1.272.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.272.2) - 2024-08-08

- [`99b1d1b`](https://github.com/alexfalkowski/go-service/commit/99b1d1bc776d518dbb45011f4a636e32a7458ade) fix(deps): update github.com/alexfalkowski/go-health to v1.18.1 (#927)

## [v1.272.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.272.1) - 2024-08-08

- [`ba39dfe`](https://github.com/alexfalkowski/go-service/commit/ba39dfe3c02312623606f5f026ccfcd259a9f657) fix(go): set to min 1.22 (#926)

## [v1.272.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.272.0) - 2024-08-08

- [`6c585c3`](https://github.com/alexfalkowski/go-service/commit/6c585c3fbcd59f3e582256cfe0f0efb49ed9ab71) feat(mem): limit allocations (#925)
- [`cd389f9`](https://github.com/alexfalkowski/go-service/commit/cd389f967d9ba11a23ebd92b442daebe44355db2) build(deps): bump bin from `cf2f550` to `2d4d510` (#924)
- [`3bbcc27`](https://github.com/alexfalkowski/go-service/commit/3bbcc2746c56a958b8379658fa402e23faa49dd9) build(deps): bump bin from `8724850` to `cf2f550` (#923)

## [v1.271.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.271.2) - 2024-08-08

- [`1837ab8`](https://github.com/alexfalkowski/go-service/commit/1837ab8e9ef9872a20f0f2b3112b8d39b9ecae2d) fix(transport): limit memory usage (#922)
- [`dad124f`](https://github.com/alexfalkowski/go-service/commit/dad124f98689e11a2aba9c96c2fb19db90359493) build(ci): add benchmarks (#921)
- [`7d58744`](https://github.com/alexfalkowski/go-service/commit/7d5874465796cc874b6433a0b74c228f07516440) build(deps): bump bin from `44badfa` to `8724850` (#920)

## [v1.271.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.271.1) - 2024-08-08

- [`7b93e47`](https://github.com/alexfalkowski/go-service/commit/7b93e47af972289b10e6e4804a285a20304e15ce) fix(deps): bump go.uber.org/fx from 1.22.1 to 1.22.2 (#918)

## [v1.271.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.271.0) - 2024-08-07

- [`7f9dcd7`](https://github.com/alexfalkowski/go-service/commit/7f9dcd7516edeb0870e2eaa829c9587421bf6cab) feat(health): upgrade github.com/alexfalkowski/go-health to v1.18.0 (#917)

## [v1.270.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.270.6) - 2024-08-07

- [`b25d985`](https://github.com/alexfalkowski/go-service/commit/b25d9856fc33cbe37ad9f00768407746736f91eb) fix(deps): update github.com/alexfalkowski/go-health to v1.17.3 (#916)

## [v1.270.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.270.5) - 2024-08-07

- [`8843208`](https://github.com/alexfalkowski/go-service/commit/884320840bbca89155846139f638dbb164a2daee) fix(go): update go to v1.22.6 (#915)

## [v1.270.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.270.4) - 2024-08-07

- [`283759d`](https://github.com/alexfalkowski/go-service/commit/283759d199eae8fd65e493c714711741e6741407) fix(deps): bump golang.org/x/net from 0.27.0 to 0.28.0 (#914)

## [v1.270.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.270.3) - 2024-08-06

- [`34dbde8`](https://github.com/alexfalkowski/go-service/commit/34dbde8e8c04da9a8e18b0458d9fdf5ebb512946) fix(logger): zap config to be nillable (#912)

## [v1.270.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.270.2) - 2024-08-06

- [`c22d5e3`](https://github.com/alexfalkowski/go-service/commit/c22d5e3874f8d8342d0386806d118e19b8587292) fix(logger): enable annotaions and stack trace (#911)

## [v1.270.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.270.1) - 2024-08-06

- [`e8f1484`](https://github.com/alexfalkowski/go-service/commit/e8f1484773b37fbdb9a56b499a48054e8d4e4698) fix(http): remove tmp variable (#910)

## [v1.270.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.270.0) - 2024-08-05

- [`25b33b1`](https://github.com/alexfalkowski/go-service/commit/25b33b121d142aa87892aa38029f47d53b69d2ad) feat(grpc): set timout like other transports (#909)

## [v1.269.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.269.0) - 2024-08-04

- [`35a9bfe`](https://github.com/alexfalkowski/go-service/commit/35a9bfeed4717be447cfd6d9039b8b37bf7dfaf8) feat(metrics): do not default to prometheus (#908)

## [v1.268.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.268.0) - 2024-08-04

- [`21d3b09`](https://github.com/alexfalkowski/go-service/commit/21d3b097c9e4ecec9f47d7e73db4fcd51154dfc0) feat(tracer): reader should be nil if disabled (#907)

## [v1.267.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.267.0) - 2024-08-03

- [`218aee8`](https://github.com/alexfalkowski/go-service/commit/218aee862680cbd5dd33402f32864ecbc00912d6) feat(feature): allow optional for feature provider (#905)

## [v1.266.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.266.0) - 2024-08-03

- [`1d3b732`](https://github.com/alexfalkowski/go-service/commit/1d3b73207bff7fb3cdec3ad63d4c384a4976a6dc) feat(net): rename to use address (#904)
- [`a8a7aa5`](https://github.com/alexfalkowski/go-service/commit/a8a7aa5e8c4ad7e8d5ba30be2b4c540b0b941a3f) build(ci): remove flipt (#903)

## [v1.265.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.265.0) - 2024-08-02

- [`b68d585`](https://github.com/alexfalkowski/go-service/commit/b68d5850e0046de17357e56daaef5cceda76aa1b) feat(flags): remove flipt to leave it to use the provider we want (#902)

## [v1.264.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.264.0) - 2024-08-01

- [`0fa91d1`](https://github.com/alexfalkowski/go-service/commit/0fa91d13651184e02b8c0c9e3259fe500b52d11b) feat(http): pass roundtripper to rpc client (#901)

## [v1.263.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.263.0) - 2024-08-01

- [`a9652e5`](https://github.com/alexfalkowski/go-service/commit/a9652e52a7ac233827edb57cf8d669cc0adffbbb) feat(http): provide default timeouts for clients (#900)

## [v1.262.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.262.1) - 2024-08-01

- [`904e790`](https://github.com/alexfalkowski/go-service/commit/904e79011379809ade4562aa99f682d36dc4f9aa) fix(deps): bump github.com/shirou/gopsutil/v4 from 4.24.6 to 4.24.7 (#899)

## [v1.262.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.262.0) - 2024-07-30

- [`25f9d24`](https://github.com/alexfalkowski/go-service/commit/25f9d24e038108f00d8fc89e8eba20b0900f82a6) feat(grpc): ignore the error for headers (#898)

## [v1.261.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.261.0) - 2024-07-30

- [`7338855`](https://github.com/alexfalkowski/go-service/commit/73388559c30ba0213c72d0c1ad5060a537216d4b) feat(http): shorten options for rpc client (#897)

## [v1.260.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.260.2) - 2024-07-29

- [`3be011e`](https://github.com/alexfalkowski/go-service/commit/3be011e45897bd523e37235b0e55eae744edd4fe) fix(http): only create client if needed for rpc (#896)

## [v1.260.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.260.1) - 2024-07-29

- [`9187f8c`](https://github.com/alexfalkowski/go-service/commit/9187f8cc30bcc2fe1680986a52646e8397238f66) fix(http): do not use default client for rpc client (#895)

## [v1.260.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.260.0) - 2024-07-29

- [`b610af8`](https://github.com/alexfalkowski/go-service/commit/b610af8b2c68fd4e338dc1d82e233965068e4f42) feat(http): pass options to rpc client (#893)

## [v1.259.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.259.0) - 2024-07-29

- [`b4869a5`](https://github.com/alexfalkowski/go-service/commit/b4869a551dba6f836d9fdadc0de1813a812a8cf1) feat(http): rpc to follow mvc (#892)

## [v1.258.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.258.0) - 2024-07-28

- [`ac4ee0b`](https://github.com/alexfalkowski/go-service/commit/ac4ee0b655cd06d4fbabf2485224dd74debf0fba) feat(http): unify errors for rpc (#891)

## [v1.257.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.257.0) - 2024-07-28

- [`f4b1838`](https://github.com/alexfalkowski/go-service/commit/f4b183800392acb6f48e183a2b473b3ee650d413) feat(http): parse patterns once for mvc (#889)

## [v1.256.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.256.1) - 2024-07-27

- [`56f08f7`](https://github.com/alexfalkowski/go-service/commit/56f08f71ce69a971874c540adf9f0bba56fe69ce) fix(grpc): set internal error for limiter (#888)

## [v1.256.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.256.0) - 2024-07-27

- [`6a7101b`](https://github.com/alexfalkowski/go-service/commit/6a7101bc6bffd97e06de0483c7693e75beb8dd3b) feat(transport): remove duplication in limiter (#887)

## [v1.255.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.255.1) - 2024-07-27

- [`290305c`](https://github.com/alexfalkowski/go-service/commit/290305c22ddde78c9126e78f98b3c060117b864f) fix(transport): calculate duration for rate limit reset (#886)

## [v1.255.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.255.0) - 2024-07-26

- [`27ee5c6`](https://github.com/alexfalkowski/go-service/commit/27ee5c611deb15bf670f98133f1159caf1ad1989) feat(http): add github.com/go-sprout/sprout v0.4.1 (#885)

## [v1.254.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.254.0) - 2024-07-26

- [`1d1ee56`](https://github.com/alexfalkowski/go-service/commit/1d1ee56b8a913965705659326802205c615d4087) feat(http): remove result from mvc (#884)

## [v1.253.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.253.0) - 2024-07-26

- [`97c34ec`](https://github.com/alexfalkowski/go-service/commit/97c34ec79d0f699b1ae6d1c87f02825b15ecdbd3) feat(http): add view function for mvc (#883)

## [v1.252.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.252.0) - 2024-07-26

- [`8bcdcd4`](https://github.com/alexfalkowski/go-service/commit/8bcdcd42fb6036d0b44e0158ae23b97f37accadc) feat(http): record the error (#882)

## [v1.251.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.251.1) - 2024-07-26

- [`5157d68`](https://github.com/alexfalkowski/go-service/commit/5157d6819ee63c6ab0ae9787d3432d6f6b9b6d3f) fix(deps): bump github.com/redis/go-redis/v9 from 9.6.0 to 9.6.1 (#881)

## [v1.251.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.251.0) - 2024-07-23

- [`4438bd5`](https://github.com/alexfalkowski/go-service/commit/4438bd50f576f47a69b87511d53a33a028e6f75b) feat(http): use context in mvc like in rpc (#880)

## [v1.250.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.250.0) - 2024-07-23

- [`0c62eef`](https://github.com/alexfalkowski/go-service/commit/0c62eefaa045a95b496dc0baf3393961472b6c80) feat(http): controller returns a result for mvc (#879)

## [v1.249.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.249.0) - 2024-07-20

- [`bb4d359`](https://github.com/alexfalkowski/go-service/commit/bb4d35965922fcbc4bd0596bb3e583b1f001ff6b) feat(deps): go: upgraded github.com/redis/go-redis/v9 v9.5.4 => v9.6.0 (#878)

## [v1.248.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.248.1) - 2024-07-15

- [`0c0a165`](https://github.com/alexfalkowski/go-service/commit/0c0a16509177bb45fbabd8ff774aac11713849f1) fix(deps): bump github.com/redis/go-redis/v9 from 9.5.3 to 9.5.4 (#877)
- [`8fa6056`](https://github.com/alexfalkowski/go-service/commit/8fa60567222e41454943406312c14f0c669b58f2) build(deps): bump bin from `897d5f1` to `44badfa` (#876)

## [v1.248.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.248.0) - 2024-07-12

- [`0a96e61`](https://github.com/alexfalkowski/go-service/commit/0a96e61216cbaf3a2823c2f96879fece1961d220) feat(http): new success view (#875)

## [v1.247.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.247.0) - 2024-07-10

- [`e09feba`](https://github.com/alexfalkowski/go-service/commit/e09febaa787ec9f7fc8ee8c3cedfe1ef350d7d22) feat(http): add ability for no controller for mvc (#874)

## [v1.246.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.246.0) - 2024-07-09

- [`c1fb821`](https://github.com/alexfalkowski/go-service/commit/c1fb821837be95f116e2635915e72684ec157dde) feat(http): add mvc (#873)

## [v1.245.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.245.0) - 2024-07-08

- [`8d7e70d`](https://github.com/alexfalkowski/go-service/commit/8d7e70de87a129b08124bd9ae39fd55b1ce94da6) feat(http): move to a func (#871)

## [v1.244.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.244.1) - 2024-07-07

- [`60dddb3`](https://github.com/alexfalkowski/go-service/commit/60dddb38f1b92edb909cf3af9698caf7f37adc6e) fix(http): use keys (#870)

## [v1.244.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.244.0) - 2024-07-07

- [`6e738c8`](https://github.com/alexfalkowski/go-service/commit/6e738c80fb1367018238d6575134d3d24355808a) feat(http): use context from the context package (#869)

## [v1.243.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.243.0) - 2024-07-06

- [`1884e54`](https://github.com/alexfalkowski/go-service/commit/1884e544edfe88aec2846f8ff68cc74f471c4b9a) feat(encoding): move to the encoding package (#868)

## [v1.242.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.242.0) - 2024-07-06

- [`e3d03fa`](https://github.com/alexfalkowski/go-service/commit/e3d03fae5a75be774772b2e0d2bd5969e0b9747f) feat(compress): move to the compress package (#867)

## [v1.241.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.241.2) - 2024-07-06

- [`0ab7a86`](https://github.com/alexfalkowski/go-service/commit/0ab7a8605efb52b51e6f6afa1601bcdddecfe340) fix(deps): bump golang.org/x/net from 0.26.0 to 0.27.0 (#866)

## [v1.241.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.241.1) - 2024-07-06

- [`a521ba7`](https://github.com/alexfalkowski/go-service/commit/a521ba703121ad758aeafdff83e69b1b20ab05e9) fix(deps): bump golang.org/x/crypto from 0.24.0 to 0.25.0 (#865)

## [v1.241.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.241.0) - 2024-07-05

- [`ae3e8f1`](https://github.com/alexfalkowski/go-service/commit/ae3e8f1a62470927af385932396d4fe30cb924ad) feat(http): move back to rpc (#864)
- [`ef504df`](https://github.com/alexfalkowski/go-service/commit/ef504dface4f6c871a97e4394407406d60ffa450) test(http): move rpc (#863)

## [v1.240.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.240.0) - 2024-07-05

- [`07f409b`](https://github.com/alexfalkowski/go-service/commit/07f409b9393dc779c956a37354de9eee387fcb07) feat(http): move back to http package (#861)

## [v1.239.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.239.1) - 2024-07-05

- [`0a39a0f`](https://github.com/alexfalkowski/go-service/commit/0a39a0fd08f7084ecb8b33a22b508761f569f139) fix(http): do not panic (#860)

## [v1.239.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.239.0) - 2024-07-05

- [`88ab19c`](https://github.com/alexfalkowski/go-service/commit/88ab19cb2373eef68694f23f05424bd7a34041f0) feat(http): move to content package (#859)

## [v1.238.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.238.0) - 2024-07-05

- [`e84b65c`](https://github.com/alexfalkowski/go-service/commit/e84b65ccd922af68247e3619bb291fc4066b2c16) feat(http): move to rpc package (#858)

## [v1.237.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.237.0) - 2024-07-05

- [`0e6ee88`](https://github.com/alexfalkowski/go-service/commit/0e6ee88b7098d74540383cb81614d83c536a1f62) feat(http): add h2c support (#857)

## [v1.236.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.236.0) - 2024-07-04

- [`4f5c3e7`](https://github.com/alexfalkowski/go-service/commit/4f5c3e76b958b873c621529b013eada92cb8ceeb) feat(metrics): add instrumentation (#855)

## [v1.235.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.235.0) - 2024-07-03

- [`f8b7c15`](https://github.com/alexfalkowski/go-service/commit/f8b7c1515129c639447a91563f004c7de03d7caf) feat(http): move to url constructor for client (#854)

## [v1.234.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.234.0) - 2024-07-03

- [`ef08319`](https://github.com/alexfalkowski/go-service/commit/ef083198aa42d302e4851ac092cd70487ef2a2f1) feat(telemetry): upgrade go.opentelemetry.io/otel to v1.28.0 (#853)

## [v1.233.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.233.0) - 2024-07-03

- [`7f4dbdc`](https://github.com/alexfalkowski/go-service/commit/7f4dbdc35534cada4316073dcc812bd0ae7a9bae) feat(deps): update google.golang.org/grpc to v1.65.0 (#852)

## [v1.232.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.232.1) - 2024-07-03

- [`d374c96`](https://github.com/alexfalkowski/go-service/commit/d374c96c8336d01e3f8f4f6d4ad3ab945da7a173) fix(go): update go to v1.22.5 (#851)

## [v1.232.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.232.0) - 2024-07-02

- [`94c9695`](https://github.com/alexfalkowski/go-service/commit/94c9695870b4cdeda84f1714444dba19dea97c25) feat(http): return code error (#845)

## [v1.231.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.231.0) - 2024-07-02

- [`991ba9f`](https://github.com/alexfalkowski/go-service/commit/991ba9f7678e70edc03ccbbe3c84ff3b8bc608e2) feat(http): use a custom error to mark the status (#844)

## [v1.230.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.230.1) - 2024-07-01

- [`dad7d55`](https://github.com/alexfalkowski/go-service/commit/dad7d5554568cbebcf1123761e07f12178c4581f) fix(transport): add meta to response (#843)

## [v1.230.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.230.0) - 2024-07-01

- [`20ec822`](https://github.com/alexfalkowski/go-service/commit/20ec822b4373df2d291d97154916ca1615b96797) feat(http): write errors as plain (#842)

## [v1.229.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.229.0) - 2024-07-01

- [`d7189cd`](https://github.com/alexfalkowski/go-service/commit/d7189cdf832b0878c862b7c8d50fa610b9bce7fe) feat(http): add client for request/response (#841)

## [v1.228.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.228.0) - 2024-07-01

- [`43ce37d`](https://github.com/alexfalkowski/go-service/commit/43ce37d62bddcd48bcfe59c671bd9c2a85039c1b) feat(limiter): change header (#840)

## [v1.227.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.227.0) - 2024-06-30

- [`014a39a`](https://github.com/alexfalkowski/go-service/commit/014a39a9e8393ef51a936efec3da98c9a552f28b) feat(limiter): move to use github.com/sethvargo/go-limiter (#839)
- [`d2568ba`](https://github.com/alexfalkowski/go-service/commit/d2568baae17f14ca79d83fd60cf38990824ed27f) docs(readme): add badges (#838)
- [`9badf23`](https://github.com/alexfalkowski/go-service/commit/9badf23ef30c3806441b1df31dc4a5e7a24ac895) build(ci): add codecov (#837)
- [`d0ee914`](https://github.com/alexfalkowski/go-service/commit/d0ee914ccbf48d4cee9ec5638940a43db7eb0ea8) docs(transport): diagrams (#836)

## [v1.226.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.226.0) - 2024-06-27

- [`34b1ec5`](https://github.com/alexfalkowski/go-service/commit/34b1ec570a03b6803c30fad2d02dc503780c1688) feat(transport): add version (#835)

## [v1.225.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.225.0) - 2024-06-27

- [`f8f899a`](https://github.com/alexfalkowski/go-service/commit/f8f899a7ccb18a0c8cca39cf47a53aae632dc608) feat(http): add specific context for handler (#834)

## [v1.224.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.224.0) - 2024-06-27

- [`03c5000`](https://github.com/alexfalkowski/go-service/commit/03c5000899b55a1c130b833873bc55a72e864439) feat(transport): add gzip compression (#833)

## [v1.223.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.223.0) - 2024-06-26

- [`18edfb3`](https://github.com/alexfalkowski/go-service/commit/18edfb3ab181c542f55bf105ae0ffab63a481f80) feat(http): set to post (#832)

## [v1.222.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.222.4) - 2024-06-26

- [`2ecefb4`](https://github.com/alexfalkowski/go-service/commit/2ecefb42ba6411c9075a056231b78df462a7a068) fix(http): remove meta as we return the error (#831)
- [`6821521`](https://github.com/alexfalkowski/go-service/commit/68215215fbc7412bcb3d5185a7a230a25abe7e9e) docs(http): remove rpc gateway (#830)

## [v1.222.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.222.3) - 2024-06-26

- [`85c14c3`](https://github.com/alexfalkowski/go-service/commit/85c14c3365cb6fe3d7927687540d11324cfb7aad) fix(deps): upgrade go.uber.org/fx to v1.22.1 (#829)

## [v1.222.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.222.2) - 2024-06-25

- [`09f8958`](https://github.com/alexfalkowski/go-service/commit/09f895804d305cf244512429e9fe554b46c76292) fix(http): invalid content type (#827)

## [v1.222.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.222.1) - 2024-06-25

- [`f7c983e`](https://github.com/alexfalkowski/go-service/commit/f7c983e4e8086be81b264fc939cec4f474b98e61) fix(http): write error to context (#826)

## [v1.222.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.222.0) - 2024-06-25

- [`e6269c8`](https://github.com/alexfalkowski/go-service/commit/e6269c8015922ca96f162c87a6bd9d72e41809ad) feat(http): move to a handler interface (#825)

## [v1.221.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.221.0) - 2024-06-25

- [`eb41d96`](https://github.com/alexfalkowski/go-service/commit/eb41d964e20a7c34757975c2243f984315a0da93) feat(http): simplify handler (#824)

## [v1.220.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.220.0) - 2024-06-24

- [`1e108e4`](https://github.com/alexfalkowski/go-service/commit/1e108e49998e1e15b02e6c3b53e06df14633c43e) feat(http): remove grpc gateway (#821)

## [v1.219.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.219.1) - 2024-06-24

- [`0b9f357`](https://github.com/alexfalkowski/go-service/commit/0b9f35780e3fa88000134cf3f9034eca84ea8e22) fix(http): parse remote addr (#823)

## [v1.219.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.219.0) - 2024-06-24

- [`17b6d14`](https://github.com/alexfalkowski/go-service/commit/17b6d14c634b7dff8fdff37055f2d993dfda3b65) feat(meta): add geolocation (#822)

## [v1.218.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.218.3) - 2024-06-23

- [`8e25222`](https://github.com/alexfalkowski/go-service/commit/8e25222d36900bbeda2fcca6bd84c7f7714a9af7) fix(http): wrap errors (#819)

## [v1.218.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.218.2) - 2024-06-23

- [`60457c1`](https://github.com/alexfalkowski/go-service/commit/60457c1c421960d898e788eeed5a6675eadda01b) fix(http): use stable errors and add the error to meta (#818)

## [v1.218.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.218.1) - 2024-06-22

- [`adfb442`](https://github.com/alexfalkowski/go-service/commit/adfb4421e49f4b6dbd27e7e18c38c2ddefe0dbf1) fix(http): pass context (#817)

## [v1.218.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.218.0) - 2024-06-22

- [`8ed8417`](https://github.com/alexfalkowski/go-service/commit/8ed84173a47d2f866b0a8644557fef5b156941f2) feat(http): add sync calls (#816)

## [v1.217.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.217.3) - 2024-06-17

- [`0d918b4`](https://github.com/alexfalkowski/go-service/commit/0d918b4fb46d12273db4a8aa55eb3f45d74e0941) fix(deps): bump github.com/spf13/cobra from 1.8.0 to 1.8.1 (#815)

## [v1.217.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.217.2) - 2024-06-13

- [`ef9961e`](https://github.com/alexfalkowski/go-service/commit/ef9961e27dc01c1476f72d97b7c0d0267855d760) fix(deps): bump github.com/klauspost/compress from 1.17.8 to 1.17.9 (#814)

## [v1.217.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.217.1) - 2024-06-11

- [`c2adba8`](https://github.com/alexfalkowski/go-service/commit/c2adba84db8c68dc1d75756a7fb45c35fd6eef98) fix(deps): bump google.golang.org/protobuf from 1.34.1 to 1.34.2 (#813)

## [v1.217.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.217.0) - 2024-06-11

- [`1a0b854`](https://github.com/alexfalkowski/go-service/commit/1a0b8543fb0fd240ca70ba0f2600142b87782b74) feat(token): move token to it's own configuration (#812)

## [v1.216.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.216.0) - 2024-06-10

- [`904d787`](https://github.com/alexfalkowski/go-service/commit/904d787bd9ef8995766507fcb0aeb364cd0551cd) feat(crypto): use string for paths (#811)

## [v1.215.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.215.0) - 2024-06-10

- [`eeffb98`](https://github.com/alexfalkowski/go-service/commit/eeffb989fc153d66ae020b2886af593594e10e8e) feat(crypto): ssh should use signature (#810)
- [`0d3a349`](https://github.com/alexfalkowski/go-service/commit/0d3a349ea10ec62d6802a8f4103b02c1c4218e6a) build(deps): bump bin from `955891d` to `897d5f1` (#809)

## [v1.214.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.214.1) - 2024-06-09

- [`e709f28`](https://github.com/alexfalkowski/go-service/commit/e709f2862ac952f557e730513bc89c88f1fe4c32) fix(crypto): move to algo interfaces (#808)

## [v1.214.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.214.0) - 2024-06-08

- [`5bd68d7`](https://github.com/alexfalkowski/go-service/commit/5bd68d7ec20cd065fbfb74a55a601c64b579b2e0) feat(crypto): add ssh (#807)
- [`0ddbf60`](https://github.com/alexfalkowski/go-service/commit/0ddbf609decbc48d9cbd62c5100381d0395455b7) build(ci): use alexfalkowski/go:1.22 (#806)

## [v1.213.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.213.3) - 2024-06-05

- [`80539ad`](https://github.com/alexfalkowski/go-service/commit/80539adb6ea2b30271141870ef151c70f7ff4404) fix(deps): upgraded github.com/alexfalkowski/go-health v1.17.1 (#805)

## [v1.213.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.213.2) - 2024-06-05

- [`876df31`](https://github.com/alexfalkowski/go-service/commit/876df316ab0cee8a1cf1da31e6d9be4f617cc706) fix(deps): bump github.com/urfave/negroni/v3 from 3.1.0 to 3.1.1 (#803)

## [v1.213.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.213.1) - 2024-06-05

- [`a78bfb9`](https://github.com/alexfalkowski/go-service/commit/a78bfb99247e364f0eb097447277f3cdeb0f521b) fix(deps): upgrade to go 1.22.4 (#804)

## [v1.213.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.213.0) - 2024-06-03

- [`23afeb9`](https://github.com/alexfalkowski/go-service/commit/23afeb98ab8affe7d5809cf32f7b6b6145c345dc) feat(token): generate (#802)

## [v1.212.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.212.2) - 2024-06-03

- [`95d3b5f`](https://github.com/alexfalkowski/go-service/commit/95d3b5f1c90ae6dc6ad0919d78a9f7a19ff75769) fix(deps): bump github.com/redis/go-redis/v9 from 9.5.1 to 9.5.2 (#801)

## [v1.212.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.212.1) - 2024-06-03

- [`00ba836`](https://github.com/alexfalkowski/go-service/commit/00ba836ae4510c3ec0a986f064e22489428be602) fix(deps): bump github.com/shirou/gopsutil/v3 from 3.24.4 to 3.24.5 (#800)

## [v1.212.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.212.0) - 2024-06-02

- [`456eb63`](https://github.com/alexfalkowski/go-service/commit/456eb6328f879c629d0552fd24b190801d340f6c) feat(crypto): use pem blocks for public/private keys (#799)

## [v1.211.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.211.0) - 2024-06-01

- [`0caf2d0`](https://github.com/alexfalkowski/go-service/commit/0caf2d0795afc741df5d1f31d633b2b58f72a2e5) feat(grpc): record ip addr kind (#798)

## [v1.210.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.210.0) - 2024-06-01

- [`8e7fc48`](https://github.com/alexfalkowski/go-service/commit/8e7fc4847934d8e9095865b067aea998cd7d2df0) feat(grpc): add ip keys (#797)

## [v1.209.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.209.2) - 2024-05-31

- [`3968281`](https://github.com/alexfalkowski/go-service/commit/39682815e837c698118abedabb0581ecb1b9bee1) fix(deps): bump go.uber.org/fx from 1.21.1 to 1.22.0 (#795)

## [v1.209.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.209.1) - 2024-05-30

- [`258bdf1`](https://github.com/alexfalkowski/go-service/commit/258bdf174ee0057f891b012ed66d92af10790c70) fix(deps): bump github.com/beevik/ntp from 1.4.2 to 1.4.3 (#794)

## [v1.209.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.209.0) - 2024-05-30

- [`03c5854`](https://github.com/alexfalkowski/go-service/commit/03c5854b30bc1d9b56c9067c6789d6ee9d4729d9) feat(cmd): allow to specify flags on individual commands (#793)

## [v1.208.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.208.2) - 2024-05-30

- [`60cb850`](https://github.com/alexfalkowski/go-service/commit/60cb850d9af014be67a80e82b34e34faf3681450) fix(deps): bump github.com/hashicorp/go-retryablehttp from 0.7.6 to 0.7.7 (#790)

## [v1.208.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.208.1) - 2024-05-30

- [`1988db4`](https://github.com/alexfalkowski/go-service/commit/1988db4951669e1a06a329149fc37b5d9c745b02) fix(telemetry): ignore stop error as there is little we can do about it (#792)
- [`79fa39e`](https://github.com/alexfalkowski/go-service/commit/79fa39ee00eee9f8aa7555ad24fe0b55c23d001e) fix(deps): bump github.com/open-feature/go-sdk from 1.11.0 to 1.12.0 (#791)

## [v1.208.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.208.0) - 2024-05-29

- [`bd0448f`](https://github.com/alexfalkowski/go-service/commit/bd0448f5beff43e131a9701fb4e72b69d12fb4f4) feat(env): strip v from user agent (#789)

## [v1.207.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.207.0) - 2024-05-29

- [`cf37a11`](https://github.com/alexfalkowski/go-service/commit/cf37a11268809a8a6adc51a1916c166cb2e8ae2c) feat(env): add name to env (#788)
- [`38ceeea`](https://github.com/alexfalkowski/go-service/commit/38ceeea87cfe8612e04aba3bc97ff710fe3dc9b1) build(make): remove newline (#787)
- [`4a97ec0`](https://github.com/alexfalkowski/go-service/commit/4a97ec0f0b00020207092d78cbcf5c13d08aace6) build(diagrams): remove diagram target (#786)
- [`033995c`](https://github.com/alexfalkowski/go-service/commit/033995c00ca238c72460f419357757c4d25d6bc5) build(deps): bump bin from `bf008ad` to `955891d` (#785)
- [`16aefa6`](https://github.com/alexfalkowski/go-service/commit/16aefa61848733304803e309a641fccaa3c07415) build(deps): bump bin from `93c02b7` to `bf008ad` (#784)
- [`7b96001`](https://github.com/alexfalkowski/go-service/commit/7b96001a7ed9cedc346954e9e02c21525059e825) build(diagrams): add dependencies (#783)
- [`a286e8b`](https://github.com/alexfalkowski/go-service/commit/a286e8b93b0f2ddcc7c1fdb5bfbacbb838aeeafd) test(transport): use missing options (#782)
- [`d08e76b`](https://github.com/alexfalkowski/go-service/commit/d08e76bfd3d68573f8079259aad8d4d57ed8bce3) test(crypto): break up large tests (#781)
- [`0d77a44`](https://github.com/alexfalkowski/go-service/commit/0d77a448a6c275e60894aa304a7ab9725292fdd5) test(crypto): add tests (#780)
- [`ffe53f0`](https://github.com/alexfalkowski/go-service/commit/ffe53f05c1708343f41c0e06c74789be86ab7459) test(http): add timeout (#779)

## [v1.206.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.206.2) - 2024-05-27

- [`7c9d536`](https://github.com/alexfalkowski/go-service/commit/7c9d536a01b5f025e878bd20198d10dd4d07023c) fix(health): gRPC observer should be optional (#778)

## [v1.206.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.206.1) - 2024-05-27

- [`740d191`](https://github.com/alexfalkowski/go-service/commit/740d1913abd4735be21a53c7f8e66d0be99ee09a) fix(deps): bump github.com/jackc/pgx/v5 from 5.5.5 to 5.6.0 (#777)

## [v1.206.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.206.0) - 2024-05-26

- [`ef67724`](https://github.com/alexfalkowski/go-service/commit/ef67724bd9c7de7819f0b56ff4b3d3869a9c5560) feat(meta): add ability to ignore meta (#776)

## [v1.205.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.205.0) - 2024-05-26

- [`853c9c9`](https://github.com/alexfalkowski/go-service/commit/853c9c974e9c96cfad39a9ab293b6b5d54085187) feat(security): add authorization to meta (#775)

## [v1.204.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.204.2) - 2024-05-26

- [`9dd92bb`](https://github.com/alexfalkowski/go-service/commit/9dd92bb6d0c34a3c808e1c8872e4ecc10a6c320d) fix(grpc): handle metadata prefix for authorization (#774)

## [v1.204.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.204.1) - 2024-05-26

- [`4592f52`](https://github.com/alexfalkowski/go-service/commit/4592f52875a20b3898ad9aa9de0ed28a429541b6) fix(limiter): do not rate limit health checks (#773)

## [v1.204.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.204.0) - 2024-05-26

- [`8c09250`](https://github.com/alexfalkowski/go-service/commit/8c092507bdb0ded6266a0191c1e0e0a07167a9d6) feat(security): add options for tokens (#772)

## [v1.203.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.203.0) - 2024-05-25

- [`0aed3dc`](https://github.com/alexfalkowski/go-service/commit/0aed3dcc4f0720980da0db88ecaa69a4c15c2f9a) feat(transport): specify default ports (#771)

## [v1.202.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.202.3) - 2024-05-24

- [`8bd0563`](https://github.com/alexfalkowski/go-service/commit/8bd05637f7e7b6c130d973b30c5393c3853be0b3) fix(deps): bump github.com/BurntSushi/toml from 1.3.2 to 1.4.0 (#770)

## [v1.202.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.202.2) - 2024-05-24

- [`76f773e`](https://github.com/alexfalkowski/go-service/commit/76f773e6fb0630641696c900aa29e88d1990e314) fix(deps): bump github.com/beevik/ntp from 1.4.1 to 1.4.2 (#769)

## [v1.202.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.202.1) - 2024-05-23

- [`f0580de`](https://github.com/alexfalkowski/go-service/commit/f0580deaed98667bdb05529a42dafab726653e55) fix(http): rename to have a more descriptive name for kind of mux (#768)

## [v1.202.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.202.0) - 2024-05-23

- [`e8e3264`](https://github.com/alexfalkowski/go-service/commit/e8e3264fbe7891a0e0b7c87b98c2d241ae2315be) feat(http): add ability to define the mux (#766)

## [v1.201.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.201.0) - 2024-05-22

- [`d8ca8b1`](https://github.com/alexfalkowski/go-service/commit/d8ca8b151bf7ef520c0185f7dbb0a6c9c3a4da23) feat(config): make sure we pass timeouts (#765)

## [v1.200.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.200.1) - 2024-05-22

- [`09bbc17`](https://github.com/alexfalkowski/go-service/commit/09bbc1753944043190d1adbe62a52f72da206da3) fix(marshaller): upgraded github.com/goccy/go-json to v0.10.3 (#763)

## [v1.200.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.200.0) - 2024-05-22

- [`d894abc`](https://github.com/alexfalkowski/go-service/commit/d894abc48ad9a32cc3188e060cf8f20d63536b74) feat(telemetry): upgraded go.opentelemetry.io/otel to v1.27.0 (#762)

## [v1.199.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.199.0) - 2024-05-22

- [`832f0c3`](https://github.com/alexfalkowski/go-service/commit/832f0c3b265fa55236b02adff3d933068d005a1b) feat(grpc): upgrade  github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.1.0 (#761)

## [v1.198.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.198.0) - 2024-05-21

- [`8b055bd`](https://github.com/alexfalkowski/go-service/commit/8b055bdaedb93aa2a06de3894ea59bb8c6fa5430) feat(client): add run for clients that run on start (#755)

## [v1.197.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.197.0) - 2024-05-21

- [`072f27f`](https://github.com/alexfalkowski/go-service/commit/072f27fe6aa98fdd1ae1876e101c8c02f53b684a) feat(deps): upgraded github.com/redis/go-redis/v9 to v9.5.1 (#754)

## [v1.196.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.196.0) - 2024-05-20

- [`619a85a`](https://github.com/alexfalkowski/go-service/commit/619a85a229b393d6c5b945e7b071938282cebc86) feat(redis): upgraded to github.com/redis/go-redis/v9 (#752)

## [v1.195.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.195.0) - 2024-05-20

- [`7d47784`](https://github.com/alexfalkowski/go-service/commit/7d47784c9d7a34216cfc431a7690831ede0c92f1) feat(os): add read file (#751)

## [v1.194.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.194.0) - 2024-05-20

- [`66faf8f`](https://github.com/alexfalkowski/go-service/commit/66faf8f94602d6f71e164c632d8de8d1165b539e) feat(secrets): make sure we load them from file system (#750)

## [v1.193.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.193.1) - 2024-05-20

- [`9ba6211`](https://github.com/alexfalkowski/go-service/commit/9ba6211d8f6cf1260b572e7c03209db313c919bb) fix(deps): bump github.com/KimMachineGun/automemlimit from 0.6.0 to 0.6.1 (#749)

## [v1.193.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.193.0) - 2024-05-19

- [`30c8755`](https://github.com/alexfalkowski/go-service/commit/30c8755459f067a9e5c266ebd12c39cf969a2e13) feat(secrets): make sure we load them from file system (#748)

## [v1.192.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.192.0) - 2024-05-19

- [`8f86afa`](https://github.com/alexfalkowski/go-service/commit/8f86afae042af997d3f14a46a2c2a4432529e290) feat(hooks): generate secret (#747)

## [v1.191.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.191.3) - 2024-05-19

- [`626c043`](https://github.com/alexfalkowski/go-service/commit/626c043ece16407b4f285d30ef9f7f1b2c39ac4b) fix(tracer): convert status code to string for http (#746)

## [v1.191.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.191.2) - 2024-05-19

- [`e9e4573`](https://github.com/alexfalkowski/go-service/commit/e9e4573902eefe8e34d036b7993b0320d2dc33d1) fix(feature): make sure we shutdown the provider (#745)

## [v1.191.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.191.1) - 2024-05-17

- [`24dd3c3`](https://github.com/alexfalkowski/go-service/commit/24dd3c3f58e786c0ee1de56e9a6aac20817a0352) fix(metrics): enable endpoint only for prometheus (#744)

## [v1.191.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.191.0) - 2024-05-17

- [`8dadb5d`](https://github.com/alexfalkowski/go-service/commit/8dadb5d660dc5292143d2f993af4d3d246899b88) feat(telemetry): pass key (#743)

## [v1.190.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.190.0) - 2024-05-16

- [`3b6b72e`](https://github.com/alexfalkowski/go-service/commit/3b6b72eb36e7e8adce5e4302dfc272742e1943e4) feat(tracer): check if trace is recording (#742)

## [v1.189.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.189.4) - 2024-05-16

- [`716665a`](https://github.com/alexfalkowski/go-service/commit/716665a1dd0f4c355277fcf0155dff99f4c56f5f) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 from 2.19.1 to 2.20.0 (#741)

## [v1.189.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.189.3) - 2024-05-15

- [`a63a4a9`](https://github.com/alexfalkowski/go-service/commit/a63a4a9f3a540da5a330abe357818625394ab544) fix(deps): bump google.golang.org/grpc from 1.63.2 to 1.64.0 (#740)

## [v1.189.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.189.2) - 2024-05-13

- [`1dc3c2a`](https://github.com/alexfalkowski/go-service/commit/1dc3c2af7b34be7a3d8bf372a596d7b267d7d8bd) fix(cmd): add more descriptive error (#739)

## [v1.189.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.189.1) - 2024-05-13

- [`ea888c0`](https://github.com/alexfalkowski/go-service/commit/ea888c0a4f85d8acbbb494f19605d50dd9309bad) fix(cmd): add more descriptive error (#738)
- [`2dd2bf5`](https://github.com/alexfalkowski/go-service/commit/2dd2bf52cf63c7d332f9ad4eb031039bc1563e3b) build(deps): bump bin from `9e5c4b7` to `93c02b7` (#737)

## [v1.189.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.189.0) - 2024-05-12

- [`167a9be`](https://github.com/alexfalkowski/go-service/commit/167a9be4654bca8be00543fcae0d2d55c50a9231) feat(crypto): use types to avoid issues (#736)

## [v1.188.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.188.4) - 2024-05-11

- [`5d62243`](https://github.com/alexfalkowski/go-service/commit/5d6224378f3e593cf7878e67cd3cf77d979abe78) fix(config): error on none (#735)

## [v1.188.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.188.3) - 2024-05-11

- [`f3ba9d4`](https://github.com/alexfalkowski/go-service/commit/f3ba9d4de65922892fc6b4a59a597f0ab5cea1db) fix(config): remove default (#734)

## [v1.188.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.188.2) - 2024-05-11

- [`8fa4566`](https://github.com/alexfalkowski/go-service/commit/8fa4566f53c0a04bc6d0a8005a57cfa2975706b7) fix(config): rename to mape from factory (#733)

## [v1.188.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.188.1) - 2024-05-11

- [`0d4b6f0`](https://github.com/alexfalkowski/go-service/commit/0d4b6f0d26a03889655399b391c52a57bb715a62) fix(config): handle empty (#732)

## [v1.188.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.188.0) - 2024-05-11

- [`67fde96`](https://github.com/alexfalkowski/go-service/commit/67fde96a84236a9989a25d9ca25252c877158515) feat(config): add none (#731)

## [v1.187.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.187.0) - 2024-05-11

- [`aa4b0d3`](https://github.com/alexfalkowski/go-service/commit/aa4b0d312b770082377073e5852cb5e4e94bd3fc) feat(cmd): allow overriding of input and output (#730)

## [v1.186.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.186.1) - 2024-05-11

- [`fd9ceb8`](https://github.com/alexfalkowski/go-service/commit/fd9ceb819e749996d3b8cb8184b70885a296350d) fix(cmd): wrap error for out config (#729)

## [v1.186.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.186.0) - 2024-05-11

- [`57056b6`](https://github.com/alexfalkowski/go-service/commit/57056b68482dee7c05235eed78f11477581a98ad) feat(cmd): add name (#728)
- [`5e4f847`](https://github.com/alexfalkowski/go-service/commit/5e4f8471f3d320ab02dede94b63b90efae0a8dfa) test(config): use output (#727)

## [v1.185.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.185.0) - 2024-05-11

- [`f087e4f`](https://github.com/alexfalkowski/go-service/commit/f087e4f7d7aa50f1219434d0d72a55f11f567a36) feat(flags): add checking if set (#726)

## [v1.184.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.184.0) - 2024-05-11

- [`2bef99d`](https://github.com/alexfalkowski/go-service/commit/2bef99d141cb7532acc02389c11483c751d857d9) feat(flags): move to package (#725)

## [v1.183.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.183.0) - 2024-05-10

- [`76b1d61`](https://github.com/alexfalkowski/go-service/commit/76b1d614abcf26b9c1921c82e50cc15adf018624) feat(deps): upgraded github.com/alexfalkowski/go-health to v1.17.0 (#724)

## [v1.182.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.182.1) - 2024-05-10

- [`645a1c9`](https://github.com/alexfalkowski/go-service/commit/645a1c9faa514970cfbf1847b72321b060c52df4) fix(deps): upgraded github.com/hashicorp/go-retryablehttp to v0.7.6 (#721)
- [`548f7f0`](https://github.com/alexfalkowski/go-service/commit/548f7f0c195648aab287a9f1d1d6e2c8cae4d08e) build(deps): bump bin from `55a2500` to `9e5c4b7` (#722)

## [v1.182.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.182.0) - 2024-05-10

- [`b84ac88`](https://github.com/alexfalkowski/go-service/commit/b84ac8893335ab3d9e52f64fb1f4ee472407defc) feat(field-alignment): run tool (#720)
- [`62e649a`](https://github.com/alexfalkowski/go-service/commit/62e649ae1fccd61fa9a829bbe6d56a78803a663d) test(cmd): verify var (#719)

## [v1.181.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.181.0) - 2024-05-09

- [`293932a`](https://github.com/alexfalkowski/go-service/commit/293932a68689b1e279af4ddacfc99eb3948ec9fe) feat(cmd): add output config (#718)

## [v1.180.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.180.1) - 2024-05-09

- [`56f45c1`](https://github.com/alexfalkowski/go-service/commit/56f45c18a6c60498ce284997f56354b197172765) fix(deps): upgraded github.com/prometheus/client_golang to v1.19.1 (#717)
- [`e38815b`](https://github.com/alexfalkowski/go-service/commit/e38815b8943bdf9b5c95709337a2de2f0649733e) build(deps): bump bin from `693b345` to `55a2500` (#716)

## [v1.180.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.180.0) - 2024-05-09

- [`d9de211`](https://github.com/alexfalkowski/go-service/commit/d9de21123b7230c5649379a43130e83f51ae6ae2) feat(crypto): add ability to handle symmetric and asymmetric keys (#715)

## [v1.179.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.179.4) - 2024-05-08

- [`07de0d9`](https://github.com/alexfalkowski/go-service/commit/07de0d983877090da0e7efa4bc48e1b4ca61e7cf) fix(deps): upgraded github.com/alexfalkowski/go-health to v1.16.2 (#714)

## [v1.179.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.179.3) - 2024-05-08

- [`61b68bb`](https://github.com/alexfalkowski/go-service/commit/61b68bb3a319176d3f5d020ee6e20b52e32333f2) fix(deps): upgraded go to v1.22.3 (#713)
- [`19e3973`](https://github.com/alexfalkowski/go-service/commit/19e39735532242d10667c3b7938b50229ee6b206) build(ci): change versions of deps (#711)

## [v1.179.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.179.2) - 2024-05-06

- [`4d1de74`](https://github.com/alexfalkowski/go-service/commit/4d1de74b34888c4a8e54e4797e1a7e72efe7f231) fix(security): remove unnecessary checks (#710)

## [v1.179.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.179.1) - 2024-05-06

- [`500f1b0`](https://github.com/alexfalkowski/go-service/commit/500f1b07c94be6c55253d9611fc91a8f5392ddcd) fix(deps): bump google.golang.org/protobuf from 1.34.0 to 1.34.1 (#708)
- [`efdc8f4`](https://github.com/alexfalkowski/go-service/commit/efdc8f473e88235d57e928b14191f28c06b0bb11) build(deps): bump bin from `253f8fe` to `693b345` (#707)
- [`7516d0c`](https://github.com/alexfalkowski/go-service/commit/7516d0c6a16ebe8b0d21d7268da5c652da17836a) build(sql): remove checks for nil (#706)

## [v1.179.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.179.0) - 2024-05-05

- [`9aec9af`](https://github.com/alexfalkowski/go-service/commit/9aec9af5d72df8c63566aa479d2334eccfd9fb59) feat(time): move time to have a network kind (#705)

## [v1.178.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.178.1) - 2024-05-05

- [`3621848`](https://github.com/alexfalkowski/go-service/commit/36218485deba40dc9cab5f38b8e42169c9720844) fix(feature): move checking to config (#704)
- [`a8d98a8`](https://github.com/alexfalkowski/go-service/commit/a8d98a8436c7470070ce8b68673310e206188663) build(docs): rename var (#703)

## [v1.178.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.178.0) - 2024-05-04

- [`98d8e16`](https://github.com/alexfalkowski/go-service/commit/98d8e16057f76c3b0df02a0058eaafec8870265b) feat(config): simplify by removing configurator (#702)
- [`fcf7edf`](https://github.com/alexfalkowski/go-service/commit/fcf7edfb8a0c6a0fbce73c06012ab1ac58736648) build(lint): remove execinquery (#701)

## [v1.177.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.177.0) - 2024-05-03

- [`888346f`](https://github.com/alexfalkowski/go-service/commit/888346fc115720337ed90d053067ec6029b09fc4) feat(config): remove enabled (#700)

## [v1.176.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.176.2) - 2024-05-03

- [`8896ba1`](https://github.com/alexfalkowski/go-service/commit/8896ba1452b3020c8fdd680e15c16c7bbdcde157) fix(sql): move module (#699)

## [v1.176.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.176.1) - 2024-05-03

- [`fb7ca4c`](https://github.com/alexfalkowski/go-service/commit/fb7ca4c4d0b411548f7b1bc6aee1fa1a3da987c6) fix(redis): stats can be nil (#698)

## [v1.176.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.176.0) - 2024-05-03

- [`d97884a`](https://github.com/alexfalkowski/go-service/commit/d97884a22ec4acb10e04e35cfeaa71d0a0b9755c) feat(time): add network time (#697)

## [v1.175.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.175.0) - 2024-05-03

- [`1262d72`](https://github.com/alexfalkowski/go-service/commit/1262d726ee7a390179be6f2c0a81d0f2db0faffc) feat(compressor): add s2 and zstd compressor (#696)

## [v1.174.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.174.0) - 2024-05-02

- [`a4e92e1`](https://github.com/alexfalkowski/go-service/commit/a4e92e1cd28eee42a6384ebdcfaa656154d20406) feat(security): load key pair from memory (#692)
- [`1733b58`](https://github.com/alexfalkowski/go-service/commit/1733b5843e35e85bd1033dae47d0f45da20855c9) build(deps): bump bin from `63b9d75` to `253f8fe` (#691)
- [`4b85341`](https://github.com/alexfalkowski/go-service/commit/4b85341b6fbb83dbf202282597df98e1909f2729) test(transport): split to different files (#690)

## [v1.173.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.173.1) - 2024-05-02

- [`0ae6a79`](https://github.com/alexfalkowski/go-service/commit/0ae6a79ea2c81c22c7d7207709fee3c39cbf3a80) fix(deps): bump github.com/shirou/gopsutil/v3 from 3.24.3 to 3.24.4 (#689)

## [v1.173.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.173.0) - 2024-05-01

- [`e390e99`](https://github.com/alexfalkowski/go-service/commit/e390e994a6189217cefe7f71ecdd87170008f55d) feat(events): restructure receiver (#688)

## [v1.172.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.10) - 2024-05-01

- [`6814853`](https://github.com/alexfalkowski/go-service/commit/681485381eb35c3761b001d8eab126f8322e2140) fix(grpc): move stream back to be an implementation detail (#687)

## [v1.172.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.9) - 2024-05-01

- [`7ce611a`](https://github.com/alexfalkowski/go-service/commit/7ce611a4138199baaa4bf32d3bc7272a1fee0db7) fix(grpc): verify stream of metrics (#686)
- [`d41698f`](https://github.com/alexfalkowski/go-service/commit/d41698fd3d9a7804692e22f2d10ddd205902cd7a) build(deps): bump bin from `cf2a686` to `63b9d75` (#682)

## [v1.172.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.8) - 2024-04-30

- [`1133fb9`](https://github.com/alexfalkowski/go-service/commit/1133fb99801c3380cbc53a7cafb9a17e42e8f17f) fix(grpc): verify metrics gRPC stream (#681)

## [v1.172.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.7) - 2024-04-30

- [`bc2f585`](https://github.com/alexfalkowski/go-service/commit/bc2f585e7cc404c57ac4c9e09d86440e201a0931) fix(deps): bump github.com/sony/gobreaker from 0.5.0 to 1.0.0 (#680)

## [v1.172.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.6) - 2024-04-30

- [`b186c2c`](https://github.com/alexfalkowski/go-service/commit/b186c2ce394d1f8b05979553156cdf8c2ec1d840) fix(deps): bump google.golang.org/protobuf from 1.33.0 to 1.34.0 (#679)
- [`c062eb8`](https://github.com/alexfalkowski/go-service/commit/c062eb83a842b6fb417ec798a82e99551ec0a323) build(deps): bump bin from `e47704b` to `cf2a686` (#678)

## [v1.172.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.5) - 2024-04-30

- [`46fa65f`](https://github.com/alexfalkowski/go-service/commit/46fa65f3f401a6f70da203f9c7c336e3ecb8d5e9) fix(errors): decorate the input error (#677)

## [v1.172.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.4) - 2024-04-30

- [`6cec24a`](https://github.com/alexfalkowski/go-service/commit/6cec24ac1b36261bc216cebb99d7ab086c8c4538) fix(errors): prefix errors for better experience (#676)

## [v1.172.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.3) - 2024-04-30

- [`0aef4be`](https://github.com/alexfalkowski/go-service/commit/0aef4be4fff96d69878d9f49c8e0d2e9e4ec34c1) fix(net): handle connection refused (#675)

## [v1.172.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.2) - 2024-04-29

- [`ff356cb`](https://github.com/alexfalkowski/go-service/commit/ff356cbf4c1eb67668d9d0c79c24c2891d1373cf) fix(grpc): remove duplication (#674)
- [`9eb49e4`](https://github.com/alexfalkowski/go-service/commit/9eb49e40b72cacf83143054feb674712c63a5088) build(ci): cover all (#673)
- [`10a0a8a`](https://github.com/alexfalkowski/go-service/commit/10a0a8abc2b0fb19c57f2dde172e92dce695920a) build(deps): bump bin from `3976e45` to `e47704b` (#672)
- [`f6d5e97`](https://github.com/alexfalkowski/go-service/commit/f6d5e9744ffece583bd124d69fae877ae771b2e7) build(ci): html coverage (#671)
- [`6a84585`](https://github.com/alexfalkowski/go-service/commit/6a8458596d3c398cdf9e9dc479fc697badbe1cbb) build(deps): bump bin from `1755d45` to `3976e45` (#670)
- [`9b00a49`](https://github.com/alexfalkowski/go-service/commit/9b00a497c9ecf11774608006accf3d57bca05f8e) build(deps): bump bin from `6fee1b8` to `1755d45` (#669)
- [`91aa5df`](https://github.com/alexfalkowski/go-service/commit/91aa5df5e3df101005def3d5711c89d5ae6601bc) test(specs): cleanup (#668)

## [v1.172.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.1) - 2024-04-27

- [`ec259d8`](https://github.com/alexfalkowski/go-service/commit/ec259d8d8d5fb767f6c1367d234cc46955a4311a) fix(tracer): remove hardcoded host (#667)
- [`c06bd98`](https://github.com/alexfalkowski/go-service/commit/c06bd98ab510ee30cc51a1ce9eefb09835aa739d) test(redis): use in memory redis (#666)
- [`97dcb90`](https://github.com/alexfalkowski/go-service/commit/97dcb9083229ffe8568e4d17a28183ba70e7fc79) build(deps): bump bin from `62c9b8d` to `6fee1b8` (#665)

## [v1.172.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.172.0) - 2024-04-27

- [`50531af`](https://github.com/alexfalkowski/go-service/commit/50531afcdc11ed384cb3869ee66d7f4dd3282bd2) feat(http): use httpsnoop (#662)
- [`ddcccfc`](https://github.com/alexfalkowski/go-service/commit/ddcccfcb3cd1e2abf3d321e8504dc6b9b9c375b0) test(config): verify config (#661)

## [v1.171.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.171.0) - 2024-04-26

- [`10ada17`](https://github.com/alexfalkowski/go-service/commit/10ada17e632dda1c50698a07c61a47c000c9e2be) feat(tracer): move start to life cycle (#660)
- [`9ea57bb`](https://github.com/alexfalkowski/go-service/commit/9ea57bb62a63234667a26d5a9637896a10e8cc45) test(server): verify grpc port (#657)

## [v1.170.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.170.0) - 2024-04-26

- [`ae7266a`](https://github.com/alexfalkowski/go-service/commit/ae7266af46275ce224b4a77347c86c6943313cd9) feat(server): move listener to net (#656)
- [`2ada426`](https://github.com/alexfalkowski/go-service/commit/2ada42657fa6f229f32f7a624753d487a7dfa922) docs(server): rename docs to suit the package (#655)
- [`e2235db`](https://github.com/alexfalkowski/go-service/commit/e2235dbe2777542c2eaceedf3b4b660fde9de23b) build(deps): bump bin from `c06391c` to `62c9b8d` (#654)

## [v1.169.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.169.0) - 2024-04-25

- [`556ff2d`](https://github.com/alexfalkowski/go-service/commit/556ff2d133c0dc0d5f9fe86ca0dac11f1c7a53dd) feat(server): move package (#653)

## [v1.168.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.168.1) - 2024-04-25

- [`f8dfa6f`](https://github.com/alexfalkowski/go-service/commit/f8dfa6fd16e16c5233743b8ec68839b7423b5f44) fix(net): simplify server (#652)
- [`5c76f4a`](https://github.com/alexfalkowski/go-service/commit/5c76f4abf874baf7a35fbd186cf9b7995810f8c9) docs(limiter): add kinds (#651)

## [v1.168.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.168.0) - 2024-04-25

- [`2f75da7`](https://github.com/alexfalkowski/go-service/commit/2f75da7183c97db2a7e4564eb9762ef7d4a26d98) feat(limiter): limit by ip (#650)

## [v1.167.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.167.0) - 2024-04-25

- [`ef9e5c9`](https://github.com/alexfalkowski/go-service/commit/ef9e5c9ba04a3ad85b1de8df2d7b56a99ffa4a6a) feat(runtime): add must (#649)

## [v1.166.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.166.3) - 2024-04-25

- [`deeb679`](https://github.com/alexfalkowski/go-service/commit/deeb679ba852f4b8f3ed2965fefc0a262e7b5466) fix(limiter): remove no key as it is not used (#648)

## [v1.166.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.166.2) - 2024-04-25

- [`f925383`](https://github.com/alexfalkowski/go-service/commit/f925383fb32942b65e6f4b4478e6007415df00be) fix(net): handle not configured server (#647)

## [v1.166.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.166.1) - 2024-04-25

- [`4b554ad`](https://github.com/alexfalkowski/go-service/commit/4b554ad0d34b7ede1c652037ffebd49e1c76e9f8) fix(http): handle closed error (#646)

## [v1.166.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.166.0) - 2024-04-25

- [`a126d3b`](https://github.com/alexfalkowski/go-service/commit/a126d3bfdc76d717068bfb12e90400bfb4a2fd13) feat(limiter): configure automatically when enabled (#645)

## [v1.165.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.165.4) - 2024-04-25

- [`9861d93`](https://github.com/alexfalkowski/go-service/commit/9861d9346701a76a241249c94e8c0d8d56001c48) fix(deps): bump github.com/rs/cors from 1.10.1 to 1.11.0 (#644)

## [v1.165.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.165.3) - 2024-04-24

- [`2629147`](https://github.com/alexfalkowski/go-service/commit/2629147ce220588ada101983bd010dcbb13ab441) fix(metrics): add target info (#643)

## [v1.165.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.165.2) - 2024-04-24

- [`1bb9891`](https://github.com/alexfalkowski/go-service/commit/1bb9891a231fb5ba1075e97f9cae8cc21e95fc22) fix(tracer): log error (#642)

## [v1.165.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.165.1) - 2024-04-24

- [`35a756f`](https://github.com/alexfalkowski/go-service/commit/35a756fcfe4c6d027da4270d1ba163faa6770a6c) fix(otel): use semconv v1.25.0 (#641)

## [v1.165.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.165.0) - 2024-04-24

- [`de24c8b`](https://github.com/alexfalkowski/go-service/commit/de24c8b5accc50ad88c08a5c4ff0b96484a832e1) feat(net): handle tls correctly (#640)

## [v1.164.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.164.1) - 2024-04-24

- [`6812832`](https://github.com/alexfalkowski/go-service/commit/68128329c44bffbaffaf8e1d56197057b6ba9c72) fix(grpc): use correct name (#639)

## [v1.164.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.164.0) - 2024-04-24

- [`25da1d3`](https://github.com/alexfalkowski/go-service/commit/25da1d3a98b7f1e21dbfb18680b56a673fd5e51a) feat(net): add server to be reused (#638)

## [v1.163.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.163.2) - 2024-04-24

- [`565d233`](https://github.com/alexfalkowski/go-service/commit/565d233886cbd09631a0a294dda518fe11a2e6da) fix(logger): set for redis (#636)

## [v1.163.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.163.1) - 2024-04-24

- [`dba3abd`](https://github.com/alexfalkowski/go-service/commit/dba3abd8f4e6a24a56d3088c613de3453be6eabc) fix(deps): bump github.com/jmoiron/sqlx from 1.3.5 to 1.4.0 (#637)

## [v1.163.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.163.0) - 2024-04-23

- [`e07e0cc`](https://github.com/alexfalkowski/go-service/commit/e07e0cc13021bca637bc5d0cc5635d9e423ef9d7) feat(http): remove error from client (#635)

## [v1.162.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.162.0) - 2024-04-23

- [`3c3a35c`](https://github.com/alexfalkowski/go-service/commit/3c3a35c49dcb2077e31a75c0195533380aac32f6) feat(tracer): standardise operation name (#634)

## [v1.161.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.161.0) - 2024-04-23

- [`e740af8`](https://github.com/alexfalkowski/go-service/commit/e740af8508b132690ee20b0e078c06b43d4bdabc) feat(metrics): must be created (#633)

## [v1.160.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.160.0) - 2024-04-23

- [`5886df1`](https://github.com/alexfalkowski/go-service/commit/5886df13d61edf93df919245e1840813060fdc2a) feat(meta): follow case conventions (#632)

## [v1.159.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.159.0) - 2024-04-22

- [`146e4aa`](https://github.com/alexfalkowski/go-service/commit/146e4aa230e91f8806c95ed5ebdd10db2ae21647) feat(debug): add devices from psutil (#631)

## [v1.158.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.158.1) - 2024-04-22

- [`703f806`](https://github.com/alexfalkowski/go-service/commit/703f806ca8052c1676ed10bc0bdec03311ddbe59) fix(transport): add http timeouts (#630)

## [v1.158.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.158.0) - 2024-04-22

- [`e375652`](https://github.com/alexfalkowski/go-service/commit/e375652360c39cb882c2258d1f996bf08b8848b2) feat(marshaller): add gob (#629)

## [v1.157.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.157.0) - 2024-04-22

- [`a64abfe`](https://github.com/alexfalkowski/go-service/commit/a64abfefaaca22722e0ca03d51f49829f0aa41b8) feat(compressor): add none (#628)
- [`c666507`](https://github.com/alexfalkowski/go-service/commit/c6665071dfc59ef16b4cfd2baa9c9b87291b6b3a) build(deps): bump bin from `d5f116d` to `c06391c` (#627)

## [v1.156.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.156.0) - 2024-04-21

- [`f22cb73`](https://github.com/alexfalkowski/go-service/commit/f22cb7380bc7edceca18e05e217c3f503345b4f9) feat(redis): make sure marshaller and compressor can be configured (#626)

## [v1.155.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.155.0) - 2024-04-20

- [`ad05af0`](https://github.com/alexfalkowski/go-service/commit/ad05af01dd63f7bd6247ac26f1686ecb273d26fd) feat(transport): expose server for gRPC (#625)

## [v1.154.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.7) - 2024-04-20

- [`ccc03a6`](https://github.com/alexfalkowski/go-service/commit/ccc03a62e75c0afedddc535730e3bdef6f982082) fix(telemetry): move metrics to http (#624)
- [`f5b4c24`](https://github.com/alexfalkowski/go-service/commit/f5b4c243415f377475fb086341ac51de10b88b5c) build(ci): find the most recently generated cache (#623)

## [v1.154.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.6) - 2024-04-17

- [`77f0774`](https://github.com/alexfalkowski/go-service/commit/77f077458fc7ce365205b593b071baa87fc75268) fix(deps): bump github.com/KimMachineGun/automemlimit from 0.5.0 to 0.6.0 (#622)
- [`1de5d13`](https://github.com/alexfalkowski/go-service/commit/1de5d13151d63c6fe9ad387c782b9a46728b2cea) build(ci): store tests (#621)
- [`92f1b30`](https://github.com/alexfalkowski/go-service/commit/92f1b302582fdbfdec85d89e27a26dc33e535962) build(deps): bump bin from `51c6ece` to `d5f116d` (#620)
- [`728eb93`](https://github.com/alexfalkowski/go-service/commit/728eb930b392311b770f1d91aea0b73b777415f5) build(ci): change cache key (#619)
- [`c1d072c`](https://github.com/alexfalkowski/go-service/commit/c1d072ce15fc65e15bd9ab5fc180a54c7b823917) build(deps): bump bin from `c322696` to `51c6ece` (#618)
- [`b5cdae0`](https://github.com/alexfalkowski/go-service/commit/b5cdae0af814dc01a93a18548c1f98157fcfe1b6) build(ci): cache go build (#617)
- [`82bdc04`](https://github.com/alexfalkowski/go-service/commit/82bdc04f6f4ec7848b454745ae47056280cddc90) build(ci): cache deps (#616)
- [`22a318f`](https://github.com/alexfalkowski/go-service/commit/22a318f2aede1317bb1fe293a55b04e78d87e573) build(deps): bump bin from `13a7302` to `c322696` (#615)

## [v1.154.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.5) - 2024-04-14

- [`79f368a`](https://github.com/alexfalkowski/go-service/commit/79f368acbf78614e323e8100684fc531e7fa51f9) fix(tracer): pass target to grpc (#614)

## [v1.154.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.4) - 2024-04-14

- [`3c0da99`](https://github.com/alexfalkowski/go-service/commit/3c0da99f240119fb8385cf36ed4a0099290d17a2) fix(tracer): pass all the url for http (#613)

## [v1.154.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.3) - 2024-04-14

- [`de05902`](https://github.com/alexfalkowski/go-service/commit/de05902b040ab5af7a5a3004f59122a8ceb9d6b8) fix(http): move mux to func (#612)

## [v1.154.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.2) - 2024-04-14

- [`2b1b2bb`](https://github.com/alexfalkowski/go-service/commit/2b1b2bbd2d6efb7c072ad00139f53485068ccf94) fix(metrics): standard labels (#611)

## [v1.154.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.1) - 2024-04-14

- [`f1c2f8d`](https://github.com/alexfalkowski/go-service/commit/f1c2f8db5f5097d74cd8406eb7ac938748642da1) fix(debug): seperate fgprof (#610)

## [v1.154.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.154.0) - 2024-04-13

- [`0965f04`](https://github.com/alexfalkowski/go-service/commit/0965f04ed6b12a9070aa70c861ae2a8042f62882) feat(telemetry): unify to be consistent (#609)

## [v1.153.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.153.2) - 2024-04-12

- [`23312a3`](https://github.com/alexfalkowski/go-service/commit/23312a337fe32c8500d073ca259a4ac99443233e) fix(metrics): use resource in provider (#607)

## [v1.153.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.153.1) - 2024-04-12

- [`b72e8e7`](https://github.com/alexfalkowski/go-service/commit/b72e8e7c105cfd7b3098d8198374a3dad823f692) fix(deps): upgraded github.com/alexfalkowski/go-health to v1.16.1 (#606)

## [v1.153.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.153.0) - 2024-04-11

- [`e5e8729`](https://github.com/alexfalkowski/go-service/commit/e5e8729c375195a002bdc66f35e98be652ca90cd) feat(debug): add fgprof (#605)

## [v1.152.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.152.0) - 2024-04-11

- [`aafee7e`](https://github.com/alexfalkowski/go-service/commit/aafee7e55fe7e01c936634c85c7a39aa6ea75c5f) feat(telemetry): add otlp metrics (#604)

## [v1.151.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.151.2) - 2024-04-10

- [`5bf1c3e`](https://github.com/alexfalkowski/go-service/commit/5bf1c3e5db261b6f7c16c761ef908ef70b64b18c) fix(deps): upgraded go.opentelemetry.io/otel v1.24.0 => v1.25.0 (#603)

## [v1.151.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.151.1) - 2024-04-10

- [`616ad0f`](https://github.com/alexfalkowski/go-service/commit/616ad0ff8ec57dd619b16d06a4e6cb81bdf965a3) fix(deps): bump google.golang.org/grpc from 1.63.0 to 1.63.2 (#601)

## [v1.151.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.151.0) - 2024-04-08

- [`6ecef04`](https://github.com/alexfalkowski/go-service/commit/6ecef044023932834e5be193c524fe0740dc061f) feat(grpc): use new client and avoid deprecation (#597)

## [v1.150.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.150.2) - 2024-04-04

- [`5ccc2a8`](https://github.com/alexfalkowski/go-service/commit/5ccc2a8e7658e2c2604ee329c7b936922f0d8733) fix(deps): upgraded google.golang.org/grpc v1.62.1 => v1.63.0 (#596)

## [v1.150.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.150.1) - 2024-04-04

- [`8b30d9d`](https://github.com/alexfalkowski/go-service/commit/8b30d9d777cca894a6c52e11deee3a5e9489f4df) fix(deps): upgraded golang.org/x/net v0.21.0 => v0.23.0 (#595)
- [`beb432a`](https://github.com/alexfalkowski/go-service/commit/beb432ab48cdfcb705f5933c44a8602e75b9ea20) build(deps): bump bin from `ed17684` to `13a7302` (#593)
- [`55666c1`](https://github.com/alexfalkowski/go-service/commit/55666c1eb406041acd373073d4b24f75baecc906) build(deps): bump bin from `024be7f` to `ed17684` (#592)

## [v1.150.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.150.0) - 2024-04-01

- [`010cde8`](https://github.com/alexfalkowski/go-service/commit/010cde82e01a987929d69c229af2ca0a060d3049) feat(time): add must parse (#591)

## [v1.149.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.149.4) - 2024-04-01

- [`1da2349`](https://github.com/alexfalkowski/go-service/commit/1da234957e4c0f35fc9b023be7f3754037e6e327) fix(deps): bump github.com/shirou/gopsutil/v3 from 3.24.2 to 3.24.3 (#590)
- [`ae29a70`](https://github.com/alexfalkowski/go-service/commit/ae29a706213867b0e465e3f6d9f36380e0ba52d9) build(deps): bump bin from `608889f` to `024be7f` (#589)
- [`8b089f2`](https://github.com/alexfalkowski/go-service/commit/8b089f23b40d731ce9d8575541a7fcff5e9a2a08) build(deps): bump bin from `a19d7ca` to `608889f` (#588)

## [v1.149.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.149.3) - 2024-03-28

- [`9e8e75a`](https://github.com/alexfalkowski/go-service/commit/9e8e75a3df77518f44e873136017202cf7f55511) fix(deps): update github.com/alexfalkowski/go-health to v1.16.0 (#587)

## [v1.149.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.149.2) - 2024-03-28

- [`e7b40cb`](https://github.com/alexfalkowski/go-service/commit/e7b40cbc5d21d62b97324fbb97d1571bde0c6cf1) fix(deps): update undetected (#586)
- [`62cef92`](https://github.com/alexfalkowski/go-service/commit/62cef92ee13135e5c7bcd0cb4424d6e6e2a38698) build(dependabot): change prefixes (#585)

## [v1.149.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.149.1) - 2024-03-28

- [`11e3eb7`](https://github.com/alexfalkowski/go-service/commit/11e3eb74bdfa265cae04916ad8d592dbd40088e5) fix(deps): bump bin from `b9b6ae3` to `a19d7ca` (#584)

## [v1.149.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.149.0) - 2024-03-27

- [`bbe0c9f`](https://github.com/alexfalkowski/go-service/commit/bbe0c9fcfe9b3df5ca187f9f0f91df26adb27a5d) feat(telemetry): remove otelconfig (#583)

## [v1.148.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.148.1) - 2024-03-27

- [`468e82e`](https://github.com/alexfalkowski/go-service/commit/468e82e090d664a95b01e8bd4bc66354a7dd3669) fix(deps): bump bin from `60071ae` to `b9b6ae3` (#582)
- [`839135a`](https://github.com/alexfalkowski/go-service/commit/839135abcecb81336aef2dfad6e8fe666b6a11b1) build(make): remove phony (#573)
- [`f56d980`](https://github.com/alexfalkowski/go-service/commit/f56d9800d68d54b5899110ae147b2329b75989cb) build: add git tasks to makefile (#569)

## [v1.148.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.148.0) - 2024-03-21

- [`c3f0984`](https://github.com/alexfalkowski/go-service/commit/c3f0984fdb37863596477c6a2fcf6a3270e34b28) feat(telemetry): conflict with schema url (#563)

## [v1.147.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.147.2) - 2024-03-21

- [`4f97e51`](https://github.com/alexfalkowski/go-service/commit/4f97e51c110aaa036f0a3d0279115fbfe48f5853) fix(meta): chaneg order to get from context first (#562)

## [v1.147.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.147.1) - 2024-03-21

- [`de12ace`](https://github.com/alexfalkowski/go-service/commit/de12ace690fffe98f117b5778efdc28f434e3028) fix(meta): append info to correct transport (#561)

## [v1.147.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.147.0) - 2024-03-20

- [`091be07`](https://github.com/alexfalkowski/go-service/commit/091be070e32b7513375e387dd973007e255dbf6d) feat(logger): add dev logger (#560)

## [v1.146.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.146.0) - 2024-03-20

- [`f82c170`](https://github.com/alexfalkowski/go-service/commit/f82c1702c452932b43b7f7ef2c1b61cd6082fdd3) feat(meta): introduce Valuer (#559)
- [`feaa01e`](https://github.com/alexfalkowski/go-service/commit/feaa01e4615f24295e8eeca537d0ed73990216f9) ci: remove depracation (#558)

## [v1.145.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.145.3) - 2024-03-20

- [`605119f`](https://github.com/alexfalkowski/go-service/commit/605119f908d87588e33c482213120782b4c44823) fix(deps): update google.golang.org/genproto/googleapis/api to v0.0.0… (#557)

## [v1.145.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.145.2) - 2024-03-19

- [`9aa1d44`](https://github.com/alexfalkowski/go-service/commit/9aa1d44db8d9d9af4eb2cd08e6a2b6b96ff2d563) fix(meta): handler error (#556)

## [v1.145.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.145.1) - 2024-03-19

- [`ea2e1f6`](https://github.com/alexfalkowski/go-service/commit/ea2e1f60186ca4aafd812f6fe30e99a69c6f05bf) fix(logger): use stringer (#555)

## [v1.145.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.145.0) - 2024-03-19

- [`eed21f7`](https://github.com/alexfalkowski/go-service/commit/eed21f74b842d9a5c5a11f7b6ea766dfddc7ba92) feat(meta): use fmt.Stringer (#554)

## [v1.144.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.144.1) - 2024-03-18

- [`abbff63`](https://github.com/alexfalkowski/go-service/commit/abbff63529abee52cfc333811336bda9aeaf18d2) fix(cmd): allow shortform flag (#553)

## [v1.144.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.144.0) - 2024-03-18

- [`337ca07`](https://github.com/alexfalkowski/go-service/commit/337ca07172b5c9b595c4983c2f049915c4cb9f9e) feat(debug): move to servers (#552)

## [v1.143.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.143.2) - 2024-03-18

- [`7376464`](https://github.com/alexfalkowski/go-service/commit/7376464137c449a927b23e2961f1aff27d0943e9) fix(config): verify enabled (#551)

## [v1.143.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.143.1) - 2024-03-18

- [`28ba256`](https://github.com/alexfalkowski/go-service/commit/28ba2565440754cd06c324e88af5da9d619b622f) fix(deps): bump github.com/alexfalkowski/go-health from 1.14.2 to 1.15.0 (#550)

## [v1.143.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.143.0) - 2024-03-17

- [`18ed017`](https://github.com/alexfalkowski/go-service/commit/18ed0178bef46948fb943de48db307c364be23e6) feat(config): allow missing values (#549)

## [v1.142.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.142.0) - 2024-03-16

- [`77c8615`](https://github.com/alexfalkowski/go-service/commit/77c861565517993cb827a1988fadf4d5880d65a0) feat(deps): add github.com/goccy/go-json v0.10.2 (#548)

## [v1.141.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.141.1) - 2024-03-16

- [`2e638b4`](https://github.com/alexfalkowski/go-service/commit/2e638b407835389ef5b61a3939a001322fa13218) fix(marshaller): indent JSON marshalling (#546)
- [`c603260`](https://github.com/alexfalkowski/go-service/commit/c6032603c3bd9686f5e7682b12bfa270fa25be63) docs: remove TOML (#545)

## [v1.141.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.141.0) - 2024-03-16

- [`931500e`](https://github.com/alexfalkowski/go-service/commit/931500e65f3192cf15edebf82f9c6b5a8c627582) feat(marshaller): add JSON (#544)

## [v1.140.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.140.1) - 2024-03-16

- [`b747b19`](https://github.com/alexfalkowski/go-service/commit/b747b19593084f82485002926a0e82da794c73e0) fix(marshaller): factory missing type (#543)

## [v1.140.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.140.0) - 2024-03-15

- [`56c71be`](https://github.com/alexfalkowski/go-service/commit/56c71be0fcf37e20158daa1b6e65c3bb619d9962) feat(transport): move hooks to http (#541)

## [v1.139.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.139.1) - 2024-03-14

- [`34af597`](https://github.com/alexfalkowski/go-service/commit/34af597f1288342bc613a4ee38b6d376ca23c3b3) fix(config): add omitempty (#539)

## [v1.139.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.139.0) - 2024-03-14

- [`b2135e4`](https://github.com/alexfalkowski/go-service/commit/b2135e43b8a4cfb6f5e7d4c4be983d4a5867525d) feat: change license to MIT (#538)

## [v1.138.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.138.5) - 2024-03-14

- [`015b761`](https://github.com/alexfalkowski/go-service/commit/015b7617c4508ed95d7d8cba1bfe94d4b88d1723) fix(cmd): add kind to config (#537)

## [v1.138.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.138.4) - 2024-03-13

- [`3e20c21`](https://github.com/alexfalkowski/go-service/commit/3e20c21941a097b62899651d708abad74d8f60b2) fix(deps): bump go.uber.org/fx from 1.20.1 to 1.21.0 (#536)

## [v1.138.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.138.3) - 2024-03-12

- [`c24906b`](https://github.com/alexfalkowski/go-service/commit/c24906be625ed1faa492bf3872bd95adb462a27c) fix(tracer): remove syncer (#535)

## [v1.138.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.138.2) - 2024-03-12

- [`3d42544`](https://github.com/alexfalkowski/go-service/commit/3d42544527567e0d9ae97f066aaa702e331352e9) fix(events): set http error (#534)

## [v1.138.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.138.1) - 2024-03-11

- [`ff686a8`](https://github.com/alexfalkowski/go-service/commit/ff686a8a9fe0c4b84bc8ea8efb0bf1d0549b61b2) fix(transport): add events config (#533)

## [v1.138.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.138.0) - 2024-03-11

- [`34da194`](https://github.com/alexfalkowski/go-service/commit/34da1947132a8ea6fb7d0fe500c8c04f874c9b3b) feat(hooks): add standard webhooks (#532)

## [v1.137.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.137.2) - 2024-03-11

- [`205b9d3`](https://github.com/alexfalkowski/go-service/commit/205b9d3c414a182ceb970459e28c5c08dc45cdd4) fix(deps): bump github.com/jackc/pgx/v5 from 5.5.4 to 5.5.5 (#531)

## [v1.137.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.137.1) - 2024-03-10

- [`95d5a51`](https://github.com/alexfalkowski/go-service/commit/95d5a51059122de0f47dfbbdc14b1ad78191034e) fix(security): remove func (#530)

## [v1.137.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.137.0) - 2024-03-09

- [`0989701`](https://github.com/alexfalkowski/go-service/commit/0989701ec626baf19c8a9055d9ce98d0a5788591) feat(transport): add cloud events (#529)

## [v1.136.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.136.0) - 2024-03-08

- [`73d2054`](https://github.com/alexfalkowski/go-service/commit/73d2054472c809ee1be68b98fef6cfbc87773b74) feat(http): allow to enable/disable (#528)

## [v1.135.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.135.0) - 2024-03-08

- [`04addcc`](https://github.com/alexfalkowski/go-service/commit/04addcca759adf8f4eab5108ee82064a35984af7) feat(cmd): remove worker (#527)

## [v1.134.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.134.2) - 2024-03-08

- [`2b28c39`](https://github.com/alexfalkowski/go-service/commit/2b28c39830cbd869ced0a692d204490ecaa33fdb) fix(deps): update github.com/alexfalkowski/go-health to v1.14.2 (#526)

## [v1.134.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.134.1) - 2024-03-08

- [`409629a`](https://github.com/alexfalkowski/go-service/commit/409629ae3f3c7c1ea3f0aaf24b203bd746dcb1df) fix(server): reuse config (#525)

## [v1.134.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.134.0) - 2024-03-08

- [`ea35aca`](https://github.com/alexfalkowski/go-service/commit/ea35acaa63cfeb89692e2f163d9ba9ed26094bf2) feat(feature): add unary client interceptors (#524)

## [v1.133.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.133.0) - 2024-03-07

- [`2fc68ee`](https://github.com/alexfalkowski/go-service/commit/2fc68ee234cd0408bb3c98d3a3731a5c3edf108c) feat(transport): remove NSQ (#523)
- [`f6190db`](https://github.com/alexfalkowski/go-service/commit/f6190db9abaf83c63a02d50f8ae99deefd68d642) test(deps): buf (#521)

## [v1.132.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.132.2) - 2024-03-06

- [`08820c6`](https://github.com/alexfalkowski/go-service/commit/08820c6208b0e44168503c7d70be4743d2ba0b73) fix(deps): bump google.golang.org/protobuf from 1.32.0 to 1.33.0 (#516)

## [v1.132.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.132.1) - 2024-03-06

- [`42e951f`](https://github.com/alexfalkowski/go-service/commit/42e951f554f4e956f76ef4200b3c8020b1a6ab20) fix(go): bump to v1.22.1 (#520)

## [v1.132.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.132.0) - 2024-03-06

- [`8965fe2`](https://github.com/alexfalkowski/go-service/commit/8965fe294c70ccfd9d460bccb7ac46346736ffe4) feat(config): add client (#518)

## [v1.131.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.131.3) - 2024-03-05

- [`5f282b8`](https://github.com/alexfalkowski/go-service/commit/5f282b86d18cc881a56d034cdbc19f542884697d) fix(deps): bump google.golang.org/grpc from 1.62.0 to 1.62.1 (#517)

## [v1.131.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.131.2) - 2024-03-05

- [`b146d39`](https://github.com/alexfalkowski/go-service/commit/b146d395761e44dd995fcb937c1fb663a857e5d5) fix(transport): listen before starting (#515)

## [v1.131.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.131.1) - 2024-03-05

- [`8ccc1d9`](https://github.com/alexfalkowski/go-service/commit/8ccc1d9c7dd92cca88ba2004f3e744b90f8db5d5) fix(deps): bump github.com/jackc/pgx/v5 from 5.5.3 to 5.5.4 (#514)

## [v1.131.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.131.0) - 2024-03-04

- [`5412e8a`](https://github.com/alexfalkowski/go-service/commit/5412e8a3fde2aed0fd4783ef33233ed37022d522) feat(security): have one cert_file and key_file (#513)

## [v1.130.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.130.0) - 2024-03-04

- [`ec5f734`](https://github.com/alexfalkowski/go-service/commit/ec5f734db79e6bb41935a3224c1e5a8d450a5162) feat(retry): move to reusable config (#512)

## [v1.129.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.129.2) - 2024-03-02

- [`b7cbbeb`](https://github.com/alexfalkowski/go-service/commit/b7cbbeb39dbb8afc4f19399b05ed3f4b484d74a7) fix(deps): bump github.com/shirou/gopsutil/v3 from 3.24.1 to 3.24.2 (#511)

## [v1.129.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.129.1) - 2024-03-02

- [`4eb523f`](https://github.com/alexfalkowski/go-service/commit/4eb523f029d217cd2de81b2dd1ecb596c48675f6) fix(deps): bump the otel group with 1 update (#510)

## [v1.129.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.129.0) - 2024-03-02

- [`ee73b8b`](https://github.com/alexfalkowski/go-service/commit/ee73b8bae77b0c516b9c081c8655188f238a51bb) feat(runtime): add github.com/KimMachineGun/automemlimit v0.5.0 (#509)

## [v1.128.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.128.2) - 2024-02-28

- [`70fddb5`](https://github.com/alexfalkowski/go-service/commit/70fddb531a2197a375016ec1dc2c61a5e9c132be) fix(deps): bump github.com/prometheus/client_golang from 1.18.0 to 1.19.0 (#508)

## [v1.128.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.128.1) - 2024-02-26

- [`a62548d`](https://github.com/alexfalkowski/go-service/commit/a62548d30fefb137aa31cccf94af5496a7035a68) fix(deps): bump the otel group with 1 update (#507)

## [v1.128.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.128.0) - 2024-02-25

- [`d82ccb2`](https://github.com/alexfalkowski/go-service/commit/d82ccb24b60a33e79b530d4cda5652eee89e06c7) feat(telemetry): add baselime (#506)

## [v1.127.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.127.0) - 2024-02-23

- [`28e6d45`](https://github.com/alexfalkowski/go-service/commit/28e6d4588b09a200af5ddb1f03b8d223bcdffc2e) feat(deps): update go.opentelemetry.io/otel v1.24.0 (#505)

## [v1.126.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.126.6) - 2024-02-22

- [`3924adf`](https://github.com/alexfalkowski/go-service/commit/3924adf6e63a45e13ce77d6d853ce187c6d8f617) fix(deps): bump github.com/klauspost/compress from 1.17.6 to 1.17.7 (#504)

## [v1.126.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.126.5) - 2024-02-22

- [`2205b58`](https://github.com/alexfalkowski/go-service/commit/2205b58cc1134e57e1d477b89ffc85436e69cc96) fix(deps): bump google.golang.org/grpc from 1.61.1 to 1.62.0 (#503)

## [v1.126.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.126.4) - 2024-02-21

- [`5325e4e`](https://github.com/alexfalkowski/go-service/commit/5325e4e9458523159558ea11626f3ba86078802f) fix(deps): bump go.uber.org/zap from 1.26.0 to 1.27.0 (#502)
- [`4b4ec39`](https://github.com/alexfalkowski/go-service/commit/4b4ec39df2b9c2ba9cb22b35f7f3f33a89f7c07a) test: verify invalid limiter (#501)

## [v1.126.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.126.3) - 2024-02-19

- [`862046e`](https://github.com/alexfalkowski/go-service/commit/862046e87f75e6544608d9dedb0adc3c5fe8b559) fix(deps): bump github.com/urfave/negroni/v3 from 3.0.0 to 3.1.0 (#499)

## [v1.126.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.126.2) - 2024-02-17

- [`8250f0c`](https://github.com/alexfalkowski/go-service/commit/8250f0c8fe63e2fe46b8b30e388782dfe9fc16d8) fix(deps): update github.com/alexfalkowski/go-health to v1.14.1 (#498)

## [v1.126.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.126.1) - 2024-02-17

- [`3c1e510`](https://github.com/alexfalkowski/go-service/commit/3c1e510afb9d3f9fbed3ae5d143eab05260355cb) fix: toolchain (#497)

## [v1.126.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.126.0) - 2024-02-16

- [`b7da165`](https://github.com/alexfalkowski/go-service/commit/b7da1652114d7ab2d12b2ee0d1d245e56d532c46) feat(go): update to v1.22 (#496)

## [v1.125.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.125.3) - 2024-02-14

- [`c72bc6f`](https://github.com/alexfalkowski/go-service/commit/c72bc6ffdec65e805b61a142848c922f4f8b392c) fix(deps): bump google.golang.org/grpc from 1.61.0 to 1.61.1 (#494)

## [v1.125.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.125.2) - 2024-02-12

- [`b0d4f3b`](https://github.com/alexfalkowski/go-service/commit/b0d4f3b0c22696aedeab3014c3c4908128f6eff0) fix(deps): bump github.com/open-feature/go-sdk from 1.9.0 to 1.10.0 (#493)

## [v1.125.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.125.1) - 2024-02-11

- [`2ffd225`](https://github.com/alexfalkowski/go-service/commit/2ffd2259c6a9f02ea26833936f55a444fdd9219a) fix(deps): use go.opentelemetry.io/otel/semconv/v1.23.0 (#492)

## [v1.125.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.125.0) - 2024-02-11

- [`c4cb994`](https://github.com/alexfalkowski/go-service/commit/c4cb9941a216fd60991a2612cfdc54feb644dd2f) feat(transport): provide multiple servers (#491)
- [`1273f57`](https://github.com/alexfalkowski/go-service/commit/1273f576d8c33b41b4c9becf7ab29273ee35a00e) style: linting (#490)

## [v1.124.17](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.17) - 2024-02-08

- [`80355ea`](https://github.com/alexfalkowski/go-service/commit/80355ea819b4f2e5dc14eb537850fae6b3ce5b8c) fix(deps): bump the otel group with 8 updates (#488)

## [v1.124.16](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.16) - 2024-02-07

- [`526bd6f`](https://github.com/alexfalkowski/go-service/commit/526bd6f60fd4aba8b1773f595ec03d406a6bc7f8) fix(deps): bump the otel group with 4 updates (#487)

## [v1.124.15](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.15) - 2024-02-06

- [`94437fa`](https://github.com/alexfalkowski/go-service/commit/94437faaee2dfbc0d34dca4b73164050669fe7b8) fix(deps): bump github.com/klauspost/compress from 1.17.5 to 1.17.6 (#486)

## [v1.124.14](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.14) - 2024-02-05

- [`5e9bbb2`](https://github.com/alexfalkowski/go-service/commit/5e9bbb2ddfa2ce6d5a83930bcb793d101093c6b2) fix(deps): bump github.com/jackc/pgx/v5 from 5.5.2 to 5.5.3 (#485)

## [v1.124.13](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.13) - 2024-02-01

- [`d578bc7`](https://github.com/alexfalkowski/go-service/commit/d578bc7d8223c65d091bdb316f22e7508bbef7a2) fix(deps): bump github.com/shirou/gopsutil/v3 from 3.23.12 to 3.24.1 (#484)

## [v1.124.12](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.12) - 2024-01-30

- [`88c918a`](https://github.com/alexfalkowski/go-service/commit/88c918a752dc5c920db6c9a86781d8cd4dd4d5d2) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 from 2.19.0 to 2.19.1 (#483)

## [v1.124.11](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.11) - 2024-01-30

- [`c0fbc4a`](https://github.com/alexfalkowski/go-service/commit/c0fbc4a30eb3ed05abf02ef2e50fb28ca8495754) fix(deps): bump github.com/klauspost/compress from 1.17.4 to 1.17.5 (#482)

## [v1.124.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.10) - 2024-01-24

- [`782b675`](https://github.com/alexfalkowski/go-service/commit/782b675d38be53162c59045f38e511cada43eb2f) fix(deps): bump google.golang.org/grpc from 1.60.1 to 1.61.0 (#480)

## [v1.124.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.9) - 2024-01-24

- [`0a96d0c`](https://github.com/alexfalkowski/go-service/commit/0a96d0c77e677abca0f7f55d9be3f27984673ea8) fix(deps): bump github.com/google/uuid from 1.5.0 to 1.6.0 (#481)

## [v1.124.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.8) - 2024-01-18

- [`e09955b`](https://github.com/alexfalkowski/go-service/commit/e09955b443eee1dba077dd1e4583505467f104aa) fix(deps): bump the otel group with 4 updates (#479)

## [v1.124.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.7) - 2024-01-15

- [`c98f553`](https://github.com/alexfalkowski/go-service/commit/c98f5538703a1d8916d789c5b08f886f5735ebcc) fix(deps): bump github.com/jackc/pgx/v5 from 5.5.1 to 5.5.2 (#478)

## [v1.124.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.6) - 2024-01-09

- [`ac83a0e`](https://github.com/alexfalkowski/go-service/commit/ac83a0eb58b7ae1f867d68f5b3d74263f17bb779) fix(deps): bump github.com/alexfalkowski/go-health from 1.13.1 to 1.13.2 (#477)

## [v1.124.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.5) - 2024-01-04

- [`31a0d61`](https://github.com/alexfalkowski/go-service/commit/31a0d61df0c39fe96409d547a34b25faf544f1f9) fix(deps): bump github.com/grpc-ecosystem/grpc-gateway/v2 from 2.18.1 to 2.19.0 (#476)

## [v1.124.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.4) - 2024-01-02

- [`584755d`](https://github.com/alexfalkowski/go-service/commit/584755ddb4fd10a98b4914b8ed819cc09992af0c) fix(deps): bump github.com/shirou/gopsutil/v3 from 3.23.11 to 3.23.12 (#475)

## [v1.124.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.3) - 2023-12-28

- [`fe87fa8`](https://github.com/alexfalkowski/go-service/commit/fe87fa868756a21d58204c9e96a8783b71dfc037) fix(deps): bump github.com/prometheus/client_golang from 1.17.0 to 1.18.0 (#474)

## [v1.124.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.2) - 2023-12-25

- [`8d68486`](https://github.com/alexfalkowski/go-service/commit/8d68486c0e70bcc0f8605f9bf429e68cbad50b1c) fix(deps): bump google.golang.org/protobuf from 1.31.0 to 1.32.0 (#473)

## [v1.124.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.1) - 2023-12-21

- [`7c688d1`](https://github.com/alexfalkowski/go-service/commit/7c688d1203a6361d68d86999973b6c8afec54078) fix(deps): bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.45.0 to 0.46.0 (#472)

## [v1.124.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.124.0) - 2023-12-21

- [`b6fc2a9`](https://github.com/alexfalkowski/go-service/commit/b6fc2a90aa510ddb53957421bc68e29be51d1d1c) feat: add open feature support (#471)

## [v1.123.0](https://github.com/alexfalkowski/go-service/releases/tag/v1.123.0) - 2023-12-21

- [`064cae1`](https://github.com/alexfalkowski/go-service/commit/064cae17b6acaa28fc7dc11d90cc7ce5dc8458c0) feat: remove migrate (#470)

## [v1.122.10](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.10) - 2023-12-19

- [`92f9350`](https://github.com/alexfalkowski/go-service/commit/92f9350dbd0a409f40014737df8b1e2bdc1e958d) fix(deps): bump google.golang.org/grpc from 1.60.0 to 1.60.1 (#467)

## [v1.122.9](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.9) - 2023-12-19

- [`6094a71`](https://github.com/alexfalkowski/go-service/commit/6094a719d0c7cdbef2426f6ec4c0e8d5abaaaa28) fix(deps): bump golang.org/x/crypto from 0.15.0 to 0.17.0 (#466)

## [v1.122.8](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.8) - 2023-12-13

- [`1d0f994`](https://github.com/alexfalkowski/go-service/commit/1d0f994c2daf52cbf5760da6cbfc43230d33168e) fix(deps): bump github.com/google/uuid from 1.4.0 to 1.5.0 (#465)

## [v1.122.7](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.7) - 2023-12-12

- [`65910dd`](https://github.com/alexfalkowski/go-service/commit/65910dd677345046edd0d94520c29a7c8b87821d) fix(deps): bump google.golang.org/grpc from 1.59.0 to 1.60.0 (#463)

## [v1.122.6](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.6) - 2023-12-12

- [`7bdd343`](https://github.com/alexfalkowski/go-service/commit/7bdd34372ab596e5ea630b7dd4626a1394b6f3a4) fix(deps): bump github.com/alexfalkowski/go-health from 1.13.0 to 1.13.1 (#464)

## [v1.122.5](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.5) - 2023-12-11

- [`84288c0`](https://github.com/alexfalkowski/go-service/commit/84288c04c874333338bda6c380cc30b06bff76b6) fix(deps): bump github.com/jackc/pgx/v5 from 5.5.0 to 5.5.1 (#462)
- [`0f43411`](https://github.com/alexfalkowski/go-service/commit/0f434114661d3c5126f931805082cf7e324b7a5a) build(deps): bump bin from `e8ef37a` to `e6271f5` (#461)

## [v1.122.4](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.4) - 2023-12-04

- [`3d6df93`](https://github.com/alexfalkowski/go-service/commit/3d6df93f869e8d6719107825d2676af10e4c3bbe) fix(deps): bump github.com/klauspost/compress from 1.17.3 to 1.17.4 (#460)
- [`26f561b`](https://github.com/alexfalkowski/go-service/commit/26f561b4ac9aa4c3bb8c5d7bd96e2e2a29b57283) build(deps): add bin (#459)

## [v1.122.3](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.3) - 2023-12-01

- [`0c1f167`](https://github.com/alexfalkowski/go-service/commit/0c1f167142a359e5c55d6f26d797d259a569f6db) fix(deps): bump github.com/shirou/gopsutil/v3 from 3.23.10 to 3.23.11 (#458)

## [v1.122.2](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.2) - 2023-11-24

- [`fe28a86`](https://github.com/alexfalkowski/go-service/commit/fe28a86f962297c2709da47c5e9adcbe57008d77) fix(cmd): use built in version (#457)

## [v1.122.1](https://github.com/alexfalkowski/go-service/releases/tag/v1.122.1) - 2023-11-24

- [`995a879`](https://github.com/alexfalkowski/go-service/commit/995a87922c3628ee5a49998941dbd914fa9bfedf) fix: remove utc (#456)

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
