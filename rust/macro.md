# macro
宏编程的主要流程就是实现若干 From和 TryFrom.

## 展开宏
ref:
- [How do I see the expanded macro code that's causing my compile error?](https://stackoverflow.com/questions/28580386/how-do-i-see-the-expanded-macro-code-thats-causing-my-compile-error)

Starting with nightly-2021-07-28, one must pass -Zunpretty=expanded instead of -Zunstable-options --pretty=expanded, like this:

`rustc -Zunpretty=expanded test.rs`

> `the option `Z` is only accepted on the nightly compiler`


其他方法:
- [`cargo-expand`](https://github.com/dtolnay/cargo-expand), 是上述方法的wrapper.
- [rust playground](https://play.rust-lang.org/)的TOOLS->Expand macros