# Digital Assistant data format 0.2

`data/tasks.json` is saved as UTF-8 without a BOM. The root is an object with a text `version` equal to `"0.2"` and a `tasks` array. Each task has exactly three required fields: a unique nonempty text `id`, a nonempty text `title`, and a Boolean `done` (`true` or `false` without quotes). Array order determines display order; object field order does not affect meaning. The learning sample has three fictional records and exactly one `true`. Do not add real names or private information.

Check: open as UTF-8; parse JSON; check version and fields; confirm unique ids, nonempty titles, and Boolean done values; compare the backup with the original byte for byte (`cmp` on Unix/macOS) or using two SHA-256 values obtained locally from a trusted source. Valid JSON alone does not imply an exact copy. Do not replace your only original with a single backup.

Checkpoint: 8/10 now and 7/10 after seven days, including transfer to a new fictional record and exact comparison of the copy. The ten checks are listed in lesson 18.
